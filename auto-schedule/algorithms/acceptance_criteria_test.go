package algorithms

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	asmodel "emcontroller/auto-schedule/model"
)

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
	testDelta := 0.0001

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

func TestInnerVmResMeetAllRestApps(t *testing.T) {
	cloud, apps, soln := cloudAppsSolnForIterTest()
	appsThisCloud := findAppsOneCloud(cloud, apps, soln) // need 0.4 CPU, 3214 Memory, 66 Storage in total.
	appsOrder := GenerateAppsOrder(apps)

	t.Log()
	t.Log("list the apps scheduled to this cloud in order")
	var appsThisCloudIter *appOneCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	var curAppName string = appsThisCloudIter.nextAppName()
	for len(curAppName) != 0 {
		t.Logf("curApp: %+v\n", appsThisCloud[curAppName])
		curAppName = appsThisCloudIter.nextAppName()
	}

	t.Log()
	t.Log("case 1")
	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()
	vm := asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 2,
			Memory:  10240,
			Storage: 100,
		},
	}
	assert.True(t, vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter, true))

	t.Log()
	t.Log("case 2")
	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.3,
			Memory:  10240,
			Storage: 100,
		},
	}
	assert.False(t, vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter, true))
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.3,
			Memory:  10240,
			Storage: 100,
		},
	}
	assert.True(t, vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter, true))
	t.Log("curAppName:", curAppName)

	t.Log()
	t.Log("case 3")
	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  3000,
			Storage: 100,
		},
	}
	assert.False(t, vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter, true))
	t.Log("curAppName:", curAppName)

	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.2,
			Memory:  3000,
			Storage: 100,
		},
	}
	assert.True(t, vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter, true))
	t.Log("curAppName:", curAppName)

	t.Log()
	t.Log("case 4: copy")
	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  3000,
			Storage: 100,
		},
	}
	assert.False(t, vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter, true))
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
	assert.True(t, vmResMeetAllRestApps(vm, apps, &curAppNameCopy1, iterCopy1, true))
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
	assert.False(t, vmResMeetAllRestApps(vm, apps, &curAppNameCopy2, iterCopy2, true))
	t.Log("curAppNameCopy2:", curAppNameCopy2)

	t.Log("copy 3, should be true")
	iterCopy3 := appsThisCloudIter.Copy()
	curAppNameCopy3 := curAppName
	vm = asmodel.K8sNode{
		Name: "vm",
		ResidualResources: asmodel.GenericResources{
			CpuCore: 0.6,
			Memory:  1000,
			Storage: 100,
		},
	}
	assert.True(t, vmResMeetAllRestApps(vm, apps, &curAppNameCopy3, iterCopy3, true))
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
	assert.False(t, vmResMeetAllRestApps(vm, apps, &curAppNameCopy4, iterCopy4, true))
	t.Log("curAppNameCopy4:", curAppNameCopy4)

	func() {
		t.Log()
		t.Log("case 5: copy without minCpu")
		appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
		curAppName = appsThisCloudIter.nextAppName()
		vm = asmodel.K8sNode{
			Name: "vm",
			ResidualResources: asmodel.GenericResources{
				CpuCore: 10.0,
				Memory:  30000,
				Storage: 1000,
			},
		}
		assert.False(t, vmResMeetAllRestApps(vm, apps, &curAppName, appsThisCloudIter, false))
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
		assert.True(t, vmResMeetAllRestApps(vm, apps, &curAppNameCopy1, iterCopy1, false))
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
		assert.False(t, vmResMeetAllRestApps(vm, apps, &curAppNameCopy2, iterCopy2, false))
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
		assert.True(t, vmResMeetAllRestApps(vm, apps, &curAppNameCopy3, iterCopy3, false))
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
		assert.False(t, vmResMeetAllRestApps(vm, apps, &curAppNameCopy4, iterCopy4, false))
		t.Log("curAppNameCopy4:", curAppNameCopy4)
	}()

}
