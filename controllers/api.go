package controllers

import (
	"github.com/astaxie/beego"
)

type ApiController struct {
	beego.Controller
}

// http://localhost:20000/api/1231asf
func (c *ApiController) Get() {
	// get the value of dynamic router
	id := c.Ctx.Input.Param(":id")
	c.Ctx.WriteString("api interface---" + id)
}
