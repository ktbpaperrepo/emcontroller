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
	cloud := Clouds[testPCloudName]
	vm, err := cloud.CreateVM("testiaasvm", 8, 16384, 150)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", vm)
}

func TestDeleteVM(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	err := cloud.DeleteVM("102")
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
