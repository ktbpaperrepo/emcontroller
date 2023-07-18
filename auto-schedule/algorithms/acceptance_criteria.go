package algorithms

import (
	asmodel "emcontroller/auto-schedule/model"
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
	for _, cloud := range clouds { // check every cloud
		if !resAccOneCloud(cloud, apps, appsOrder, soln) {
			return false
		}
	}

	// all clouds passed
	return true
}

// Check whether the resources of one cloud are enough for the applications scheduled to it.
func resAccOneCloud(cloud asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) bool {

	// To check whether the resources are enough, we try to simulate deploying all applications scheduled to this cloud on this cloud.
	appsThisCloud := findAppsOneCloud(cloud, apps, soln)

	// To ensure the definiteness/uniqueness of the result, to ensure that the check use the same order with the actual deploy, we use the appsOrder as the fixed order to deploy the applications.
	appsThisCloudIter := newAppOneCloudIter(appsThisCloud, appsOrder)
	curAppName := appsThisCloudIter.nextAppName()
	if len(curAppName) == 0 {
		return true // this means that no applications are scheduled to this cloud, so the resources are certainly enough.
	}

	// TODO: test cases are finished

	// Step 1. use up the resources of existing VMs
	for _, vm := range cloud.K8sNodes { // this vm is only the copy, so the change of it will not affect the original cloud
		// For every VM, if the resources of this vm can meet all rest applications, it means that the resources of this cloud is enough for the applications scheduled to it.
		if vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter) {
			return true
		}
	}

	// Now, we have tried all existing VMs, and they are not enough for the applications. If this cloud does not support creating a new VM, we directly return false.
	if !cloud.SupportCreateNewVM() {
		return false
	}

	// Step 2. create a new VM
	if !cloud.Resources.Limit.AllMoreThan(cloud.Resources.InUse) {
		return false // The cloud does not have more resources to create new VMs.
	}

	/* On a cloud, when existing VMs do not have enough resources for the applications scheduled here, we have 3 possible choices to create a new VM:
	1. a VM with 50% resources; if 50% is not enough, try 3.
	2. a VM with 30% resources; if 30% is not enough, try 3.
	3. a VM with all rest resources; if all rest is not enough, it means this solution is not acceptable.
	*/
	cloudLeastResPct := cloud.Resources.LeastRemainPct()

	if cloudLeastResPct > biggerVmResPct { // try 1
		vmToCreate := cloud.GetInfoVmToCreate(biggerVmResPct)

		// if the 1 does not work, we should do 3, so we should use copy the iter and curAppName, avoiding changing the environment.
		iterCopy := appsThisCloudIter.Copy()
		curAppNameCopy := curAppName

		// If vmToCreate have enough resources for the applications, it means this cloud has enough resources.
		if vmResMeetAllRestApps(vmToCreate, apps, &curAppNameCopy, iterCopy) {
			return true
		}
	} else if cloudLeastResPct > smallerVmResPct { // try 2
		vmToCreate := cloud.GetInfoVmToCreate(smallerVmResPct)

		// if the 2 does not work, we should do 3, so we should use copy the iter and curAppName, avoiding changing the environment.
		iterCopy := appsThisCloudIter.Copy()
		curAppNameCopy := curAppName

		if vmResMeetAllRestApps(vmToCreate, apps, &curAppNameCopy, iterCopy) {
			return true
		}
	}

	// If the rest percentage of this cloud's resources are less than 30%,
	// or if the vmToCreate in the above tried 1 or 2 cannot meet all rest applications scheduled to this cloud,
	// we try 3.
	vmToCreate := cloud.GetInfoVmToCreate(allRestVmResPct)
	if vmResMeetAllRestApps(vmToCreate, apps, &curAppName, appsThisCloudIter) {
		return true
	}

	// the vmToCreate in the above tried 3 cannot meet all rest applications scheduled to this cloud, which means that this cloud does not have enough resources.
	return false
}

// check whether the residual resources of a VM can support an application
func isResEnough(vm asmodel.K8sNode, app asmodel.Application) bool {
	// CPU is a soft resource, so we think that cpuCoreStep CPU is the minimum requirement for each application
	return vm.ResidualResources.CpuCore >= cpuCoreStep &&
		vm.ResidualResources.Memory >= app.Resources.Memory &&
		vm.ResidualResources.Storage >= app.Resources.Storage
}

// subtract the resources required by an application from a VM
func subRes(vm *asmodel.K8sNode, app asmodel.Application) {
	// CPU is a soft resource, so we think that cpuCoreStep CPU is the minimum requirement for each application
	vm.ResidualResources.CpuCore -= cpuCoreStep
	vm.ResidualResources.Memory -= app.Resources.Memory
	vm.ResidualResources.Storage -= app.Resources.Storage
}

// check whether the resources of the input VM can support all rest applications.
func vmResMeetAllRestApps(vm asmodel.K8sNode, apps map[string]asmodel.Application, curAppName *string, appsThisCloudIter *appOneCloudIter) bool {
	// we loop until the resources of this VM is used up.
	for isResEnough(vm, apps[*curAppName]) {
		subRes(&vm, apps[*curAppName]) // simulate deploying this application on this VM.
		// After the current applications is deployed, we go to the next application.
		*curAppName = appsThisCloudIter.nextAppName() // curAppName and appsThisCloudIter are pointers, so the value change will affect the variables outside this function.
		if len(*curAppName) == 0 {
			return true // this means that all applications are already deployed, so the resources are enough.
		}
	}
	return false // this means that the rest of the VM's resources can not meet the next application, so the resources are not enough
}
