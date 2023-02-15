package models

import (
	"fmt"
	"testing"
)

const testPCloudName = "NOKIA8"

func TestNodeStatus(t *testing.T) {
	InitClouds()
	for _, cloud := range Clouds {
		switch cloud.(type) {
		case *Proxmox:
			nodeStatus, _ := cloud.(*Proxmox).NodeStatus()
			fmt.Printf("%+v\n", string(nodeStatus))
		}
	}
}

func TestListQemus(t *testing.T) {
	InitClouds()
	for _, cloud := range Clouds {
		switch cloud.(type) {
		case *Proxmox:
			nodeStatus, _ := cloud.(*Proxmox).ListQemus()
			fmt.Printf("%+v\n", string(nodeStatus))
		}
	}
}
