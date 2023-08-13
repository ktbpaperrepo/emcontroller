package models

import (
	"fmt"

	"github.com/astaxie/beego"
	apiv1 "k8s.io/api/core/v1"
)

type K8sNodeInfo struct {
	Name   string
	IP     string
	Status string
}

// create the new VMs and add them to Kubernetes
func AddNewVms(vmsToCreate []IaasVm) ([]IaasVm, error) {
	beego.Info(fmt.Sprintf("Create new VMs [%s].", JsonString(vmsToCreate)))
	createdVms, err := CreateVms(vmsToCreate)
	if err != nil {
		outErr := fmt.Errorf("Create new VMs [%s], Error: [%w]", JsonString(vmsToCreate), err)
		beego.Error(outErr)
		return []IaasVm{}, outErr
	}

	beego.Info(fmt.Sprintf("Add new VMs [%s] to Kubernetes cluster.", JsonString(createdVms)))
	errs := AddNodes(createdVms)
	if len(errs) != 0 {
		sumErr := HandleErrSlice(errs)
		outErr := fmt.Errorf("Add new VMs [%s] to Kubernetes cluster, Error: [%w]", JsonString(createdVms), sumErr)
		beego.Error(outErr)
		return createdVms, outErr
	}

	return createdVms, nil
}

// remove a Kubernetes node with a appointed name from a list
func RemoveNodeFromList(nodeList *[]apiv1.Node, nodeNameToRemove string) {
	nodeIdx := FindIdxNodeInList(*nodeList, nodeNameToRemove)
	if nodeIdx >= 0 {
		*nodeList = append((*nodeList)[:nodeIdx], (*nodeList)[nodeIdx+1:]...)
	}
}

// find the index of a Kubernetes node with a appointed name in a list. return -1 if not found.
func FindIdxNodeInList(nodeList []apiv1.Node, nodeNameToFind string) int {
	for idx, node := range nodeList {
		if node.Name == nodeNameToFind {
			return idx
		}
	}
	return -1
}
