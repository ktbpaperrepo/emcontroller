package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListVMsNamePrefix(t *testing.T) {
	InitSomeThing()
	vms, err := ListVMsNamePrefix("auto-sched-")
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	} else {
		for _, vm := range vms {
			t.Log(vm.Name, vm.IPs)
		}
	}
}

var vmsForTest = []IaasVm{
	IaasVm{
		Name: "vm1",
	},
	IaasVm{
		Name: "vm2",
	},
	IaasVm{
		Name: "vm3",
	},
	IaasVm{
		Name: "vm4",
	},
}

func TestFindIdxVmInList(t *testing.T) {

	testCases := []struct {
		name           string
		vmList         []IaasVm
		vmNameToFind   string
		expectedResult int
	}{
		{
			name:           "case found 1",
			vmList:         vmsForTest,
			vmNameToFind:   "vm1",
			expectedResult: 0,
		},
		{
			name:           "case found 2",
			vmList:         vmsForTest,
			vmNameToFind:   "vm2",
			expectedResult: 1,
		},
		{
			name:           "case found 3",
			vmList:         vmsForTest,
			vmNameToFind:   "vm3",
			expectedResult: 2,
		},
		{
			name:           "case found 4",
			vmList:         vmsForTest,
			vmNameToFind:   "vm4",
			expectedResult: 3,
		},
		{
			name:           "case not found 1",
			vmList:         vmsForTest,
			vmNameToFind:   "vm5",
			expectedResult: -1,
		},
		{
			name:           "case not found 2",
			vmList:         vmsForTest,
			vmNameToFind:   "asdfasf",
			expectedResult: -1,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := FindIdxVmInList(testCase.vmList, testCase.vmNameToFind)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}

func TestRemoveVmFromList(t *testing.T) {

	testCases := []struct {
		name           string
		vmList         []IaasVm
		vmNameToRemove string
		expectedResult []IaasVm
	}{
		{
			name: "case exist 1",
			vmList: []IaasVm{
				IaasVm{
					Name: "vm1",
				},
				IaasVm{
					Name: "vm2",
				},
				IaasVm{
					Name: "vm3",
				},
				IaasVm{
					Name: "vm4",
				},
			},
			vmNameToRemove: "vm1",
			expectedResult: []IaasVm{
				IaasVm{
					Name: "vm2",
				},
				IaasVm{
					Name: "vm3",
				},
				IaasVm{
					Name: "vm4",
				},
			},
		},
		{
			name: "case exist 2",
			vmList: []IaasVm{
				IaasVm{
					Name: "vm1",
				},
				IaasVm{
					Name: "vm2",
				},
				IaasVm{
					Name: "vm3",
				},
				IaasVm{
					Name: "vm4",
				},
			},
			vmNameToRemove: "vm2",
			expectedResult: []IaasVm{
				IaasVm{
					Name: "vm1",
				},
				IaasVm{
					Name: "vm3",
				},
				IaasVm{
					Name: "vm4",
				},
			},
		},
		{
			name: "case exist 3",
			vmList: []IaasVm{
				IaasVm{
					Name: "vm1",
				},
				IaasVm{
					Name: "vm2",
				},
				IaasVm{
					Name: "vm3",
				},
				IaasVm{
					Name: "vm4",
				},
			},
			vmNameToRemove: "vm3",
			expectedResult: []IaasVm{
				IaasVm{
					Name: "vm1",
				},
				IaasVm{
					Name: "vm2",
				},
				IaasVm{
					Name: "vm4",
				},
			},
		},
		{
			name: "case exist 4",
			vmList: []IaasVm{
				IaasVm{
					Name: "vm1",
				},
				IaasVm{
					Name: "vm2",
				},
				IaasVm{
					Name: "vm3",
				},
				IaasVm{
					Name: "vm4",
				},
			},
			vmNameToRemove: "vm4",
			expectedResult: []IaasVm{
				IaasVm{
					Name: "vm1",
				},
				IaasVm{
					Name: "vm2",
				},
				IaasVm{
					Name: "vm3",
				},
			},
		},
		{
			name: "case not exist 1",
			vmList: []IaasVm{
				IaasVm{
					Name: "vm1",
				},
				IaasVm{
					Name: "vm2",
				},
				IaasVm{
					Name: "vm3",
				},
				IaasVm{
					Name: "vm4",
				},
			},
			vmNameToRemove: "vm5",
			expectedResult: []IaasVm{
				IaasVm{
					Name: "vm1",
				},
				IaasVm{
					Name: "vm2",
				},
				IaasVm{
					Name: "vm3",
				},
				IaasVm{
					Name: "vm4",
				},
			},
		},
		{
			name: "case not exist 2",
			vmList: []IaasVm{
				IaasVm{
					Name: "vm1",
				},
				IaasVm{
					Name: "vm2",
				},
				IaasVm{
					Name: "vm3",
				},
				IaasVm{
					Name: "vm4",
				},
			},
			vmNameToRemove: "asdfasf",
			expectedResult: []IaasVm{
				IaasVm{
					Name: "vm1",
				},
				IaasVm{
					Name: "vm2",
				},
				IaasVm{
					Name: "vm3",
				},
				IaasVm{
					Name: "vm4",
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		RemoveVmFromList(&(testCase.vmList), testCase.vmNameToRemove)
		assert.Equal(t, testCase.expectedResult, testCase.vmList, fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}
