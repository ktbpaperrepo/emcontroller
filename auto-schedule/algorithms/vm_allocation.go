/*
The functions to allocate Virtual Machines to applications inside clouds.
*/

package algorithms

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"

	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

type VmAllocType int

const (
	SharedVm     VmAllocType = iota // This cloud only create at most 1 shared VM for the applications.
	DedicatedVms                    // This cloud create dedicated VMs for the applications with the Max Priority.
	UnAcceptable                    // This cloud does not have enough resources for the applications allocated to it.
)

// allocate VMs to complete this solution.
func allocateVms(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) (asmodel.Solution, bool) {

	// Get the AppsSolution from the input solution. We use the copy function to avoid changing the original solution.
	solnWithVm := asmodel.SolutionCopy(soln)
	// If there is original VmsToCreate in the input solution, we ignore them, as this function will generate a new VM allocation scheme from zero.
	solnWithVm.VmsToCreate = nil

	// We should allocate VMs cloud by cloud.
	for _, cloud := range clouds {
		solnWithVmsThisCloud, allocType := allocateVmsOneCloud(cloud, apps, appsOrder, soln)
		if allocType == UnAcceptable { // if any cloud cannot accept the scheduled applications, this whole solution is not acceptable.
			return asmodel.Solution{}, false
		}

		// update the part of solution of this cloud to the total solution.
		solnWithVm.Absorb(solnWithVmsThisCloud)
	}

	return solnWithVm, true
}

// allocate VMs in one cloud
func allocateVmsOneCloud(cloud asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) (asmodel.Solution, VmAllocType) {
	// check shared allocation at first
	solnWithSharedVm, sharedAcceptable := resAccOneCloudSharedVm(cloud, apps, appsOrder, soln, AllPriApps)
	if !sharedAcceptable {
		return asmodel.Solution{}, UnAcceptable // shared not acceptable, we cannot accept.
	}

	// Then, check dedicated allocation
	solnWithDedVms, dedAcceptable := resAccOneCloudDedicatedVms(cloud, apps, appsOrder, soln)
	if !dedAcceptable {
		return solnWithSharedVm, SharedVm // shared acceptable but dedicated not acceptable, we can only use shared.
	}

	// Actually, I think, theoretically, there is a possibility that some solutions can meet dedicated allocation but cannot meet shared allocation, because of the order of applications, but we still check shared scheme first and then dedicated, and if a solution cannot meet the shared scheme, we will simply consider that it cannot meet the dedicated scheme without a check of it. This can make the algorithm less complicated.

	return solnWithDedVms, DedicatedVms // shared acceptable and dedicated acceptable, we use dedicated.
}

// Check whether the resources of one cloud are enough for the applications scheduled to it when only creating a shared VM.
// Also return the solution with the vm allocation scheme.
func resAccOneCloudSharedVm(cloud asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, solnWithoutVm asmodel.Solution, appTypeToHandle AppType) (asmodel.Solution, bool) {

	// To check whether the resources are enough, we try to simulate deploying all applications scheduled to this cloud on this cloud.
	appsThisCloud := findAppsOneCloud(cloud, apps, solnWithoutVm)

	// filter the type of applications that we will handle
	var appsToHandle map[string]asmodel.Application
	switch appTypeToHandle {
	case AllPriApps:
		appsToHandle = appsThisCloud
	case MaxPriApps:
		appsToHandle = filterMaxPriApps(appsThisCloud)
	case NotMaxPriApps:
		appsToHandle = filterOutMaxPriApps(appsThisCloud)
	default:
		panic(fmt.Sprintf("invalid appTypeToHandle: %v", appTypeToHandle))
	}

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

	/**
	NOTE:
	On a cloud, when existing VMs do not have enough resources for the applications scheduled here, we have 3 possible choices to create a new VM:
	1. a VM with 50% resources; if 50% is not enough, try 3.
	2. a VM with 30% resources; if 30% is not enough, try 3.
	3. a VM with all rest resources; if all rest is not enough, it means this solution is not acceptable.
	*/
	cloudLeastResPct := cloud.Resources.LeastRemainPct()

	if cloudLeastResPct > biggerVmResPct { // try 1
		vmToCreate := cloud.GetSharedVmToCreate(biggerVmResPct, false)
		k8sNodeToCreate := asmodel.GenK8sNodeFromPods(vmToCreate, []apiv1.Pod{})

		// if the 1 does not work, we should do 3, so we should copy the iter and curAppName, avoiding changing the environment.
		iterCopy := appsIter.Copy()
		curAppNameCopy := curAppName

		appNamesToThisVm, meetAllRest := vmResMeetAllRestApps(k8sNodeToCreate, apps, &curAppNameCopy, iterCopy.nextAppName, true)

		// Only if this VM can meet all rest applications, we apply this vm scheme and put the vm allocation information into the solution.
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

	} else if cloudLeastResPct > smallerVmResPct { // try 2
		vmToCreate := cloud.GetSharedVmToCreate(smallerVmResPct, false)
		k8sNodeToCreate := asmodel.GenK8sNodeFromPods(vmToCreate, []apiv1.Pod{})

		// if the 2 does not work, we should do 3, so we should copy the iter and curAppName, avoiding changing the environment.
		iterCopy := appsIter.Copy()
		curAppNameCopy := curAppName

		appNamesToThisVm, meetAllRest := vmResMeetAllRestApps(k8sNodeToCreate, apps, &curAppNameCopy, iterCopy.nextAppName, true)

		// Only if this VM can meet all rest applications, we apply this vm scheme and put the vm allocation information into the solution.
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

	}

	// If the rest percentage of this cloud's resources are less than 30%,
	// or if the vmToCreate in the above tried 1 or 2 cannot meet all rest applications scheduled to this cloud,
	// we try 3.
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

// Check whether the resources of one cloud are enough for the applications scheduled to it when creating dedicated VMs.
// Also return the solution with the vm allocation scheme.
func resAccOneCloudDedicatedVms(cloud asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, solnWithoutVm asmodel.Solution) (asmodel.Solution, bool) {
	// To use dedicated VMs, creating new VMs is necessary.
	if !cloud.SupportCreateNewVM() {
		return asmodel.Solution{}, false
	}

	appsThisCloud := findAppsOneCloud(cloud, apps, solnWithoutVm)

	// We put the result to return into this variable.
	solnWithVm := asmodel.GenEmptySoln()
	for appName, _ := range appsThisCloud { // the result should only include the solutions for the applications to handle
		solnWithVm.AppsSolution[appName] = asmodel.SasCopy(solnWithoutVm.AppsSolution[appName])
	}

	// 1. try to create dedicated VMs for Max-Priority applications

	// find out the Max-Priority applications, because the dedicated VMs are for them.
	// At here, we do not need to keep the order of the applications, so we do not need to use iterator
	var maxPriApps map[string]asmodel.Application = filterMaxPriApps(appsThisCloud)

	// group the max-priority applications according to their dependencies. The applications with dependencies should be in the same group, and we will create one dedicated VM for one group.
	maxPriAppsGroups := groupByDep(maxPriApps)
	simulatedCloud := asmodel.CloudCopy(cloud) // avoid changing the original cloud variable
	dedicatedVmsToCreate := getDedicatedVmsToCreate(&simulatedCloud, apps, maxPriAppsGroups)

	// put the app scheduling vm information into the solution
	for i := 0; i < len(maxPriAppsGroups); i++ {
		vmToCreateName := dedicatedVmsToCreate[i].Name // this group of apps are scheduled to this VM
		for j := 0; j < len(maxPriAppsGroups[i]); j++ {
			appName := maxPriAppsGroups[i][j]
			thisAppSoln := solnWithVm.AppsSolution[appName]
			thisAppSoln.K8sNodeName = vmToCreateName
			solnWithVm.AppsSolution[appName] = thisAppSoln
		}
	}
	// put the information of VMs to create into the solution
	solnWithVm.VmsToCreate = append(solnWithVm.VmsToCreate, dedicatedVmsToCreate...)

	// after creating dedicated vms, if any type of resources overflows, it means that this cloud does not have enough resources to create the dedicated VMs, so we return false.
	if simulatedCloud.Resources.Overflow() {
		return asmodel.Solution{}, false
	}

	// 2. try to deploy the rest of the applications on existing VMs and a shared VM.
	solnNotMaxPri, sharedAcceptable := resAccOneCloudSharedVm(simulatedCloud, apps, appsOrder, solnWithoutVm, NotMaxPriApps)
	if !sharedAcceptable {
		return asmodel.Solution{}, false
	}
	solnWithVm.Absorb(solnNotMaxPri) // combing the solution about max-priority and non-max-priority applications.

	return solnWithVm, true
}

// In some conditions, we create a dedicated VM for each application group.
// NOTE: This function changes the cloud, and we need to pass a copy to avoid the change.
func getDedicatedVmsToCreate(cloud *asmodel.Cloud, apps map[string]asmodel.Application, appGroups [][]string) []models.IaasVm {
	var vmsToCreate []models.IaasVm = make([]models.IaasVm, len(appGroups))

	for i := 0; i < len(appGroups); i++ {
		vmsToCreate[i] = getDedVmOneGroup(*cloud, apps, appGroups[i]) // The apps in appGroups[i] are scheduled to vmsToCreate[i].
		simulateCreateVm(cloud, vmsToCreate[i], apps, appGroups[i])   // subtract the resources and add the new vm record to this cloud
	}

	return vmsToCreate
}

// get the dedicated vm to create for an app group
func getDedVmOneGroup(cloud asmodel.Cloud, apps map[string]asmodel.Application, appGroup []string) models.IaasVm {

	// calculate the total resources that are needed by this group of applications.
	// for max-priority applications on dedicated VMs, we regulate that it occupy its requested CPU cores rather than minCpu
	var neededAvailRes asmodel.AppResources = calcNeededRes(apps, appGroup, false)

	// As every VM has some reserved resources, so according to the reserved resources, we calculate the needed available resources of the VM to create.
	var deDVmToCreate models.IaasVm = models.IaasVm{
		Name:    cloud.GetNameVmToCreate(),
		Cloud:   cloud.Name,
		VCpu:    models.CalcVmTotalVcpu(neededAvailRes.CpuCore),
		Ram:     models.CalcVmTotalRamMiB(neededAvailRes.Memory),
		Storage: models.CalcVmTotalStorGiB(neededAvailRes.Storage),
	}

	return deDVmToCreate
}

// simulate to create a vm on a simulated cloud
func simulateCreateVm(simCloud *asmodel.Cloud, vmToCreate models.IaasVm, apps map[string]asmodel.Application, appGroup []string) {
	// Add the resources to the InUse
	simCloud.Resources.InUse.VCpu += vmToCreate.VCpu
	simCloud.Resources.InUse.Ram += vmToCreate.Ram
	simCloud.Resources.InUse.Storage += vmToCreate.Storage

	// add this K8sNode into the simulated cloud.
	convertedK8sNode := asmodel.GenK8sNodeFromApps(vmToCreate, apps, appGroup)
	simCloud.K8sNodes = append(simCloud.K8sNodes, convertedK8sNode)
}

// check whether the residual resources of a VM can support an application
func isResEnough(vm asmodel.K8sNode, app asmodel.Application, minCpu bool) bool {
	// In some conditions, the occupied CPU can be considered as the original requirement of the application.
	var cpuToOccupy float64 = app.Resources.CpuCore

	// CPU is a soft resource, so we think that cpuCoreStep CPU is the minimum requirement for each application.
	// In some conditions, the occupied CPU can be considered as the minimum requirement.
	if minCpu {
		cpuToOccupy = cpuCoreStep
	}

	return vm.ResidualResources.CpuCore >= cpuToOccupy &&
		vm.ResidualResources.Memory >= app.Resources.Memory &&
		vm.ResidualResources.Storage >= app.Resources.Storage
}

// subtract the resources required by an application from a VM
func subRes(vm *asmodel.K8sNode, app asmodel.Application, minCpu bool) {
	// In some conditions, the occupied CPU can be considered as the original requirement of the application.
	var cpuToOccupy float64 = app.Resources.CpuCore

	// CPU is a soft resource, so we think that cpuCoreStep CPU is the minimum requirement for each application.
	// In some conditions, the occupied CPU can be considered as the minimum requirement.
	if minCpu {
		cpuToOccupy = cpuCoreStep
	}
	vm.ResidualResources.CpuCore -= cpuToOccupy
	vm.ResidualResources.Memory -= app.Resources.Memory
	vm.ResidualResources.Storage -= app.Resources.Storage
}

// check whether the resources of the input VM can support all rest applications. Also return the applications that are scheduled to this VM.
func vmResMeetAllRestApps(vm asmodel.K8sNode, apps map[string]asmodel.Application, curAppName *string, nextAppNameFunc func() string, minCpu bool) ([]string, bool) {

	var appNamesToThisVm []string // the application names that are scheduled to this VM

	// we loop until the resources of this VM is used up.
	for isResEnough(vm, apps[*curAppName], minCpu) {
		// simulate deploying this application on this VM.
		subRes(&vm, apps[*curAppName], minCpu)
		appNamesToThisVm = append(appNamesToThisVm, *curAppName)

		// After the current applications is deployed, we go to the next application.
		// curAppName is a pointer, so the value change will affect the variables outside this function.
		// nextAppNameFunc should be a method of *iterForApps, which is also a pointer, so the value will also affect it outside this function.
		*curAppName = nextAppNameFunc()
		if len(*curAppName) == 0 {
			return appNamesToThisVm, true // this means that all applications are already deployed, so the resources are enough.
		}
	}

	return appNamesToThisVm, false // this means that the rest of the VM's resources can not meet the next application, so the resources are not enough
}
