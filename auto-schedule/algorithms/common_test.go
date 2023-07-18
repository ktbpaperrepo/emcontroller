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

func cloudAppsSolnForIterTest() (asmodel.Cloud, map[string]asmodel.Application, asmodel.Solution) {
	cloud := asmodel.Cloud{
		Name: "cloud1",
	}
	apps := map[string]asmodel.Application{
		"app1": asmodel.Application{
			Name: "app1",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 5.0,
					Memory:  1024,
					Storage: 10,
				},
			},
		},
		"app2": asmodel.Application{
			Name: "app2",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 5.0,
					Memory:  990,
					Storage: 15,
				},
			},
		},
		"app3": asmodel.Application{
			Name: "app3",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 5.0,
					Memory:  990,
					Storage: 15,
				},
			},
		},
		"app4": asmodel.Application{
			Name: "app4",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 5.0,
					Memory:  660,
					Storage: 6,
				},
			},
		},
		"app5": asmodel.Application{
			Name: "app5",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 5.0,
					Memory:  990,
					Storage: 15,
				},
			},
		},
		"app6": asmodel.Application{
			Name: "app6",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 5.0,
					Memory:  990,
					Storage: 15,
				},
			},
		},
		"app7": asmodel.Application{
			Name: "app7",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 5.0,
					Memory:  540,
					Storage: 35,
				},
			},
		},
		"app8": asmodel.Application{
			Name: "app8",
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 5.0,
					Memory:  540,
					Storage: 15,
				},
			},
		},
	}
	soln := asmodel.Solution{
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
			TargetCloudName: "cloud1",
		},
		"app5": asmodel.SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "cloud2",
		},
		"app6": asmodel.SingleAppSolution{
			Accepted: false,
		},
		"app7": asmodel.SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "cloud1",
		},
		"app8": asmodel.SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "cloud3",
		},
	}
	return cloud, apps, soln
}

func TestInnerAppOneCloudIter(t *testing.T) {
	cloud, apps, soln := cloudAppsSolnForIterTest()

	appsThisCloud := findAppsOneCloud(cloud, apps, soln)
	t.Logf("appsThisCloud:\n%+v\n", appsThisCloud)

	appsOrder := GenerateAppsOrder(apps)
	t.Logf("appsOrder:\n%+v\n", appsOrder)

	var appsThisCloudIter *appOneCloudIter
	var curAppName string

	t.Log()
	t.Log("The following 3 orders of curAppName should be the same.")
	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	for {
		curAppName = appsThisCloudIter.nextAppName()
		t.Logf("curAppName: %s\n", curAppName)
		if len(curAppName) == 0 {
			break
		}
	}

	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	for {
		curAppName = appsThisCloudIter.nextAppName()
		t.Logf("curAppName: %s\n", curAppName)
		if len(curAppName) == 0 {
			break
		}
	}

	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	for {
		curAppName = appsThisCloudIter.nextAppName()
		t.Logf("curAppName: %s\n", curAppName)
		if len(curAppName) == 0 {
			break
		}
	}

	t.Log()
	t.Log("Then, we test copying an iterator.")
	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()
	t.Logf("curAppName: %s\n", curAppName)
	curAppName = appsThisCloudIter.nextAppName()
	t.Logf("curAppName: %s\n", curAppName)

	t.Log("Now, we copy an iterator")
	iterCopy := appsThisCloudIter.Copy()
	curAppNameCopy := curAppName

	t.Log("Now, we run the original iterator.")
	for {
		curAppName = appsThisCloudIter.nextAppName()
		t.Logf("curAppName: %s\n", curAppName)
		if len(curAppName) == 0 {
			break
		}
	}

	t.Log("Now, we run the copied iterator, and it should work from the point where we copied.")
	for {
		curAppNameCopy = iterCopy.nextAppName()
		t.Logf("curAppNameCopy: %s\n", curAppNameCopy)
		if len(curAppNameCopy) == 0 {
			break
		}
	}

	t.Log()
	t.Log("Test the scenario in which no applications are scheduled to this cloud.")
	solnNon := asmodel.Solution{
		"app1": asmodel.SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "cloud2",
		},
		"app2": asmodel.SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "cloud2",
		},
		"app3": asmodel.SingleAppSolution{
			Accepted: false,
		},
		"app4": asmodel.SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "cloud2",
		},
		"app5": asmodel.SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "cloud2",
		},
		"app6": asmodel.SingleAppSolution{
			Accepted: false,
		},
		"app7": asmodel.SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "cloud3",
		},
		"app8": asmodel.SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "cloud3",
		},
	}
	appsThisCloudNon := findAppsOneCloud(cloud, apps, solnNon)
	appsThisCloudIter = newAppOneCloudIter(appsThisCloudNon, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()
	assert.Equal(t, 0, len(curAppName))

}
