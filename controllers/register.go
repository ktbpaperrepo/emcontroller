package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Get() {
	c.TplName = "register.html"
}

func (c *RegisterController) DoRegister() {
	username := c.GetString("username")
	password := c.GetString("password")
	rpassword := c.GetString("rpassword")
	logs.Info(username, password, rpassword)

	c.TplName = "success.html"
}
