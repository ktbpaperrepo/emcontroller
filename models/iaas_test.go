package models

import (
	"fmt"
	"testing"
)

func TestCheckResources(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		resourceStatus := Clouds[i].CheckResources()
		fmt.Printf("Limit: %#v\n", resourceStatus.Limit)
		fmt.Printf("InUse: %#v\n", resourceStatus.InUse)
	}
}
