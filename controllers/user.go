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
	c.TplName = "userAdd.html"
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

func (c *UserController) EditUser() {
	c.TplName = "userEdit.html"
}

type User struct {
	ID       string   `form:"id" json:"id1"`
	Username string   `form:"username" json:"username1"`
	Password string   `form:"password" json:"password1"`
	Hobby    []string `form:"hobby" json:"hobby1"`
}

func (c *UserController) DoEditUser() {
	u := User{}
	if err := c.ParseForm(&u); err != nil {
		c.Ctx.WriteString("post error")
		return
	}
	logs.Info("%#v", u)
	c.Ctx.WriteString("parse form successful")
}

func (c *UserController) GetUser() {
	u := User{
		ID:       "3232",
		Username: "abc",
		Password: "1234",
		Hobby:    []string{"1", "2"},
	}
	c.Data["json"] = u
	c.ServeJSON()
}
