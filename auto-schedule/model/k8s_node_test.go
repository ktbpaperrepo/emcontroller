package model

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"emcontroller/models"
)

func TestGenK8sNodeFromPods(t *testing.T) {
	testCases := []struct {
		name           string
		vm             models.IaasVm
		podsOnNode     []apiv1.Pod
		expectedResult K8sNode
	}{
		{
			name: "case 2 pods 1 has resources",
			vm: models.IaasVm{
				Name:    "n8test",
				Cloud:   "NOKIA8",
				VCpu:    4,
				Ram:     8192,
				Storage: 200,
			},
			podsOnNode: []apiv1.Pod{
				apiv1.Pod{
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Resources: apiv1.ResourceRequirements{},
							},
						},
					},
				},
				apiv1.Pod{
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Resources: apiv1.ResourceRequirements{
									Limits: map[apiv1.ResourceName]resource.Quantity{
										apiv1.ResourceCPU:              resource.MustParse("100m"),
										apiv1.ResourceMemory:           resource.MustParse("500Mi"),
										apiv1.ResourceEphemeralStorage: resource.MustParse("10Gi"),
									},
									Requests: map[apiv1.ResourceName]resource.Quantity{
										apiv1.ResourceCPU:              resource.MustParse("100m"),
										apiv1.ResourceMemory:           resource.MustParse("500Mi"),
										apiv1.ResourceEphemeralStorage: resource.MustParse("10Gi"),
									},
								},
							},
						},
					},
				},
			},
			expectedResult: K8sNode{
				Name: "n8test",
				ResidualResources: GenericResources{
					CpuCore: math.Max(models.CalcVmAvailVcpu(4)-0.1, 0),
					Memory:  math.Max(models.CalcVmAvailRamMiB(8192)-500, 0),
					Storage: math.Max(models.CalcVmAvailStorGiB(200)-10, 0),
				},
			},
		},
		{
			name: "case 3 pods 2 have resources",
			vm: models.IaasVm{
				Name:    "n8test",
				Cloud:   "NOKIA8",
				VCpu:    4,
				Ram:     8192,
				Storage: 200,
			},
			podsOnNode: []apiv1.Pod{
				apiv1.Pod{
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Resources: apiv1.ResourceRequirements{},
							},
						},
					},
				},
				apiv1.Pod{
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Resources: apiv1.ResourceRequirements{
									Limits: map[apiv1.ResourceName]resource.Quantity{
										apiv1.ResourceCPU:              resource.MustParse("100m"),
										apiv1.ResourceMemory:           resource.MustParse("500Mi"),
										apiv1.ResourceEphemeralStorage: resource.MustParse("10Gi"),
									},
									Requests: map[apiv1.ResourceName]resource.Quantity{
										apiv1.ResourceCPU:              resource.MustParse("100m"),
										apiv1.ResourceMemory:           resource.MustParse("500Mi"),
										apiv1.ResourceEphemeralStorage: resource.MustParse("10Gi"),
									},
								},
							},
						},
					},
				},
				apiv1.Pod{
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Resources: apiv1.ResourceRequirements{
									Limits: map[apiv1.ResourceName]resource.Quantity{
										apiv1.ResourceCPU:              resource.MustParse("1.2"),
										apiv1.ResourceMemory:           resource.MustParse("666Mi"),
										apiv1.ResourceEphemeralStorage: resource.MustParse("7Gi"),
									},
									Requests: map[apiv1.ResourceName]resource.Quantity{
										apiv1.ResourceCPU:              resource.MustParse("1.2"),
										apiv1.ResourceMemory:           resource.MustParse("666Mi"),
										apiv1.ResourceEphemeralStorage: resource.MustParse("7Gi"),
									},
								},
							},
						},
					},
				},
			},
			expectedResult: K8sNode{
				Name: "n8test",
				ResidualResources: GenericResources{
					CpuCore: math.Max(models.CalcVmAvailVcpu(4)-0.1-1.2, 0),
					Memory:  math.Max(models.CalcVmAvailRamMiB(8192)-500-666, 0),
					Storage: math.Max(models.CalcVmAvailStorGiB(200)-10-7, 0),
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := GenK8sNodeFromPods(testCase.vm, testCase.podsOnNode)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}

func TestGenK8sNodeFromApps(t *testing.T) {
	testCases := []struct {
		name           string
		vm             models.IaasVm
		apps           map[string]Application
		appGroup       []string
		expectedResult K8sNode
	}{
		{
			name: "case 2 apps",
			vm: models.IaasVm{
				Name:    "n8test",
				Cloud:   "NOKIA8",
				VCpu:    20,
				Ram:     8192,
				Storage: 230,
			},
			apps: map[string]Application{
				"app1": Application{
					Name:     "app1",
					Priority: 5,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 2.4,
							Memory:  1024,
							Storage: 10,
						},
					},
				},
				"app2": Application{
					Name:     "app2",
					Priority: 10,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 3.9,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app3": Application{
					Name:     "app3",
					Priority: 1,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 1.4,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app4": Application{
					Name:     "app4",
					Priority: 3,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 5.0,
							Memory:  660,
							Storage: 6,
						},
					},
				},
				"app5": Application{
					Name:     "app5",
					Priority: 2,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 5.0,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app6": Application{
					Name:     "app6",
					Priority: 7,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 4.0,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app7": Application{
					Name:     "app7",
					Priority: 10,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 3.0,
							Memory:  540,
							Storage: 35,
						},
					},
				},
				"app8": Application{
					Name:     "app8",
					Priority: 8,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 2.0,
							Memory:  540,
							Storage: 15,
						},
					},
				},
			},
			appGroup: []string{"app1", "app2"},
			expectedResult: K8sNode{
				Name: "n8test",
				ResidualResources: GenericResources{
					CpuCore: math.Max(models.CalcVmAvailVcpu(20)-3.9-2.4, 0),
					Memory:  math.Max(models.CalcVmAvailRamMiB(8192)-990-1024, 0),
					Storage: math.Max(models.CalcVmAvailStorGiB(230)-10-15, 0),
				},
			},
		},
		{
			name: "case 3 apps",
			vm: models.IaasVm{
				Name:    "n8test",
				Cloud:   "NOKIA8",
				VCpu:    20,
				Ram:     8192,
				Storage: 230,
			},
			apps: map[string]Application{
				"app1": Application{
					Name:     "app1",
					Priority: 5,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 2.4,
							Memory:  1024,
							Storage: 10,
						},
					},
				},
				"app2": Application{
					Name:     "app2",
					Priority: 10,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 3.9,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app3": Application{
					Name:     "app3",
					Priority: 1,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 1.4,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app4": Application{
					Name:     "app4",
					Priority: 3,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 5.0,
							Memory:  660,
							Storage: 6,
						},
					},
				},
				"app5": Application{
					Name:     "app5",
					Priority: 2,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 5.0,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app6": Application{
					Name:     "app6",
					Priority: 7,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 4.0,
							Memory:  990,
							Storage: 15,
						},
					},
				},
				"app7": Application{
					Name:     "app7",
					Priority: 10,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 3.0,
							Memory:  540,
							Storage: 35,
						},
					},
				},
				"app8": Application{
					Name:     "app8",
					Priority: 8,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 2.0,
							Memory:  540,
							Storage: 15,
						},
					},
				},
			},
			appGroup: []string{"app1", "app2", "app8"},
			expectedResult: K8sNode{
				Name: "n8test",
				ResidualResources: GenericResources{
					CpuCore: math.Max(models.CalcVmAvailVcpu(20)-3.9-2.4-2.0, 0),
					Memory:  math.Max(models.CalcVmAvailRamMiB(8192)-990-1024-540, 0),
					Storage: math.Max(models.CalcVmAvailStorGiB(230)-10-15-15, 0),
				},
			},
		},
	}

	testDelta := 0.0001

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := GenK8sNodeFromApps(testCase.vm, testCase.apps, testCase.appGroup)
		assert.Equal(t, testCase.expectedResult.Name, actualResult.Name, fmt.Sprintf("%s: result is not expected", testCase.name))
		assert.InDelta(t, testCase.expectedResult.ResidualResources.CpuCore, actualResult.ResidualResources.CpuCore, testDelta, fmt.Sprintf("%s: result is not expected", testCase.name))
		assert.InDelta(t, testCase.expectedResult.ResidualResources.Memory, actualResult.ResidualResources.Memory, testDelta, fmt.Sprintf("%s: result is not expected", testCase.name))
		assert.InDelta(t, testCase.expectedResult.ResidualResources.Storage, actualResult.ResidualResources.Storage, testDelta, fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}

func TestK8sNodeCopy(t *testing.T) {
	testCases := []struct {
		name string
		src  K8sNode
	}{
		{
			name: "case 1",
			src: K8sNode{
				Name: "node1",
				ResidualResources: GenericResources{
					CpuCore: 16,
					Memory:  10240,
					Storage: 100,
				},
			},
		},
		{
			name: "case 2",
			src: K8sNode{
				Name: "node2",
				ResidualResources: GenericResources{
					CpuCore: 4,
					Memory:  4096,
					Storage: 50,
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		dst := K8sNodeCopy(testCase.src)
		assert.Equal(t, testCase.src, dst)
	}
}
