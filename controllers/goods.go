package controllers

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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
	Title   string `form:"title" xml:"title" json:"title"`
	Content string `form:"content" xml:"content" json:"content"`
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

// need conf/app.conf: copyrequestbody = true
func (c *GoodsController) Xml() { // delete
	str := string(c.Ctx.Input.RequestBody)
	logs.Info(str)

	p := Product{}

	if err := xml.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Ctx.WriteString(fmt.Sprintf("error: %s\n", err.Error()))
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = p
	}
	c.ServeJSON()
}
