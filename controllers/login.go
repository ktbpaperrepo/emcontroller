package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
}

func (c *LoginController) DoLogin() {
	username := c.GetString("username")
	password := c.GetString("password")
	logs.Info(username, password)

	// c.Redirect("/", 302)
	//c.Redirect("/api/"+username, 302)
	c.Ctx.Redirect(302, "/api/"+username)
}
