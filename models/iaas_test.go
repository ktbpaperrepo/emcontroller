package models

import (
	"fmt"
	"testing"
)

func TestCheckResources(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		resourceStatus, _ := Clouds[i].CheckResources()
		fmt.Printf("Limit: %#v\n", resourceStatus.Limit)
		fmt.Printf("InUse: %#v\n", resourceStatus.InUse)
	}
}

func TestCreateVM(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		vm, err := Clouds[i].CreateVM("testiaasvm2", 4, 16384, 150)
		if err != nil {
			t.Errorf("error: %s\n", err.Error())
		}
		fmt.Printf("%+v\n", vm)
	}
}
