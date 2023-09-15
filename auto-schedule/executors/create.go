package executors

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/astaxie/beego"

	"emcontroller/auto-schedule/algorithms"
	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

// algoName is the name of the scheduling algorithm to use.
func CreateAutoScheduleApps(apps []models.K8sApp, algoName string) ([]models.AppInfo, error, int) {

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

	// Whether this order is fixed or random does not affect the performance of algorithms, because the applications are generated randomly, which will not be changed by a fixed order. However, when we fix the order here, the comparison between different algorithms can have the same input, because apps order is one input parameter.
	sort.Strings(appsOrder)

	// the parameters for genetic algorithms
	var chromosomesCount int = 200
	var iterationCount int = 5000
	var crossoverProbability float64 = 0.3
	var mutationProbability float64 = 0.019
	var stopNoUpdateIteration int = 200

	// call the Scheduling method according to the input parameter "algo"

	// create algorithm instances, and put them in a map
	mcssgaInstance := algorithms.NewMcssga(chromosomesCount, iterationCount, crossoverProbability, mutationProbability, stopNoUpdateIteration)
	var allAlgos map[string]algorithms.SchedulingAlgorithm = make(map[string]algorithms.SchedulingAlgorithm)
	allAlgos[algorithms.McssgaName] = mcssgaInstance
	allAlgos[algorithms.CompRandName] = algorithms.NewCompRand()
	allAlgos[algorithms.BERandName] = algorithms.NewBERand()
	allAlgos[algorithms.AmpgaName] = algorithms.NewApaga(chromosomesCount, iterationCount, crossoverProbability, mutationProbability, stopNoUpdateIteration)

	// select the algorithm to use according to the input parameter algoName
	beego.Info(fmt.Sprintf("Looking for the algorithm \"%s\".", algoName))
	var algoToUse algorithms.SchedulingAlgorithm
	var algoNameToUse string = algoName
	if algo, exist := allAlgos[algoName]; exist {
		beego.Info(fmt.Sprintf("Algorithm \"%s\" is found.", algoName))
		algoToUse = algo
	} else { // if we cannot find the input algoName, we use MCASSGA algorithm by default.
		algoNameToUse = algorithms.McssgaName
		beego.Info(fmt.Sprintf("Algorithm \"%s\" is not found, so we use \"%s\" by default.", algoName, algoNameToUse))
		algoToUse = mcssgaInstance
	}

	solution, err := algoToUse.Schedule(cloudsForScheduling, appsForScheduling, appsOrder)
	if err != nil {
		outErr := fmt.Errorf("Run the Schedule method of %s, Error: [%w]", algoNameToUse, err)
		beego.Error(outErr)
		return []models.AppInfo{}, outErr, http.StatusInternalServerError
	}

	// If we did not use Mcssga to schedule apps, now its max rtt has not been set, so we should set it now to calculate the fitness value in the following log.
	mcssgaInstance.SetMaxReaRtt(cloudsForScheduling)
	mcssgaInstance.SetAvgDepNum(appsForScheduling)
	beego.Info(fmt.Sprintf("The algorithm works out the solution: %s\nIts fitness value is %g.", models.JsonString(solution), mcssgaInstance.Fitness(cloudsForScheduling, appsForScheduling, solution)))

	//// This part is for debug ----------------------------
	//
	//// draw evolution chart
	//if mcssgaAlgo, ok := algoToUse.(*algorithms.Mcssga); ok {
	//	mcssgaAlgo.DrawEvoChart()
	//}
	//switch realAlgo := algoToUse.(type) {
	//case *algorithms.Mcssga:
	//	realAlgo.DrawEvoChart()
	//case *algorithms.Ampga:
	//	realAlgo.DrawEvoChart()
	//default:
	//}
	//
	//// return accepted names for experiments
	//var acceptedApps []models.AppInfo
	//for appName, app := range appsForScheduling {
	//	if solution.AppsSolution[appName].Accepted {
	//		acceptedApps = append(acceptedApps, models.AppInfo{AppName: appName, Priority: app.Priority})
	//	}
	//}
	//return acceptedApps, nil, http.StatusCreated
	//// This part is for debug ----------------------------

	/**
	TODO:
	migration: I set a lock, migration and deployment (or multiple deployments) cannot be done at the same time. When doing migration, we skip the resources occupied by the applications to be migrated, and count them as the VM resources. When the resources are not enough, the rolling update may be blocked, because the new pods cannot be created. Maybe I can make a dependency topo-sort to avoid it.
	I will put the migration into the next paper.
	*/

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
		app.Containers[0].Resources.Requests.CPU = fmt.Sprintf("%.0f", scheSoln.AppsSolution[app.Name].AllocatedCpuCore)
		app.Containers[0].Resources.Limits.CPU = fmt.Sprintf("%.0f", scheSoln.AppsSolution[app.Name].AllocatedCpuCore)

		appsWithScheInfo = append(appsWithScheInfo, app)
	}

	return appsWithScheInfo
}
