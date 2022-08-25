package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
)

type GoodsController struct {
	beego.Controller
}

func (c *GoodsController) Get() { // get
	c.Data["mytitle"] = "Hello Beego!"
	c.Data["num"] = 12
	c.TplName = "goods.tpl"
}

func (c *GoodsController) DoAdd() { // post
	c.Ctx.WriteString("execute add operation")
}

type Product struct {
	Title   string `form:"title"`
	Content string `form:"content"`
}

func (c *GoodsController) DoEdit() { // put
	//title := c.GetString("title")
	p := Product{}
	if err := c.ParseForm(&p); err != nil {
		c.Ctx.WriteString("error")
		return
	}
	fmt.Printf("%#v\n", p)

	c.Ctx.WriteString("execute edit operation" + p.Title + p.Content)
}

func (c *GoodsController) DoDelete() { // delete
	id, err := c.GetInt("id")
	if err != nil {
		c.Ctx.WriteString("parameter error")
	}
	c.Ctx.WriteString("execute delete operation----" + strconv.Itoa(id))
}
