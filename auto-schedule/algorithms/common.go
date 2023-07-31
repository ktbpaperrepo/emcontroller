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

	floatDelta float64 = 0.00001 // binary-floating-point data is not accurate, so we need to allow a delta when checking whether 2 float values are equal
)

// SchedulingAlgorithm is the interface that all algorithms should implement
type SchedulingAlgorithm interface {
	Schedule(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) (asmodel.Solution, error)
}

// After scheduling applications to clouds, we get a coarse solution. Then, we use this function to refine the solution, do 3 things:
// 1. schedule applications to VMs inside clouds;
// 2. allocate CPUs to applications inside VMs;
// 3. Check whether this solution is acceptable.
// If this solution passes the above 3 things, we return the refined solution.
func RefineSoln(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) (asmodel.Solution, bool) {
	// 1. give the solution node names
	solnWithVm, vmAcceptable := allocateVms(clouds, apps, appsOrder, soln)
	if !vmAcceptable {
		return asmodel.Solution{}, false
	}
	// 2. Allocate CPU cores
	solnWithCpu, cpuAcceptable := allocateCpus(clouds, apps, appsOrder, solnWithVm)
	if !cpuAcceptable {
		return asmodel.Solution{}, false
	}
	// 3. Check whether this solution is acceptable.
	if !Acceptable(clouds, apps, appsOrder, solnWithCpu) {
		return asmodel.Solution{}, false
	}
	return solnWithCpu, true
}
