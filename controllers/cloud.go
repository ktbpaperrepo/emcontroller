package controllers

import (
	"fmt"
	"github.com/astaxie/beego"

	"emcontroller/models"
)

type CloudController struct {
	beego.Controller
}

func (c *CloudController) Get() {
	cloudList, errs := models.ListClouds()
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

	cloudInfo, vmList, _, _ := models.GetCloud(cloudName)

	c.Data["cloudInfo"] = cloudInfo
	c.Data["vmList"] = vmList

	c.TplName = "singleCloud.tpl"
}
