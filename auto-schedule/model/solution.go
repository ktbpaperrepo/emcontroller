package model

/*
For an application, the information for its scheduling solution only needs to include the name of the cloud and does not need to show how to schedule the application to a VM, because in my algorithm when an application is scheduled to a cloud, the scheduling to a VM will already be determined, following these rules:
1. If the resources are enough, we create dedicated/exclusive VMs for all applications with the priority 10 (highest). If the resources are not enough, no need for this;
2. If the resources are enough, If multiple 10-priority applications have dependency relationships between them, they should be scheduled to the same dedicated/exclusive VMs. If the resources are not enough, no need for this; (If App1 depends on App2, App2.Priority should >= App1.Priority).
3. For the applications whose priorities are less than 10, if the existing VMs do not have enough resources to deploy them, we need to create new VMs. In this condition, currently, I think we can:
	(1) if the rest of the most-used (Which type of resources has the largest used proportion?) resources are higher than 50%, we create a VM with 50% of the total amounts of every type of resources;
	(2) if the rest of ... are higher than 30%, we create a VM with 30% ....
	(3) Otherwise, we create a VM with all rest resources.
	(4) We set a periodic Garbage Collection mechanism, in our periodic check, if a VM does not have any applications, we delete it.
*/

// Solution for scheduling applications to clouds. The keys of the map is the name of applications.
type Solution map[string]SingleAppSolution

// The scheduling scheme for a single application
type SingleAppSolution struct {
	TargetCloudName string `json:"targetCloudName"`

	// number of CPU logical cores allocated to this application.
	// CPU core is a soft requirement, which means we do not have to allocate all required CPU cores to an application.
	// For example, if an application requires 4 CPU cores, but we only allocate 2 CPU cores to it, in the containerSpec of it, we will set the required CPU is 2 and the Limit CPU is 4.
	AllocatedCpuCore float64 `json:"allocatedCpuCore"`
}
