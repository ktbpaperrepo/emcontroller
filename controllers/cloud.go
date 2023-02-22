package controllers

import (
	"emcontroller/models"
	"fmt"
	"github.com/astaxie/beego"
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
	for _, cloud := range models.Clouds {
		var thisCloud CloudInfo
		thisCloud.Name = cloud.ShowName()
		thisCloud.Type = cloud.ShowType()
		thisCloud.WebUrl = cloud.ShowWebUrl()
		resources, err := cloud.CheckResources()
		if err != nil {
			beego.Error(fmt.Sprintf("Check resources for cloud Name [%s] Type [%s], error: %s", cloud.ShowType(), cloud.ShowType(), err.Error()))
		}
		thisCloud.Resources = resources
		cloudList = append(cloudList, thisCloud)
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
