package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type ArticleController struct {
	beego.Controller
}

func (c *ArticleController) Get() {
	c.Data["title"] = "hello article"
	c.TplName = "article.html"
}

func (c *ArticleController) AddArticle() {
	c.Ctx.WriteString("Add News")
}

func (c *ArticleController) EditArticle() {
	id := c.GetString("id")
	logs.Info("The id is: ", id)
	logs.Info("Type: %T", id)
	id2, err := c.GetInt("id2")
	if err != nil {
		logs.Info(err)
		c.Ctx.WriteString("error")
		return
	}
	logs.Info("The id2 is: %T %v", id2, id2)
	c.Ctx.WriteString(fmt.Sprintf("Edit News: %v %d", id, id2))
}
