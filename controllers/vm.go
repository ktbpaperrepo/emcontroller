package controllers

import (
	"fmt"
	"sync"

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
	var errs []error

	// the slice in golang is not safe for concurrent read/write
	var allVmsMu sync.Mutex
	var errsMu sync.Mutex

	// List VMs in every cloud in parallel
	var wg sync.WaitGroup

	for _, cloud := range models.Clouds {
		wg.Add(1)
		go func(c models.Iaas) {
			defer wg.Done()
			vms, err := c.ListAllVMs()
			if err != nil {
				outErr := fmt.Errorf("List vms in cloud [%s] type [%s], error %s.", c.ShowName(), c.ShowType(), err)
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

	// create vms
	if err = models.CreateVms(vms); err != nil {
		outErr := fmt.Errorf("DoNewVms error: %w", err)
		beego.Error(outErr)
		c.Data["errorMessage"] = outErr.Error()
		c.TplName = "error.tpl"
		return
	}

	c.TplName = "newVmsSuccess.tpl"
}
