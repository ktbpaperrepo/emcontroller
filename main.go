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
	fmt.Printf("Build time: [%s]. Git commit: [%s]\n", models.BuildDate, models.GitCommit)
}

func main() {
	// It seems that the -X parameter of go build -ldflags can only set the values for the variables in the main package
	models.BuildDate = buildDate
	models.GitCommit = gitCommit

	versionFlag := flag.Bool("v", false, "Print the current version and exit.")
	flag.Parse()
	switch {
	case *versionFlag:
		printVersion()
		return
	default:
		// continue with other stuff
	}

	models.InitSomeThing()

	if netTestOn, err := beego.AppConfig.Bool("TurnOnNetTest"); err == nil && netTestOn {
		beego.Info("Network performance test function is on.")
		if err := models.InitNetPerfDB(); err != nil {
			outErr := fmt.Errorf("Initialize the database [%s] in MySQL failed, error: [%w]", models.NetPerfDbName, err)
			beego.Error(outErr)
			panic(outErr)
		}
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

		models.NetTestFuncOn = true
		models.NetTestPeriodSec = netTestPeriodSec
	} else if err != nil {
		beego.Error(fmt.Sprintf("Read \"TurnOnNetTest\" in app.conf, error: [%s]. We turn off the network performance test function.", err.Error()))
	} else {
		beego.Info("Network performance test function is off.")
	}

	beego.Run()
}
