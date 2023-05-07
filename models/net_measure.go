package models

import (
	"fmt"
	"strings"
	"sync"

	"github.com/astaxie/beego"
)

const (
	// The names of McmNetTestVMs are "CloudName+suffix"
	mcmNetTestServerSuffix string = "-mcm-net-test-server"
	mcmNetTestClientSuffix string = "-mcm-net-test-client"

	mcmNetTestCpu  int = 2    // number of logical CPU cores
	mcmNetTestMem  int = 2048 // memory size unit: MB
	mcmNetTestDisk int = 10   // storage size unit: GB

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
			// 1. network test VMs exist
			if err := ensureTestPreC(c); err != nil {
				outErr := fmt.Errorf("Cannot ensure the network test preconditions for cloud [%s], error: [%w]", c.ShowName(), err)
				beego.Error(outErr)
				errsMu.Lock()
				errs = append(errs, outErr)
				errsMu.Unlock()
				return
			}
			// TODO
			// 2. network test VMs are managed by Kubernetes

			// TODO
			// 3. network test VMs have the NetTestTaint in Kubernetes

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
	if err := ensureVms(cloud); err != nil {
		return err
	}
	beego.Info(fmt.Sprintf("On cloud [%s], network test VMs exist.", cloud.ShowName()))

	beego.Info(fmt.Sprintf("On cloud [%s], The precondition is OK for network test.", cloud.ShowName()))
	return nil
}

// Ensure that the VMs for network test exists
func ensureVms(cloud Iaas) error {
	var netTestVmNames []string = getNetTestVMNames(cloud)
	vms, err := cloud.ListAllVMs()
	if err != nil {
		outErr := fmt.Errorf("List vms in cloud [%s] type [%s], error %s.", cloud.ShowName(), cloud.ShowType(), err)
		beego.Error(outErr)
		return outErr
	}

	// Here we cannot check each VM in parallel, because if we create more than one VM in proxmox, there will be problems.
	// It seems that proxmox does not support creating VMs in parallel.

	// Check each VM in serial
	for _, netTestVmName := range netTestVmNames {
		if !VmNameExist(netTestVmName, vms) {
			if _, err := cloud.CreateVM(netTestVmName, mcmNetTestCpu, mcmNetTestMem, mcmNetTestDisk); err != nil {
				outErr := fmt.Errorf("Create vm [%s] in cloud [%s] type [%s], error %w.", netTestVmName, cloud.ShowName(), cloud.ShowType(), err)
				beego.Error(outErr)
				return outErr
			}
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
	//		if !VmNameExist(name, vms) {
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

	return nil
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
