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
	beego.Router("/cloud/:cloudName/vm/:vmID", &controllers.VmController{}, "get:GetVM")

	beego.Router("/vm", &controllers.VmController{}, "get:ListVMsAllClouds")
	beego.Router("/vm/new", &controllers.VmController{}, "get:NewVms")
	beego.Router("/vm/doNew", &controllers.VmController{}, "post:DoNewVms")

	beego.Router("/image", &controllers.ImageController{}, "get:Get")
	beego.Router("/image/:repo", &controllers.ImageController{}, "delete:DeleteRepo")
	beego.Router("/upload", &controllers.ImageController{}, "post:Upload")

	beego.Router("/application", &controllers.ApplicationController{}, "get:Get")
	beego.Router("/application/:appName", &controllers.ApplicationController{}, "delete:DeleteApp")
	beego.Router("/application/:appName", &controllers.ApplicationController{}, "get:GetApp")
	beego.Router("/newApplication", &controllers.ApplicationController{}, "get:NewApplication")
	beego.Router("/doNewApplication", &controllers.ApplicationController{}, "post:DoNewApplication")

	beego.Router("/k8sNode", &controllers.K8sNodeController{}, "get:Get")
	beego.Router("/k8sNode/:nodeName", &controllers.K8sNodeController{}, "delete:DeleteNode")
	beego.Router("/k8sNode/add", &controllers.K8sNodeController{}, "get:AddNodes")
	beego.Router("/k8sNode/doAdd", &controllers.K8sNodeController{}, "post:DoAddNodes")

	beego.Router("/netState", &controllers.NetStateController{}, "get:Get")
}
