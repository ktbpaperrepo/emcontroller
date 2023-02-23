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

func TestGetQemu(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	qemu, err := cloud.(*Proxmox).GetQemu("100")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", string(qemu))
}

func TestGetNetInterfaces(t *testing.T) {
	InitClouds()
	for _, cloud := range Clouds {
		switch cloud.(type) {
		case *Proxmox:
			netIfs, _ := cloud.(*Proxmox).GetNetInterfaces("100")
			fmt.Printf("%+v\n", string(netIfs))
		}
	}
}
