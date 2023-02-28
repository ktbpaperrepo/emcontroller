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

func TestShutdownQemu(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	resp, err := cloud.(*Proxmox).ShutdownQemu("102")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", string(resp))
}

func TestDeleteQemu(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	resp, err := cloud.(*Proxmox).DeleteQemu("102")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", string(resp))
}

func TestGetTaskStatus(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	taskStatus, err := cloud.(*Proxmox).GetTaskStatus("UPID:NOKIA8:00355E3A:0A280A27:63FC6BB4:qmclone:103:root@pam!multicloud:")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", string(taskStatus))
}

func TestGetQemuConfig(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	config, err := cloud.(*Proxmox).GetQemuConfig("102")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	if err := cloud.(*Proxmox).CheckErrInResp(config); err != nil {
		t.Errorf("error in resp: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", string(config))
}

func TestInnerGetDiskName(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	diskName, err := cloud.(*Proxmox).getDiskName("102")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%s\n", diskName)
}

func TestCloneQemu(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	resp, err := cloud.(*Proxmox).CloneQemu("104", "testcode")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", string(resp))
}

func TestConfigCoreRam(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	resp, err := cloud.(*Proxmox).ConfigCoreRam(102, 16384, 8)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	if err := cloud.(*Proxmox).CheckErrInResp(resp); err != nil {
		t.Errorf("error in resp: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", string(resp))
}

func TestResizeDisk(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	resp, err := cloud.(*Proxmox).ResizeDisk(102, "scsi0", "150G")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", string(resp))
}

func TestStartQemu(t *testing.T) {
	InitClouds()
	cloud := Clouds[testPCloudName]
	resp, err := cloud.(*Proxmox).StartQemu(104)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", string(resp))
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
