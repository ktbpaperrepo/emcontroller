package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"emcontroller/models"
)

func TestCloudCopy(t *testing.T) {
	testCases := []struct {
		name string
		src  Cloud
	}{
		{
			name: "case 1",
			src: Cloud{
				Name: "cloud1",
				K8sNodes: []K8sNode{
					K8sNode{
						Name: "node1",
					},
					K8sNode{
						Name: "node2",
					},
				},
			},
		},
		{
			name: "case 2",
			src: Cloud{
				Name: "cloud2",
				K8sNodes: []K8sNode{
					K8sNode{
						Name: "node3",
					},
					K8sNode{
						Name: "node4",
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		dst := CloudCopy(testCase.src)
		assert.Equal(t, testCase.src, dst)
		dst.K8sNodes[0].Name = "node100"
		assert.NotEqual(t, testCase.src.K8sNodes[0], dst.K8sNodes[0], fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}

func TestCloudMapCopy(t *testing.T) {
	testCases := []struct {
		name string
		src  map[string]Cloud
	}{
		{
			name: "case1",
			src: map[string]Cloud{
				"cloud1": Cloud{
					Name: "cloud1",
					K8sNodes: []K8sNode{
						K8sNode{
							Name: "node1",
						},
						K8sNode{
							Name: "node2",
						},
					},
				},
				"cloud2": Cloud{
					Name: "cloud2",
					K8sNodes: []K8sNode{
						K8sNode{
							Name: "node3",
						},
						K8sNode{
							Name: "node4",
						},
					},
				},
			},
		},
		{
			name: "case2",
			src: map[string]Cloud{
				"cloud1": Cloud{
					Name: "cloud1",
					K8sNodes: []K8sNode{
						K8sNode{
							Name: "node10",
						},
						K8sNode{
							Name: "node20",
						},
					},
				},
				"cloud2": Cloud{
					Name: "cloud2",
					K8sNodes: []K8sNode{
						K8sNode{
							Name: "node30",
						},
						K8sNode{
							Name: "node40",
						},
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		var dst map[string]Cloud

		// test extra key
		dst = CloudMapCopy(testCase.src)
		assert.Equal(t, testCase.src, dst)

		extraKey := "cloud100"
		dst[extraKey] = Cloud{
			Name: extraKey,
			K8sNodes: []K8sNode{
				K8sNode{
					Name: "node300",
				},
				K8sNode{
					Name: "node400",
				},
			},
		}
		if app, exist := testCase.src[extraKey]; exist {
			t.Errorf("key \"%s\" should not exist in testCase.src, the value is: %+v", extraKey, app)
		} else {
			t.Logf("As expected, key \"%s\" does not exist in testCase.src", extraKey)
		}

		// test change
		dst = CloudMapCopy(testCase.src)
		assert.Equal(t, testCase.src, dst)

		changeKey := "cloud1"
		dst[changeKey] = Cloud{
			Name: "cloud1",
			K8sNodes: []K8sNode{
				K8sNode{
					Name: "node150",
				},
				K8sNode{
					Name: "node250",
				},
			},
		}

		t.Logf("testCase.src[changeKey]:\n%+v\ndst[changeKey]:\n%+v", testCase.src[changeKey], dst[changeKey])
		assert.NotEqual(t, testCase.src[changeKey], dst[changeKey], fmt.Sprintf("testCase.src[changeKey] [%+v] and dst[changeKey] [%+v], show not be equal.", testCase.src[changeKey], dst[changeKey]))

	}

}

func TestGetSharedVmToCreate(t *testing.T) {
	testCases := []struct {
		name           string
		cloud          Cloud
		resPct         float64
		allRest        bool
		expectedResult models.IaasVm
	}{
		{
			name: "case 33%",
			cloud: Cloud{
				Name: "NOKIA4",
				Resources: models.ResourceStatus{
					Limit: models.ResSet{
						VCpu:    56,
						Ram:     128796.75390625,
						Storage: 1396.5185890197754,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
					InUse: models.ResSet{
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []K8sNode{},
			},
			resPct:  0.33,
			allRest: false,
			expectedResult: models.IaasVm{
				Name:    "auto-sched-nokia4-0",
				Cloud:   "NOKIA4",
				VCpu:    18,
				Ram:     42502,
				Storage: 460,
			},
		},
		{
			name: "case 50%",
			cloud: Cloud{
				Name: "NOKIA4",
				Resources: models.ResourceStatus{
					Limit: models.ResSet{
						VCpu:    56,
						Ram:     128796.75390625,
						Storage: 1396.5185890197754,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
					InUse: models.ResSet{
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []K8sNode{{Name: "auto-sched-nokia4-0"}},
			},
			resPct:  0.5,
			allRest: false,
			expectedResult: models.IaasVm{
				Name:    "auto-sched-nokia4-1",
				Cloud:   "NOKIA4",
				VCpu:    28,
				Ram:     64398,
				Storage: 698,
			},
		},
		{
			name: "all rest",
			cloud: Cloud{
				Name: "NOKIA4",
				Resources: models.ResourceStatus{
					Limit: models.ResSet{
						VCpu:    56,
						Ram:     128796.75390625,
						Storage: 1396.5185890197754,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
					InUse: models.ResSet{
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []K8sNode{
					K8sNode{Name: "auto-sched-nokia4-0"},
					K8sNode{Name: "auto-sched-nokia4-1"},
					K8sNode{Name: "auto-sched-nokia4-2"},
					K8sNode{Name: "auto-sched-nokia4-4"},
				},
			},
			resPct:  0,
			allRest: true,
			expectedResult: models.IaasVm{
				Name:    "auto-sched-nokia4-3",
				Cloud:   "NOKIA4",
				VCpu:    30,
				Ram:     69404,
				Storage: 767,
			},
		},
		{
			name: "all rest 2",
			cloud: Cloud{
				Name: "NOKIA4",
				Resources: models.ResourceStatus{
					Limit: models.ResSet{
						VCpu:    56,
						Ram:     128796.75390625,
						Storage: 1396.5185890197754,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
					InUse: models.ResSet{
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []K8sNode{
					K8sNode{Name: "auto-sched-nokia4-0"},
					K8sNode{Name: "auto-sched-nokia4-1"},
					K8sNode{Name: "auto-sched-nokia4-2"},
					K8sNode{Name: "auto-sched-nokia4-4"},
				},
			},
			resPct:  0.34,
			allRest: true,
			expectedResult: models.IaasVm{
				Name:    "auto-sched-nokia4-3",
				Cloud:   "NOKIA4",
				VCpu:    30,
				Ram:     69404,
				Storage: 767,
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := testCase.cloud.GetSharedVmToCreate(testCase.resPct, testCase.allRest)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestSupportCreateNewVM(t *testing.T) {
	testCases := []struct {
		name           string
		cloud          Cloud
		expectedResult bool
	}{
		{
			name: "case-proxmox",
			cloud: Cloud{
				Name: "NOKIA4",
				Type: models.ProxmoxIaas,
			},
			expectedResult: true,
		},
		{
			name: "case-openstack",
			cloud: Cloud{
				Name: "CLAAUDIAweifan",
				Type: models.OpenstackIaas,
			},
			expectedResult: false,
		},
		{
			name: "case-other",
			cloud: Cloud{
				Name: "HPE2",
				Type: "proxmox (powered off)",
			},
			expectedResult: false,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := testCase.cloud.SupportCreateNewVM()
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestGetNameVmToCreate(t *testing.T) {
	testCases := []struct {
		name           string
		cloud          Cloud
		expectedResult string
	}{
		{
			name: "case1",
			cloud: Cloud{
				Name: "NOKIA4",
				Resources: models.ResourceStatus{
					Limit: models.ResSet{
						VCpu:    56,
						Ram:     128796.75390625,
						Storage: 1396.5185890197754,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
					InUse: models.ResSet{
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []K8sNode{},
			},
			expectedResult: "auto-sched-nokia4-0",
		},
		{
			name: "case2",
			cloud: Cloud{
				Name: "NOKIA4",
				Resources: models.ResourceStatus{
					Limit: models.ResSet{
						VCpu:    56,
						Ram:     128796.75390625,
						Storage: 1396.5185890197754,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
					InUse: models.ResSet{
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []K8sNode{{Name: "auto-sched-nokia4-0"}},
			},
			expectedResult: "auto-sched-nokia4-1",
		},
		{
			name: "case3",
			cloud: Cloud{
				Name: "NOKIA4",
				Resources: models.ResourceStatus{
					Limit: models.ResSet{
						VCpu:    56,
						Ram:     128796.75390625,
						Storage: 1396.5185890197754,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
					InUse: models.ResSet{
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []K8sNode{
					K8sNode{Name: "auto-sched-nokia4-0"},
					K8sNode{Name: "auto-sched-nokia4-2"},
					K8sNode{Name: "auto-sched-nokia4-4"},
				},
			},
			expectedResult: "auto-sched-nokia4-1",
		},
		{
			name: "case4",
			cloud: Cloud{
				Name: "NOKIA4",
				Resources: models.ResourceStatus{
					Limit: models.ResSet{
						VCpu:    56,
						Ram:     128796.75390625,
						Storage: 1396.5185890197754,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
					InUse: models.ResSet{
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []K8sNode{
					K8sNode{Name: "auto-sched-nokia4-0"},
					K8sNode{Name: "auto-sched-nokia4-1"},
					K8sNode{Name: "auto-sched-nokia4-2"},
					K8sNode{Name: "auto-sched-nokia4-4"},
				},
			},
			expectedResult: "auto-sched-nokia4-3",
		},
		{
			name: "case5",
			cloud: Cloud{
				Name: "NOKIA4",
				Resources: models.ResourceStatus{
					Limit: models.ResSet{
						VCpu:    56,
						Ram:     128796.75390625,
						Storage: 1396.5185890197754,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
					InUse: models.ResSet{
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []K8sNode{
					K8sNode{Name: "auto-sched-nokia4-0"},
					K8sNode{Name: "auto-sched-nokia4-1"},
					K8sNode{Name: "auto-sched-nokia4-2"},
					K8sNode{Name: "auto-sched-nokia4-3"},
					K8sNode{Name: "auto-sched-nokia4-4"},
				},
			},
			expectedResult: "auto-sched-nokia4-5",
		},
		{
			name: "case6",
			cloud: Cloud{
				Name: "NOKIA4",
				Resources: models.ResourceStatus{
					Limit: models.ResSet{
						VCpu:    56,
						Ram:     128796.75390625,
						Storage: 1396.5185890197754,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
					InUse: models.ResSet{
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []K8sNode{
					K8sNode{Name: "asdfa"},
					K8sNode{Name: "dfgdf"},
					K8sNode{Name: "werwe"},
					K8sNode{Name: "asfdasf"},
					K8sNode{Name: "asdfaf"},
				},
			},
			expectedResult: "auto-sched-nokia4-0",
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := testCase.cloud.GetNameVmToCreate()
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}
