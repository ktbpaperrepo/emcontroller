/*
The functions to allocate Virtual Machines to applications inside clouds.
*/

package algorithms

import asmodel "emcontroller/auto-schedule/model"

// allocate VMs to complete this solution.
func allocateVms(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) (asmodel.Solution, bool) {

	// avoiding changing the original solution
	solnWithVm := asmodel.SolutionCopy(soln)

	// We should allocate VMs cloud by cloud.
	for _, cloud := range clouds {
		solnWithVmsThisCloud, acceptable := allocateVmsOneCLoud(cloud, apps, appsOrder, soln)
		if !acceptable { // if any cloud cannot accept the scheduled applications, this whole solution is not acceptable.
			return nil, false
		}

		// add the solution of this cloud to the total solution
		for appName, appSoln := range solnWithVmsThisCloud {
			solnWithVm[appName] = appSoln
		}
	}

	return solnWithVm, true
}

// allocate VMs in one cloud
func allocateVmsOneCLoud(cloud asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) (map[string]asmodel.SingleAppSolution, bool) {
	// For every cloud, at first, we find out the applications scheduled on it.
	//appsThisCloud := findAppsOneCloud(cloud, apps, soln)
	_ = findAppsOneCloud(cloud, apps, soln)
	return soln, true
}
