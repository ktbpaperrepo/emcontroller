package controllers

import (
	"github.com/astaxie/beego"
)

type CmsController struct {
	beego.Controller
}

func (c *CmsController) Get() {
	id := c.Ctx.Input.Param(":id")
	c.Ctx.WriteString("CMS details---" + id)
}
