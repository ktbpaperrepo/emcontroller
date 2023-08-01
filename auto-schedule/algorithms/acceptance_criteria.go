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

// Check whether a solution is acceptable
func Acceptable(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) bool {
	// check resources
	if !depAcc(clouds, apps, soln) {
		return false
	}

	// Maybe there will be other aspects in the future

	// all checks passed
	return true
}

// Check whether a solution is acceptable in terms of dependency
func depAcc(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, soln asmodel.Solution) bool {

	for appName, app := range apps {
		for _, dep := range app.Dependencies {
			depAppName := dep.AppName

			// If an application is accepted, all its dependent applications should be accepted.
			if soln.AppsSolution[appName].Accepted && !soln.AppsSolution[depAppName].Accepted {
				return false
			}

			// We presume that every application needs to send network requests to all its dependent applications, so the network RTT should not be too large.
			// This check is only needed when this pair of applications are accepted.
			if soln.AppsSolution[appName].Accepted && soln.AppsSolution[depAppName].Accepted {
				srcVmName := soln.AppsSolution[appName].K8sNodeName
				dstVmName := soln.AppsSolution[depAppName].K8sNodeName
				// If 2 applications are deployed on the same VM, we think that the RTT between them is 0, so this check will not be needed in that condition.
				if srcVmName != dstVmName {
					srcCloudName := soln.AppsSolution[appName].TargetCloudName
					dstCloudName := soln.AppsSolution[depAppName].TargetCloudName

					// If the RTT is too large, this solution is not acceptable.
					if clouds[srcCloudName].NetState[dstCloudName].Rtt > minAccRttMs {
						return false
					}

				}

			}

		}
	}

	// all clouds passed
	return true
}
