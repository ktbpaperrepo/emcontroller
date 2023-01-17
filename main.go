package main

import (
	"emcontroller/models"
	"github.com/astaxie/beego"

	_ "emcontroller/routers"
)

func main() {
	models.InitDockerClient()
	models.InitKubernetesClient()

	beego.Run()
}
