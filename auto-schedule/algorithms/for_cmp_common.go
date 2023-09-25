package algorithms

import (
	"fmt"
	"math"

	"github.com/KeepTheBeats/routing-algorithms/mymath"
	"github.com/KeepTheBeats/routing-algorithms/random"
	"github.com/astaxie/beego"

	asmodel "emcontroller/auto-schedule/model"
	apiv1 "k8s.io/api/core/v1"
)

// the function to refine solutions in the algorithms for comparison.
func CmpRefineSoln(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) (asmodel.Solution, bool) {
	// 1. give the solution node names
	solnWithVm, vmAcceptable := cmpAllocateVms(clouds, apps, appsOrder, soln)
	if !vmAcceptable {
		return asmodel.Solution{}, false
	}
	// 2. Allocate CPU cores
	solnWithCpu, cpuAcceptable := cmpAllocateCpus(clouds, apps, appsOrder, solnWithVm)
	if !cpuAcceptable {
		return asmodel.Solution{}, false
	}
	// 3. Check whether this solution is acceptable.
	if !Acceptable(clouds, apps, appsOrder, solnWithCpu) {
		return asmodel.Solution{}, false
	}
	return solnWithCpu, true
}

// allocate VMs to complete this solution in the algorithms for comparison.
func cmpAllocateVms(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) (asmodel.Solution, bool) {

	// Get the AppsSolution from the input solution. We use the copy function to avoid changing the original solution.
	solnWithVm := asmodel.SolutionCopy(soln)
	// If there is original VmsToCreate in the input solution, we ignore them, as this function will generate a new VM allocation scheme from zero.
	solnWithVm.VmsToCreate = nil

	// We should allocate VMs cloud by cloud.
	for _, cloud := range clouds {
		solnWithVmsThisCloud, allocType := cmpAllocateVmsOneCloud(cloud, apps, appsOrder, soln)
		if allocType == UnAcceptable { // if any cloud cannot accept the scheduled applications, this whole solution is not acceptable.
			return asmodel.Solution{}, false
		}

		// update the part of solution of this cloud to the total solution.
		solnWithVm.Absorb(solnWithVmsThisCloud)
	}

	return solnWithVm, true
}

// allocate VMs in one cloud in the algorithms for comparison.
func cmpAllocateVmsOneCloud(cloud asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) (asmodel.Solution, VmAllocType) {
	// check acceptablity and get the allocation scheme
	solnWithVm, acceptable := cmpResAccOneCloud(cloud, apps, appsOrder, soln)
	if !acceptable {
		return asmodel.Solution{}, UnAcceptable // not acceptable, we cannot accept.
	}

	return solnWithVm, SharedVm // In the algorithms for comparison, we only create a shared VM.
}

func cmpResAccOneCloud(cloud asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, solnWithoutVm asmodel.Solution) (asmodel.Solution, bool) {

	// To check whether the resources are enough, we try to simulate deploying all applications scheduled to this cloud on this cloud.
	appsThisCloud := findAppsOneCloud(cloud, apps, solnWithoutVm)
	var appsToHandle map[string]asmodel.Application = appsThisCloud

	// We put the result to return into this variable.
	solnWithVm := asmodel.GenEmptySoln()
	for appName, _ := range appsToHandle { // the result should only include the solutions for the applications to handle
		solnWithVm.AppsSolution[appName] = asmodel.SasCopy(solnWithoutVm.AppsSolution[appName])
	}

	// To ensure the definiteness/uniqueness of the result, to ensure that the check use the same order with the actual deploy, we use the appsOrder as the fixed order to deploy the applications.
	appsIter := newIterForApps(appsToHandle, appsOrder)

	curAppName := appsIter.nextAppName()
	if len(curAppName) == 0 {
		return solnWithVm, true // this means that no applications are scheduled to this cloud, so the resources are certainly enough.
	}

	// Step 1. use up the resources of existing VMs
	for _, vm := range cloud.K8sNodes { // this vm is only the copy, so the change of it will not affect the original cloud
		// For every VM, if the resources of this vm can meet all rest applications, it means that the resources of this cloud is enough for the applications scheduled to it.
		appNamesToThisVm, meetAllRest := vmResMeetAllRestApps(vm, apps, &curAppName, appsIter.nextAppName, true)

		// put the vm allocation information into the solution.
		for _, appName := range appNamesToThisVm {
			thisAppSoln := solnWithVm.AppsSolution[appName]
			thisAppSoln.K8sNodeName = vm.Name
			solnWithVm.AppsSolution[appName] = thisAppSoln
		}

		if meetAllRest {
			return solnWithVm, true
		}
	}

	// Now, we have tried all existing VMs, and they are not enough for the applications. If this cloud does not support creating a new VM, we directly return false.
	if !cloud.SupportCreateNewVM() {
		return solnWithVm, false
	}

	// Step 2. create a new shared VM
	if !cloud.Resources.Limit.AllMoreThan(cloud.Resources.InUse) {
		return solnWithVm, false // The cloud does not have more resources to create new VMs.
	}

	// On a cloud, when existing VMs do not have enough resources for the applications scheduled here, we try to create a VM with all rest resources.
	vmToCreate := cloud.GetSharedVmToCreate(0, true)
	k8sNodeToCreate := asmodel.GenK8sNodeFromPods(vmToCreate, []apiv1.Pod{})

	appNamesToThisVm, meetAllRest := vmResMeetAllRestApps(k8sNodeToCreate, apps, &curAppName, appsIter.nextAppName, true)

	if meetAllRest {
		// modify single app solutions
		for _, appName := range appNamesToThisVm {
			thisAppSoln := solnWithVm.AppsSolution[appName]
			thisAppSoln.K8sNodeName = vmToCreate.Name
			solnWithVm.AppsSolution[appName] = thisAppSoln
		}
		// modify VmsToCreate
		solnWithVm.VmsToCreate = append(solnWithVm.VmsToCreate, vmToCreate)

		// If vmToCreate have enough resources for the applications, it means this cloud has enough resources.
		return solnWithVm, true
	}

	// the vmToCreate in the above tried 3 cannot meet all rest applications scheduled to this cloud, which means that this cloud does not have enough resources.
	return solnWithVm, false
}

// allocate cpus to complete this solution in the algorithms for comparison.
func cmpAllocateCpus(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, solnWithVm asmodel.Solution) (asmodel.Solution, bool) {

	// This method ues an incremental way to set allocate CPU cores to applications, so before the incremental way, in the base solution, all values of AllocatedCpuCore must be 0.
	for appName, _ := range apps {
		thisAppSoln := solnWithVm.AppsSolution[appName]
		thisAppSoln.AllocatedCpuCore = 0
		solnWithVm.AppsSolution[appName] = thisAppSoln
	}

	// avoiding changing the original solution
	solnWithCpu := asmodel.SolutionCopy(solnWithVm)

	// We should allocate CPUs cloud by cloud.
	for _, cloud := range clouds {
		solnWithCpuThisCloud, acceptable := cmpAllocateCpusOneCloud(cloud, apps, appsOrder, solnWithVm)
		if !acceptable { // if any cloud cannot accept the scheduled applications, this whole solution is not acceptable.
			return asmodel.Solution{}, false
		}

		// update the solution of this cloud to the total solution.
		solnWithCpu.Absorb(solnWithCpuThisCloud)
	}

	return solnWithCpu, true
}

// allocate cpus in one cloud in the algorithms for comparison.
func cmpAllocateCpusOneCloud(cloud asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, solnWithVm asmodel.Solution) (asmodel.Solution, bool) {
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
		vm := getVmByName(vmName, cloud, solnWithVm)                                         // handle this vm
		solnWithCpuThisVm := cmpAllocateCpusOneVm(vm, apps, appsOrder, appNames, solnWithVm) // allocate CPUs to applications on this VM

		solnWithCpu.Absorb(solnWithCpuThisVm) // combine the solution of this VM into the solution of this cloud.
	}

	return solnWithCpu, true
}

// allocate CPUs to the applications on a VM in the algorithms for comparison.
func cmpAllocateCpusOneVm(vm asmodel.K8sNode, apps map[string]asmodel.Application, appsOrder []string, appNamesThisVm []string, solnWithVm asmodel.Solution) asmodel.Solution {
	appsThisVm := filterAppsByNames(appNamesThisVm, apps) // get the applications scheduled to this VM
	return cmpVmCpuAllocation(vm, appsThisVm, appNamesThisVm, appsOrder, solnWithVm)
}

// allocate the CPUs of a VM to the applications scheduled to it, in a completely random way.
func cmpVmCpuAllocation(vm asmodel.K8sNode, appsThisVm map[string]asmodel.Application, appNamesThisVm []string, appsOrder []string, solnWithVm asmodel.Solution) asmodel.Solution {
	var solnWithCpuThisVm asmodel.Solution = asmodel.GenEmptySoln()
	for appName, _ := range appsThisVm { // the result should only include the solutions for the applications to handle
		solnWithCpuThisVm.AppsSolution[appName] = asmodel.SasCopy(solnWithVm.AppsSolution[appName])
	}

	// copy and avoid changing the original variable.
	vmCopy := asmodel.K8sNodeCopy(vm)

	// First round, we allocate one CPU core to every application to meet the minimum requirement
	for appName, _ := range appsThisVm {
		if vmCopy.ResidualResources.CpuCore >= cpuCoreStep { // if this VM has residual CPUs, we allocate more CPUs to this app.
			thisAppSoln := solnWithCpuThisVm.AppsSolution[appName]
			thisAppSoln.AllocatedCpuCore += cpuCoreStep
			solnWithCpuThisVm.AppsSolution[appName] = thisAppSoln
			vmCopy.ResidualResources.CpuCore -= cpuCoreStep
		}
	}

	// Second round, we allocate a random number of CPU cores to every application. The application order is also random.
	appNamesCopy := make([]string, len(appNamesThisVm)) // copy to avoid changing the original variable.
	copy(appNamesCopy, appNamesThisVm)
	for len(appNamesCopy) > 0 && vmCopy.ResidualResources.CpuCore > floatDelta { // when the CPU cores are used up, we stop this allocation
		randIdx := random.RandomInt(0, len(appNamesCopy)-1)
		randAppName := appNamesCopy[randIdx]
		appNamesCopy = append(appNamesCopy[:randIdx], appNamesCopy[randIdx+1:]...)

		// randomly choose CPU to allocate
		randCpu := float64(random.RandomInt(0, int(mymath.UnitRound(vmCopy.ResidualResources.CpuCore, 1))))

		thisAppSoln := solnWithCpuThisVm.AppsSolution[randAppName]
		thisAppSoln.AllocatedCpuCore += randCpu
		solnWithCpuThisVm.AppsSolution[randAppName] = thisAppSoln
		vmCopy.ResidualResources.CpuCore -= randCpu

	}

	// fix the inaccuracy of float
	for appName, _ := range appsThisVm {
		thisAppSoln := solnWithCpuThisVm.AppsSolution[appName]
		if math.Abs(thisAppSoln.AllocatedCpuCore-mymath.UnitRound(thisAppSoln.AllocatedCpuCore, cpuCoreStep)) < floatDelta {
			thisAppSoln.AllocatedCpuCore = mymath.UnitRound(thisAppSoln.AllocatedCpuCore, cpuCoreStep)
		} else {
			thisAppSoln.AllocatedCpuCore = mymath.UnitFloor(thisAppSoln.AllocatedCpuCore, cpuCoreStep)
		}
		solnWithCpuThisVm.AppsSolution[appName] = thisAppSoln
	}

	return solnWithCpuThisVm
}

// Deprecated: this method is not random enough.
// allocate the CPUs of a VM to the applications scheduled to it, in a completely random way. The difference between this and oldCmpVmCpuAllocation3 is that this allocation method retain some resources for the possible future applications, which is more reasonable.
func oldCmpVmCpuAllocation4(vm asmodel.K8sNode, appsThisVm map[string]asmodel.Application, appNamesThisVm []string, appsOrder []string, solnWithVm asmodel.Solution) asmodel.Solution {
	var solnWithCpuThisVm asmodel.Solution = asmodel.GenEmptySoln()
	for appName, _ := range appsThisVm { // the result should only include the solutions for the applications to handle
		solnWithCpuThisVm.AppsSolution[appName] = asmodel.SasCopy(solnWithVm.AppsSolution[appName])
	}

	// copy and avoid changing the original variable.
	vmCopy := asmodel.K8sNodeCopy(vm)

	// firstly, we allocate one CPU core to every application to meet the minimum requirement
	for appName, _ := range appsThisVm {
		if vmCopy.ResidualResources.CpuCore >= cpuCoreStep { // if this VM has residual CPUs, we allocate more CPUs to this app.
			thisAppSoln := solnWithCpuThisVm.AppsSolution[appName]
			thisAppSoln.AllocatedCpuCore += cpuCoreStep
			solnWithCpuThisVm.AppsSolution[appName] = thisAppSoln
			vmCopy.ResidualResources.CpuCore -= cpuCoreStep
		}
	}

	// Then randomly allocate CPU cores to applications until all CPU cores are allocated
	// we only allocate at most len(appNamesThisVm) times, because we should not use up all CPUs, because we should retain some for the possible future applications.
	for i := 0; i < len(appNamesThisVm) && vmCopy.ResidualResources.CpuCore > floatDelta; i++ {
		// randomly pick an app to allocate CPU
		randAppName := appNamesThisVm[random.RandomInt(0, len(appNamesThisVm)-1)]
		// randomly choose CPU to allocate
		randCpu := float64(random.RandomInt(1, int(mymath.UnitRound(vmCopy.ResidualResources.CpuCore, 1))))

		thisAppSoln := solnWithCpuThisVm.AppsSolution[randAppName]
		thisAppSoln.AllocatedCpuCore += randCpu
		solnWithCpuThisVm.AppsSolution[randAppName] = thisAppSoln
		vmCopy.ResidualResources.CpuCore -= randCpu

	}

	// fix the inaccuracy of float
	for appName, _ := range appsThisVm {
		thisAppSoln := solnWithCpuThisVm.AppsSolution[appName]
		if math.Abs(thisAppSoln.AllocatedCpuCore-mymath.UnitRound(thisAppSoln.AllocatedCpuCore, cpuCoreStep)) < floatDelta {
			thisAppSoln.AllocatedCpuCore = mymath.UnitRound(thisAppSoln.AllocatedCpuCore, cpuCoreStep)
		} else {
			thisAppSoln.AllocatedCpuCore = mymath.UnitFloor(thisAppSoln.AllocatedCpuCore, cpuCoreStep)
		}
		solnWithCpuThisVm.AppsSolution[appName] = thisAppSoln
	}

	return solnWithCpuThisVm
}

// Deprecated: this method uses up all CPUs of a VM, but it should still retain some CPUs for the possible future applications.
// allocate the CPUs of a VM to the applications scheduled to it, in a completely random way.
func oldCmpVmCpuAllocation3(vm asmodel.K8sNode, appsThisVm map[string]asmodel.Application, appNamesThisVm []string, appsOrder []string, solnWithVm asmodel.Solution) asmodel.Solution {
	var solnWithCpuThisVm asmodel.Solution = asmodel.GenEmptySoln()
	for appName, _ := range appsThisVm { // the result should only include the solutions for the applications to handle
		solnWithCpuThisVm.AppsSolution[appName] = asmodel.SasCopy(solnWithVm.AppsSolution[appName])
	}

	// copy and avoid changing the original variable.
	vmCopy := asmodel.K8sNodeCopy(vm)

	// firstly, we allocate one CPU core to every application to meet the minimum requirement
	for appName, _ := range appsThisVm {
		if vmCopy.ResidualResources.CpuCore >= cpuCoreStep { // if this VM has residual CPUs, we allocate more CPUs to this app.
			thisAppSoln := solnWithCpuThisVm.AppsSolution[appName]
			thisAppSoln.AllocatedCpuCore += cpuCoreStep
			solnWithCpuThisVm.AppsSolution[appName] = thisAppSoln
			vmCopy.ResidualResources.CpuCore -= cpuCoreStep
		}
	}

	// Then randomly allocate CPU cores to applications until all CPU cores are allocated
	for vmCopy.ResidualResources.CpuCore > floatDelta {
		// randomly pick an app to allocate CPU
		randAppName := appNamesThisVm[random.RandomInt(0, len(appNamesThisVm)-1)]
		// randomly choose CPU to allocate
		randCpu := float64(random.RandomInt(1, int(mymath.UnitRound(vmCopy.ResidualResources.CpuCore, 1))))

		thisAppSoln := solnWithCpuThisVm.AppsSolution[randAppName]
		thisAppSoln.AllocatedCpuCore += randCpu
		solnWithCpuThisVm.AppsSolution[randAppName] = thisAppSoln
		vmCopy.ResidualResources.CpuCore -= randCpu

	}

	// fix the inaccuracy of float
	for appName, _ := range appsThisVm {
		thisAppSoln := solnWithCpuThisVm.AppsSolution[appName]
		if math.Abs(thisAppSoln.AllocatedCpuCore-mymath.UnitRound(thisAppSoln.AllocatedCpuCore, cpuCoreStep)) < floatDelta {
			thisAppSoln.AllocatedCpuCore = mymath.UnitRound(thisAppSoln.AllocatedCpuCore, cpuCoreStep)
		} else {
			thisAppSoln.AllocatedCpuCore = mymath.UnitFloor(thisAppSoln.AllocatedCpuCore, cpuCoreStep)
		}
		solnWithCpuThisVm.AppsSolution[appName] = thisAppSoln
	}

	return solnWithCpuThisVm
}

// Deprecated: this method allocates CPU cores averagely, which is good for some scenarios but bad for some others, so we should not use this, and we should do it completely randomly.
// allocate the CPUs of a VM to the applications scheduled to it, in an average way.
func oldCmpVmCpuAllocation2(vm asmodel.K8sNode, appsThisVm map[string]asmodel.Application, appsOrder []string, solnWithVm asmodel.Solution) asmodel.Solution {
	var solnWithCpuThisVm asmodel.Solution = asmodel.GenEmptySoln()
	for appName, _ := range appsThisVm { // the result should only include the solutions for the applications to handle
		solnWithCpuThisVm.AppsSolution[appName] = asmodel.SasCopy(solnWithVm.AppsSolution[appName])
	}

	// copy and avoid changing the original variable.
	vmCopy := asmodel.K8sNodeCopy(vm)

	var allocated bool = true
	for allocated { // if this VM does not have residual CPUs or all applications cannot use more CPUs, this will be false.
		allocated = false

		// in every loop, we try to allocate cpuCoreStep CPUs to every application on this VM.
		for appName, _ := range appsThisVm {

			if vmCopy.ResidualResources.CpuCore >= cpuCoreStep { // if this VM has residual CPUs, we allocate more CPUs to this app.
				thisAppSoln := solnWithCpuThisVm.AppsSolution[appName]
				thisAppSoln.AllocatedCpuCore += cpuCoreStep
				solnWithCpuThisVm.AppsSolution[appName] = thisAppSoln

				vmCopy.ResidualResources.CpuCore -= cpuCoreStep

				allocated = true
			}
		}
	}

	// fix the inaccuracy of float
	for appName, _ := range appsThisVm {
		thisAppSoln := solnWithCpuThisVm.AppsSolution[appName]
		if math.Abs(thisAppSoln.AllocatedCpuCore-mymath.UnitRound(thisAppSoln.AllocatedCpuCore, cpuCoreStep)) < floatDelta {
			thisAppSoln.AllocatedCpuCore = mymath.UnitRound(thisAppSoln.AllocatedCpuCore, cpuCoreStep)
		} else {
			thisAppSoln.AllocatedCpuCore = mymath.UnitFloor(thisAppSoln.AllocatedCpuCore, cpuCoreStep)
		}
		solnWithCpuThisVm.AppsSolution[appName] = thisAppSoln
	}

	return solnWithCpuThisVm
}

// Deprecated: this CPU allocation is too intelligent, because it considers the CPU requirement of applications.
// allocate the CPUs of a VM to the applications scheduled to it, in a random and best-effort way.
func oldCmpVmCpuAllocation(vm asmodel.K8sNode, appsThisVm map[string]asmodel.Application, appsOrder []string, solnWithVm asmodel.Solution) asmodel.Solution {
	var solnWithCpuThisVm asmodel.Solution = asmodel.GenEmptySoln()
	for appName, _ := range appsThisVm { // the result should only include the solutions for the applications to handle
		solnWithCpuThisVm.AppsSolution[appName] = asmodel.SasCopy(solnWithVm.AppsSolution[appName])
	}

	// copy and avoid changing the original variable.
	vmCopy := asmodel.K8sNodeCopy(vm)

	var allocated bool = true
	for allocated { // if this VM does not have residual CPUs or all applications cannot use more CPUs, this will be false.
		allocated = false

		// in every loop, we try to allocate cpuCoreStep CPUs to every application on this VM.
		for appName, app := range appsThisVm {

			if solnWithCpuThisVm.AppsSolution[appName].AllocatedCpuCore+cpuCoreStep <= app.Resources.CpuCore && vmCopy.ResidualResources.CpuCore >= cpuCoreStep { // if this VM has residual CPUs, and this application can use more CPUs, we allocate more CPUs to this app.
				thisAppSoln := solnWithCpuThisVm.AppsSolution[appName]
				thisAppSoln.AllocatedCpuCore += cpuCoreStep
				solnWithCpuThisVm.AppsSolution[appName] = thisAppSoln

				vmCopy.ResidualResources.CpuCore -= cpuCoreStep

				allocated = true
			}
		}
	}

	// fix the inaccuracy of float
	for appName, _ := range appsThisVm {
		thisAppSoln := solnWithCpuThisVm.AppsSolution[appName]
		if math.Abs(thisAppSoln.AllocatedCpuCore-mymath.UnitRound(thisAppSoln.AllocatedCpuCore, cpuCoreStep)) < floatDelta {
			thisAppSoln.AllocatedCpuCore = mymath.UnitRound(thisAppSoln.AllocatedCpuCore, cpuCoreStep)
		} else {
			thisAppSoln.AllocatedCpuCore = mymath.UnitFloor(thisAppSoln.AllocatedCpuCore, cpuCoreStep)
		}
		solnWithCpuThisVm.AppsSolution[appName] = thisAppSoln
	}

	return solnWithCpuThisVm
}

// for comparison
// Randomly explore all possibilities of 2-point crossover, to try to get an acceptable solution. If this function cannot find an acceptable solution after trying all possibilities, it will return the original 2 chromosomes without doing crossover. The only difference with AllPossTwoPointCrossover is that this function use CmpRefineSoln instead of RefineSoln.
func CmpAllPossTwoPointCrossover(firstChromosome asmodel.Solution, secondChromosome asmodel.Solution, clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) (asmodel.Solution, asmodel.Solution) {
	// in our unit tests, we will set both the input cloud and apps as nil
	var testMode bool = clouds == nil && apps == nil

	if len(firstChromosome.AppsSolution) != len(secondChromosome.AppsSolution) || len(firstChromosome.AppsSolution) != len(appsOrder) {
		panic(fmt.Sprintf("len(firstChromosome.AppsSolution) is %d; len(secondChromosome.AppsSolution) is %d; len(appsOrder) is %d. They should be equal.", len(firstChromosome.AppsSolution), len(secondChromosome.AppsSolution), len(appsOrder)))
	}

	// the number of genes in a chromosome, also the number of applications to schedule.
	geneNumber := len(firstChromosome.AppsSolution)

	/**
	We exchange the genes of the 2 chromosomes in the closed interval [point1, point2]. The width of this closed interval ranges from 1 to geneNumber-1.
	The following loop randomly traverse all possibility of different width and point1 which determine point2.
	*/

	// build an array to help select point widths randomly
	var possiblePointWidths []int
	for pointWidth := 1; pointWidth <= geneNumber-1; pointWidth++ {
		possiblePointWidths = append(possiblePointWidths, pointWidth)
	}
	for len(possiblePointWidths) > 0 {
		// randomly select a possible point width, and then remove it from the array, in order not to select it again.
		widthIdx := random.RandomInt(0, len(possiblePointWidths)-1)
		pointWidth := possiblePointWidths[widthIdx]
		possiblePointWidths = append(possiblePointWidths[:widthIdx], possiblePointWidths[widthIdx+1:]...)
		if testMode {
			beego.Info("pointWidth is:", pointWidth) // for debug
		}

		// build an array to help select point1 randomly
		var possiblePoint1 []int
		for point1 := 0; calcPoint2(point1, pointWidth) < geneNumber; point1++ {
			possiblePoint1 = append(possiblePoint1, point1)
		}
		for len(possiblePoint1) > 0 {
			// randomly select a possible point1, and then remove it from the array, in order not to select it again.
			pointIdx := random.RandomInt(0, len(possiblePoint1)-1)
			point1 := possiblePoint1[pointIdx]
			possiblePoint1 = append(possiblePoint1[:pointIdx], possiblePoint1[pointIdx+1:]...)

			// calculate point 2 by point 1
			point2 := calcPoint2(point1, pointWidth)

			if testMode {
				beego.Info("point1, point2:", point1, point2) // for debug
			} else {

				/**
				Then we do crossover with the randomly selected point1 and point2.
				We set the tryFunc here, because with this the unit tests will be easier to make.
				*/

				// if in this possibility the 2 crossovered chromosomes are acceptable, return them.
				crossoveredChromosome1, crossoveredChromosome2 := twoPointCrossover(firstChromosome, secondChromosome, appsOrder, point1, point2)

				// refine the 2 crossovered chromosomes and check whether they are acceptable. If both of them are acceptable, we return them as the result.
				if crossoveredChromosome1, acceptable1 := CmpRefineSoln(clouds, apps, appsOrder, crossoveredChromosome1); acceptable1 {
					if crossoveredChromosome2, acceptable2 := CmpRefineSoln(clouds, apps, appsOrder, crossoveredChromosome2); acceptable2 {
						return crossoveredChromosome1, crossoveredChromosome2
					}
				}
			}

		}
		if testMode {
			fmt.Println() // for debug
		}
	}

	return firstChromosome, secondChromosome
}
