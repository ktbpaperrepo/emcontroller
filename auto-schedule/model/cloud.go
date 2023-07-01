package model

import "emcontroller/models"

type Cloud struct {
	Name      string                         `json:"name"`
	Resources models.ResourceStatus          `json:"resources"` // used and all resources of this cloud. Here we start with struct defined in "models" package, and in the future if we find that this cannot meet the needs here, we can define new structs.
	NetState  map[string]models.NetworkState `json:"netState"`  // the network state from this cloud to every cloud
	K8sNodes  []K8sNode                      `json:"k8sNodes"`  // all existing Kubernetes nodes whose VMs are on this cloud
}
