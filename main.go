package main

import (
	"emcontroller/models"
	"github.com/astaxie/beego"

	_ "emcontroller/routers"
)

func main() {
	models.InitDockerClient()
	models.InitKubernetesClient()

	// viper is case-insensitive, so all keys in iaas.json should be lowercase
	models.InitClouds()

	beego.Run()
}
