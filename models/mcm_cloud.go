package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"sync"
)

type CloudInfo struct {
	Name      string
	Type      string
	WebUrl    string
	Resources ResourceStatus
}

func ListClouds() ([]CloudInfo, []error) {
	var cloudList []CloudInfo
	var errs []error

	// the slice in golang is not safe for concurrent read/write
	var cloudListMu sync.Mutex
	var errsMu sync.Mutex

	// Handle every cloud in parallel
	var wg sync.WaitGroup

	for _, cloud := range Clouds {
		wg.Add(1)
		go func(cloud Iaas) {
			defer wg.Done()
			var thisCloud CloudInfo
			thisCloud.Name = cloud.ShowName()
			thisCloud.Type = cloud.ShowType()
			thisCloud.WebUrl = cloud.ShowWebUrl()
			resources, err := cloud.CheckResources()
			if err != nil {
				outErr := fmt.Errorf("Check resources for cloud Name [%s] Type [%s], error: %w", cloud.ShowType(), cloud.ShowType(), err)
				beego.Error(outErr)

				errsMu.Lock()
				errs = append(errs, outErr)
				errsMu.Unlock()
			}
			thisCloud.Resources = resources

			beego.Info(fmt.Sprintf("Cloud [%s], type [%s], resources [%+v]", thisCloud.Name, thisCloud.Type, thisCloud.Resources))
			cloudListMu.Lock()
			cloudList = append(cloudList, thisCloud)
			cloudListMu.Unlock()
			beego.Info(fmt.Sprintf("Cloud [%s], type [%s] added to cloudList", thisCloud.Name, thisCloud.Type))
		}(cloud)
	}
	wg.Wait()

	return cloudList, errs
}

func GetCloud(cloudName string) (CloudInfo, []IaasVm, error, error) {
	cloud := Clouds[cloudName]
	resources, errRes := cloud.CheckResources()
	if errRes != nil {
		beego.Error(fmt.Sprintf("Check resources for cloud Name [%s] Type [%s], error: %s", cloud.ShowType(), cloud.ShowType(), errRes.Error()))
	}
	cloudInfo := CloudInfo{
		Name:      cloud.ShowName(),
		Type:      cloud.ShowType(),
		WebUrl:    cloud.ShowWebUrl(),
		Resources: resources,
	}

	// show all VMs of this cloud on the web
	vmList, errVMs := cloud.ListAllVMs()
	if errVMs != nil {
		beego.Error(fmt.Sprintf("List VMs in cloud Name [%s] Type [%s], error: %s", cloud.ShowType(), cloud.ShowType(), errVMs.Error()))
	}

	return cloudInfo, vmList, errRes, errVMs
}
