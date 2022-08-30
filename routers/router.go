package routers

import (
	"emcontroller/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/image", &controllers.ImageController{}, "get:Get")
	beego.Router("/upload", &controllers.ImageController{}, "post:Upload")

	beego.Router("/service", &controllers.ServiceController{}, "get:Get")
}
