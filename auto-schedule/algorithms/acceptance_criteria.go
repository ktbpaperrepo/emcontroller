package algorithms

import (
	asmodel "emcontroller/auto-schedule/model"
)

type AppType int

const (
	AllPriApps    AppType = iota // all applications.
	MaxPriApps                   // only max-priority applications.
	NotMaxPriApps                // only not-max-priority applications.
)

// Check whether a solution is acceptable (whether the resources are enough or not)
func Acceptable(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) bool {
	// check resources
	if !resAcc(clouds, apps, appsOrder, soln) {
		return false
	}

	// TODO: check other aspects

	// all checks passed
	return true
}

// Check whether a solution is acceptable in terms of resources
func resAcc(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) bool {
	//for _, cloud := range clouds { // check every cloud
	//	if !resAccOneCloudSharedVm(cloud, apps, appsOrder, soln, AllPriApps) {
	//		return false
	//	}
	//}

	// all clouds passed
	return true
}
