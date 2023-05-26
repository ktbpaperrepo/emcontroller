package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"

	"emcontroller/models"
	_ "emcontroller/routers"
)

var gitCommit, buildDate string

func printVersion() {
	fmt.Printf("Build time: [%s]. Git commit: [%s]\n", buildDate, gitCommit)
}

func main() {
	versionFlag := flag.Bool("v", false, "Print the current version and exit.")
	flag.Parse()
	switch {
	case *versionFlag:
		printVersion()
		return
	default:
		// continue with other stuff
	}

	models.InitDockerClient()
	models.InitKubernetesClient()

	// viper is case-insensitive, so all keys in iaas.json should be lowercase
	models.InitClouds()

	if NetTestOn, err := beego.AppConfig.Bool("TurnOnNetTest"); err == nil && NetTestOn {
		beego.Info("Network performance test function is on.")
		// When multi-cloud manager starts up, we do measure the network performance once instantly, because we need the network performance information to schedule applications when deploying them.
		models.MeasNetPerf()
		// periodically measure the network performance
		netTestPeriodSec, err := strconv.Atoi(beego.AppConfig.String("NetTestPeriodSec"))
		if err != nil {
			beego.Error(fmt.Sprintf("Read config \"NetTestPeriodSec\" error: %s, set the period as the DefaultNetTestPeriodSec", err.Error()))
			netTestPeriodSec = models.DefaultNetTestPeriodSec
		}
		beego.Info(fmt.Sprintf("The period of measuring network performance is %d seconds.", netTestPeriodSec))
		go models.CronTaskTimer(models.MeasNetPerf, time.Duration(netTestPeriodSec)*time.Second)
	} else if err != nil {
		beego.Error(fmt.Sprintf("Read \"TurnOnNetTest\" in app.conf, error: [%s]. We turn off the network performance test function.", err.Error()))
	} else {
		beego.Info("Network performance test function is off.")
	}

	beego.Run()
}
