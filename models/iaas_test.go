package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckResources(t *testing.T) {
	InitSomeThing()
	cloud := Clouds[testPCloudName]
	resourceStatus, _ := cloud.CheckResources()
	fmt.Printf("Limit: %#v\n", resourceStatus.Limit)
	fmt.Printf("InUse: %#v\n", resourceStatus.InUse)
}

func TestGetVM(t *testing.T) {
	InitSomeThing()
	cloud := Clouds[testPCloudName]
	vm, err := cloud.GetVM("100")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", vm)
}

func TestListAllVMs(t *testing.T) {
	InitSomeThing()
	cloud := Clouds[testPCloudName]
	vms, err := cloud.ListAllVMs()
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	t.Logf("The result is: [%+v]\n", vms)
}

func TestCreateVM(t *testing.T) {
	InitSomeThing()
	cloud := Clouds[testOsCloudName]
	vm, err := cloud.CreateVM("testiaasvm", 8, 16384, 150)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%+v\n", vm)
}

func TestCreateVms(t *testing.T) {
	InitSomeThing()
	var vmsToCreate []IaasVm = []IaasVm{
		{Cloud: "NOKIA8", Name: "n8-node1", VCpu: 4, Ram: 8192, Storage: 3},
		//{Cloud: "NOKIA8", Name: "n8-node2", VCpu: 4, Ram: 8192, Storage: 100},
		//{Cloud: "NOKIA7", Name: "n7-node1", VCpu: 4, Ram: 8192, Storage: 100},
		//{Cloud: "NOKIA7", Name: "n7-node2", VCpu: 4, Ram: 8192, Storage: 100},
		//{Cloud: "CLAAUDIAweifan", Name: "cnode1", VCpu: 4, Ram: 32768, Storage: 100},
		//{Cloud: "CLAAUDIAweifan", Name: "cnode2", VCpu: 4, Ram: 32768, Storage: 100},
	}
	if vms, err := CreateVms(vmsToCreate); err != nil {
		t.Errorf("Create VMs error: %s", err.Error())
	} else {
		t.Logf("Create VMs successfully.")
		t.Logf("Created VMs are: [%v].", vms)
	}
}

func TestDeleteVM(t *testing.T) {
	InitSomeThing()
	cloud := Clouds["HPE1"]
	err := cloud.DeleteVM("107")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
}

func TestDeleteBatchVms(t *testing.T) {
	InitSomeThing()

	vmsToDelete, err := ListVMsNamePrefix("auto-sched-")
	if err != nil {
		t.Errorf("get VMs to delete error: %s", err.Error())
	}

	t.Log("VMs to delete:")
	for _, vm := range vmsToDelete {
		t.Logf("%s/%s", vm.Cloud, vm.Name)
	}

	errs := DeleteBatchVms(vmsToDelete)
	if errs != nil {
		t.Errorf("uninstall nodes error: %s", HandleErrSlice(errs).Error())
	} else {
		t.Logf("Uninstall Kubernetes nodes successfully")
	}
}

func TestIsCreatedByMcm(t *testing.T) {
	InitSomeThing()
	cloud := Clouds[testOsCloudName]
	is, err := cloud.IsCreatedByMcm("d2076789-f289-4ae1-b599-b8e20e7658b3")
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}
	fmt.Printf("%t\n", is)
}

func TestAllMoreThan(t *testing.T) {
	testCases := []struct {
		name           string
		resStatus      ResourceStatus
		expectedResult bool
	}{
		{
			name: "all more",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    6,
					Ram:     2048,
					Storage: 2048,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: true,
		},
		{
			name: "cpuLess",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    5,
					Ram:     4096,
					Storage: 4096,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    6,
					Ram:     2048,
					Storage: 2048,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: false,
		},
		{
			name: "RamEqual",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    6,
					Ram:     4096,
					Storage: 2048,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: false,
		},
		{
			name: "vmEqual",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    6,
					Ram:     2048,
					Storage: 2048,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: false,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := testCase.resStatus.Limit.AllMoreThan(testCase.resStatus.InUse)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestLeastRemainPct(t *testing.T) {
	testCases := []struct {
		name           string
		resStatus      ResourceStatus
		expectedResult float64
	}{
		{
			name: "cpu least",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    6,
					Ram:     2048,
					Storage: 2048,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: 1.0 / 7.0,
		},
		{
			name: "cpu Ram both least ",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    8,
					Ram:     4096,
					Storage: 4096,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    4,
					Ram:     2048,
					Storage: 1024,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: 0.5,
		},
		{
			name: "ram used up",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    3,
					Ram:     4096,
					Storage: 2048,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: 0.0,
		},
		{
			name: "storage least",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    2,
					Ram:     2048,
					Storage: 4000,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: 96.0 / 4096.0,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := testCase.resStatus.LeastRemainPct()
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestOverflow(t *testing.T) {
	testCases := []struct {
		name           string
		resStatus      ResourceStatus
		expectedResult bool
	}{
		{
			name: "no overflow",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    6,
					Ram:     2048,
					Storage: 2048,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: false,
		},
		{
			name: "cpu overflows",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    8,
					Ram:     2048,
					Storage: 2048,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: true,
		},
		{
			name: "cpu equal",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    8,
					Ram:     4096,
					Storage: 4096,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    8,
					Ram:     2048,
					Storage: 1024,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: false,
		},
		{
			name: "ram overflow",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    3,
					Ram:     5000,
					Storage: 2048,
					Vm:      -1,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: true,
		},
		{
			name: "storage overflow",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    2,
					Ram:     2048,
					Storage: 6000,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: true,
		},
		{
			name: "memory and storage overflow",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    2,
					Ram:     5001,
					Storage: 6000,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: true,
		},
		{
			name: "all overflow",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    10,
					Ram:     5001,
					Storage: 6000,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: true,
		},
		{
			name: "all equal",
			resStatus: ResourceStatus{
				Limit: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
				InUse: ResSet{
					VCpu:    7,
					Ram:     4096,
					Storage: 4096,
					Vm:      5,
					Volume:  -1,
					Port:    -1,
				},
			},
			expectedResult: false,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := testCase.resStatus.Overflow()
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestGroupVmsByCloud(t *testing.T) {
	testCases := []struct {
		name           string
		vms            []IaasVm
		expectedResult map[string][]IaasVm
	}{
		{
			name: "case1",
			vms: []IaasVm{
				IaasVm{
					Name:    "auto-sched-nokia7-0",
					VCpu:    3,
					Ram:     9462,
					Storage: 238,
					Cloud:   "NOKIA7",
				},
				IaasVm{
					Name:    "auto-sched-nokia7-1",
					VCpu:    15,
					Ram:     52154,
					Storage: 516,
					Cloud:   "NOKIA7",
				},
				IaasVm{
					Name:    "auto-sched-nokia8-0",
					VCpu:    38,
					Ram:     98076,
					Storage: 482,
					Cloud:   "NOKIA8",
				},
				IaasVm{
					Name:    "auto-sched-hpe1-0",
					VCpu:    8,
					Ram:     2432,
					Storage: 45,
					Cloud:   "HPE1",
				},
				IaasVm{
					Name:    "auto-sched-hpe1-1",
					VCpu:    66,
					Ram:     453675,
					Storage: 363,
					Cloud:   "HPE1",
				},
				IaasVm{
					Name:    "auto-sched-nokia4-0",
					VCpu:    34,
					Ram:     73500,
					Storage: 867,
					Cloud:   "NOKIA4",
				},
			},
			expectedResult: map[string][]IaasVm{
				"NOKIA7": []IaasVm{
					IaasVm{
						Name:    "auto-sched-nokia7-0",
						VCpu:    3,
						Ram:     9462,
						Storage: 238,
						Cloud:   "NOKIA7",
					},
					IaasVm{
						Name:    "auto-sched-nokia7-1",
						VCpu:    15,
						Ram:     52154,
						Storage: 516,
						Cloud:   "NOKIA7",
					},
				},
				"HPE1": []IaasVm{
					IaasVm{
						Name:    "auto-sched-hpe1-0",
						VCpu:    8,
						Ram:     2432,
						Storage: 45,
						Cloud:   "HPE1",
					},
					IaasVm{
						Name:    "auto-sched-hpe1-1",
						VCpu:    66,
						Ram:     453675,
						Storage: 363,
						Cloud:   "HPE1",
					},
				},
				"NOKIA8": []IaasVm{
					IaasVm{
						Name:    "auto-sched-nokia8-0",
						VCpu:    38,
						Ram:     98076,
						Storage: 482,
						Cloud:   "NOKIA8",
					},
				},
				"NOKIA4": []IaasVm{
					IaasVm{
						Name:    "auto-sched-nokia4-0",
						VCpu:    34,
						Ram:     73500,
						Storage: 867,
						Cloud:   "NOKIA4",
					},
				},
			},
		},
		{
			name: "case2",
			vms: []IaasVm{
				IaasVm{
					Name:    "auto-sched-nokia4-0",
					VCpu:    3,
					Ram:     9462,
					Storage: 238,
					Cloud:   "NOKIA4",
				},
				IaasVm{
					Name:    "auto-sched-nokia6-0",
					VCpu:    15,
					Ram:     52154,
					Storage: 516,
					Cloud:   "NOKIA6",
				},
				IaasVm{
					Name:    "auto-sched-nokia8-0",
					VCpu:    38,
					Ram:     98076,
					Storage: 482,
					Cloud:   "NOKIA8",
				},
				IaasVm{
					Name:    "auto-sched-nokia6-1",
					VCpu:    8,
					Ram:     2432,
					Storage: 45,
					Cloud:   "NOKIA6",
				},
				IaasVm{
					Name:    "auto-sched-nokia6-2",
					VCpu:    66,
					Ram:     453675,
					Storage: 363,
					Cloud:   "NOKIA6",
				},
				IaasVm{
					Name:    "auto-sched-nokia4-0",
					VCpu:    34,
					Ram:     73500,
					Storage: 867,
					Cloud:   "NOKIA4",
				},
				IaasVm{
					Name:    "auto-sched-nokia8-1",
					VCpu:    34,
					Ram:     73500,
					Storage: 867,
					Cloud:   "NOKIA8",
				},
			},
			expectedResult: map[string][]IaasVm{
				"NOKIA4": []IaasVm{
					IaasVm{
						Name:    "auto-sched-nokia4-0",
						VCpu:    3,
						Ram:     9462,
						Storage: 238,
						Cloud:   "NOKIA4",
					},
					IaasVm{
						Name:    "auto-sched-nokia4-0",
						VCpu:    34,
						Ram:     73500,
						Storage: 867,
						Cloud:   "NOKIA4",
					},
				},
				"NOKIA6": []IaasVm{
					IaasVm{
						Name:    "auto-sched-nokia6-0",
						VCpu:    15,
						Ram:     52154,
						Storage: 516,
						Cloud:   "NOKIA6",
					},
					IaasVm{
						Name:    "auto-sched-nokia6-1",
						VCpu:    8,
						Ram:     2432,
						Storage: 45,
						Cloud:   "NOKIA6",
					},
					IaasVm{
						Name:    "auto-sched-nokia6-2",
						VCpu:    66,
						Ram:     453675,
						Storage: 363,
						Cloud:   "NOKIA6",
					},
				},
				"NOKIA8": []IaasVm{
					IaasVm{
						Name:    "auto-sched-nokia8-0",
						VCpu:    38,
						Ram:     98076,
						Storage: 482,
						Cloud:   "NOKIA8",
					},
					IaasVm{
						Name:    "auto-sched-nokia8-1",
						VCpu:    34,
						Ram:     73500,
						Storage: 867,
						Cloud:   "NOKIA8",
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := GroupVmsByCloud(testCase.vms)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}
