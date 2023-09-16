package algorithms

import (
	"fmt"
	"math"
	"testing"

	"github.com/KeepTheBeats/routing-algorithms/mymath"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"

	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

func TestInnerGroupAppsByVm(t *testing.T) {
	cloud, apps, soln := cloudAppsSolnForTest()

	testCases := []struct {
		name           string
		appsToGroup    map[string]asmodel.Application
		appsOrder      []string
		solnWithVm     asmodel.Solution
		expectedResult map[string][]string
	}{
		{
			name:        "case1",
			appsToGroup: findAppsOneCloud(cloud, apps, soln),
			appsOrder:   appOrdersForTest()[0],
			solnWithVm:  soln,
			expectedResult: map[string][]string{
				"auto-sched-cloud1-1": []string{"app1", "app2"},
				"auto-sched-cloud1-2": []string{"app4"},
				"auto-sched-cloud1-3": []string{"app7"},
			},
		},
		{
			name:        "case2",
			appsToGroup: findAppsOneCloud(cloud, apps, soln),
			appsOrder:   appOrdersForTest()[0],
			solnWithVm:  solnsForTest()[0],
			expectedResult: map[string][]string{
				"auto-sched-cloud1-1": []string{"app1", "app2"},
				"auto-sched-cloud1-2": []string{"app4"},
				"auto-sched-cloud1-3": []string{"app7"},
			},
		},
		{
			name:        "case3",
			appsToGroup: findAppsOneCloud(cloud, apps, solnsForTest()[0]),
			appsOrder:   appOrdersForTest()[0],
			solnWithVm:  solnsForTest()[0],
			expectedResult: map[string][]string{
				"auto-sched-cloud1-1": []string{"app1", "app2"},
				"auto-sched-cloud1-2": []string{"app4", "app6"},
				"auto-sched-cloud1-3": []string{"app7"},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := groupAppsByVm(testCase.appsToGroup, testCase.appsOrder, testCase.solnWithVm)

		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))

		// This unit test can ensure the fixed application order, so we can directly use assert.Equal, and do not need the following way.
		/**
		assert.Equal(t, len(testCase.expectedResult), len(actualResult), fmt.Sprintf("%s: result length is not expected", testCase.name))
		for key, _ := range testCase.expectedResult {
			assert.ElementsMatch(t, testCase.expectedResult[key], actualResult[key], fmt.Sprintf("%s: result is not expected", testCase.name))
		}
		*/

	}

}

func TestInnerGetVmByName(t *testing.T) {
	testCases := []struct {
		name           string
		tgtVmName      string
		cloud          asmodel.Cloud
		solnWithVm     asmodel.Solution
		expectedResult asmodel.K8sNode
		expectedPanic  bool
	}{
		{
			name:       "found ori",
			tgtVmName:  "nokia4-ori-node1",
			cloud:      cloudsForTest()["nokia4With2OriNode"],
			solnWithVm: solnsForTest()[1],
			expectedResult: asmodel.K8sNode{
				Name: "nokia4-ori-node1",
			},
			expectedPanic: false,
		},
		{
			name:       "found create",
			tgtVmName:  "auto-sched-nokia4-4",
			cloud:      cloudsForTest()["nokia4With2OriNode"],
			solnWithVm: solnsForTest()[1],
			expectedResult: asmodel.GenK8sNodeFromPods(models.IaasVm{
				Name:    "auto-sched-nokia4-4",
				Cloud:   "NOKIA4",
				VCpu:    8,
				Ram:     3122,
				Storage: 44,
			}, []apiv1.Pod{}),
			expectedPanic: false,
		},
		{
			name:           "cannot found incorrect cloud",
			tgtVmName:      "auto-sched-nokia6-3",
			cloud:          cloudsForTest()["nokia4With2OriNode"],
			solnWithVm:     solnsForTest()[1],
			expectedResult: asmodel.K8sNode{},
			expectedPanic:  true,
		},
		{
			name:           "cannot found target VM name does not exist",
			tgtVmName:      "asdfasf",
			cloud:          cloudsForTest()["nokia4With2OriNode"],
			solnWithVm:     solnsForTest()[1],
			expectedResult: asmodel.K8sNode{},
			expectedPanic:  true,
		},
		{
			name:           "both Found",
			tgtVmName:      "auto-sched-nokia4-4",
			cloud:          cloudsForTest()["nokia4With2SchedNode"],
			solnWithVm:     solnsForTest()[1],
			expectedResult: asmodel.K8sNode{},
			expectedPanic:  true,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		var actualResult asmodel.K8sNode
		runFunc := func() {
			actualResult = getVmByName(testCase.tgtVmName, testCase.cloud, testCase.solnWithVm)
		}

		if testCase.expectedPanic {
			assert.Panics(t, runFunc)
		} else {
			assert.NotPanics(t, runFunc)
			assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
		}

	}
}

func TestInnerCalcAppWeight(t *testing.T) {
	_, apps, _ := cloudAppsSolnForTest()

	testCases := []struct {
		name           string
		app            asmodel.Application
		expectedResult float64
	}{
		{
			name:           "case1",
			app:            apps["app1"],
			expectedResult: 5,
		},
		{
			name:           "case2",
			app:            apps["app2"],
			expectedResult: 10,
		},
		{
			name:           "case3",
			app:            apps["app3"],
			expectedResult: 1,
		},
		{
			name:           "case4",
			app:            apps["app4"],
			expectedResult: 3,
		},
		{
			name:           "case5",
			app:            apps["app5"],
			expectedResult: 2,
		},
		{
			name:           "case6",
			app:            apps["app6"],
			expectedResult: 7,
		},
		{
			name:           "case7",
			app:            apps["app7"],
			expectedResult: 10,
		},
		{
			name:           "case8",
			app:            apps["app8"],
			expectedResult: 8,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		actualResult := calcAppWeight(testCase.app)
		assert.InDelta(t, testCase.expectedResult, actualResult, testDelta, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerCalcAppsSumWeight(t *testing.T) {
	_, apps, _ := cloudAppsSolnForTest()

	testCases := []struct {
		name           string
		apps           map[string]asmodel.Application
		expectedResult float64
	}{
		{
			name:           "case1",
			apps:           apps,
			expectedResult: 46,
		},
		{
			name:           "case2",
			apps:           appsForTest()[0],
			expectedResult: 46,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		actualResult := calcAppsSumWeight(testCase.apps)
		assert.InDelta(t, testCase.expectedResult, actualResult, testDelta, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerCalcCpuOneApp(t *testing.T) {
	_, apps, _ := cloudAppsSolnForTest()

	testCases := []struct {
		name           string
		vm             asmodel.K8sNode
		apps           map[string]asmodel.Application
		thisAppName    string
		expectedResult float64
	}{
		{
			name: "case1",
			vm: asmodel.K8sNode{
				Name: "node1",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 14.4,
					Memory:  6668,
					Storage: 150,
				},
			},
			thisAppName:    "app2",
			apps:           apps,
			expectedResult: 14.4 * calcAppWeight(apps["app2"]) / calcAppsSumWeight(apps),
		},
		{
			name: "case2",
			vm: asmodel.K8sNode{
				Name: "node1",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 14.4,
					Memory:  6668,
					Storage: 150,
				},
			},
			thisAppName:    "app5",
			apps:           apps,
			expectedResult: 14.4 * calcAppWeight(apps["app5"]) / calcAppsSumWeight(apps),
		},
		{
			name: "case3",
			vm: asmodel.K8sNode{
				Name: "node2",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 7.6,
					Memory:  4096,
					Storage: 80,
				},
			},
			thisAppName:    "app3",
			apps:           appsForTest()[0],
			expectedResult: 7.6 * calcAppWeight(appsForTest()[0]["app3"]) / calcAppsSumWeight(appsForTest()[0]),
		},
		{
			name: "case3",
			vm: asmodel.K8sNode{
				Name: "node2",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 7.6,
					Memory:  4096,
					Storage: 80,
				},
			},
			thisAppName:    "app6",
			apps:           appsForTest()[0],
			expectedResult: 7.6 * calcAppWeight(appsForTest()[0]["app6"]) / calcAppsSumWeight(appsForTest()[0]),
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		actualResult := calcCpuOneApp(testCase.vm, testCase.apps, testCase.thisAppName)
		assert.InDelta(t, testCase.expectedResult, actualResult, testDelta, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerCheckValidDistriCpu(t *testing.T) {
	_, apps, _ := cloudAppsSolnForTest()

	testCases := []struct {
		name          string
		vm            asmodel.K8sNode
		apps          map[string]asmodel.Application
		expectedPanic bool
	}{
		{
			name: "case panic",
			vm: asmodel.K8sNode{
				Name: "node1",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 7,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps:          apps,
			expectedPanic: true,
		},
		{
			name: "case no panic1",
			vm: asmodel.K8sNode{
				Name: "node1",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 8,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps:          apps,
			expectedPanic: false,
		},
		{
			name: "case no panic2",
			vm: asmodel.K8sNode{
				Name: "node1",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 9,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps:          apps,
			expectedPanic: false,
		},
		{
			name: "case panic2",
			vm: asmodel.K8sNode{
				Name: "node1",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 5,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps:          appsForTest()[1],
			expectedPanic: true,
		},
		{
			name: "case no panic3",
			vm: asmodel.K8sNode{
				Name: "node1",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 9,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps:          appsForTest()[1],
			expectedPanic: false,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		runFunc := func() {
			checkValidDistriCpu(testCase.vm, testCase.apps)
		}

		if testCase.expectedPanic {
			assert.Panics(t, runFunc)
		} else {
			assert.NotPanics(t, runFunc)
		}

	}
}

func TestInnerDistrCpuNextApp(t *testing.T) {
	_, apps, _ := cloudAppsSolnForTest()
	apps = asmodel.AppMapCopy(apps)
	appOrder := appOrdersForTest()[1]
	vm := asmodel.K8sNode{
		Name: "node1",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 200,
			Memory:  4096,
			Storage: 80,
		},
	}

	t.Log("Round 1")
	thisAppName, allocatedCpu := distrCpuNextApp(vm, apps, appOrder)
	assert.Equal(t, "app3", thisAppName)
	assert.InDelta(t, vm.ResidualResources.CpuCore*calcAppWeight(apps["app3"])/calcAppsSumWeight(apps), allocatedCpu, testDelta)

	var actualAllocatedCpu float64
	if math.Abs(allocatedCpu-mymath.UnitRound(allocatedCpu, cpuCoreStep)) < floatDelta {
		// e.g., because of the inaccuracy of binary-floating-point data, a value like 2 may be represented as 1.999999999999, so its floor will be 1 rather than 2, but we need its floor to be 2, so we do this.
		actualAllocatedCpu = mymath.UnitRound(allocatedCpu, cpuCoreStep)
	} else {
		actualAllocatedCpu = mymath.UnitFloor(allocatedCpu, cpuCoreStep)
	}

	if actualAllocatedCpu > apps[thisAppName].Resources.CpuCore {
		actualAllocatedCpu = apps[thisAppName].Resources.CpuCore
	}

	delete(apps, thisAppName)
	vm.ResidualResources.CpuCore -= actualAllocatedCpu
	assert.Equal(t, 7, len(apps))

	t.Log("Round 2")
	thisAppName, allocatedCpu = distrCpuNextApp(vm, apps, appOrder)
	assert.Equal(t, "app2", thisAppName)
	assert.InDelta(t, vm.ResidualResources.CpuCore*calcAppWeight(apps["app2"])/calcAppsSumWeight(apps), allocatedCpu, testDelta)

	if math.Abs(allocatedCpu-mymath.UnitRound(allocatedCpu, cpuCoreStep)) < floatDelta {
		// e.g., because of the inaccuracy of binary-floating-point data, a value like 2 may be represented as 1.999999999999, so its floor will be 1 rather than 2, but we need its floor to be 2, so we do this.
		actualAllocatedCpu = mymath.UnitRound(allocatedCpu, cpuCoreStep)
	} else {
		actualAllocatedCpu = mymath.UnitFloor(allocatedCpu, cpuCoreStep)
	}

	if actualAllocatedCpu > apps[thisAppName].Resources.CpuCore {
		actualAllocatedCpu = apps[thisAppName].Resources.CpuCore
	}

	delete(apps, thisAppName)
	vm.ResidualResources.CpuCore -= actualAllocatedCpu
	assert.Equal(t, 6, len(apps))

	t.Log("Round 3")
	thisAppName, allocatedCpu = distrCpuNextApp(vm, apps, appOrder)
	assert.Equal(t, "app5", thisAppName)
	assert.InDelta(t, vm.ResidualResources.CpuCore*calcAppWeight(apps["app5"])/calcAppsSumWeight(apps), allocatedCpu, testDelta)

	if math.Abs(allocatedCpu-mymath.UnitRound(allocatedCpu, cpuCoreStep)) < floatDelta {
		// e.g., because of the inaccuracy of binary-floating-point data, a value like 2 may be represented as 1.999999999999, so its floor will be 1 rather than 2, but we need its floor to be 2, so we do this.
		actualAllocatedCpu = mymath.UnitRound(allocatedCpu, cpuCoreStep)
	} else {
		actualAllocatedCpu = mymath.UnitFloor(allocatedCpu, cpuCoreStep)
	}

	if actualAllocatedCpu > apps[thisAppName].Resources.CpuCore {
		actualAllocatedCpu = apps[thisAppName].Resources.CpuCore
	}

	delete(apps, thisAppName)
	vm.ResidualResources.CpuCore -= actualAllocatedCpu
	assert.Equal(t, 5, len(apps))

	t.Log("Round 4")
	thisAppName, allocatedCpu = distrCpuNextApp(vm, apps, appOrder)
	assert.Equal(t, "app1", thisAppName)
	assert.InDelta(t, vm.ResidualResources.CpuCore*calcAppWeight(apps["app1"])/calcAppsSumWeight(apps), allocatedCpu, testDelta)

	if math.Abs(allocatedCpu-mymath.UnitRound(allocatedCpu, cpuCoreStep)) < floatDelta {
		// e.g., because of the inaccuracy of binary-floating-point data, a value like 2 may be represented as 1.999999999999, so its floor will be 1 rather than 2, but we need its floor to be 2, so we do this.
		actualAllocatedCpu = mymath.UnitRound(allocatedCpu, cpuCoreStep)
	} else {
		actualAllocatedCpu = mymath.UnitFloor(allocatedCpu, cpuCoreStep)
	}

	if actualAllocatedCpu > apps[thisAppName].Resources.CpuCore {
		actualAllocatedCpu = apps[thisAppName].Resources.CpuCore
	}

	delete(apps, thisAppName)
	vm.ResidualResources.CpuCore -= actualAllocatedCpu
	assert.Equal(t, 4, len(apps))
}

func TestInnerDistrCpuApps(t *testing.T) {
	_, apps, _ := cloudAppsSolnForTest()

	testCases := []struct {
		name           string
		vm             asmodel.K8sNode
		apps           map[string]asmodel.Application
		expectedResult map[string]float64
	}{
		{
			name: "case1",
			vm: asmodel.K8sNode{
				Name: "node1",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 20,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps: apps,
			expectedResult: map[string]float64{
				"app1": 20.0 * calcAppWeight(apps["app1"]) / calcAppsSumWeight(apps),
				"app2": 20.0 * calcAppWeight(apps["app2"]) / calcAppsSumWeight(apps),
				"app3": 20.0 * calcAppWeight(apps["app3"]) / calcAppsSumWeight(apps),
				"app4": 20.0 * calcAppWeight(apps["app4"]) / calcAppsSumWeight(apps),
				"app5": 20.0 * calcAppWeight(apps["app5"]) / calcAppsSumWeight(apps),
				"app6": 20.0 * calcAppWeight(apps["app6"]) / calcAppsSumWeight(apps),
				"app7": 20.0 * calcAppWeight(apps["app7"]) / calcAppsSumWeight(apps),
				"app8": 20.0 * calcAppWeight(apps["app8"]) / calcAppsSumWeight(apps),
			},
		},
		{
			name: "case2",
			vm: asmodel.K8sNode{
				Name: "node2",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 7.3,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps: appsForTest()[1],
			expectedResult: map[string]float64{
				"app1": 7.3 * calcAppWeight(appsForTest()[1]["app1"]) / calcAppsSumWeight(appsForTest()[1]),
				"app2": 7.3 * calcAppWeight(appsForTest()[1]["app2"]) / calcAppsSumWeight(appsForTest()[1]),
				"app3": 7.3 * calcAppWeight(appsForTest()[1]["app3"]) / calcAppsSumWeight(appsForTest()[1]),
				"app4": 7.3 * calcAppWeight(appsForTest()[1]["app4"]) / calcAppsSumWeight(appsForTest()[1]),
				"app5": 7.3 * calcAppWeight(appsForTest()[1]["app5"]) / calcAppsSumWeight(appsForTest()[1]),
				"app6": 7.3 * calcAppWeight(appsForTest()[1]["app6"]) / calcAppsSumWeight(appsForTest()[1]),
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		actualResult := distrCpuApps(testCase.vm, testCase.apps)
		assert.InDeltaMapValues(t, testCase.expectedResult, actualResult, testDelta, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerAllocateCpuAsRequest(t *testing.T) {
	testCases := []struct {
		name           string
		appsToHandle   map[string]asmodel.Application
		solnWithVm     asmodel.Solution
		expectedResult asmodel.Solution
	}{
		{
			name:         "case1",
			appsToHandle: appsForTest()[2],
			solnWithVm:   solnsForTest()[2],
			expectedResult: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA4",
						K8sNodeName:      "nokia4-ori-node1",
						AllocatedCpuCore: 5.2,
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA4",
						K8sNodeName:      "nokia4-ori-node1",
						AllocatedCpuCore: 1.3,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA4",
						K8sNodeName:      "nokia4-ori-node1",
						AllocatedCpuCore: 5.6,
					},
				},
			},
		},
		{
			name:         "case2",
			appsToHandle: appsForTest()[3],
			solnWithVm:   solnsForTest()[3],
			expectedResult: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app2": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1.3,
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2.2,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 0.5,
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		actualResult := allocateCpuAsRequest(testCase.appsToHandle, testCase.solnWithVm)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerVmCpuWeightedAllocation(t *testing.T) {
	testCases := []struct {
		name           string
		vm             asmodel.K8sNode
		appsThisVm     map[string]asmodel.Application
		appsOrder      []string
		solnWithVm     asmodel.Solution
		expectedResult asmodel.Solution
	}{
		{
			name: "case4",
			vm: asmodel.K8sNode{
				Name: "auto-sched-nokia6-3",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 20,
					Memory:  4096,
					Storage: 80,
				},
			},
			appsThisVm: appsForTest()[7],
			appsOrder:  appOrdersForTest()[0],
			solnWithVm: solnsForTest()[5],
			expectedResult: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 3,
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 5,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
				},
			},
		},
		{
			name: "case4-2",
			vm: asmodel.K8sNode{
				Name: "auto-sched-nokia6-3",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 13,
					Memory:  4096,
					Storage: 80,
				},
			},
			appsThisVm: appsForTest()[7],
			appsOrder:  appOrdersForTest()[0],
			solnWithVm: solnsForTest()[5],
			expectedResult: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 5,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
				},
			},
		},
		{
			name: "case4-3",
			vm: asmodel.K8sNode{
				Name: "auto-sched-nokia6-3",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 10,
					Memory:  4096,
					Storage: 80,
				},
			},
			appsThisVm: appsForTest()[7],
			appsOrder:  appOrdersForTest()[0],
			solnWithVm: solnsForTest()[5],
			expectedResult: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		actualResult := vmCpuWeightedAllocation(testCase.vm, testCase.appsThisVm, testCase.appsOrder, testCase.solnWithVm)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerAllocateCpusOneVm(t *testing.T) {
	testCases := []struct {
		name           string
		vm             asmodel.K8sNode
		apps           map[string]asmodel.Application
		appsOrder      []string
		appNamesThisVm []string
		solnWithVm     asmodel.Solution
		expectedResult asmodel.Solution
	}{
		{
			name: "case vm less than requested",
			vm: asmodel.K8sNode{
				Name: "auto-sched-nokia6-3",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 20,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps:           appsForTest()[8],
			appsOrder:      appOrdersForTest()[0],
			appNamesThisVm: []string{"app1", "app2", "app3", "app4", "app6", "app7", "app8"},
			solnWithVm:     solnsForTest()[5],
			expectedResult: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 3,
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 11,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
				},
			},
		},
		{
			name: "case vm less than requested 2",
			vm: asmodel.K8sNode{
				Name: "auto-sched-nokia6-3",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 10,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps:           appsForTest()[8],
			appsOrder:      appOrdersForTest()[0],
			appNamesThisVm: []string{"app1", "app2", "app3", "app4", "app6", "app7", "app8"},
			solnWithVm:     solnsForTest()[5],
			expectedResult: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
				},
			},
		},
		{
			name: "case vm less than requested 3",
			vm: asmodel.K8sNode{
				Name: "auto-sched-nokia6-3",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 15,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps:           appsForTest()[8],
			appsOrder:      appOrdersForTest()[0],
			appNamesThisVm: []string{"app1", "app2", "app3", "app4", "app6", "app7", "app8"},
			solnWithVm:     solnsForTest()[5],
			expectedResult: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 3,
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 6,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
				},
			},
		},
		{
			name: "case vm less than requested 4",
			vm: asmodel.K8sNode{
				Name: "auto-sched-nokia6-3",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 32,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps:           appsForTest()[8],
			appsOrder:      appOrdersForTest()[0],
			appNamesThisVm: []string{"app1", "app2", "app3", "app4", "app6", "app7", "app8"},
			solnWithVm:     solnsForTest()[5],
			expectedResult: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 3,
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 17,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 4,
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: cpuCoreStep,
					},
				},
			},
		},
		{
			name: "case vm more than requested",
			vm: asmodel.K8sNode{
				Name: "auto-sched-nokia6-3",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 33,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps:           appsForTest()[8],
			appsOrder:      appOrdersForTest()[0],
			appNamesThisVm: []string{"app1", "app2", "app3", "app4", "app6", "app7", "app8"},
			solnWithVm:     solnsForTest()[5],
			expectedResult: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 3,
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 20,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 4,
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
				},
			},
		},
		{
			name: "case vm more than requested 2",
			vm: asmodel.K8sNode{
				Name: "auto-sched-nokia6-3",
				ResidualResources: asmodel.GenericResources{
					CpuCore: 34,
					Memory:  4096,
					Storage: 80,
				},
			},
			apps:           appsForTest()[8],
			appsOrder:      appOrdersForTest()[0],
			appNamesThisVm: []string{"app1", "app2", "app3", "app4", "app6", "app7", "app8"},
			solnWithVm:     solnsForTest()[5],
			expectedResult: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 3,
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 2,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 20,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 4,
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 1,
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		actualResult := allocateCpusOneVm(testCase.vm, testCase.apps, testCase.appsOrder, testCase.appNamesThisVm, testCase.solnWithVm)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}
