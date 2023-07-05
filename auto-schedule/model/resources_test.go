package model

import (
	"fmt"
	"testing"

	apiv1 "k8s.io/api/core/v1"

	"emcontroller/models"
)

func tryReqValue(pod apiv1.Pod) {
	for _, container := range pod.Spec.Containers {
		fmt.Println(container.Name)
		fmt.Println(container.Resources.Requests.Cpu().Value())              // not accurate
		fmt.Println(container.Resources.Requests.Cpu().MilliValue())         // accurate, unit m
		fmt.Println(container.Resources.Requests.Memory().Value())           // accurate, unit Byte
		fmt.Println(container.Resources.Requests.StorageEphemeral().Value()) // accurate, unit Byte
	}
}

func TestInnerTryReqValue(t *testing.T) {
	models.InitKubernetesClient()
	pods, err := models.ListPodsOnNode("", "n4test")
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	}
	for _, pod := range pods {
		tryReqValue(pod)
	}
}

func TestGetOccupiedResByPod(t *testing.T) {
	models.InitKubernetesClient()
	pods, err := models.ListPodsOnNode("", "n8test")
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	}
	for _, pod := range pods {
		occupiedRes := GetResOccupiedByPod(pod)
		t.Logf("%s/%s, occupied: %+v", pod.Namespace, pod.Name, occupiedRes)
	}
}
