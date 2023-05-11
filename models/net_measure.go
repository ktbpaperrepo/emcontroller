package models

import (
	"fmt"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// The names of McmNetTestVMs are "CloudName+suffix"
	mcmNetTestServerSuffix string = "-mcm-net-test-server"
	mcmNetTestClientSuffix string = "-mcm-net-test-client"

	mcmNetTestCpu  int = 2    // number of logical CPU cores
	mcmNetTestMem  int = 2048 // memory size unit: MB
	mcmNetTestDisk int = 10   // storage size unit: GB

	taintValue   string                   = "net-test"
	taintEffect  apiv1.TaintEffect        = apiv1.TaintEffectNoSchedule
	tolerationOp apiv1.TolerationOperator = apiv1.TolerationOpEqual
)

// The function to measure network performance between every two clouds
// This function should be executed every time period
func MeasNetPerf() {
	beego.Info("Start to ensure the network test preconditions for each cloud.")

	// Ensure the network test preconditions for each cloud in parallel
	var wg sync.WaitGroup
	var errsMu sync.Mutex // the slice in golang is not safe for concurrent read/write
	var errs []error
	for name, cloud := range Clouds {
		beego.Info(fmt.Sprintf("Ensure the network test preconditions for cloud %s", name))
		wg.Add(1)
		go func(c Iaas) {
			defer wg.Done()
			if err := ensureTestPreC(c); err != nil {
				outErr := fmt.Errorf("Cannot ensure the network test preconditions for cloud [%s], error: [%w]", c.ShowName(), err)
				beego.Error(outErr)
				errsMu.Lock()
				errs = append(errs, outErr)
				errsMu.Unlock()
				return
			}
		}(cloud)
	}
	wg.Wait()

	if len(errs) != 0 {
		sumErr := HandleErrSlice(errs)
		outErr := fmt.Errorf("Cannot ensure the network test preconditions, Error: %w", sumErr)
		beego.Error(outErr)
		return
	}

	beego.Info("Finish ensuring the network test preconditions for each cloud.")

	// TODO
	// Execute network test between every two clouds
	beego.Info("Start to measure network performance between every two clouds.")

	beego.Info("Finish measuring network performance between every two clouds.")
}

// Ensure the preconditions for network test, including VMs, K8s nodes, and K8s taints.
func ensureTestPreC(cloud Iaas) error {
	// 1. network test VMs exist
	netTestVmMap, err := ensureVms(cloud)
	if err != nil {
		outErr := fmt.Errorf("Cannot ensure that network test VMs exist, Error: %w", err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("On cloud [%s], network test VMs exist.", cloud.ShowName()))

	// 2. network test VMs are managed by Kubernetes
	err = ensureK8sMg(netTestVmMap)
	if err != nil {
		outErr := fmt.Errorf("Cannot ensure that network test VMs are managed by Kubernetes, Error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	// 3. network test VMs have the NetTestTaint in Kubernetes
	err = ensureTaints(netTestVmMap)
	if err != nil {
		outErr := fmt.Errorf("Cannot ensure that VMs have the NetTestTaint in Kubernetes, Error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	beego.Info(fmt.Sprintf("On cloud [%s], The precondition is OK for network test.", cloud.ShowName()))
	return nil
}

// Ensure that the VMs for network test exists
func ensureVms(cloud Iaas) (map[string]*IaasVm, error) {
	var netTestVmNames []string = getNetTestVMNames(cloud)
	vms, err := cloud.ListAllVMs()
	if err != nil {
		outErr := fmt.Errorf("List vms in cloud [%s] type [%s], error %s.", cloud.ShowName(), cloud.ShowType(), err)
		beego.Error(outErr)
		return nil, outErr
	}

	// We should use the information of the network performance test VMs in the following steps, so we use a map to return the VMs information
	var netTestVmMap map[string]*IaasVm = make(map[string]*IaasVm)

	// Here we cannot check each VM in parallel, because if we create more than one VM in proxmox, there will be problems.
	// It seems that proxmox does not support creating VMs in parallel.

	// Check each VM in serial
	for _, netTestVmName := range netTestVmNames {
		if vm, found := FindVm(netTestVmName, vms); found { // If a test vm exists, we return it.
			netTestVmMap[netTestVmName] = vm
		} else { // If a test vm does not exist, we create and return a new one.
			newVm, err := cloud.CreateVM(netTestVmName, mcmNetTestCpu, mcmNetTestMem, mcmNetTestDisk)
			if err != nil {
				outErr := fmt.Errorf("Create vm [%s] in cloud [%s] type [%s], error %w.", netTestVmName, cloud.ShowName(), cloud.ShowType(), err)
				beego.Error(outErr)
				return nil, outErr
			}
			netTestVmMap[netTestVmName] = newVm
		}
	}

	//// Check each VM in parallel
	//var wg sync.WaitGroup
	//var errsMu sync.Mutex // the slice in golang is not safe for concurrent read/write
	//var errs []error
	//for _, netTestVmName := range netTestVmNames {
	//	wg.Add(1)
	//	go func(name string) {
	//		defer wg.Done()
	//		if !FindVm(name, vms) {
	//			if _, err := cloud.CreateVM(name, mcmNetTestCpu, mcmNetTestMem, mcmNetTestDisk); err != nil {
	//				outErr := fmt.Errorf("Create vm [%s] in cloud [%s] type [%s], error %w.", name, cloud.ShowName(), cloud.ShowType(), err)
	//				beego.Error(outErr)
	//				errsMu.Lock()
	//				errs = append(errs, outErr)
	//				errsMu.Unlock()
	//			}
	//		}
	//	}(netTestVmName)
	//}
	//wg.Wait()
	//
	//if len(errs) != 0 {
	//	sumErr := HandleErrSlice(errs)
	//	outErr := fmt.Errorf("Ensure network test VMs in cloud [%s], Error: %w", cloud.ShowName(), sumErr)
	//	beego.Error(outErr)
	//	return outErr
	//}

	return netTestVmMap, nil
}

func getNetTestVMNames(cloud Iaas) []string {
	return []string{
		getNetTestServerName(cloud),
		getNetTestClientName(cloud),
	}
}

func getNetTestServerName(cloud Iaas) string {
	return strings.ToLower(cloud.ShowName()) + mcmNetTestServerSuffix
}

func getNetTestClientName(cloud Iaas) string {
	return strings.ToLower(cloud.ShowName()) + mcmNetTestClientSuffix
}

// Ensure that the VMs for network test are managed by Kubernetes
func ensureK8sMg(vmMap map[string]*IaasVm) error {
	for vmName, vm := range vmMap {
		_, err := GetNode(vmName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				beego.Info(fmt.Sprintf("Node [%s] is not managed by Kubernetes. Then we add it.", vmName))

				joinCmd, err := GetJoinCmd()
				if err != nil {
					outErr := fmt.Errorf("Cannot add node [%s] to Kubernetes cluster, GetJoinCmd error: %w", vmName, err)
					beego.Error(outErr)
					return outErr
				}
				if err := AddNode(*vm, joinCmd); err != nil {
					outErr := fmt.Errorf("Cannot add node [%s] to Kubernetes cluster, error: %w.", vmName, err)
					beego.Error(outErr)
					return outErr
				}

				beego.Info(fmt.Sprintf("Node [%s] is already added to Kubernetes cluster.", vmName))
			} else {
				outErr := fmt.Errorf("Get node [%s] in Kubernetes, error %w.", vmName, err)
				beego.Error(outErr)
				return outErr
			}
		} else {
			beego.Info(fmt.Sprintf("Node [%s] is already managed by Kubernetes.", vmName))
		}
	}
	return nil
}

// Ensure that the nodes of VMs in Kubernetes have the network test taints.
func ensureTaints(vmMap map[string]*IaasVm) error {
	netTestTaint := &apiv1.Taint{
		Key:    McmKey,
		Effect: taintEffect,
		Value:  taintValue,
	}
	for nodeName, _ := range vmMap {
		if err := TaintNode(nodeName, netTestTaint); err != nil {
			outErr := fmt.Errorf("ensureTaints, node [%s], error %w.", nodeName, err)
			beego.Error(outErr)
			return outErr
		}
	}

	return nil
}
