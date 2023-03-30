package controllers

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego"

	"emcontroller/models"
)

type CloudController struct {
	beego.Controller
}

type CloudInfo struct {
	Name      string
	Type      string
	WebUrl    string
	Resources models.ResourceStatus
}

func (c *CloudController) Get() {
	var cloudList []CloudInfo
	var errs []error

	// the slice in golang is not safe for concurrent read/write
	var cloudListMu sync.Mutex
	var errsMu sync.Mutex

	// Handle every cloud in parallel
	var wg sync.WaitGroup

	for _, cloud := range models.Clouds {
		wg.Add(1)
		go func(cloud models.Iaas) {
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

			cloudListMu.Lock()
			cloudList = append(cloudList, thisCloud)
			cloudListMu.Unlock()
		}(cloud)
	}
	wg.Wait()

	if len(errs) != 0 {
		sumErr := models.HandleErrSlice(errs)
		beego.Error(fmt.Sprintf("Get Clouds, Error: %s", sumErr.Error()))
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
		c.Data["errorMessage"] = sumErr.Error()
		c.TplName = "error.tpl"
		return
	}

	c.Data["cloudList"] = cloudList
	c.TplName = "cloud.tpl"
}

func (c *CloudController) GetSingleCloud() {
	cloudName := c.Ctx.Input.Param(":cloudName")
	cloud := models.Clouds[cloudName]
	resources, err := cloud.CheckResources()
	if err != nil {
		beego.Error(fmt.Sprintf("Check resources for cloud Name [%s] Type [%s], error: %s", cloud.ShowType(), cloud.ShowType(), err.Error()))
	}
	cloudInfo := CloudInfo{
		Name:      cloud.ShowName(),
		Type:      cloud.ShowType(),
		WebUrl:    cloud.ShowWebUrl(),
		Resources: resources,
	}

	// show all VMs of this cloud on the web
	vmList, err := cloud.ListAllVMs()
	if err != nil {
		beego.Error(fmt.Sprintf("List VMs in cloud Name [%s] Type [%s], error: %s", cloud.ShowType(), cloud.ShowType(), err.Error()))
	}

	c.Data["cloudInfo"] = cloudInfo
	c.Data["vmList"] = vmList

	c.TplName = "singleCloud.tpl"
}
