package controllers

import (
	"github.com/astaxie/beego"
)

type ServiceController struct {
	beego.Controller
}

func (c *ServiceController) Get() {

	c.TplName = "service.tpl"
}
