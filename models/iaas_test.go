package models

import (
	"fmt"
	"testing"
)

func TestCheckResources(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	resourceStatus, _ := cloud.CheckResources()
	fmt.Printf("Limit: %#v\n", resourceStatus.Limit)
	fmt.Printf("InUse: %#v\n", resourceStatus.InUse)

}

func TestCreateVM(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	vm, err := cloud.CreateVM("testiaasvm2", 4, 16384, 150)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", vm)

}

func TestDeleteVM(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	err := cloud.DeleteVM("55d6b25a-1f91-4bd2-9230-130e3078e92c")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

}
