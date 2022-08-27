package main

import (
	_ "emcontroller/routers"
	"github.com/astaxie/beego"

	"emcontroller/models"
)

func main() {
	//beego.BConfig.WebConfig.TemplateLeft = "<<"
	//beego.BConfig.WebConfig.TemplateRight = ">>"
	beego.AddFuncMap("hi", models.Hello)
	beego.AddFuncMap("unixToDate", models.UnixToDate)

	// configure static resource paths
	beego.SetStaticPath("/down", "download")
	beego.Run()
}
