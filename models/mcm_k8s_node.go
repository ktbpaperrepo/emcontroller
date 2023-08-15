package models

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/astaxie/beego"
	apiv1 "k8s.io/api/core/v1"
)

type K8sNodeInfo struct {
	Name           string     `json:"name"`
	IP             string     `json:"ip"`
	Status         string     `json:"status"`
	TotalResources K8sNodeRes `json:"totalResources"` // the total available resources of this Kubernetes node.
	UsedResources  K8sNodeRes `json:"UsedResources"`  // the resources used by all Kubernetes pods running on this node.
}

type K8sNodeRes struct {
	CpuCore float64 `json:"cpuCore"` // number of CPU logical cores
	Memory  float64 `json:"memory"`  // unit Mebibyte (MiB)
	Storage float64 `json:"storage"` // unit Gibibyte (GiB)
}

func GetResOccupiedByPod(pod apiv1.Pod) K8sNodeRes {
	var occupied K8sNodeRes
	for _, container := range pod.Spec.Containers {
		// convert unit to 1 CPU
		occupied.CpuCore += float64(container.Resources.Requests.Cpu().MilliValue()) / 1000
		// convert unit to Mi
		occupied.Memory += float64(container.Resources.Requests.Memory().Value()) / 1024 / 1024
		// convert unit to Gi
		occupied.Storage += float64(container.Resources.Requests.StorageEphemeral().Value()) / 1024 / 1024 / 1024
	}
	return occupied
}

func ListK8sNodes() []K8sNodeInfo {
	// TODO: This code does not work, I do not know the reason.
	//K8sMasterSelector := labels.NewSelector()
	//K8sMasterReq, err := labels.NewRequirement(models.K8sMasterNodeRole, selection.NotEquals, []string{""})
	//if err != nil {
	//	beego.Error(fmt.Sprintf("Construct Kubernetes Master requirement, error: %s", err.Error()))
	//}
	//K8sMasterSelector.Add(*K8sMasterReq)
	//beego.Info(fmt.Sprintf("List nodes with selector: %v", K8sMasterSelector))
	//beego.Info(fmt.Sprintf("List nodes with selector: %s", K8sMasterSelector.String()))
	//nodes, err := models.ListNodes(metav1.ListOptions{LabelSelector: K8sMasterSelector.String()})

	// get all VMs, for calculate the total available resources of node.
	allVms, errs := ListVMsAllClouds()
	if len(errs) != 0 {
		sumErr := HandleErrSlice(errs)
		outErr := fmt.Errorf("Get Kubernetes Nodes, List VMs in all clouds, Error: %w", sumErr)
		beego.Error(outErr)
	}

	nodes, err := ListNodes(metav1.ListOptions{})
	if err != nil {
		beego.Error(fmt.Sprintf("List Kubernetes nodes, error: %s", err.Error()))
	}

	selectorControlPlane := labels.SelectorFromSet(labels.Set(map[string]string{
		K8sMasterNodeRole: "",
	}))

	var k8sNodeList []K8sNodeInfo
	for _, node := range nodes {
		if selectorControlPlane.Matches(labels.Set(node.Labels)) {
			beego.Info(fmt.Sprintf("node %s is a Master node, so we do not show it.", node.Name))
			continue
		}

		// I think there is no need to hide these network test nodes, or else there will be other troubles.
		//if models.NodeHasTaint(&node, models.NetTestTaint) {
		//	beego.Info(fmt.Sprintf("node %s is a network performance test node, so we do not show it.", node.Name))
		//	continue
		//}

		var thisOutNode K8sNodeInfo
		thisOutNode.Name = node.Name
		thisOutNode.IP = GetNodeInternalIp(node)
		thisOutNode.Status = ExtractNodeStatus(node)

		// calculate the resources occupied by pods
		var resInUse K8sNodeRes
		podsOnNode, err := ListPodsOnNode(KubernetesNamespace, node.Name)
		if err != nil {
			outErr := fmt.Errorf("List pods on Kubernetes node [%s], error: %w", node.Name, err)
			beego.Error(outErr)
			// when there is an error, we set the occupied resources as -1.
			resInUse.CpuCore = -1
			resInUse.Memory = -1
			resInUse.Storage = -1
		} else {
			for _, pod := range podsOnNode {
				usedByPod := GetResOccupiedByPod(pod)

				resInUse.CpuCore += usedByPod.CpuCore
				resInUse.Memory += usedByPod.Memory
				resInUse.Storage += usedByPod.Storage
			}
		}
		thisOutNode.UsedResources = resInUse

		// calculate the total available resources of this node
		var availRes K8sNodeRes
		// if the above models.ListVMsAllClouds() has errors, or the VM name and node name are different, we will not find a VM for this node. In this case, we will set the avail resources as -1.
		var found bool = false
		for _, vm := range allVms {
			// find the VM of this Kubernetes node
			if len(vm.IPs) == 0 {
				continue
			}
			if vm.IPs[0] == thisOutNode.IP && vm.Name == thisOutNode.Name {
				found = true
				availRes.CpuCore = CalcVmAvailVcpu(vm.VCpu)
				availRes.Memory = CalcVmAvailRamMiB(vm.Ram)
				availRes.Storage = CalcVmAvailStorGiB(vm.Storage)
				break
			}
		}
		if !found {
			availRes.CpuCore = -1
			availRes.Memory = -1
			availRes.Storage = -1
		}
		thisOutNode.TotalResources = availRes

		k8sNodeList = append(k8sNodeList, thisOutNode)
	}

	return k8sNodeList
}

// create the new VMs and add them to Kubernetes
func AddNewVms(vmsToCreate []IaasVm) ([]IaasVm, error) {
	beego.Info(fmt.Sprintf("Create new VMs [%s].", JsonString(vmsToCreate)))
	createdVms, err := CreateVms(vmsToCreate)
	if err != nil {
		outErr := fmt.Errorf("Create new VMs [%s], Error: [%w]", JsonString(vmsToCreate), err)
		beego.Error(outErr)
		return []IaasVm{}, outErr
	}

	beego.Info(fmt.Sprintf("Add new VMs [%s] to Kubernetes cluster.", JsonString(createdVms)))
	errs := AddNodes(createdVms)
	if len(errs) != 0 {
		sumErr := HandleErrSlice(errs)
		outErr := fmt.Errorf("Add new VMs [%s] to Kubernetes cluster, Error: [%w]", JsonString(createdVms), sumErr)
		beego.Error(outErr)
		return createdVms, outErr
	}

	return createdVms, nil
}

// remove a Kubernetes node with a appointed name from a list
func RemoveNodeFromList(nodeList *[]apiv1.Node, nodeNameToRemove string) {
	nodeIdx := FindIdxNodeInList(*nodeList, nodeNameToRemove)
	if nodeIdx >= 0 {
		*nodeList = append((*nodeList)[:nodeIdx], (*nodeList)[nodeIdx+1:]...)
	}
}

// find the index of a Kubernetes node with a appointed name in a list. return -1 if not found.
func FindIdxNodeInList(nodeList []apiv1.Node, nodeNameToFind string) int {
	for idx, node := range nodeList {
		if node.Name == nodeNameToFind {
			return idx
		}
	}
	return -1
}
