package controllers

import (
	"emcontroller/models"
	"fmt"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = models.ControllerName
	c.Data["VersionInfo"] = fmt.Sprintf("Build time: [%s]. Git commit: [%s]\n", models.BuildDate, models.GitCommit)

	c.TplName = "index.tpl"
}
