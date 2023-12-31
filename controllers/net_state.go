package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"emcontroller/models"
	"github.com/astaxie/beego"
)

type NetStateController struct {
	beego.Controller
}

func (c *NetStateController) Get() {
	contentType := c.Ctx.Request.Header.Get("Content-Type")
	beego.Info(fmt.Sprintf("The header \"Content-Type\" is [%s]", contentType))

	switch {
	case strings.Contains(strings.ToLower(contentType), JsonContentType):
		beego.Info(fmt.Sprintf("The input body should be json"))
		c.GetJson()
	default:
		beego.Info(fmt.Sprintf("The input body should be form"))
		c.GetForm()
	}
}

// test command:
// curl -i -X GET -H Content-Type:application/json http://172.27.15.31:20000/netState
func (c *NetStateController) GetJson() {
	if !models.NetTestFuncOn {
		beego.Info(models.NetTestFuncOffMsg)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusServiceUnavailable)
		c.Ctx.WriteString(models.NetTestFuncOffMsg)
		return
	}

	netState, err := models.GetNetState()
	if err != nil {
		outErr := fmt.Errorf("Check network state from MySQL Error: %w", err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString(outErr.Error())
		return
	}

	//c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Ctx.Output.Status = http.StatusOK
	c.Data["json"] = netState
	c.ServeJSON()
}

func (c *NetStateController) GetForm() {
	if models.NetTestFuncOn {
		netState, err := models.GetNetState()
		if err != nil {
			outErr := fmt.Errorf("Check network state from MySQL Error: %w", err)
			beego.Error(outErr)
			c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
			c.Data["errorMessage"] = outErr.Error()
			c.TplName = "error.tpl"
			return
		}

		// we should give an order to the frontend to read the information in the map
		var netStateKeys []string
		for key, _ := range netState {
			netStateKeys = append(netStateKeys, key)
		}

		c.Data["NetTestPeriodSec"] = models.NetTestPeriodSec
		c.Data["netStateKeys"] = netStateKeys
		c.Data["netState"] = netState
	} else {
		c.Data["netTestFuncOffMsg"] = models.NetTestFuncOffMsg
	}
	c.Data["NetTestFuncOn"] = models.NetTestFuncOn
	c.TplName = "netState.tpl"
}
