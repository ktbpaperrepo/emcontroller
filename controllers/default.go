package controllers

import (
	"emcontroller/models"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = models.ControllerName

	c.TplName = "index.tpl"
}
