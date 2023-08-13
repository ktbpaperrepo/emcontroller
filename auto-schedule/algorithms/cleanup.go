package algorithms

import (
	"fmt"

	"github.com/astaxie/beego"
	apiv1 "k8s.io/api/core/v1"

	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

// garbage collection of auto-scheduling VMs and K8s nodes.
func GcASVms() {
	// scheduling, migration, and cleanup cannot be done at the same time
	ScheMu.Lock()
	defer ScheMu.Unlock()

	beego.Info("Start to do the periodical cleanup of auto-scheduling VMs and Kubernetes nodes.")

	// list all Kubernetes nodes with the auto-scheduling prefix
	autoK8sNodes, err := models.ListNodesNamePrefix(asmodel.ASVmNamePrefix)
	if err != nil {
		outErr := fmt.Errorf("cleanup auto-scheduling VMs, List Kubernetes Nodes with name prefix [%s] Error: %w", asmodel.ASVmNamePrefix, err)
		beego.Error(outErr)
		return
	}

	// list all VMs with the auto-scheduling prefix
	autoVms, err := models.ListVMsNamePrefix(asmodel.ASVmNamePrefix)
	if err != nil {
		outErr := fmt.Errorf("cleanup auto-scheduling VMs, List VMs with name prefix [%s] Error: %w", asmodel.ASVmNamePrefix, err)
		beego.Error(outErr)
		return
	}

	// look for the Kubernetes nodes and vms to delete
	k8sNodesToDelete := make([]apiv1.Node, len(autoK8sNodes))
	copy(k8sNodesToDelete, autoK8sNodes)
	vmsToDelete := make([]models.IaasVm, len(autoVms))
	copy(vmsToDelete, autoVms)

	// We delete all auto-scheduling kubernetes nodes and auto-scheduling VMs which do not have any Kubernetes applications running.
	for _, node := range autoK8sNodes {
		podsOnNode, err := models.ListPodsOnNode(models.KubernetesNamespace, node.Name)
		if err != nil {
			outErr := fmt.Errorf("List pods on Kubernetes node [%s], error: %w", node.Name, err)
			beego.Error(outErr)
			return
		}
		if len(podsOnNode) > 0 { // if this node has applications running, we remove it from the delete list.
			models.RemoveNodeFromList(&k8sNodesToDelete, node.Name)
			models.RemoveVmFromList(&vmsToDelete, node.Name)
		}
	}

	// make the slice of the names of the VMs and nodes to delete, for log and deletion.
	var k8sNodeNamesToDelete, vmNamesToDelete []string
	for _, node := range k8sNodesToDelete {
		k8sNodeNamesToDelete = append(k8sNodeNamesToDelete, node.Name)
	}
	for _, vm := range vmsToDelete {
		vmNamesToDelete = append(vmNamesToDelete, vm.Name)
	}

	beego.Info(fmt.Sprintf("Delete Kubernetes nodes %v from the cluster.", k8sNodeNamesToDelete))
	if errs := models.UninstallBatchNodes(k8sNodeNamesToDelete); errs != nil {
		outErr := fmt.Errorf("Delete Kubernetes nodes from the cluster, error: %w", models.HandleErrSlice(errs))
		beego.Error(outErr)
	}
	beego.Info(fmt.Sprintf("Successfully, Delete Kubernetes nodes %v from the cluster.", k8sNodeNamesToDelete))

	beego.Info(fmt.Sprintf("Delete Virtual Machines %v.", vmNamesToDelete))
	if errs := models.DeleteBatchVms(vmsToDelete); errs != nil {
		outErr := fmt.Errorf("Delete Virtual Machines, error: %w", models.HandleErrSlice(errs))
		beego.Error(outErr)
	}
	beego.Info(fmt.Sprintf("Successfully, Delete Virtual Machines %v.", vmNamesToDelete))

	beego.Info("Finished this period of the cleanup of auto-scheduling VMs and Kubernetes nodes.")
}
