package executors

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"

	"emcontroller/auto-schedule/algorithms"
	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

func CreateAutoScheduleApps(apps []models.K8sApp) ([]models.AppInfo, error, int) {

	// we only accept the valid applications, or otherwise we will have too much unnecessary workload
	if errs := ValidateAutoScheduleApps(apps); len(errs) != 0 {
		outErr := fmt.Errorf("The input applicatios are invalid, Error: [%w]", models.HandleErrSlice(errs))
		beego.Error(outErr)
		return []models.AppInfo{}, outErr, http.StatusBadRequest
	}

	// make the asmodel.Cloud structure as the input of Schedule function
	cloudsForScheduling, err := asmodel.GenerateClouds(models.Clouds)
	if err != nil {
		outErr := fmt.Errorf("Generate input clouds for auto-scheduling, Error: [%w]", err)
		beego.Error(outErr)
		return []models.AppInfo{}, outErr, http.StatusInternalServerError
	}

	// make the asmodel.Application structure as the input of Schedule function
	appsForScheduling, err := asmodel.GenerateApplications(apps)
	if err != nil {
		outErr := fmt.Errorf("Generate input applications for auto-scheduling, Error: [%w]", err)
		beego.Error(outErr)
		return []models.AppInfo{}, outErr, http.StatusInternalServerError
	}

	// call the Schedule method in mcasga.go
	mcssgaInstance := algorithms.Mcssga{}
	solution, err := mcssgaInstance.Schedule(cloudsForScheduling, appsForScheduling)
	if err != nil {
		outErr := fmt.Errorf("Run the Schedule method of Mcssga, Error: [%w]", err)
		beego.Error(outErr)
		return []models.AppInfo{}, outErr, http.StatusInternalServerError
	}
	fmt.Println(solution)

	// call models.CreateApplication to deploy

	// wait for application running

	// get applications and return then AppInfo

	return []models.AppInfo{}, nil, http.StatusCreated
}
