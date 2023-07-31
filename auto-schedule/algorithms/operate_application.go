package algorithms

import asmodel "emcontroller/auto-schedule/model"

// With a solution, Find out the applications scheduled on this cloud
func findAppsOneCloud(cloud asmodel.Cloud, apps map[string]asmodel.Application, soln asmodel.Solution) map[string]asmodel.Application {
	appsThisCloud := make(map[string]asmodel.Application)
	for appName, appSoln := range soln.AppsSolution {
		if appSoln.Accepted && appSoln.TargetCloudName == cloud.Name {
			appsThisCloud[appName] = asmodel.AppCopy(apps[appName])
		}
	}
	return appsThisCloud
}

// filter the max-priority applications
func filterMaxPriApps(apps map[string]asmodel.Application) map[string]asmodel.Application {
	var maxPriApps map[string]asmodel.Application = make(map[string]asmodel.Application)
	for appName, app := range apps {
		if app.Priority == asmodel.MaxPriority {
			maxPriApps[appName] = asmodel.AppCopy(app)
		}
	}
	return maxPriApps
}

// filter out the max-priority applications
func filterOutMaxPriApps(apps map[string]asmodel.Application) map[string]asmodel.Application {
	var maxPriApps map[string]asmodel.Application = make(map[string]asmodel.Application)
	for appName, app := range apps {
		if app.Priority != asmodel.MaxPriority {
			maxPriApps[appName] = asmodel.AppCopy(app)
		}
	}
	return maxPriApps
}

// filter the applications whose names are in the input target application names
func filterAppsByNames(tgtAppNames []string, apps map[string]asmodel.Application) map[string]asmodel.Application {
	outApps := make(map[string]asmodel.Application)
	for _, appName := range tgtAppNames {
		outApps[appName] = asmodel.AppCopy(apps[appName])
	}
	return outApps
}

// In Golang, the iteration order of map is random, but in some steps of scheduling, we need a fixed order of applications, so we make this function to randomly generate a order of applications and use it as the fixed application order in scheduling.
func GenerateAppsOrder(apps map[string]asmodel.Application) []string {
	var appsOrder []string
	for appName, _ := range apps {
		appsOrder = append(appsOrder, appName)
	}
	return appsOrder
}

// calculate the resources needed by a group of applications
func calcNeededRes(apps map[string]asmodel.Application, appNames []string, minCpu bool) asmodel.AppResources {
	var neededRes asmodel.AppResources = asmodel.AppResources{
		GenericResources: asmodel.GenericResources{
			CpuCore: 0,
			Memory:  0,
			Storage: 0,
		},
	}

	for _, appName := range appNames {
		var cpuToOccupy float64
		if minCpu {
			cpuToOccupy = cpuCoreStep
		} else {
			cpuToOccupy = apps[appName].Resources.CpuCore
		}

		neededRes.CpuCore += cpuToOccupy
		neededRes.Memory += apps[appName].Resources.Memory
		neededRes.Storage += apps[appName].Resources.Storage
	}
	return neededRes
}
