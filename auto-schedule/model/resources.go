/*
NOTE:

There are 2 types of resources: (If no people defined them previously, we can define them here.)
1. soft resource: we can choose to meet the requirement or not, if we do not meet the requirement, the response time will be longer, but the application will still work, such as CPU and RTT.
2. hard resource: If we accept an application, we have to meet its hard requirements, or otherwise this application will not work, such as Memory and Storage (If memory is not enough, there will be OOM (Out-of memory) error).
*/

package model

import (
	apiv1 "k8s.io/api/core/v1"
)

const (
	CpuResReg         string = `^[0-9]+(\.[0-9]+)?$`
	MemUnitSuffix     string = "Mi"
	StorageUnitSuffix string = "Gi"
)

// We use a different Object for the applications resources, in case of some special scenarios.
type AppResources struct {
	GenericResources `json:",inline"`
}

type GenericResources struct {
	CpuCore float64 `json:"cpuCore"` // number of CPU logical cores that this application needs, this is a soft requirement
	Memory  float64 `json:"memory"`  // unit Mebibyte (MiB), this is a hard requirement
	Storage float64 `json:"storage"` // unit Gibibyte (GiB), this is a hard requirement
}

func GetResOccupiedByPod(pod apiv1.Pod) GenericResources {
	var occupied GenericResources
	for _, container := range pod.Spec.Containers {
		// convert unit to 1 CPU
		occupied.CpuCore += float64(container.Resources.Requests.Cpu().MilliValue()) / 1000
		// convert unit to Mi
		occupied.Memory += float64(container.Resources.Requests.Memory().Value()) / 1024 / 1024
		// convert unit to Gi
		occupied.Storage += float64(container.Resources.Requests.StorageEphemeral().Value()) / 1024 / 1024 / 1024
	}
	return occupied
}
