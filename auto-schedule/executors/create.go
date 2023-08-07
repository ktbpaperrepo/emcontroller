package executors

import (
	"fmt"
	"github.com/astaxie/beego"
	"net/http"

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
	// In some steps of scheduling, we need a fixed order of applications.
	appsOrder := algorithms.GenerateAppsOrder(appsForScheduling)

	//// for debug, sometimes, we need the fixed apps order to do some comparison.
	//sort.Strings(appsOrder)

	// call the Schedule method in mcasga.go
	mcssgaInstance := algorithms.NewMcssga(100, 5000, 0.7, 0.2, 200)
	solution, err := mcssgaInstance.Schedule(cloudsForScheduling, appsForScheduling, appsOrder)
	if err != nil {
		outErr := fmt.Errorf("Run the Schedule method of Mcssga, Error: [%w]", err)
		beego.Error(outErr)
		return []models.AppInfo{}, outErr, http.StatusInternalServerError
	}
	beego.Info(fmt.Sprintf("The algorithm works out the solution: %s\nIts fitness value is %g.", models.JsonString(solution), mcssgaInstance.Fitness(cloudsForScheduling, appsForScheduling, solution)))

	// for debug
	return []models.AppInfo{}, nil, http.StatusCreated

	// create the VMs and add them to Kubernetes
	if _, err := models.AddNewVms(solution.VmsToCreate); err != nil {
		outErr := fmt.Errorf("Add new auto-scheduling VMs, Error: [%w]", err)
		beego.Error(outErr)
		return []models.AppInfo{}, outErr, http.StatusInternalServerError
	}

	// add the auto-scheduling information into the applications to deploy.
	appsToDeploy := addScheInfoToApps(apps, solution)

	// deploy applications, and wait for them running.
	createdAppsInfo, err := models.CreateAppsWait(appsToDeploy)
	if err != nil {
		outErr := fmt.Errorf("Create auto-scheduling applications [%s], Error: [%w]", models.JsonString(appsToDeploy), err)
		beego.Error(outErr)
		return []models.AppInfo{}, outErr, http.StatusInternalServerError
	}

	return createdAppsInfo, nil, http.StatusCreated
}

// After scheduling applications, we should use this functions to add the scheduling information to applications.
func addScheInfoToApps(apps []models.K8sApp, scheSoln asmodel.Solution) []models.K8sApp {
	var appsWithScheInfo []models.K8sApp

	for _, app := range apps {
		// Remove the rejected applications from the array.
		if !scheSoln.AppsSolution[app.Name].Accepted {
			continue
		}

		// add node name
		app.NodeName = scheSoln.AppsSolution[app.Name].K8sNodeName
		// configure allocated CPU
		app.Containers[0].Resources.Requests.CPU = fmt.Sprintf("%.1f", scheSoln.AppsSolution[app.Name].AllocatedCpuCore)

		appsWithScheInfo = append(appsWithScheInfo, app)
	}

	return appsWithScheInfo
}
