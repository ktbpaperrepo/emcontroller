/*
The functions to allocate Virtual Machines to applications inside clouds.
*/

package algorithms

import (
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

	// avoiding changing the original solution
	solnWithVm := asmodel.SolutionCopy(soln)

	// We should allocate VMs cloud by cloud.
	for _, cloud := range clouds {
		solnWithVmsThisCloud, allocType := allocateVmsOneCLoud(cloud, apps, appsOrder, soln)
		if allocType == UnAcceptable { // if any cloud cannot accept the scheduled applications, this whole solution is not acceptable.
			return asmodel.Solution{}, false
		}

		// update the part of solution of this cloud to the total solution.
		solnWithVm.Absorb(solnWithVmsThisCloud)
	}

	return solnWithVm, true
}

// allocate VMs in one cloud
func allocateVmsOneCLoud(cloud asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, soln asmodel.Solution) (asmodel.Solution, VmAllocType) {
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
	var neededAvailRes asmodel.AppResources = asmodel.AppResources{
		GenericResources: asmodel.GenericResources{
			CpuCore: 0,
			Memory:  0,
			Storage: 0,
		},
	}

	for _, appName := range appGroup {
		// for max-priority applications on dedicated VMs, we regulate that it occupy its requested CPU cores rather than minCpu
		neededAvailRes.CpuCore += apps[appName].Resources.CpuCore
		neededAvailRes.Memory += apps[appName].Resources.Memory
		neededAvailRes.Storage += apps[appName].Resources.Storage
	}

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
