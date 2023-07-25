/*
The functions to allocate CPU cores to applications inside Virtual Machines.
*/

package algorithms

import asmodel "emcontroller/auto-schedule/model"

// allocate cpus to complete this solution.
func allocateCpus(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, solnWithVm asmodel.Solution) (asmodel.Solution, bool) {

	// avoiding changing the original solution
	solnWithCpu := asmodel.SolutionCopy(solnWithVm)

	// We should allocate VMs cloud by cloud.
	for _, cloud := range clouds {
		solnWithCpuThisCloud, acceptable := allocateCpusOneCLoud(cloud, apps, appsOrder, solnWithVm)
		if !acceptable { // if any cloud cannot accept the scheduled applications, this whole solution is not acceptable.
			return asmodel.Solution{}, false
		}

		// add the solution of this cloud to the total solution
		for appName, appSoln := range solnWithCpuThisCloud {
			solnWithCpu.AppsSolution[appName] = appSoln
		}
	}

	return solnWithCpu, true
}

// allocate cpus in one cloud
func allocateCpusOneCLoud(cloud asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, solnWithVm asmodel.Solution) (map[string]asmodel.SingleAppSolution, bool) {
	// For every cloud, at first, we find out the applications scheduled on it.
	// appsThisCloud := findAppsOneCloud(cloud, apps, solnWithVm)
	_ = findAppsOneCloud(cloud, apps, solnWithVm)
	return solnWithVm.AppsSolution, true
}
