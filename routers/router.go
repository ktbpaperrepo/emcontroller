package routers

import (
	"emcontroller/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/goods", &controllers.GoodsController{})
	beego.Router("/article", &controllers.ArticleController{})                       // get method
	beego.Router("/article/add", &controllers.ArticleController{}, "get:AddArticle") // customized
	beego.Router("/article/edit", &controllers.ArticleController{}, "get:EditArticle")
	beego.Router("/user", &controllers.UserController{})
	beego.Router("/user/add", &controllers.UserController{}, "get:AddUser")
	beego.Router("/user/doAdd", &controllers.UserController{}, "post:DoAddUser")
}
