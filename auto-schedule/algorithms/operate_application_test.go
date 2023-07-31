package algorithms

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	asmodel "emcontroller/auto-schedule/model"
)

func TestInnerFindAppsOneCloud(t *testing.T) {
	testCases := []struct {
		name           string
		cloud          asmodel.Cloud
		apps           map[string]asmodel.Application
		soln           asmodel.Solution
		expectedResult map[string]asmodel.Application
	}{
		{
			name: "caseUnacceptedWithCloud",
			cloud: asmodel.Cloud{
				Name: "cloud1",
			},
			apps: map[string]asmodel.Application{
				"app1": asmodel.Application{
					Name: "app1",
				},
				"app2": asmodel.Application{
					Name: "app2",
				},
				"app3": asmodel.Application{
					Name: "app3",
				},
				"app4": asmodel.Application{
					Name: "app4",
				},
			},
			soln: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "cloud1",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "cloud3",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:        false,
						TargetCloudName: "cloud1",
					},
				},
			},
			expectedResult: map[string]asmodel.Application{
				"app1": asmodel.Application{
					Name: "app1",
				},
			},
		},
		{
			name: "caseTwoResults",
			cloud: asmodel.Cloud{
				Name: "cloud1",
			},
			apps: map[string]asmodel.Application{
				"app1": asmodel.Application{
					Name: "app1",
				},
				"app2": asmodel.Application{
					Name: "app2",
				},
				"app3": asmodel.Application{
					Name: "app3",
				},
				"app4": asmodel.Application{
					Name: "app4",
				},
			},
			soln: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "cloud1",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "cloud1",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "cloud12",
					},
				},
			},
			expectedResult: map[string]asmodel.Application{
				"app2": asmodel.Application{
					Name: "app2",
				},
				"app1": asmodel.Application{
					Name: "app1",
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := findAppsOneCloud(testCase.cloud, testCase.apps, testCase.soln)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerFilterMaxPriApps(t *testing.T) {
	_, testApps, _ := cloudAppsSolnForTest()

	testCases := []struct {
		name           string
		apps           map[string]asmodel.Application
		expectedResult map[string]asmodel.Application
	}{
		{
			name: "case 1",
			apps: testApps,
			expectedResult: map[string]asmodel.Application{
				"app2": asmodel.Application{
					Name:     "app2",
					Priority: 10,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 3.9,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app7": asmodel.Application{
					Name:     "app7",
					Priority: 10,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 3.0,
							Memory:  540,
							Storage: 35,
						},
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := filterMaxPriApps(testCase.apps)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerFilterOutMaxPriApps(t *testing.T) {
	_, testApps, _ := cloudAppsSolnForTest()

	testCases := []struct {
		name           string
		apps           map[string]asmodel.Application
		expectedResult map[string]asmodel.Application
	}{
		{
			name: "case 1",
			apps: testApps,
			expectedResult: map[string]asmodel.Application{
				"app1": asmodel.Application{
					Name:     "app1",
					Priority: 5,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 2.4,
							Memory:  1024,
							Storage: 10,
						},
					},
				},
				"app3": asmodel.Application{
					Name:     "app3",
					Priority: 1,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 1.4,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app4": asmodel.Application{
					Name:     "app4",
					Priority: 3,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 5.0,
							Memory:  660,
							Storage: 6,
						},
					},
				},
				"app5": asmodel.Application{
					Name:     "app5",
					Priority: 2,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 5.0,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app6": asmodel.Application{
					Name:     "app6",
					Priority: 7,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 4.0,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app8": asmodel.Application{
					Name:     "app8",
					Priority: 8,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 2.0,
							Memory:  540,
							Storage: 15,
						},
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := filterOutMaxPriApps(testCase.apps)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerFilterAppsByNames(t *testing.T) {
	_, testApps, _ := cloudAppsSolnForTest()

	testCases := []struct {
		name           string
		tgtAppNames    []string
		apps           map[string]asmodel.Application
		expectedResult map[string]asmodel.Application
	}{
		{
			name:        "case 1",
			tgtAppNames: []string{"app2", "app7"},
			apps:        testApps,
			expectedResult: map[string]asmodel.Application{
				"app2": asmodel.Application{
					Name:     "app2",
					Priority: 10,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 3.9,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app7": asmodel.Application{
					Name:     "app7",
					Priority: 10,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 3.0,
							Memory:  540,
							Storage: 35,
						},
					},
				},
			},
		},
		{
			name:        "case 2",
			tgtAppNames: []string{"app1", "app3", "app5", "app7"},
			apps:        testApps,
			expectedResult: map[string]asmodel.Application{
				"app1": asmodel.Application{
					Name:     "app1",
					Priority: 5,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 2.4,
							Memory:  1024,
							Storage: 10,
						},
					},
				},
				"app3": asmodel.Application{
					Name:     "app3",
					Priority: 1,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 1.4,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app5": asmodel.Application{
					Name:     "app5",
					Priority: 2,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 5.0,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app7": asmodel.Application{
					Name:     "app7",
					Priority: 10,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 3.0,
							Memory:  540,
							Storage: 35,
						},
					},
				},
			},
		},
		{
			name:        "case 3",
			tgtAppNames: []string{"app1", "app2", "app3", "app6", "app8"},
			apps:        testApps,
			expectedResult: map[string]asmodel.Application{
				"app1": asmodel.Application{
					Name:     "app1",
					Priority: 5,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 2.4,
							Memory:  1024,
							Storage: 10,
						},
					},
				},
				"app2": asmodel.Application{
					Name:     "app2",
					Priority: 10,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 3.9,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app3": asmodel.Application{
					Name:     "app3",
					Priority: 1,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 1.4,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app6": asmodel.Application{
					Name:     "app6",
					Priority: 7,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 4.0,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app8": asmodel.Application{
					Name:     "app8",
					Priority: 8,
					Resources: asmodel.AppResources{
						GenericResources: asmodel.GenericResources{
							CpuCore: 2.0,
							Memory:  540,
							Storage: 15,
						},
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := filterAppsByNames(testCase.tgtAppNames, testCase.apps)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerCalcNeededRes(t *testing.T) {
	_, testApps, _ := cloudAppsSolnForTest()

	testCases := []struct {
		name           string
		apps           map[string]asmodel.Application
		appNames       []string
		minCpu         bool
		expectedResult asmodel.AppResources
	}{
		{
			name:     "case1 not minCpu",
			apps:     testApps,
			appNames: []string{"app1", "app2"},
			minCpu:   false,
			expectedResult: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 6.3,
					Memory:  2014,
					Storage: 25,
				},
			},
		},
		{
			name:     "case1 minCpu",
			apps:     testApps,
			appNames: []string{"app1", "app2"},
			minCpu:   true,
			expectedResult: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 0.2,
					Memory:  2014,
					Storage: 25,
				},
			},
		},
		{
			name:     "case2 not minCpu",
			apps:     testApps,
			appNames: []string{"app3", "app8", "app7", "app5"},
			minCpu:   false,
			expectedResult: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 11.4,
					Memory:  3060,
					Storage: 80,
				},
			},
		},
		{
			name:     "case2 minCpu",
			apps:     testApps,
			appNames: []string{"app3", "app8", "app7", "app5"},
			minCpu:   true,
			expectedResult: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 0.4,
					Memory:  3060,
					Storage: 80,
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := calcNeededRes(testCase.apps, testCase.appNames, testCase.minCpu)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}
