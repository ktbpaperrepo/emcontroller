package algorithms

import (
	asmodel "emcontroller/auto-schedule/model"
)

const (
	cpuCoreStep float64 = 0.1 // also named as stride. This can be seen as the unit that we allocate CPU cores in our algorithm.

	// On a cloud, when existing VMs do not have enough resources for the applications scheduled here, we have 3 possible choices to create a new VM:
	// 1. a VM with 50% resources; if 50% is not enough, do 3.
	// 2. a VM with 30% resources; if 30% is not enough, do 3.
	// 3. a VM with all rest resources; if all rest is not enough, it means this solution is not acceptable.
	biggerVmResPct  float64 = 0.5
	smallerVmResPct float64 = 0.3
)

// SchedulingAlgorithm is the interface that all algorithms should implement
type SchedulingAlgorithm interface {
	Schedule(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) (asmodel.Solution, error)
}

// With a solution, Find out the applications scheduled on this cloud
func findAppsOneCloud(cloud asmodel.Cloud, apps map[string]asmodel.Application, soln asmodel.Solution) map[string]asmodel.Application {
	appsThisCloud := make(map[string]asmodel.Application)
	for appName, appSoln := range soln.AppsSolution {
		if appSoln.Accepted && appSoln.TargetCloudName == cloud.Name {
			appsThisCloud[appName] = apps[appName]
		}
	}
	return appsThisCloud
}

// filter the max-priority applications
func filterMaxPriApps(apps map[string]asmodel.Application) map[string]asmodel.Application {
	var maxPriApps map[string]asmodel.Application = make(map[string]asmodel.Application)
	for appName, app := range apps {
		if app.Priority == asmodel.MaxPriority {
			maxPriApps[appName] = app
		}
	}
	return maxPriApps
}

// filter out the max-priority applications
func filterOutMaxPriApps(apps map[string]asmodel.Application) map[string]asmodel.Application {
	var maxPriApps map[string]asmodel.Application = make(map[string]asmodel.Application)
	for appName, app := range apps {
		if app.Priority != asmodel.MaxPriority {
			maxPriApps[appName] = app
		}
	}
	return maxPriApps
}

// In Golang, the iteration order of map is random, but in some steps of scheduling, we need a fixed order of applications, so we make this function to randomly generate a order of applications and use it as the fixed application order in scheduling.
func GenerateAppsOrder(apps map[string]asmodel.Application) []string {
	var appsOrder []string
	for appName, _ := range apps {
		appsOrder = append(appsOrder, appName)
	}
	return appsOrder
}
