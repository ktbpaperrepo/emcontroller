package algorithms

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

func TestInnerSimulateCreateVm(t *testing.T) {
	_, apps, _ := cloudAppsSolnForTest()
	clouds := cloudsForTest()

	type caseStruct struct {
		name                  string
		simCloud              *asmodel.Cloud
		vmToCreate            models.IaasVm
		apps                  map[string]asmodel.Application
		appGroup              []string
		expectedSimCloudAfter asmodel.Cloud
	}
	var testCases []caseStruct

	func() {
		thisSimCloudCopy := asmodel.CloudCopy(clouds["nokia4WithOneNode"])
		thisVm := models.IaasVm{
			Name:    "auto-sched-nokia4-1",
			Cloud:   "NOKIA4",
			VCpu:    12,
			Ram:     8192,
			Storage: 100,
		}
		thisAppGroup := []string{"app3", "app6", "app7"}
		convertedK8sNode := asmodel.GenK8sNodeFromApps(thisVm, apps, thisAppGroup)

		testCases = append(testCases, caseStruct{
			name:       "cloud with no original nodes",
			simCloud:   &thisSimCloudCopy,
			vmToCreate: thisVm,
			apps:       apps,
			appGroup:   thisAppGroup,
			expectedSimCloudAfter: asmodel.Cloud{
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
						VCpu:    38,
						Ram:     67584,
						Storage: 729,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{
					{
						Name: "auto-sched-nokia4-0",
					},
					convertedK8sNode,
				}},
		})
	}()

	func() {
		thisSimCloudCopy := asmodel.CloudCopy(clouds["nokia4"])
		thisVm := models.IaasVm{
			Name:    "auto-sched-nokia4-0",
			Cloud:   "NOKIA4",
			VCpu:    8,
			Ram:     4096,
			Storage: 120,
		}
		var thisAppGroup []string
		convertedK8sNode := asmodel.GenK8sNodeFromApps(thisVm, apps, thisAppGroup)

		testCases = append(testCases, caseStruct{
			name:       "cloud with 1 original node",
			simCloud:   &thisSimCloudCopy,
			vmToCreate: thisVm,
			apps:       apps,
			appGroup:   thisAppGroup,
			expectedSimCloudAfter: asmodel.Cloud{
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
						VCpu:    34,
						Ram:     63488,
						Storage: 749,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{
					convertedK8sNode,
				}},
		})
	}()

	func() {
		thisSimCloudCopy := asmodel.CloudCopy(clouds["nokia4"])
		thisVm := models.IaasVm{
			Name:    "auto-sched-nokia4-0",
			Cloud:   "NOKIA4",
			VCpu:    40,
			Ram:     100000,
			Storage: 1000,
		}
		var thisAppGroup []string
		convertedK8sNode := asmodel.GenK8sNodeFromApps(thisVm, apps, thisAppGroup)

		testCases = append(testCases, caseStruct{
			name:       "overflow",
			simCloud:   &thisSimCloudCopy,
			vmToCreate: thisVm,
			apps:       apps,
			appGroup:   thisAppGroup,
			expectedSimCloudAfter: asmodel.Cloud{
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
						VCpu:    66,
						Ram:     159392,
						Storage: 1629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{
					convertedK8sNode,
				}},
		})
	}()

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		simulateCreateVm(testCase.simCloud, testCase.vmToCreate, testCase.apps, testCase.appGroup)
		assert.Equal(t, testCase.expectedSimCloudAfter, *testCase.simCloud, fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}

func TestInnerGetDedVmOneGroup(t *testing.T) {
	_, apps, _ := cloudAppsSolnForTest()
	clouds := cloudsForTest()

	testCases := []struct {
		name           string
		cloud          asmodel.Cloud
		apps           map[string]asmodel.Application
		appGroup       []string
		expectedResult models.IaasVm
	}{
		{
			name:     "case1",
			cloud:    asmodel.CloudCopy(clouds["nokia4WithOneNode"]),
			apps:     apps,
			appGroup: []string{"app3", "app6", "app7"},
			expectedResult: models.IaasVm{
				Name:    "auto-sched-nokia4-1",
				Cloud:   "NOKIA4",
				VCpu:    10,
				Ram:     3544,
				Storage: 89,
			},
		},
		{
			name:     "case2",
			cloud:    asmodel.CloudCopy(clouds["nokia4"]),
			apps:     apps,
			appGroup: []string{"app5", "app8", "app4"},
			expectedResult: models.IaasVm{
				Name:    "auto-sched-nokia4-0",
				Cloud:   "NOKIA4",
				VCpu:    13,
				Ram:     3214,
				Storage: 55,
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := getDedVmOneGroup(testCase.cloud, testCase.apps, testCase.appGroup)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}

func TestInnerGetDedicatedVmsToCreate(t *testing.T) {
	_, apps, _ := cloudAppsSolnForTest()
	clouds := cloudsForTest()

	type caseStruct struct {
		name               string
		cloud              *asmodel.Cloud
		apps               map[string]asmodel.Application
		appGroups          [][]string
		expectedVms        []models.IaasVm
		expectedCloudAfter asmodel.Cloud
	}
	var testCases []caseStruct

	func() {
		thisCloudCopy := asmodel.CloudCopy(clouds["nokia4WithOneNode"])
		var thisAppGroups [][]string = [][]string{
			[]string{
				"app3", "app4", "app1",
			},
			[]string{
				"app2", "app5", "app7",
			},
			[]string{
				"app6", "app4", "app8",
			},
		}

		testCases = append(testCases, caseStruct{
			name:      "case1",
			cloud:     &thisCloudCopy,
			apps:      apps,
			appGroups: thisAppGroups,
			expectedVms: []models.IaasVm{
				models.IaasVm{
					Name:    "auto-sched-nokia4-1",
					Cloud:   "NOKIA4",
					VCpu:    10,
					Ram:     3698,
					Storage: 49,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia4-2",
					Cloud:   "NOKIA4",
					VCpu:    13,
					Ram:     3544,
					Storage: 89,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia4-3",
					Cloud:   "NOKIA4",
					VCpu:    12,
					Ram:     3214,
					Storage: 55,
				},
			},
			expectedCloudAfter: asmodel.Cloud{
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
						VCpu:    61,
						Ram:     69848,
						Storage: 822,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{
					{
						Name: "auto-sched-nokia4-0",
					},
					asmodel.GenK8sNodeFromApps(models.IaasVm{
						Name:    "auto-sched-nokia4-1",
						Cloud:   "NOKIA4",
						VCpu:    10,
						Ram:     3698,
						Storage: 49,
					}, apps, []string{
						"app3", "app4", "app1",
					}),
					asmodel.GenK8sNodeFromApps(models.IaasVm{
						Name:    "auto-sched-nokia4-2",
						Cloud:   "NOKIA4",
						VCpu:    13,
						Ram:     3544,
						Storage: 89,
					}, apps, []string{
						"app2", "app5", "app7",
					}),
					asmodel.GenK8sNodeFromApps(models.IaasVm{
						Name:    "auto-sched-nokia4-3",
						Cloud:   "NOKIA4",
						VCpu:    12,
						Ram:     3214,
						Storage: 55,
					}, apps, []string{
						"app6", "app4", "app8",
					}),
				},
			},
		})
	}()

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualVms := getDedicatedVmsToCreate(testCase.cloud, testCase.apps, testCase.appGroups)
		assert.Equal(t, testCase.expectedVms, actualVms, fmt.Sprintf("%s: result is not expected", testCase.name))
		assert.Equal(t, testCase.expectedCloudAfter, *testCase.cloud, fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}
