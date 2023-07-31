package model

import (
	"emcontroller/models"
	apiv1 "k8s.io/api/core/v1"
)

const (
	ASVmNamePrefix string = "auto-sched-"
)

// There may be some existing VMs on a cloud, which we should consider when scheduling applications, so we use this struct to represent them.
type K8sNode struct {
	Name              string           `json:"name"`
	ResidualResources GenericResources `json:"residualResources"` // The rest of the resources on the VM of this Kubernetes Node (The resources occupied by pods are already subtracted.)

	/* We do not need to put the information of running applications here, because:
	1. When deploying new applications, we do not need to consider running ones;
	2. When migrating old applications, we can assume that they are new applications, which means we can simulate to remove them from the clouds and then pass them to the scheduling functions.
	3. We do not support to migrating old applications and deploying new applications in one scheduling.
	*/

}

// copy a K8sNode to generate a new and same one
func K8sNodeCopy(src K8sNode) K8sNode {
	var dst K8sNode = src
	return dst
}

// generate a K8sNode variable from a VM and the pods deployed on it.
func GenK8sNodeFromPods(vm models.IaasVm, podsOnNode []apiv1.Pod) K8sNode {
	// Get available resources of this VM
	residualCpuCore := models.CalcVmAvailVcpu(vm.VCpu)
	residualRamMiB := models.CalcVmAvailRamMiB(vm.Ram)
	residualStorGiB := models.CalcVmAvailStorGiB(vm.Storage)

	// subtract the resources occupied by pods
	for _, pod := range podsOnNode {
		occupied := GetResOccupiedByPod(pod)
		residualCpuCore -= occupied.CpuCore
		residualRamMiB -= occupied.Memory
		residualStorGiB -= occupied.Storage
	}

	// handle possible negative results
	if residualCpuCore < 0 {
		residualCpuCore = 0
	}
	if residualRamMiB < 0 {
		residualRamMiB = 0
	}
	if residualStorGiB < 0 {
		residualStorGiB = 0
	}

	// we put the information needed by auto-scheduling to the K8sNode structure
	var thisNode K8sNode
	thisNode.Name = vm.Name
	thisNode.ResidualResources.CpuCore = residualCpuCore
	thisNode.ResidualResources.Memory = residualRamMiB
	thisNode.ResidualResources.Storage = residualStorGiB
	return thisNode
}

// generate a K8sNode variable from a VM and the applications deployed on it.
func GenK8sNodeFromApps(vm models.IaasVm, apps map[string]Application, appGroup []string) K8sNode {
	// Get available resources of this VM
	residualCpuCore := models.CalcVmAvailVcpu(vm.VCpu)
	residualRamMiB := models.CalcVmAvailRamMiB(vm.Ram)
	residualStorGiB := models.CalcVmAvailStorGiB(vm.Storage)

	// subtract the resources occupied by applications
	for _, appName := range appGroup {
		residualCpuCore -= apps[appName].Resources.CpuCore
		residualRamMiB -= apps[appName].Resources.Memory
		residualStorGiB -= apps[appName].Resources.Storage
	}

	// handle possible negative results
	if residualCpuCore < 0 {
		residualCpuCore = 0
	}
	if residualRamMiB < 0 {
		residualRamMiB = 0
	}
	if residualStorGiB < 0 {
		residualStorGiB = 0
	}

	// we put the information needed by auto-scheduling to the K8sNode structure
	var thisNode K8sNode
	thisNode.Name = vm.Name
	thisNode.ResidualResources.CpuCore = residualCpuCore
	thisNode.ResidualResources.Memory = residualRamMiB
	thisNode.ResidualResources.Storage = residualStorGiB
	return thisNode
}
