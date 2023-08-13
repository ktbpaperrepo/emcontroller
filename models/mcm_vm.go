package models

import (
	"fmt"
	"strings"
	"sync"

	"github.com/astaxie/beego"
)

func ListVMsAllClouds() ([]IaasVm, []error) {
	var allVms []IaasVm
	var errs []error

	// the slice in golang is not safe for concurrent read/write
	var allVmsMu sync.Mutex
	var errsMu sync.Mutex

	// List VMs in every cloud in parallel
	var wg sync.WaitGroup

	for _, cloud := range Clouds {
		wg.Add(1)
		go func(c Iaas) {
			defer wg.Done()
			vms, err := c.ListAllVMs()
			if err != nil {
				outErr := fmt.Errorf("List vms in cloud [%s] type [%s], error %w.", c.ShowName(), c.ShowType(), err)
				beego.Error(outErr)
				errsMu.Lock()
				errs = append(errs, outErr)
				errsMu.Unlock()
			}
			allVmsMu.Lock()
			allVms = append(allVms, vms...)
			allVmsMu.Unlock()
		}(cloud)
	}
	wg.Wait()
	return allVms, errs
}

// list all VMs with a name prefix
func ListVMsNamePrefix(prefix string) ([]IaasVm, error) {
	// get all VMs
	allVms, errs := ListVMsAllClouds()
	if len(errs) != 0 {
		sumErr := HandleErrSlice(errs)
		outErr := fmt.Errorf("cleanup auto-scheduling VMs, List VMs in all clouds, Error: %w", sumErr)
		beego.Error(outErr)
		return nil, outErr
	}

	// filter the VMs with the name prefix
	var outVms []IaasVm
	for _, vm := range allVms {
		if strings.HasPrefix(vm.Name, prefix) {
			outVms = append(outVms, vm)
		}
	}

	return outVms, nil
}

// remove a Virtual Machine with a appointed name from a list
func RemoveVmFromList(vmList *[]IaasVm, vmNameToRemove string) {
	vmIdx := FindIdxVmInList(*vmList, vmNameToRemove)
	if vmIdx >= 0 {
		*vmList = append((*vmList)[:vmIdx], (*vmList)[vmIdx+1:]...)
	}
}

// find the index of a Virtual Machine with a appointed name in a list. return -1 if not found.
func FindIdxVmInList(vmList []IaasVm, vmNameToFind string) int {
	for idx, vm := range vmList {
		if vm.Name == vmNameToFind {
			return idx
		}
	}
	return -1
}
