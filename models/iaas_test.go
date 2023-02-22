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
	cloud := Clouds[testOsCloudName]
	vm, err := cloud.GetVM("e5df0ca9-1f33-4d56-a3d5-7387d66bac6e")
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
	vm, err := cloud.CreateVM("testiaasvm", 4, 16384, 150)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", vm)
}

func TestDeleteVM(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	err := cloud.DeleteVM("9b8e9722-7ea4-4152-9f42-d2ac3bad720e")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
}
