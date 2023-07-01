package model

// There may be some existing VMs on a cloud, which we should consider when scheduling applications, so we use this struct to represent them.
type K8sNode struct {
	Name               string           `json:"name"`
	AvailableResources GenericResources `json:"availableResources"` // available resources on the VM of this Kubernetes Node

	/* We do not need to put the information of running applications here, because:
	1. When deploying new applications, we do not need to consider running ones;
	2. When migrating old applications, we can assume that they are new applications, which means we can simulate to remove them from the clouds and then pass them to the scheduling functions.
	3. We do not support to migrating old applications and deploying new applications in one scheduling.
	*/

}
