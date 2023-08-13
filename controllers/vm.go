package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/astaxie/beego"

	"emcontroller/models"
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

// delete multiple VMs
// test command:
// curl -i -X DELETE -H Content-Type:application/json http://localhost:20000/vm -d '[{"name":"auto-sched-nokia8-0","cloud":"NOKIA8","id":"106"},{"name":"auto-sched-nokia8-1","cloud":"NOKIA8","id":"107"},{"name":"auto-sched-nokia7-0","cloud":"NOKIA7","id":"104"}]'
func (c *VmController) DeleteVMs() {
	var vms []models.IaasVm
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vms); err != nil {
		outErr := fmt.Errorf("json.Unmarshal the vms in RequestBody, error: %w", err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		//c.Ctx.WriteString(outErr.Error())
		if result, err := c.Ctx.ResponseWriter.Write([]byte(outErr.Error())); err != nil {
			beego.Error(fmt.Sprintf("Write Error to response, error: %s, result: %d", err.Error(), result))
		}
		return
	}

	beego.Info(fmt.Sprintf("Delete VMs %v.", vms))

	// Use the parsed vms as the input information to delete VMs
	if errs := models.DeleteBatchVms(vms); len(errs) != 0 {
		outErr := models.HandleErrSlice(errs)
		beego.Error(fmt.Sprintf("DeleteBatchVms Error: %s", outErr.Error()))
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString(outErr.Error())
		return
	}

	beego.Info(fmt.Sprintf("Successful! Delete VMs %v.", vms))

	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
}

// test command:
// curl -i -X GET http://localhost:20000/cloud/NOKIA7/vm/102
// curl -i -X GET http://localhost:20000/cloud/CLAAUDIAweifan/vm/8117edb9-0fdc-4334-a1d9-51779e79f377
func (c *VmController) GetVM() {
	cloudName := c.Ctx.Input.Param(":cloudName")
	vmID := c.Ctx.Input.Param(":vmID")

	vm, err := models.Clouds[cloudName].GetVM(vmID)
	if err != nil {
		beego.Error(fmt.Sprintf("Get VM %s on cloud %s, error: %s.", vmID, cloudName, err.Error()))
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	c.Ctx.Output.Status = http.StatusOK
	c.Data["json"] = vm
	c.ServeJSON()
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
	allVms, errs := models.ListVMsAllClouds()
	if len(errs) != 0 {
		sumErr := models.HandleErrSlice(errs)
		beego.Error(fmt.Sprintf("List VMs in all clouds, Error: %s", sumErr.Error()))
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
		c.Data["errorMessage"] = sumErr.Error()
		c.TplName = "error.tpl"
		return
	}

	c.Data["allVms"] = allVms
	c.TplName = "vm.tpl"
}

func (c *VmController) NewVms() {
	c.TplName = "newVms.tpl"
}

func (c *VmController) DoNewVms() {
	contentType := c.Ctx.Request.Header.Get("Content-Type")
	beego.Info(fmt.Sprintf("The header \"Content-Type\" is [%s]", contentType))

	switch {
	case strings.Contains(strings.ToLower(contentType), JsonContentType):
		beego.Info(fmt.Sprintf("The input body should be json"))
		c.DoNewVmsJson()
	default:
		beego.Info(fmt.Sprintf("The input body should be form"))
		c.DoNewVmsForm()
	}
}

func (c *VmController) DoNewVmsForm() {
	vmNum, err := c.GetInt("newVmNumber")
	if err != nil {
		outErr := fmt.Errorf("Get newVmNumber error: %w", err)
		beego.Error(outErr)
		c.Data["errorMessage"] = outErr.Error()
		c.TplName = "error.tpl"
		return
	}
	beego.Info(fmt.Sprintf("%d vms need to be created", vmNum))

	// prepare the information of the vms to add
	vms := make([]models.IaasVm, vmNum, vmNum)

	for i := 0; i < vmNum; i++ {
		vms[i].Name = c.GetString(fmt.Sprintf("vm%dName", i))
		vms[i].Cloud = c.GetString(fmt.Sprintf("vm%dCloudName", i))
		if vms[i].VCpu, err = c.GetFloat(fmt.Sprintf("vm%dVCpu", i)); err != nil {
			outErr := fmt.Errorf("Get vms[%d].VCpu, error: %w", i, err)
			beego.Error(outErr)
			c.Data["errorMessage"] = outErr.Error()
			c.TplName = "error.tpl"
			return
		}
		if vms[i].Ram, err = c.GetFloat(fmt.Sprintf("vm%dRam", i)); err != nil {
			outErr := fmt.Errorf("Get vms[%d].Ram, error: %w", i, err)
			beego.Error(outErr)
			c.Data["errorMessage"] = outErr.Error()
			c.TplName = "error.tpl"
			return
		}
		if vms[i].Storage, err = c.GetFloat(fmt.Sprintf("vm%dStorage", i)); err != nil {
			outErr := fmt.Errorf("Get vms[%d].Storage, error: %w", err)
			beego.Error(outErr)
			c.Data["errorMessage"] = outErr.Error()
			c.TplName = "error.tpl"
			return
		}
	}

	logContent := "VMs to create:"
	for i := 0; i < vmNum; i++ {
		logContent += fmt.Sprintf("\n\r%d. Name: %s\tCloud: %s\tVCpu: %f\tRam: %f\tStorage: %f", i+1, vms[i].Name, vms[i].Cloud, vms[i].VCpu, vms[i].Ram, vms[i].Storage)
	}
	beego.Info(logContent)

	vmsJson, err := json.Marshal(vms)
	if err != nil {
		outErr := fmt.Errorf("json Marshal this: %v, error: %w", vms, err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
		c.Data["errorMessage"] = outErr.Error()
		c.TplName = "error.tpl"
		return
	}
	beego.Info(fmt.Sprintf("VMs json is\n%s", string(vmsJson)))

	// create vms
	if _, err = models.CreateVms(vms); err != nil {
		outErr := fmt.Errorf("DoNewVms error: %w", err)
		beego.Error(outErr)
		c.Data["errorMessage"] = outErr.Error()
		c.TplName = "error.tpl"
		return
	}

	c.TplName = "newVmsSuccess.tpl"
}

// test command:
// curl -i -X POST -H Content-Type:application/json -d '[{"name":"cnode1","vcpu":4,"ram":32768,"storage":400,"cloud":"CLAAUDIAweifan"},{"name":"cnode2","vcpu":4,"ram":32768,"storage":400,"cloud":"CLAAUDIAweifan"},{"name":"hpe1","vcpu":4,"ram":8192,"storage":100,"cloud":"HPE1"},{"name":"nokia7","vcpu":4,"ram":8192,"storage":100,"cloud":"NOKIA7"}]' http://localhost:20000/vm/doNew
func (c *VmController) DoNewVmsJson() {
	var vms []models.IaasVm
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vms); err != nil {
		outErr := fmt.Errorf("json.Unmarshal the vms in RequestBody, error: %w", err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		//c.Ctx.WriteString(outErr.Error())
		if result, err := c.Ctx.ResponseWriter.Write([]byte(outErr.Error())); err != nil {
			beego.Error(fmt.Sprintf("Write Error to response, error: %s, result: %d", err.Error(), result))
		}
		return
	}

	beego.Info(fmt.Sprintf("From json input, we successfully parsed vms [%v]", vms))

	// Use the parsed vms to create VMs
	outVms, err := models.CreateVms(vms)
	if err != nil {
		outErr := fmt.Errorf("Create VMs %v, error: %w", vms, err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString(outErr.Error())
		return
	}

	// If the user uses Json as the HTTP input body, he may need the information of the VMs such as IP addresses and VM IDs in the response body, so we put this information into the response body.

	//c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Ctx.Output.Status = http.StatusCreated
	c.Data["json"] = outVms
	c.ServeJSON()
}
