package main

import (
	_ "emcontroller/routers"
	"github.com/astaxie/beego"
)

func main() {
	//beego.BConfig.WebConfig.TemplateLeft = "<<"
	//beego.BConfig.WebConfig.TemplateRight = ">>"
	beego.Run()
}
