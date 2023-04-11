package models

import (
	"fmt"
	"testing"
)

func TestCheckResources(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	resourceStatus, _ := cloud.CheckResources()
	fmt.Printf("Limit: %#v\n", resourceStatus.Limit)
	fmt.Printf("InUse: %#v\n", resourceStatus.InUse)
}

func TestGetVM(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	vm, err := cloud.GetVM("100")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", vm)
}

func TestListAllVMs(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	vm, err := cloud.ListAllVMs()
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", vm)
}

func TestCreateVM(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	vm, err := cloud.CreateVM("testiaasvm", 8, 16384, 150)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", vm)
}

func TestCreateVms(t *testing.T) {
	InitClouds()
	var vmsToCreate []IaasVm = []IaasVm{
		{Cloud: "NOKIA10", Name: "node1", VCpu: 4, Ram: 32768, Storage: 100},
		{Cloud: "NOKIA8", Name: "node2", VCpu: 4, Ram: 32768, Storage: 100},
		{Cloud: "CLAAUDIAweifan", Name: "cnode1", VCpu: 4, Ram: 32768, Storage: 100},
		{Cloud: "CLAAUDIAweifan", Name: "cnode2", VCpu: 4, Ram: 32768, Storage: 100},
	}
	if vms, err := CreateVms(vmsToCreate); err != nil {
		t.Errorf("Create VMs error: %s", err.Error())
	} else {
		t.Logf("Create VMs successfully.")
		t.Logf("Created VMs are: [%v].", vms)
	}
}

func TestDeleteVM(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	err := cloud.DeleteVM("a3e02a3a-7213-462a-bbe4-5411a6f92be2")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
}

func TestIsCreatedByMcm(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	is, err := cloud.IsCreatedByMcm("d2076789-f289-4ae1-b599-b8e20e7658b3")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%t\n", is)
}
