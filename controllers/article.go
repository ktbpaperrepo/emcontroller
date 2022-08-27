package controllers

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"emcontroller/models"
)

type ArticleController struct {
	beego.Controller
}

func (c *ArticleController) Get() {
	c.Data["title"] = "hello article"

	now := time.Now()
	fmt.Println(now)
	c.Data["now"] = now

	c.Data["title"] = "This is a list of the article"

	c.Data["html"] = "<h2>This h2 is rendered by the backend</h2>"

	userinfo := make(map[string]interface{})
	userinfo["username"] = "Zhang San"
	userinfo["age"] = 20
	userinfo["a"] = map[string]float64{
		"c": 4,
	}
	c.Data["userinfo"] = userinfo
	c.Data["unix"] = 1587880013

	// cut string
	var str string = "This is an article list"
	slice := []rune(str)
	fmt.Println(slice)
	fmt.Println(string(slice[5:]))
	c.Data["str"] = string(slice[5:10])

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

func (c *ArticleController) Add() {
	unix := 1599880013
	str := models.UnixToDate(unix)
	c.Ctx.WriteString("news add --- " + str)

}
