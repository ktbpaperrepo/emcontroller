package controllers

import (
	"emcontroller/models"
	"fmt"
	"github.com/astaxie/beego"
)

type VmController struct {
	beego.Controller
}

func (c *VmController) DeleteVM() {
	cloudName := c.Ctx.Input.Param(":cloudName")
	vmID := c.Ctx.Input.Param(":vmID")

	beego.Info(fmt.Sprintf("Delete VM %s on cloud %s.", vmID, cloudName))
	err := models.Clouds[cloudName].DeleteVM(vmID)
	if err != nil {
		beego.Error(fmt.Sprintf("Delete VM %s on cloud %s, error: %s.", vmID, cloudName, err.Error()))
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	beego.Info(fmt.Sprintf("Successful! Delete VM %s on cloud %s.", vmID, cloudName))

	c.Ctx.ResponseWriter.WriteHeader(200)
}

func (c *VmController) CreateVM() {
	cloudName := c.Ctx.Input.Param(":cloudName")

	vmName := c.GetString("newVmName")
	vcpu, err := c.GetInt("newVmVCpu")
	if err != nil {
		beego.Error(fmt.Sprintf("read vcpu in int error: %s", err.Error()))
		return
	}
	ram, err := c.GetInt("newVmRam")
	if err != nil {
		beego.Error(fmt.Sprintf("read ram in int error: %s", err.Error()))
		return
	}
	storage, err := c.GetInt("newVmStorage")
	if err != nil {
		beego.Error(fmt.Sprintf("read storage in int error: %s", err.Error()))
		return
	}

	beego.Info(fmt.Sprintf("Start to create vm."))
	createdVM, err := models.Clouds[cloudName].CreateVM(vmName, vcpu, ram, storage)
	if err != nil {
		beego.Error(fmt.Sprintf("Create vm error %s.", err.Error()))
		return
	}
	beego.Info(fmt.Sprintf("Successful! Create vm:\n%+v\n", createdVM))

	c.Data["cloudName"] = cloudName
	c.TplName = "createVMSuccess.tpl"
}

// List VMs in all clouds
func (c *VmController) ListVMsAllClouds() {
	var allVms []models.IaasVm
	for _, cloud := range models.Clouds {
		vms, err := cloud.ListAllVMs()
		if err != nil {
			outErr := fmt.Errorf("List vms in cloud [%s] type [%s], error %s.", cloud.ShowName(), cloud.ShowType(), err)
			beego.Error(outErr)
			c.Data["errorMessage"] = outErr.Error()
			c.TplName = "error.tpl"
			return
		}
		allVms = append(allVms, vms...)
	}

	c.Data["allVms"] = allVms
	c.TplName = "vm.tpl"
}
