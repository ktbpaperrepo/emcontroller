package model

import (
	"fmt"
	"math"
	"strings"

	"github.com/astaxie/beego"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"emcontroller/models"
)

type Cloud struct {
	Name      string                         `json:"name"`
	Type      string                         `json:"type"`
	Resources models.ResourceStatus          `json:"resources"` // used and all resources of this cloud. Here we start with struct defined in "models" package, and in the future if we find that this cannot meet the needs here, we can define new structs.
	NetState  map[string]models.NetworkState `json:"netState"`  // the network state from this cloud to every cloud
	K8sNodes  []K8sNode                      `json:"k8sNodes"`  // all existing Kubernetes nodes whose VMs are on this cloud
}

// the set of all cloud types that support creating new VMs when auto-scheduling
var typesCanCreateNewVM map[string]struct{} = map[string]struct{}{
	models.ProxmoxIaas: struct{}{},
}

// Not all cloud types support creating new VMs.
// For example, CLAAUDIA does not allow users to create flavors in Openstack, so if we want to create new VMs in auto-scheduling, this will be more complicated, so we do not support creating new VMs in auto-scheduling.
func (c Cloud) SupportCreateNewVM() bool {
	if _, exist := typesCanCreateNewVM[c.Type]; exist {
		return true
	}
	return false
}

// According to the inpur resource percentage, this function can generate the information of the VM to create.
func (c Cloud) GetInfoVmToCreate(resPct float64) K8sNode {
	return K8sNode{
		Name: c.getNameVmToCreate(),
		ResidualResources: GenericResources{
			CpuCore: math.Floor(resPct * c.Resources.Limit.VCpu),
			Memory:  math.Floor(resPct * c.Resources.Limit.Ram),
			Storage: math.Floor(resPct * c.Resources.Limit.Storage),
		},
	}
}

// auto-schedule vms should have special prefixes
func (c Cloud) getNameVmToCreate() string {
	var vmName string
OUTLOOP:
	for i := 0; i < math.MaxInt; i++ {
		vmName = fmt.Sprintf("%s%s-%d", ASVmNamePrefix, strings.ToLower(c.Name), i)
		for _, existingVm := range c.K8sNodes {
			if existingVm.Name == vmName {
				continue OUTLOOP // the intended vmName is already used, we continue to try the next possible name.
			}
		}
		return vmName // the intended vmName is not used, so we can use it.
	}
	panic(fmt.Sprintf("Cloud [%s], all available auto-schedule vm names are used up. There are [%d] existing auto-schedule vms on this cloud.", c.Name, len(c.K8sNodes)))
}

func CloudMapCopy(src map[string]Cloud) map[string]Cloud {
	var dst map[string]Cloud = make(map[string]Cloud)
	for name, cloud := range src {
		dst[name] = cloud
	}
	return dst
}

// generate a group of Cloud from a group of models.Iaas
func GenerateClouds(inputClouds map[string]models.Iaas) (map[string]Cloud, error) {
	var outputClouds map[string]Cloud = make(map[string]Cloud)

	netStates, err := models.GetNetState()
	if err != nil {
		outErr := fmt.Errorf("Check network state from MySQL Error: %w", err)
		beego.Error(outErr)
		return nil, outErr
	}

	allK8sNodes, err := models.ListNodes(metav1.ListOptions{})
	if err != nil {
		outErr := fmt.Errorf("List Kubernetes Nodes Error: %w", err)
		beego.Error(outErr)
		return nil, outErr
	}

	for name, inCloud := range inputClouds {
		outputClouds[name], err = GenerateOneCloud(inCloud, netStates[name], allK8sNodes)
		if err != nil {
			outErr := fmt.Errorf("generate the Cloud  [%s] from models.Iaas, Error: %w", name, err)
			beego.Error(outErr)
			return nil, outErr
		}
	}

	return outputClouds, nil
}

// generate the Cloud from models.Iaas, the network states of this cloud, and all Kubernetes Nodes
func GenerateOneCloud(inCloud models.Iaas, cloudNetStates map[string]models.NetworkState, allK8sNodes []apiv1.Node) (Cloud, error) {

	resources, err := inCloud.CheckResources()
	if err != nil {
		outErr := fmt.Errorf("Check resources for cloud Name [%s] Type [%s], error: %w", inCloud.ShowType(), inCloud.ShowType(), err)
		beego.Error(outErr)
		return Cloud{}, outErr
	}

	k8sNodesOnCloud, err := getK8sNodesOnCloud(inCloud, allK8sNodes)
	if err != nil {
		outErr := fmt.Errorf("Get Kubernetes Nodes on Cloud [%s] Type [%s], error: %w", inCloud.ShowType(), inCloud.ShowType(), err)
		beego.Error(outErr)
		return Cloud{}, outErr
	}

	var outCloud Cloud = Cloud{
		Name:      inCloud.ShowName(),
		Type:      inCloud.ShowType(),
		NetState:  cloudNetStates,
		Resources: resources,
		K8sNodes:  k8sNodesOnCloud,
	}

	return outCloud, nil
}

func getK8sNodesOnCloud(cloud models.Iaas, allK8sNodes []apiv1.Node) ([]K8sNode, error) {
	var k8sNodes []K8sNode

	// We should find the VMs that meet the 2 conditions:
	// 1. They are on this cloud;
	// 2. They are Kubernetes Nodes;

	vms, err := cloud.ListAllVMs()
	if err != nil {
		outErr := fmt.Errorf("List vms in Cloud [%s] Type [%s], error: %w", cloud.ShowType(), cloud.ShowType(), err)
		beego.Error(outErr)
		return []K8sNode{}, outErr
	}

	for _, vm := range vms {
		// If we cannot get the IP of a VM, we do not consider scheduling applications to it.
		if len(vm.IPs) == 0 {
			continue
		}
		// We do not schedule applications to the VMs that are not created by multi-cloud manager
		if createdByMcm, err := cloud.IsCreatedByMcm(vm.ID); err != nil {
			printErr := fmt.Errorf("Cloud [%s] Type [%s] check IsCreatedByMcm [%s/%s], error: %w. So we jump it.", cloud.ShowType(), cloud.ShowType(), vm.Name, vm.ID, err)
			beego.Error(printErr)
			continue
		} else if !createdByMcm {
			continue
		}

		for _, node := range allK8sNodes {
			if len(node.Spec.Taints) > 0 {
				// If len(node.Spec.Taints) > 0, do not auto schedule applications to this node. This can ensure the auto-schedule will not put applications on network state test VMs, and also provide a way to users to avoid auto-schedule from putting applications on their VMs.
				continue
			}

			// If the IPs match, we find one Kubernetes node on this cloud
			if vm.IPs[0] == models.GetNodeInternalIp(node) {
				// Get available resources of this VM
				residualCpuCore := models.CalcVmAvailVcpu(vm.VCpu)
				residualRamMiB := models.CalcVmAvailRamMiB(vm.Ram)
				residualStorGiB := models.CalcVmAvailStorGiB(vm.Storage)

				// subtract the resources occupied by pods
				podsOnNode, err := models.ListPodsOnNode(models.KubernetesNamespace, node.Name)
				if err != nil {
					outErr := fmt.Errorf("List pods on Kubernetes node [%s], error: %w", node.Name, err)
					beego.Error(outErr)
					return []K8sNode{}, outErr
				}
				for _, pod := range podsOnNode {
					occupied := GetResOccupiedByPod(pod)
					residualCpuCore -= occupied.CpuCore
					residualRamMiB -= occupied.Memory
					residualStorGiB -= occupied.Storage
				}

				// handle possible negative results
				if residualCpuCore < 0 {
					residualCpuCore = 0
				}
				if residualRamMiB < 0 {
					residualRamMiB = 0
				}
				if residualStorGiB < 0 {
					residualStorGiB = 0
				}

				// we put the information needed by auto-scheduling to the K8sNode structure
				var thisNode K8sNode
				thisNode.Name = node.Name
				thisNode.ResidualResources.CpuCore = residualCpuCore
				thisNode.ResidualResources.Memory = residualRamMiB
				thisNode.ResidualResources.Storage = residualStorGiB
				k8sNodes = append(k8sNodes, thisNode)

				break // When we find a match, we can break to search for the next VM.
			}
		}
	}

	return k8sNodes, nil
}
