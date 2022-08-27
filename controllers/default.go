package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type Article struct {
	Title   string
	Content string
}

func (c *MainController) Get() {
	// 1. bind basic data types in the template: string, num, bool
	c.Data["website"] = "beego.me"
	c.Data["title"] = "hello beego"
	c.Data["num"] = 12
	c.Data["flag"] = true

	// 2. bind the data of structures in the template
	article := Article{
		Title:   "I am a golang course",
		Content: "beego fake xiaomi project",
	}
	c.Data["article"] = article

	// 3. customize variables in the template

	// 4. loop in template, range, loop slice
	c.Data["sliceList"] = []string{"php", "java", "golang"}

	// 5. loop in template, range, loop Map
	userinfo := make(map[string]interface{})
	userinfo["username"] = "Zhang San"
	userinfo["age"] = 20
	userinfo["sex"] = "male"
	c.Data["userinfo"] = userinfo

	// 6. loop in template, loop slice with the type of structure
	// usually used for database data
	c.Data["articleList"] = []Article{
		{
			Title:   "news 111",
			Content: "news content 111",
		},
		{
			Title:   "news 222",
			Content: "news content 222",
		},
		{
			Title:   "news 333",
			Content: "news content 333",
		},
		{
			Title:   "news 444",
			Content: "news content 444",
		},
	}

	// 7. loop in template, another definition method of slice with the type of structure
	/*
		anonymous structure, it is a type
			struct {
				Title string
			}
	*/
	c.Data["cmdList"] = []struct {
		Title string
	}{
		{
			Title: "news 111111111",
		},
		{
			Title: "news 222222222222",
		},
		{
			Title: "news 3333333333333",
		},
		{
			Title: "news 4444444444444444",
		},
	}

	// 8. conditions in template
	c.Data["isLogin"] = true
	c.Data["isHome"] = false
	c.Data["isAbout"] = true

	// 9. if condition eq / ne / lt / le / gt / ge
	c.Data["n1"] = 12
	c.Data["n2"] = 6

	// 12. use self-defined template functions
	c.Data["unix"] = 1587880013

	c.TplName = "index.tpl"
}
