package models

import (
	"fmt"
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
