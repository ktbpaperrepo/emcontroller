package routers

import (
	"emcontroller/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/cloud", &controllers.CloudController{}, "get:Get")
	beego.Router("/cloud/:cloudName", &controllers.CloudController{}, "get:GetSingleCloud")

	beego.Router("/cloud/:cloudName/vm/:vmID", &controllers.VmController{}, "delete:DeleteVM")
	beego.Router("/cloud/:cloudName/vm", &controllers.VmController{}, "post:CreateVM")

	beego.Router("/image", &controllers.ImageController{}, "get:Get")
	beego.Router("/image/:repo", &controllers.ImageController{}, "delete:DeleteRepo")
	beego.Router("/upload", &controllers.ImageController{}, "post:Upload")

	beego.Router("/application", &controllers.ApplicationController{}, "get:Get")
	beego.Router("/application/:appName", &controllers.ApplicationController{}, "delete:DeleteApp")
	beego.Router("/newApplication", &controllers.ApplicationController{}, "get:NewApplication")
	beego.Router("/doNewApplication", &controllers.ApplicationController{}, "post:DoNewApplication")

}
