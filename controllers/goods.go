package controllers

import (
	"github.com/astaxie/beego"
)

type GoodsController struct {
	beego.Controller
}

func (c *GoodsController) Get() {
	c.Data["mytitle"] = "Hello Beego!"
	c.Data["num"] = 12
	c.TplName = "goods.tpl"
}
