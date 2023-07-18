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
	allRestVmResPct float64 = 1.0
)

// SchedulingAlgorithm is the interface that all algorithms should implement
type SchedulingAlgorithm interface {
	Schedule(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) (asmodel.Solution, error)
}

// With a solution, Find out the applications scheduled on this cloud
func findAppsOneCloud(cloud asmodel.Cloud, apps map[string]asmodel.Application, soln asmodel.Solution) map[string]asmodel.Application {
	appsThisCloud := make(map[string]asmodel.Application)
	for appName, appSoln := range soln {
		if appSoln.Accepted && appSoln.TargetCloudName == cloud.Name {
			appsThisCloud[appName] = apps[appName]
		}
	}
	return appsThisCloud
}

// In Golang, the iteration order of map is random, but in some steps of scheduling, we need a fixed order of applications, so we make this function to randomly generate a order of applications and use it as the fixed application order in scheduling.
func GenerateAppsOrder(apps map[string]asmodel.Application) []string {
	var appsOrder []string
	for appName, _ := range apps {
		appsOrder = append(appsOrder, appName)
	}
	return appsOrder
}

// Iterator to iterate the applications on one cloud in a fixed order
type appOneCloudIter struct {
	appsThisCloud map[string]asmodel.Application
	appsOrder     []string // a slice of app names
	nextAppIdx    int
}

func newAppOneCloudIter(appsThisCloud map[string]asmodel.Application, appsOrder []string) *appOneCloudIter {
	return &appOneCloudIter{
		appsThisCloud: appsThisCloud,
		appsOrder:     appsOrder,
		nextAppIdx:    0,
	}
}

func (it *appOneCloudIter) nextAppName() string {
	for it.nextAppIdx < len(it.appsOrder) {
		curAppName := it.appsOrder[it.nextAppIdx]
		if _, exist := it.appsThisCloud[curAppName]; exist {
			it.nextAppIdx++
			return curAppName
		}
		it.nextAppIdx++
	}
	return ""
}

func (it *appOneCloudIter) Copy() *appOneCloudIter {
	// copy the map (not completely deep, but enough for our scenarios)
	appsThisCloudCopy := asmodel.AppMapCopy(it.appsThisCloud)

	// deep copy the slice
	appsOrderCopy := make([]string, len(it.appsOrder))
	copy(appsOrderCopy, it.appsOrder)

	return &appOneCloudIter{
		appsThisCloud: appsThisCloudCopy,
		appsOrder:     appsOrderCopy,
		nextAppIdx:    it.nextAppIdx,
	}
}
