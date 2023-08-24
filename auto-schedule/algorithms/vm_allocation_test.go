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

func TestInnerIsResEnough(t *testing.T) {
	testCases := []struct {
		name           string
		vm             asmodel.K8sNode
		app            asmodel.Application
		minCpu         bool
		expectedResult bool
	}{
		{
			name: "case enough equal 1",
			vm: asmodel.K8sNode{
				Name: "test-vm",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 1.0,
					Memory:  10240,
					Storage: 100,
				},
			},
			app: asmodel.Application{
				Name: "test-app",
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 4,
						Memory:  10240,
						Storage: 100,
					},
				},
			},
			minCpu:         true,
			expectedResult: true,
		},
		{
			name: "case enough equal 2",
			vm: asmodel.K8sNode{
				Name: "test-vm",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 1.0,
					Memory:  10243,
					Storage: 120,
				},
			},
			app: asmodel.Application{
				Name: "test-app",
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 1.0,
						Memory:  10243,
						Storage: 120,
					},
				},
			},
			minCpu:         true,
			expectedResult: true,
		},
		{
			name: "case enough more",
			vm: asmodel.K8sNode{
				Name: "test-vm",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 6.5,
					Memory:  10240,
					Storage: 100,
				},
			},
			app: asmodel.Application{
				Name: "test-app",
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.0,
						Memory:  1024,
						Storage: 10,
					},
				},
			},
			minCpu:         true,
			expectedResult: true,
		},
		{
			name: "case not enough less cpu 1",
			vm: asmodel.K8sNode{
				Name: "test-vm",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 0,
					Memory:  10240,
					Storage: 100,
				},
			},
			app: asmodel.Application{
				Name: "test-app",
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.0,
						Memory:  1024,
						Storage: 10,
					},
				},
			},
			minCpu:         true,
			expectedResult: false,
		},
		{
			name: "case not enough less cpu 2",
			vm: asmodel.K8sNode{
				Name: "test-vm",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 0.07,
					Memory:  10240,
					Storage: 100,
				},
			},
			app: asmodel.Application{
				Name: "test-app",
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.0,
						Memory:  1024,
						Storage: 10,
					},
				},
			},
			minCpu:         true,
			expectedResult: false,
		},
		{
			name: "case not enough memory",
			vm: asmodel.K8sNode{
				Name: "test-vm",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 6.5,
					Memory:  10240,
					Storage: 100,
				},
			},
			app: asmodel.Application{
				Name: "test-app",
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.0,
						Memory:  10241,
						Storage: 10,
					},
				},
			},
			minCpu:         true,
			expectedResult: false,
		},
		{
			name: "case not enough storage",
			vm: asmodel.K8sNode{
				Name: "test-vm",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 6.5,
					Memory:  10240,
					Storage: 100,
				},
			},
			app: asmodel.Application{
				Name: "test-app",
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.0,
						Memory:  1024,
						Storage: 1101,
					},
				},
			},
			minCpu:         true,
			expectedResult: false,
		},
		{
			name: "case not enough storage not minCpu",
			vm: asmodel.K8sNode{
				Name: "test-vm",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 6.5,
					Memory:  10240,
					Storage: 100,
				},
			},
			app: asmodel.Application{
				Name: "test-app",
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.0,
						Memory:  1024,
						Storage: 1101,
					},
				},
			},
			minCpu:         false,
			expectedResult: false,
		},
		{
			name: "case not enough CPU not minCpu",
			vm: asmodel.K8sNode{
				Name: "test-vm",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 4.5,
					Memory:  10240,
					Storage: 100,
				},
			},
			app: asmodel.Application{
				Name: "test-app",
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.0,
						Memory:  1024,
						Storage: 50,
					},
				},
			},
			minCpu:         false,
			expectedResult: false,
		},
		{
			name: "control group of last one, with minCpu",
			vm: asmodel.K8sNode{
				Name: "test-vm",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 4.5,
					Memory:  10240,
					Storage: 100,
				},
			},
			app: asmodel.Application{
				Name: "test-app",
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.0,
						Memory:  1024,
						Storage: 50,
					},
				},
			},
			minCpu:         true,
			expectedResult: true,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := isResEnough(testCase.vm, testCase.app, testCase.minCpu)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerSubRes(t *testing.T) {

	func() {
		t.Log("test case with minCpu.")

		var vm *asmodel.K8sNode = &asmodel.K8sNode{
			Name: "test-vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 6.5,
				Memory:  10240,
				Storage: 100,
			},
		}

		var app asmodel.Application

		t.Log("First time.")
		app = asmodel.Application{
			Name: "test-app",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 5.0,
					Memory:  1024,
					Storage: 9,
				},
			},
		}
		subRes(vm, app, true)
		assert.InDelta(t, 6.4, vm.ResidualResources.CpuCore, testDelta)
		assert.InDelta(t, 9216, vm.ResidualResources.Memory, testDelta)
		assert.InDelta(t, 91, vm.ResidualResources.Storage, testDelta)

		t.Log("Second time.")
		app = asmodel.Application{
			Name: "test-app",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 3.0,
					Memory:  2000,
					Storage: 0,
				},
			},
		}
		subRes(vm, app, true)
		assert.InDelta(t, 6.3, vm.ResidualResources.CpuCore, testDelta)
		assert.InDelta(t, 7216, vm.ResidualResources.Memory, testDelta)
		assert.InDelta(t, 91, vm.ResidualResources.Storage, testDelta)
	}()

	func() {
		t.Log("test case without minCpu.")

		var vm *asmodel.K8sNode = &asmodel.K8sNode{
			Name: "test-vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 6.5,
				Memory:  10240,
				Storage: 100,
			},
		}

		var app asmodel.Application

		t.Log("First time.")
		app = asmodel.Application{
			Name: "test-app",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 3.1,
					Memory:  1024,
					Storage: 9,
				},
			},
		}
		subRes(vm, app, false)
		assert.InDelta(t, 3.4, vm.ResidualResources.CpuCore, testDelta)
		assert.InDelta(t, 9216, vm.ResidualResources.Memory, testDelta)
		assert.InDelta(t, 91, vm.ResidualResources.Storage, testDelta)

		t.Log("Second time.")
		app = asmodel.Application{
			Name: "test-app",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 2.0,
					Memory:  2000,
					Storage: 0,
				},
			},
		}
		subRes(vm, app, false)
		assert.InDelta(t, 1.4, vm.ResidualResources.CpuCore, testDelta)
		assert.InDelta(t, 7216, vm.ResidualResources.Memory, testDelta)
		assert.InDelta(t, 91, vm.ResidualResources.Storage, testDelta)
	}()
}

func TestInnerVmResMeetAllRestAppsAll(t *testing.T) {
	TestInnerVmResMeetAllRestAppsNoPriLimit(t)
	TestInnerVmResMeetAllRestAppsMaxPri(t)
	TestInnerVmResMeetAllRestAppsNotMaxPri(t)
}

func TestInnerVmResMeetAllRestAppsNoPriLimit(t *testing.T) {
	cloud, apps, soln := cloudAppsSolnForTest()
	appsThisCloud := findAppsOneCloud(cloud, apps, soln) // need 0.4 CPU, 3214 Memory, 66 Storage in total.
	appsOrder := appOrdersForTest()[0]

	t.Log()
	t.Log("list the apps scheduled to this cloud in order")
	var appsThisCloudIter *iterForApps = newIterForApps(appsThisCloud, appsOrder)
	var curAppName string = appsThisCloudIter.nextAppName()
	for len(curAppName) != 0 {
		t.Logf("curApp: %+v\n", appsThisCloud[curAppName])
		curAppName = appsThisCloudIter.nextAppName()
	}

	var appNamesToThisVm []string
	var meetAllRest bool

	t.Log()
	t.Log("case 1")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()
	vm := asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 2,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app1", "app2", "app4", "app7"})

	t.Log()
	t.Log("case 2")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.2,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app1", "app2"})
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.1,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app4"})
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.3,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app7"})
	t.Log("curAppName:", curAppName)

	t.Log()
	t.Log("case 3")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  2000,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app1"})
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1200,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app2"})
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.2,
			Memory:  3000,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app4", "app7"})
	t.Log("curAppName:", curAppName)

	t.Log()
	t.Log("case 4: copy")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  3000,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app1", "app2", "app4"})
	t.Log("curAppName:", curAppName)

	t.Log("copy 1, should be true")
	iterCopy1 := appsThisCloudIter.Copy()
	curAppNameCopy1 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  3000,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy1, iterCopy1.nextAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app7"})
	t.Log("curAppNameCopy1:", curAppNameCopy1)

	t.Log("copy 2, should be false")
	iterCopy2 := appsThisCloudIter.Copy()
	curAppNameCopy2 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  10,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy2, iterCopy2.nextAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{})
	t.Log("curAppNameCopy2:", curAppNameCopy2)

	t.Log("copy 3, should be true")
	iterCopy3 := appsThisCloudIter.Copy()
	curAppNameCopy3 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1200,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy3, iterCopy3.nextAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app7"})
	t.Log("curAppNameCopy3:", curAppNameCopy3)

	t.Log("copy 4, should be false")
	iterCopy4 := appsThisCloudIter.Copy()
	curAppNameCopy4 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  10,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy4, iterCopy4.nextAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{})
	t.Log("curAppNameCopy4:", curAppNameCopy4)

	func() {
		t.Log()
		t.Log("case 5: copy without minCpu")
		appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
		curAppName = appsThisCloudIter.nextAppName()
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 10.0,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextAppName, false)
		assert.False(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{"app1", "app2"})
		t.Log("curAppName:", curAppName)

		t.Log("copy 1, should be true")
		iterCopy1 := appsThisCloudIter.Copy()
		curAppNameCopy1 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 10,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy1, iterCopy1.nextAppName, false)
		assert.True(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{"app4", "app7"})
		t.Log("curAppNameCopy1:", curAppNameCopy1)

		t.Log("copy 2, should be false")
		iterCopy2 := appsThisCloudIter.Copy()
		curAppNameCopy2 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 3,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy2, iterCopy2.nextAppName, false)
		assert.False(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{})
		t.Log("curAppNameCopy2:", curAppNameCopy2)

		t.Log("copy 3, should be true")
		iterCopy3 := appsThisCloudIter.Copy()
		curAppNameCopy3 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 10,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy3, iterCopy3.nextAppName, false)
		assert.True(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{"app4", "app7"})
		t.Log("curAppNameCopy3:", curAppNameCopy3)

		t.Log("copy 4, should be false")
		iterCopy4 := appsThisCloudIter.Copy()
		curAppNameCopy4 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 3,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy4, iterCopy4.nextAppName, false)
		assert.False(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{})
		t.Log("curAppNameCopy4:", curAppNameCopy4)
	}()

}

func TestInnerVmResMeetAllRestAppsMaxPri(t *testing.T) {
	cloud, apps, soln := cloudAppsSolnForTest()
	appsThisCloud := findAppsOneCloud(cloud, apps, soln) // need 0.4 CPU, 3214 Memory, 66 Storage in total.
	appsOrder := appOrdersForTest()[0]

	t.Log()
	t.Log("list the MaxPri apps scheduled to this cloud in order")
	var appsThisCloudIter *iterForApps = newIterForApps(appsThisCloud, appsOrder)
	var curAppName string = appsThisCloudIter.nextMaxPriAppName()
	for len(curAppName) != 0 {
		t.Logf("curApp: %+v\n", appsThisCloud[curAppName])
		curAppName = appsThisCloudIter.nextMaxPriAppName()
	}

	var appNamesToThisVm []string
	var meetAllRest bool

	t.Log()
	t.Log("case 1")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextMaxPriAppName()
	vm := asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 2,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextMaxPriAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app2", "app7"})

	t.Log()
	t.Log("case 2")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextMaxPriAppName()

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.1,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app2"})
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.05,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{})
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.3,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextMaxPriAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app7"})
	t.Log("curAppName:", curAppName)

	t.Log()
	t.Log("case 3")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextMaxPriAppName()

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1000,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app2"})
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1000,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextMaxPriAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app7"})
	t.Log("curAppName:", curAppName)

	t.Log()
	t.Log("case 4: copy")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextMaxPriAppName()
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1000,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app2"})
	t.Log("curAppName:", curAppName)

	t.Log("copy 1, should be true")
	iterCopy1 := appsThisCloudIter.Copy()
	curAppNameCopy1 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1000,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy1, iterCopy1.nextMaxPriAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app7"})
	t.Log("curAppNameCopy1:", curAppNameCopy1)

	t.Log("copy 2, should be false")
	iterCopy2 := appsThisCloudIter.Copy()
	curAppNameCopy2 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  10,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy2, iterCopy2.nextMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{})
	t.Log("curAppNameCopy2:", curAppNameCopy2)

	t.Log("copy 3, should be true")
	iterCopy3 := appsThisCloudIter.Copy()
	curAppNameCopy3 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1200,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy3, iterCopy3.nextMaxPriAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app7"})
	t.Log("curAppNameCopy3:", curAppNameCopy3)

	t.Log("copy 4, should be false")
	iterCopy4 := appsThisCloudIter.Copy()
	curAppNameCopy4 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  10,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy4, iterCopy4.nextMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{})
	t.Log("curAppNameCopy4:", curAppNameCopy4)

	func() {
		t.Log()
		t.Log("case 5: copy without minCpu")
		appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
		curAppName = appsThisCloudIter.nextMaxPriAppName()
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 4.0,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextMaxPriAppName, false)
		assert.False(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{"app2"})
		t.Log("curAppName:", curAppName)

		t.Log("copy 1, should be true")
		iterCopy1 := appsThisCloudIter.Copy()
		curAppNameCopy1 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 4,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy1, iterCopy1.nextMaxPriAppName, false)
		assert.True(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{"app7"})
		t.Log("curAppNameCopy1:", curAppNameCopy1)

		t.Log("copy 2, should be false")
		iterCopy2 := appsThisCloudIter.Copy()
		curAppNameCopy2 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 2.5,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy2, iterCopy2.nextMaxPriAppName, false)
		assert.False(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{})
		t.Log("curAppNameCopy2:", curAppNameCopy2)

		t.Log("copy 3, should be true")
		iterCopy3 := appsThisCloudIter.Copy()
		curAppNameCopy3 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 4,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy3, iterCopy3.nextMaxPriAppName, false)
		assert.True(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{"app7"})
		t.Log("curAppNameCopy3:", curAppNameCopy3)

		t.Log("copy 4, should be false")
		iterCopy4 := appsThisCloudIter.Copy()
		curAppNameCopy4 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 2.5,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy4, iterCopy4.nextMaxPriAppName, false)
		assert.False(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{})
		t.Log("curAppNameCopy4:", curAppNameCopy4)
	}()

}

func TestInnerVmResMeetAllRestAppsNotMaxPri(t *testing.T) {
	cloud, apps, soln := cloudAppsSolnForTest()
	appsThisCloud := findAppsOneCloud(cloud, apps, soln) // need 0.4 CPU, 3214 Memory, 66 Storage in total.
	appsOrder := appOrdersForTest()[0]

	t.Log()
	t.Log("list the NotMaxPri apps scheduled to this cloud in order")
	var appsThisCloudIter *iterForApps = newIterForApps(appsThisCloud, appsOrder)
	var curAppName string = appsThisCloudIter.nextNotMaxPriAppName()
	for len(curAppName) != 0 {
		t.Logf("curApp: %+v\n", appsThisCloud[curAppName])
		curAppName = appsThisCloudIter.nextNotMaxPriAppName()
	}

	var appNamesToThisVm []string
	var meetAllRest bool

	t.Log()
	t.Log("case 1")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextNotMaxPriAppName()
	vm := asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 2,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextNotMaxPriAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app1", "app4"})

	t.Log()
	t.Log("case 2")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextNotMaxPriAppName()

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.1,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextNotMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app1"})
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.05,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextNotMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{})
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.3,
			Memory:  10240,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextNotMaxPriAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app4"})
	t.Log("curAppName:", curAppName)

	t.Log()
	t.Log("case 3")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextNotMaxPriAppName()

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1100,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextNotMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app1"})
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  500,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextNotMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{})
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1100,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextNotMaxPriAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app4"})
	t.Log("curAppName:", curAppName)

	t.Log()
	t.Log("case 4: copy")
	appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextNotMaxPriAppName()
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1100,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextNotMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app1"})
	t.Log("curAppName:", curAppName)

	t.Log("copy 1, should be true")
	iterCopy1 := appsThisCloudIter.Copy()
	curAppNameCopy1 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1100,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy1, iterCopy1.nextNotMaxPriAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app4"})
	t.Log("curAppNameCopy1:", curAppNameCopy1)

	t.Log("copy 2, should be false")
	iterCopy2 := appsThisCloudIter.Copy()
	curAppNameCopy2 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  10,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy2, iterCopy2.nextNotMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{})
	t.Log("curAppNameCopy2:", curAppNameCopy2)

	t.Log("copy 3, should be true")
	iterCopy3 := appsThisCloudIter.Copy()
	curAppNameCopy3 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1200,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy3, iterCopy3.nextNotMaxPriAppName, true)
	assert.True(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{"app4"})
	t.Log("curAppNameCopy3:", curAppNameCopy3)

	t.Log("copy 4, should be false")
	iterCopy4 := appsThisCloudIter.Copy()
	curAppNameCopy4 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  10,
			Storage: 100,
		},
	}
	appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy4, iterCopy4.nextNotMaxPriAppName, true)
	assert.False(t, meetAllRest)
	assert.ElementsMatch(t, appNamesToThisVm, []string{})
	t.Log("curAppNameCopy4:", curAppNameCopy4)

	func() {
		t.Log()
		t.Log("case 5: copy without minCpu")
		appsThisCloudIter = newIterForApps(appsThisCloud, appsOrder)
		curAppName = appsThisCloudIter.nextNotMaxPriAppName()
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 5.2,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter.nextNotMaxPriAppName, false)
		assert.False(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{"app1"})
		t.Log("curAppName:", curAppName)

		t.Log("copy 1, should be true")
		iterCopy1 := appsThisCloudIter.Copy()
		curAppNameCopy1 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 5.5,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy1, iterCopy1.nextNotMaxPriAppName, false)
		assert.True(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{"app4"})
		t.Log("curAppNameCopy1:", curAppNameCopy1)

		t.Log("copy 2, should be false")
		iterCopy2 := appsThisCloudIter.Copy()
		curAppNameCopy2 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 2.3,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy2, iterCopy2.nextNotMaxPriAppName, false)
		assert.False(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{})
		t.Log("curAppNameCopy2:", curAppNameCopy2)

		t.Log("copy 3, should be true")
		iterCopy3 := appsThisCloudIter.Copy()
		curAppNameCopy3 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 5.6,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy3, iterCopy3.nextNotMaxPriAppName, false)
		assert.True(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{"app4"})
		t.Log("curAppNameCopy3:", curAppNameCopy3)

		t.Log("copy 4, should be false")
		iterCopy4 := appsThisCloudIter.Copy()
		curAppNameCopy4 := curAppName
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 2.1,
				Memory:  30000,
				Storage: 1000,
			},
		}
		appNamesToThisVm, meetAllRest = vmResMeetAllRestApps(vm, apps, &curAppNameCopy4, iterCopy4.nextNotMaxPriAppName, false)
		assert.False(t, meetAllRest)
		assert.ElementsMatch(t, appNamesToThisVm, []string{})
		t.Log("curAppNameCopy4:", curAppNameCopy4)
	}()

}
