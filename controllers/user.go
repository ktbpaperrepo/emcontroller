package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) Get() {
	c.Ctx.WriteString("User center")
}

func (c *UserController) AddUser() {
	c.TplName = "user.html"
}

func (c *UserController) DoAddUser() {
	id, err := c.GetInt("id")
	if err != nil {
		c.Ctx.WriteString("id should be int")
		return
	}
	logs.Info("%v-----%T", id, id)
	username := c.GetString("username")
	password := c.GetString("password")
	hobby := c.GetStrings("hobby")
	logs.Info("value: %v-----type: %T", hobby, hobby)

	c.Ctx.WriteString("User center" + strconv.Itoa(id) + username + password)
}
