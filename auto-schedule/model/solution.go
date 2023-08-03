package model

import (
	"emcontroller/models"
)

/**
NOTE:

For an application, the information for its scheduling solution only needs to include the name of the cloud and does not need to show how to schedule the application to a VM, because in my algorithm when an application is scheduled to a cloud, the scheduling to a VM will already be determined, following these rules:
1. If the resources are enough, we create dedicated/exclusive VMs for all applications with the priority 10 (highest). If multiple 10-priority applications have dependency relationships between them, they should be scheduled to the same dedicated/exclusive VMs (If App1 depends on App2, App2.Priority should >= App1.Priority). If the resources are not enough, no need for this;
2. If the resources of a cloud are not enough when using dedicated/exclusive VMs but enough for creating shared VM, we create one VM using all rest resources.
	(1) When checking the resources for dedicated/exclusive VMs, we regulate that one application occupies the number of CPU cores in its requests;
	(2) When checking the resources for a shared VM or an existing VM, we regulate that one application occupies 0.1 CPU cores.
3. For the applications whose priorities are less than 10, if the existing VMs do not have enough resources to deploy them, we need to create new shared VMs. If the resources are not enough for dedicated/exclusive VMs, we also need to deploy 10-priority applications on the shared VMs. When we create the new VM, I think we can:
	(1) if the rest of the most-used (Which type of resources has the largest used proportion?) resources are higher than 50%, we create a VM with 50% of the total amounts of every type of resources; If 50% resources are not enough for the applications, do (3).
	(2) if the rest of ... are higher than 30%, we create a VM with 30% ..... If 30% resources are not enough for the applications, do (3).
	(3) Otherwise, we create a VM with all rest resources. If the all rest resources are not enough for the applications, it means that the solution is not acceptable.
	(4) We set a periodic Garbage Collection mechanism, in our periodic check, if a VM does not have any applications, we delete it. For this, we should add an annotation to the VMs created by auto-scheduling, maybe name prefix "auto-schedule-" is more convenient. With the annotation or name prefix, we can choose only to delete the VMs created by auto-scheduling.
*/

/**
NOTE:

In our algorithm, we use some methods (random, mutation, crossover, etc.) to choose the target clouds of applications. Then the allocated CPU cores of each application will be determined:
1. The scheme of creating VMs will be determined;
2. When the CPU cores are not enough, we will only create 1 VM with all rest resources;
3. In this one VM, we can use "weighted average" with upper-bound and lower-bound to allocate the CPU cores to VMs. The weight is calculated by app.Priority * app.RequestedCpuCores. The upper-bound of the CPU cores allocated to an app is app.RequestedCpuCores. The lower-bound of the CPU cores allocated to an app is 0.1. Because of the upper-bound and lower-bound, the CPU cores may be left or not enough:
	(1) if there is left, we do another time of allocation, loop until no left.
	(2) if the CPU cores are not enough for the "weighted average" allocation, we reduce the allocated CPU cores from the lower-weight applications whose allocated CPU cores are more than 0.1.
4. In the fitness function, we still consider both of computation time (related to CPU) and communication time (related to network RTT).
*/

// Solution for scheduling applications to clouds. The keys of the map is the name of applications.
type Solution struct {
	AppsSolution map[string]SingleAppSolution `json:"appsSolution"` // key: application name
	VmsToCreate  []models.IaasVm              `json:"vmsToCreate"`
}

func (absorber *Solution) Absorb(absorbate Solution) {
	// add the solutions of the absorbate into the absorber
	for appName, appSoln := range absorbate.AppsSolution {
		absorber.AppsSolution[appName] = appSoln
	}
	// add the vms to create of the absorbate into the absorber
	absorber.VmsToCreate = append(absorber.VmsToCreate, absorbate.VmsToCreate...)
}

// The scheduling scheme for a single application
type SingleAppSolution struct {
	// This member variable means whether this application is accepted or rejected.
	// other member variables are meaningful only when this Accepted is true.
	// When Accepted is set to true, TargetCloudName must be set together.
	// In the mutation operator, when we generate a new gene, each application should (can?) have 50% possibility to be accepted and 50% to be rejected.
	Accepted bool `json:"accepted"`

	// The name of the cloud where this application is scheduled to.
	TargetCloudName string `json:"targetCloudName"`

	// The name of the kubernetes node used to deploy this application.
	K8sNodeName string `json:"k8sNodeName"`

	// number of CPU logical cores allocated to this application.
	// CPU core is a soft requirement, which means we do not have to allocate all required CPU cores to an application.
	// For example, if an application requires 4 CPU cores, but we only allocate 2 CPU cores to it, in the containerSpec of it, we will set the required CPU is 2 and the Limit CPU is 4.
	AllocatedCpuCore float64 `json:"allocatedCpuCore"`
}

// single app solution copy
func SasCopy(src SingleAppSolution) SingleAppSolution {
	var dst SingleAppSolution = src
	return dst
}

var (
	// the solution to reject an application
	RejSoln SingleAppSolution = SingleAppSolution{
		Accepted: false,
	}
)

// generate an empty solution
func GenEmptySoln() Solution {
	return Solution{
		AppsSolution: make(map[string]SingleAppSolution),
	}
}

func SolutionCopy(src Solution) Solution {
	// copy AppsSolution
	var dst Solution = Solution{
		AppsSolution: make(map[string]SingleAppSolution),
	}
	for name, singleSoln := range src.AppsSolution {
		dst.AppsSolution[name] = SasCopy(singleSoln)
	}

	// copy VmsToCreate
	if src.VmsToCreate != nil {
		dst.VmsToCreate = make([]models.IaasVm, len(src.VmsToCreate))
		copy(dst.VmsToCreate, src.VmsToCreate)
	}

	return dst
}
