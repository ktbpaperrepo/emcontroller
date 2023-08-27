/*
The functions to allocate CPU cores to applications inside Virtual Machines.
*/

package algorithms

import (
	"fmt"
	"math"

	asmodel "emcontroller/auto-schedule/model"
	"github.com/KeepTheBeats/routing-algorithms/mymath"
	apiv1 "k8s.io/api/core/v1"
)

// allocate cpus to complete this solution.
func allocateCpus(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, solnWithVm asmodel.Solution) (asmodel.Solution, bool) {

	// avoiding changing the original solution
	solnWithCpu := asmodel.SolutionCopy(solnWithVm)

	// We should allocate CPUs cloud by cloud.
	for _, cloud := range clouds {
		solnWithCpuThisCloud, acceptable := allocateCpusOneCloud(cloud, apps, appsOrder, solnWithVm)
		if !acceptable { // if any cloud cannot accept the scheduled applications, this whole solution is not acceptable.
			return asmodel.Solution{}, false
		}

		// update the solution of this cloud to the total solution.
		solnWithCpu.Absorb(solnWithCpuThisCloud)
	}

	return solnWithCpu, true
}

// allocate cpus in one cloud
func allocateCpusOneCloud(cloud asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, solnWithVm asmodel.Solution) (asmodel.Solution, bool) {
	// For every cloud, at first, we find out the applications scheduled on it.
	appsThisCloud := findAppsOneCloud(cloud, apps, solnWithVm)

	// we put the result to return into this variable
	solnWithCpu := asmodel.GenEmptySoln()
	for appName, _ := range appsThisCloud { // the result should only include the solutions for the applications to handle
		solnWithCpu.AppsSolution[appName] = asmodel.SasCopy(solnWithVm.AppsSolution[appName])
	}

	// group applications by the VMs on which they are scheduled
	vmAppGroups := groupAppsByVm(appsThisCloud, appsOrder, solnWithVm)

	// allocate CPUs to applications on each VM
	for vmName, appNames := range vmAppGroups {
		vm := getVmByName(vmName, cloud, solnWithVm)                                      // handle this vm
		solnWithCpuThisVm := allocateCpusOneVm(vm, apps, appsOrder, appNames, solnWithVm) // allocate CPUs to applications on this VM

		solnWithCpu.Absorb(solnWithCpuThisVm) // combine the solution of this VM into the solution of this cloud.
	}

	return solnWithCpu, true
}

// group applications, putting those on a same VM to a same group.
func groupAppsByVm(appsToGroup map[string]asmodel.Application, appsOrder []string, solnWithVm asmodel.Solution) map[string][]string {

	// We put the result in this map. The key is the VM name, and the value is a slice of the applications names.
	var vmAppGroups map[string][]string = make(map[string][]string)

	appIter := newIterForApps(appsToGroup, appsOrder)
	for curAppName := appIter.nextAppName(); len(curAppName) > 0; curAppName = appIter.nextAppName() {
		if solnWithVm.AppsSolution[curAppName].Accepted {
			schedVmName := solnWithVm.AppsSolution[curAppName].K8sNodeName // the name of the VM that this application is scheduled to.
			vmAppGroups[schedVmName] = append(vmAppGroups[schedVmName], curAppName)
		}
	}

	return vmAppGroups
}

// get a K8sNode variable from a target VM name, the cloud, and a solution.
func getVmByName(tgtVmName string, cloud asmodel.Cloud, solnWithVm asmodel.Solution) asmodel.K8sNode {
	// the VM can be got from either cloud.K8sNodes or solnWithVm.VmsToCreate.

	// try cloud.K8sNodes
	var oriNode asmodel.K8sNode
	var oriFound bool = false
	for _, node := range cloud.K8sNodes {
		if node.Name == tgtVmName {
			oriNode = node
			oriFound = true
			break
		}
	}

	// if we cannot find this node in cloud.K8sNodes, we try solnWithVm.VmsToCreate.
	var createNode asmodel.K8sNode
	var createFound bool = false
	for _, vm := range solnWithVm.VmsToCreate {
		if vm.Cloud == cloud.Name && vm.Name == tgtVmName { // the cloud must be specified correctly, or otherwise the VM cannot be found
			createNode = asmodel.GenK8sNodeFromPods(vm, []apiv1.Pod{})
			createFound = true
			break
		}
	}

	// To be safe, we write check the panic conditions.
	switch {
	case oriFound && !createFound:
		return oriNode
	case !oriFound && createFound:
		return createNode
	case oriFound && createFound:
		panic(fmt.Sprintf("The target node is found both in cloud.K8sNodes and in solnWithVm.VmsToCreate. oriFound is [%t], createFound is [%t]. oriNode is [%v], createNode is createNode [%v].", oriFound, createFound, oriNode, createNode))
	case !oriFound && !createFound:
		panic(fmt.Sprintf("The target node is not found either in cloud.K8sNodes or in solnWithVm.VmsToCreate. oriFound is [%t], createFound is [%t]. oriNode is [%v], createNode is createNode [%v].", oriFound, createFound, oriNode, createNode))
	default:
		panic(fmt.Sprintf("Condition should not exist. oriFound is [%t], createFound is [%t]. oriNode is [%v], createNode is createNode [%v].", oriFound, createFound, oriNode, createNode))
	}

}

// allocate CPUs to the applications on a VM.
func allocateCpusOneVm(vm asmodel.K8sNode, apps map[string]asmodel.Application, appsOrder []string, appNamesThisVm []string, solnWithVm asmodel.Solution) asmodel.Solution {
	appsThisVm := filterAppsByNames(appNamesThisVm, apps) // get the applications scheduled to this VM

	neededRes := calcNeededRes(apps, appNamesThisVm, false)
	if vm.ResidualResources.CpuCore >= neededRes.CpuCore { // If the code can reach here, it means that Memory and Storage of this VM are certainly enough for all applications, so we only need to check CPU.
		// condition 1: If the CPUs of this VM can meet all applications' requested CPUs, allocate CPUs as app's requests.
		return allocateCpuAsRequest(appsThisVm, solnWithVm)
	} else {
		// condition 2: If the CPUs of this VM are not enough for all applications' requested CPUs, allocate CPUs according to the "product of their requested CPUs and priorities".
		return vmCpuWeightedAllocation(vm, appsThisVm, appsOrder, solnWithVm)
	}
}

// allocate CPUs to some applications as their requests
func allocateCpuAsRequest(appsToHandle map[string]asmodel.Application, solnWithVm asmodel.Solution) asmodel.Solution {
	var theseAppsSolnWithCpu asmodel.Solution = asmodel.GenEmptySoln()
	for appName, _ := range appsToHandle { // the result should only include the solutions for the applications to handle
		thisAppSoln := asmodel.SasCopy(solnWithVm.AppsSolution[appName])       // get original single solution with VM
		thisAppSoln.AllocatedCpuCore = appsToHandle[appName].Resources.CpuCore // add CPU information in this single solution.
		theseAppsSolnWithCpu.AppsSolution[appName] = thisAppSoln               // put this single solution into the output solution.
	}
	return theseAppsSolnWithCpu
}

// allocate the CPUs of a VM to the applications scheduled to it weighted by "App requested CPU * App priority"
func vmCpuWeightedAllocation(vm asmodel.K8sNode, appsThisVm map[string]asmodel.Application, appsOrder []string, solnWithVm asmodel.Solution) asmodel.Solution {
	var solnWithCpuThisVm asmodel.Solution = asmodel.GenEmptySoln()
	for appName, _ := range appsThisVm { // the result should only include the solutions for the applications to handle
		solnWithCpuThisVm.AppsSolution[appName] = asmodel.SasCopy(solnWithVm.AppsSolution[appName])
	}

	// copy and avoid changing the original variable.
	vmCopy := asmodel.K8sNodeCopy(vm)
	remainingApps := asmodel.AppMapCopy(appsThisVm)

	// the Step 1.2 described in the following comments
	for {
		cpuAllocScheme := distrCpuApps(vmCopy, remainingApps)

		minCpuFound := false

		// do the Step 1.1 described in the following comments
		for appName, allocatedCpu := range cpuAllocScheme {
			if allocatedCpu <= cpuCoreStep {

				minCpuFound = true

				// the CPU allocation of this application is decided
				thisAppSoln := solnWithCpuThisVm.AppsSolution[appName]
				thisAppSoln.AllocatedCpuCore = cpuCoreStep
				solnWithCpuThisVm.AppsSolution[appName] = thisAppSoln

				delete(remainingApps, appName)                  // no need to handle this application later
				vmCopy.ResidualResources.CpuCore -= cpuCoreStep // subtract the CPU allocated to this app from the VM
			}
		}

		// When no applications are allocated less than cpuCoreStep CPU, we stop Step 1.2.
		if !minCpuFound {
			break
		}

	}

	// do Step 2, allocate CPUs to the remaining Applications one by one.
	for len(remainingApps) > 0 {
		thisAppName, allocatedCpu := distrCpuNextApp(vmCopy, remainingApps, appsOrder)

		// do "floor" with the unit cpuCoreStep
		var actualAllocatedCpu float64
		if math.Abs(allocatedCpu-mymath.UnitRound(allocatedCpu, cpuCoreStep)) < floatDelta {
			// e.g., because of the inaccuracy of binary-floating-point data, a value like 2 may be represented as 1.999999999999, so its floor will be 1 rather than 2, but we need its floor to be 2, so we do this.
			actualAllocatedCpu = mymath.UnitRound(allocatedCpu, cpuCoreStep)
		} else {
			actualAllocatedCpu = mymath.UnitFloor(allocatedCpu, cpuCoreStep)
		}

		// For every application, if the allocated CPUs are more than its request, we reduce them to its request.
		if actualAllocatedCpu > remainingApps[thisAppName].Resources.CpuCore {
			actualAllocatedCpu = remainingApps[thisAppName].Resources.CpuCore
		}

		// the CPU allocation of this application is decided
		thisAppSoln := solnWithCpuThisVm.AppsSolution[thisAppName]
		thisAppSoln.AllocatedCpuCore = actualAllocatedCpu
		solnWithCpuThisVm.AppsSolution[thisAppName] = thisAppSoln

		delete(remainingApps, thisAppName)                     // no need to handle this application later
		vmCopy.ResidualResources.CpuCore -= actualAllocatedCpu // subtract the CPU allocated to this app from the VM
	}

	return solnWithCpuThisVm
}

// For vmCpuWeightedAllocation
/**
In a VM, when we allocate CPUs with weights, the following 2 ways are equivalent:
(1) allocate CPUs for all applications in one round.
(2) allocate CPUs for the applications in different rounds.

*********** Our method to allocate CPUs. *************
(Now, our minCPU is set as 0.1 CPU, but in the future there may be a possibility to change this value.)
My original plan was (1), but then I found that the accuracy is 0.1 CPU, so we should do some round/ceil/floor to 0.1 CPU.
In order to allocate the CPUs weighted-averagely and avoid the overflow after round/ceil/floor, I decide to:
Step 1.
	1.1. Use (1) to allocate CPUs to all applications without round/ceil/floor, picking out the applications that are allocated less than 0.1 CPU.
         Allocate 0.1 CPU to them, which means we do "ceil" to them.
    1.2. Do 1.1 repeatedly in a loop, until no applications are allocated less than 0.1 CPU.
Step 2. Use (2) to allocate CPUs one by one to other applications, and do "floor" to them,
        and for every application, if the allocated CPUs are more than its request, we reduce them to its request.

Proof of the rationality:
We have these facts:
Fact 1. My algorithm can guarantee that the VM's CPUs are certainly enough if every application is allocated 0.1 CPU.
Fact 2. When we do "ceil" to an application, other applications have fewer CPUs.
Fact 3. When we do "floor" to an application, other applications have more CPUs.

Because of Fact 2, after we do Step 1.1, new less-0.1-CPU applications may emerge.
Because of Fact 1, if we do Step 1.2 repeatedly, there will come a time when no new less-0.1-CPU applications emerge.

After finishing Step 1.2, there may be 2 conditions:
Condition 1. All applications are allocated 0.1 CPU and there are no applications left.
Condition 2. There are some applications left, and all of them should be allocated more than 0.1 CPU.
For Condition 1, we have finished the CPU allocation at this time.
For Condition 2, we continue to do Step 2.

Because of Fact 3, in Step 2, every time when we allocate CPUs to one application using "floor", other applications will have more applications, so there will not be an overflow.
Moreover, although the "ceil" and "floor" have impacts on the "weighted-average allocation", in Step 2, we allocate CPUs to applications one by one, so the "weighted-average allocation" will be adjusted every time after a "floor", which can guarantee the weighted-averageness as far as possible.

In conclusion:
Conclusion 1. There will not be an overflow.
Conclusion 2. This method can guarantee the weighted-averageness as far as possible.

*********** Proof: (1) and (2) are equivalent. *************

We need to allocate s CPUs to 4 applications A, B, C, D with the weights a, b, c, d.

If we do (1), allocating CPUs for all applications in one round.
- A will get "s * a / (a+b+c+d)" CPUs.
- B will get "s * b / (a+b+c+d)" CPUs.
- C will get "s * c / (a+b+c+d)" CPUs.
- D will get "s * d / (a+b+c+d)" CPUs.

If we do (2), allocating CPUs for all applications in 4 rounds (each round one application):
- Round 1, we allocate CPUs for A.
  A will get "s * a / (a+b+c+d)" CPUs.
- Round 2, we allocate CPUs for B.
  B will get "(s - s * a / (a+b+c+d)) * b / (b+c+d)" = "(s * (1 - a/(a+b+c+d)) * b / (b+c+d)" = "(s * (b+c+d)/(a+b+c+d) * b / (b+c+d)" = "s * b / (a+b+c+d)" CPUs.
- Round 3, we allocate CPUs for C.
  C will get "(s - s * a / (a+b+c+d) - s * b / (a+b+c+d)) * c / (c+d)" = "(s * (c+d) / (a+b+c+d) * c / (c+d)" = "s * d / (a+b+c+d)" CPUs.
- Round 4, we allocate CPUs for D.
  C will get "(s - s * a / (a+b+c+d) - s * b / (a+b+c+d) - s * d / (a+b+c+d)) * d / d" = "s * d / (a+b+c+d)" CPUs.

Additionally, if we allocate CPUs for all applications in 3 rounds (Round 1 for A, Round 2 for B and C, Round 3 for D):
- Round 1, we allocate CPUs for A.
  A will get "s * a / (a+b+c+d)" CPUs.
- Round 2, we allocate CPUs for B and C.
  B will get "(s - s * a / (a+b+c+d)) * b / (b+c+d)" = "s * b / (a+b+c+d)" CPUs;
  C will get "(s - s * a / (a+b+c+d)) * c / (b+c+d)" = "s * c / (a+b+c+d)" CPUs;
- Round 3, we allocate CPUs for D CPUs.
  D will get "(s - s * a / (a+b+c+d) - s * b / (a+b+c+d) - s * d / (a+b+c+d)) * d / d" = "s * d / (a+b+c+d)" CPUs.

Lastly, if we allocate CPUs for all applications in 3 rounds (Round 1 for A and B, Round 2 for C and D):
- Round 1, we allocate CPUs for A and B.
  A will get "s * a / (a+b+c+d)" CPUs.
  B will get "s * b / (a+b+c+d)" CPUs.
- Round 2, we allocate CPUs for C and D.
  C will get "(s - s * a / (a+b+c+d) - s * b / (a+b+c+d)) * c / (c+d)" = "s * c / (a+b+c+d)" CPUs.
  D will get "(s - s * a / (a+b+c+d) - s * b / (a+b+c+d)) * d / (c+d)" = "s * d / (a+b+c+d)" CPUs.

To sum up, We can see that the results of allocating CPUs in different number of rounds are the same.
*/

// distribute CPUs to some applications weighted by apps' requested CPUs and their priorities.
func distrCpuApps(vm asmodel.K8sNode, apps map[string]asmodel.Application) map[string]float64 {
	checkValidDistriCpu(vm, apps)

	// result to output, key: App names, value: number of CPU cores to allocate to this App.
	var allocationScheme map[string]float64 = make(map[string]float64)

	sumWeight := calcAppsSumWeight(apps)
	for appName, app := range apps {
		thisWeight := calcAppWeight(app)
		allocationScheme[appName] = vm.ResidualResources.CpuCore * (thisWeight / sumWeight)
	}

	return allocationScheme
}

// distribute CPUs to the next application weighted by apps' requested CPUs and their priorities.
func distrCpuNextApp(vm asmodel.K8sNode, apps map[string]asmodel.Application, appsOrder []string) (string, float64) {
	checkValidDistriCpu(vm, apps)

	// get one App to handle
	appsIter := newIterForApps(apps, appsOrder)
	thisAppName := appsIter.nextAppName()

	allocatedCpus := calcCpuOneApp(vm, apps, thisAppName)

	return thisAppName, allocatedCpus
}

// When the code reaches here, the VM's CPUs should be certainly enough if every application is allocated cpuCoreStep CPU.
func checkValidDistriCpu(vm asmodel.K8sNode, apps map[string]asmodel.Application) {
	if vm.ResidualResources.CpuCore < float64(len(apps))*cpuCoreStep {
		panic(fmt.Sprintf("vm.ResidualResources.CpuCore [%f] is not enough for the minimum CPU of [%d] applications.", vm.ResidualResources.CpuCore, len(apps)))
	}
}

// calculate the CPUs to allocate to this application weighted by its requested CPU and its priority
func calcCpuOneApp(vm asmodel.K8sNode, apps map[string]asmodel.Application, thisAppName string) float64 {
	sumWeight := calcAppsSumWeight(apps)
	thisWeight := calcAppWeight(apps[thisAppName])
	return vm.ResidualResources.CpuCore * (thisWeight / sumWeight) // allocate the CPUs proportionally
}

// calculate the sum of the weights of some applications for CPU allocation
func calcAppsSumWeight(apps map[string]asmodel.Application) float64 {
	var sumWeight float64
	for _, app := range apps {
		sumWeight += calcAppWeight(app)
	}
	return sumWeight
}

// calculate the weight of an application for CPU allocation
func calcAppWeight(app asmodel.Application) float64 {
	return app.Resources.CpuCore * float64(app.Priority)
}
