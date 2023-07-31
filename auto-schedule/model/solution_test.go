package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"emcontroller/models"
)

func TestChangeSolution(t *testing.T) {
	testCases := []struct {
		name string
		src  Solution
	}{
		{
			name: "case1",
			src: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app1": SingleAppSolution{
						Accepted: false,
					},
					"app2": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "cloud1",
					},
				},
				VmsToCreate: []models.IaasVm{
					models.IaasVm{
						Name:    "aaa",
						VCpu:    10,
						Ram:     4096,
						Storage: 10,
					},
					models.IaasVm{
						Name:    "bbb",
						VCpu:    10,
						Ram:     4096,
						Storage: 10,
					},
				},
			},
		},
		{
			name: "case2",
			src: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app2": SingleAppSolution{
						Accepted: false,
					},
					"app20": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "cloud2",
					},
				},
				VmsToCreate: []models.IaasVm{
					models.IaasVm{
						Name:    "aaa",
						VCpu:    10,
						Ram:     4096,
						Storage: 10,
					},
					models.IaasVm{
						Name:    "bbb",
						VCpu:    10,
						Ram:     4096,
						Storage: 10,
					},
				},
			},
		},
		{
			name: "case no VMs to create not nil",
			src: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app2": SingleAppSolution{
						Accepted: false,
					},
					"app21": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "cloud3",
					},
				},
				VmsToCreate: []models.IaasVm{},
			},
		},
		{
			name: "case no VMs to create nil",
			src: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app2": SingleAppSolution{
						Accepted: false,
					},
					"app21": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "cloud3",
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		var dst Solution

		t.Log("test AppsSolution extra key")
		dst = SolutionCopy(testCase.src)
		assert.Equal(t, testCase.src, dst)

		extraKey := "app100"
		dst.AppsSolution[extraKey] = SingleAppSolution{
			Accepted: false,
		}
		if app, exist := testCase.src.AppsSolution[extraKey]; exist {
			t.Errorf("key \"%s\" should not exist in testCase.src, the value is: %+v", extraKey, app)
		} else {
			t.Logf("As expected, key \"%s\" does not exist in testCase.src", extraKey)
		}

		t.Log("test AppsSolution change")
		dst = SolutionCopy(testCase.src)
		assert.Equal(t, testCase.src, dst)

		changeKey := "app2"
		dst.AppsSolution[changeKey] = SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "cloud20",
		}

		t.Logf("testCase.src[changeKey]:\n%+v\ndst[changeKey]:\n%+v", testCase.src.AppsSolution[changeKey], dst.AppsSolution[changeKey])
		assert.NotEqual(t, testCase.src.AppsSolution[changeKey], dst.AppsSolution[changeKey], fmt.Sprintf("testCase.src[changeKey] [%+v] and dst[changeKey] [%+v], show not be equal.", testCase.src.AppsSolution[changeKey], dst.AppsSolution[changeKey]))

		if len(testCase.src.VmsToCreate) > 0 {
			t.Log("test VmsToCreate change")
			dst.VmsToCreate[0] = models.IaasVm{
				Name:    "ddd",
				VCpu:    100,
				Ram:     4096,
				Storage: 10,
			}
			assert.NotEqual(t, dst.VmsToCreate[0], testCase.src.VmsToCreate[0])
		}

		t.Log("test VmsToCreate append")
		dst.VmsToCreate = append(dst.VmsToCreate, models.IaasVm{
			Name:    "ccc",
			VCpu:    12,
			Ram:     5000,
			Storage: 100,
		})
		assert.NotEqual(t, len(dst.VmsToCreate), len(testCase.src.VmsToCreate))

	}
}

func TestAbsorb(t *testing.T) {
	testCases := []struct {
		name           string
		absorber       Solution
		absorbates     []Solution
		expectedResult Solution
	}{
		{
			name: "same cloud",
			absorber: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app1": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app2": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-1",
					},
					"app3": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app4": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app5": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-2",
					},
					"app6": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-2",
					},
					"app7": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
				},
				VmsToCreate: []models.IaasVm{
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
			},
			absorbates: []Solution{
				Solution{
					AppsSolution: map[string]SingleAppSolution{
						"app1": SingleAppSolution{
							Accepted:        true,
							TargetCloudName: "NOKIA4",
							K8sNodeName:     "ori-node1",
						},
						"app3": SingleAppSolution{
							Accepted:        true,
							TargetCloudName: "NOKIA4",
							K8sNodeName:     "ori-node2",
						},
						"app4": SingleAppSolution{
							Accepted:        true,
							TargetCloudName: "NOKIA4",
							K8sNodeName:     "auto-sched-nokia4-4",
						},
					},
					VmsToCreate: []models.IaasVm{
						models.IaasVm{
							Name:    "auto-sched-nokia4-4",
							Cloud:   "NOKIA4",
							VCpu:    8,
							Ram:     3122,
							Storage: 44,
						},
					},
				},
			},
			expectedResult: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app1": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "ori-node1",
					},
					"app2": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-1",
					},
					"app3": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "ori-node2",
					},
					"app4": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-4",
					},
					"app5": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-2",
					},
					"app6": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-2",
					},
					"app7": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
				},
				VmsToCreate: []models.IaasVm{
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
					models.IaasVm{
						Name:    "auto-sched-nokia4-4",
						Cloud:   "NOKIA4",
						VCpu:    8,
						Ram:     3122,
						Storage: 44,
					},
				},
			},
		},
		{
			name: "different clouds",
			absorber: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app1": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app2": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA5",
					},
					"app3": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app4": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
					},
					"app5": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA8",
					},
					"app6": SingleAppSolution{
						Accepted: false,
					},
					"app7": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
					},
				},
				VmsToCreate: []models.IaasVm{},
			},
			absorbates: []Solution{
				Solution{
					AppsSolution: map[string]SingleAppSolution{
						"app1": SingleAppSolution{
							Accepted:        true,
							TargetCloudName: "NOKIA4",
							K8sNodeName:     "nokia4-ori-node1",
						},
						"app3": SingleAppSolution{
							Accepted:        true,
							TargetCloudName: "NOKIA4",
							K8sNodeName:     "auto-sched-nokia4-4",
						},
					},
					VmsToCreate: []models.IaasVm{
						models.IaasVm{
							Name:    "auto-sched-nokia4-4",
							Cloud:   "NOKIA4",
							VCpu:    8,
							Ram:     3122,
							Storage: 44,
						},
					},
				},
				Solution{
					AppsSolution: map[string]SingleAppSolution{
						"app2": SingleAppSolution{
							Accepted:        true,
							TargetCloudName: "NOKIA5",
							K8sNodeName:     "nokia5-ori-node1",
						},
					},
					VmsToCreate: []models.IaasVm{},
				},
				Solution{
					AppsSolution: map[string]SingleAppSolution{
						"app4": SingleAppSolution{
							Accepted:        true,
							TargetCloudName: "NOKIA6",
							K8sNodeName:     "nokia6-ori-node1",
						},
						"app7": SingleAppSolution{
							Accepted:        true,
							TargetCloudName: "NOKIA6",
							K8sNodeName:     "auto-sched-nokia6-3",
						},
					},
					VmsToCreate: []models.IaasVm{
						models.IaasVm{
							Name:    "auto-sched-nokia6-3",
							Cloud:   "NOKIA6",
							VCpu:    10,
							Ram:     2567,
							Storage: 80,
						},
					},
				},
				Solution{
					AppsSolution: map[string]SingleAppSolution{
						"app5": SingleAppSolution{
							Accepted:        true,
							TargetCloudName: "NOKIA8",
							K8sNodeName:     "auto-sched-nokia8-0",
						},
					},
					VmsToCreate: []models.IaasVm{
						models.IaasVm{
							Name:    "auto-sched-nokia8-0",
							Cloud:   "NOKIA8",
							VCpu:    7,
							Ram:     4000,
							Storage: 100,
						},
					},
				},
			},
			expectedResult: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app1": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "nokia4-ori-node1",
					},
					"app2": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA5",
						K8sNodeName:     "nokia5-ori-node1",
					},
					"app3": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-4",
					},
					"app4": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "nokia6-ori-node1",
					},
					"app5": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA8",
						K8sNodeName:     "auto-sched-nokia8-0",
					},
					"app6": SingleAppSolution{
						Accepted: false,
					},
					"app7": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "auto-sched-nokia6-3",
					},
				},
				VmsToCreate: []models.IaasVm{
					models.IaasVm{
						Name:    "auto-sched-nokia4-4",
						Cloud:   "NOKIA4",
						VCpu:    8,
						Ram:     3122,
						Storage: 44,
					},
					models.IaasVm{
						Name:    "auto-sched-nokia6-3",
						Cloud:   "NOKIA6",
						VCpu:    10,
						Ram:     2567,
						Storage: 80,
					},
					models.IaasVm{
						Name:    "auto-sched-nokia8-0",
						Cloud:   "NOKIA8",
						VCpu:    7,
						Ram:     4000,
						Storage: 100,
					},
				},
			},
		},
		{
			name: "different clouds only allocate CPU",
			absorber: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app1": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "nokia4-ori-node1",
					},
					"app2": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA5",
						K8sNodeName:     "nokia5-ori-node1",
					},
					"app3": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-4",
					},
					"app4": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "nokia6-ori-node1",
					},
					"app5": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA8",
						K8sNodeName:     "auto-sched-nokia8-0",
					},
					"app6": SingleAppSolution{
						Accepted: false,
					},
					"app7": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "auto-sched-nokia6-3",
					},
				},
				VmsToCreate: []models.IaasVm{
					models.IaasVm{
						Name:    "auto-sched-nokia4-4",
						Cloud:   "NOKIA4",
						VCpu:    8,
						Ram:     3122,
						Storage: 44,
					},
					models.IaasVm{
						Name:    "auto-sched-nokia6-3",
						Cloud:   "NOKIA6",
						VCpu:    10,
						Ram:     2567,
						Storage: 80,
					},
					models.IaasVm{
						Name:    "auto-sched-nokia8-0",
						Cloud:   "NOKIA8",
						VCpu:    7,
						Ram:     4000,
						Storage: 100,
					},
				},
			},
			absorbates: []Solution{
				Solution{
					AppsSolution: map[string]SingleAppSolution{
						"app1": SingleAppSolution{
							Accepted:         true,
							TargetCloudName:  "NOKIA4",
							K8sNodeName:      "nokia4-ori-node1",
							AllocatedCpuCore: 0.6,
						},
						"app3": SingleAppSolution{
							Accepted:         true,
							TargetCloudName:  "NOKIA4",
							K8sNodeName:      "auto-sched-nokia4-4",
							AllocatedCpuCore: 2.3,
						},
					},
					VmsToCreate: []models.IaasVm{},
				},
				Solution{
					AppsSolution: map[string]SingleAppSolution{
						"app2": SingleAppSolution{
							Accepted:         true,
							TargetCloudName:  "NOKIA5",
							K8sNodeName:      "nokia5-ori-node1",
							AllocatedCpuCore: 0.4,
						},
					},
					VmsToCreate: []models.IaasVm{},
				},
				Solution{
					AppsSolution: map[string]SingleAppSolution{
						"app4": SingleAppSolution{
							Accepted:         true,
							TargetCloudName:  "NOKIA6",
							K8sNodeName:      "nokia6-ori-node1",
							AllocatedCpuCore: 2.4,
						},
						"app7": SingleAppSolution{
							Accepted:         true,
							TargetCloudName:  "NOKIA6",
							K8sNodeName:      "auto-sched-nokia6-3",
							AllocatedCpuCore: 3.4,
						},
					},
					VmsToCreate: []models.IaasVm{},
				},
				Solution{
					AppsSolution: map[string]SingleAppSolution{
						"app5": SingleAppSolution{
							Accepted:         true,
							TargetCloudName:  "NOKIA8",
							K8sNodeName:      "auto-sched-nokia8-0",
							AllocatedCpuCore: 1.7,
						},
					},
					VmsToCreate: []models.IaasVm{},
				},
			},
			expectedResult: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app1": SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA4",
						K8sNodeName:      "nokia4-ori-node1",
						AllocatedCpuCore: 0.6,
					},
					"app2": SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA5",
						K8sNodeName:      "nokia5-ori-node1",
						AllocatedCpuCore: 0.4,
					},
					"app3": SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA4",
						K8sNodeName:      "auto-sched-nokia4-4",
						AllocatedCpuCore: 2.3,
					},
					"app4": SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "nokia6-ori-node1",
						AllocatedCpuCore: 2.4,
					},
					"app5": SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA8",
						K8sNodeName:      "auto-sched-nokia8-0",
						AllocatedCpuCore: 1.7,
					},
					"app6": SingleAppSolution{
						Accepted: false,
					},
					"app7": SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA6",
						K8sNodeName:      "auto-sched-nokia6-3",
						AllocatedCpuCore: 3.4,
					},
				},
				VmsToCreate: []models.IaasVm{
					models.IaasVm{
						Name:    "auto-sched-nokia4-4",
						Cloud:   "NOKIA4",
						VCpu:    8,
						Ram:     3122,
						Storage: 44,
					},
					models.IaasVm{
						Name:    "auto-sched-nokia6-3",
						Cloud:   "NOKIA6",
						VCpu:    10,
						Ram:     2567,
						Storage: 80,
					},
					models.IaasVm{
						Name:    "auto-sched-nokia8-0",
						Cloud:   "NOKIA8",
						VCpu:    7,
						Ram:     4000,
						Storage: 100,
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		for _, absorbate := range testCase.absorbates {
			testCase.absorber.Absorb(absorbate)
		}
		assert.Equal(t, testCase.expectedResult, testCase.absorber, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestSolutionCopy(t *testing.T) {
	testCases := []struct {
		name string
		src  Solution
	}{
		{
			name: "case with VMs to create",
			src: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app1": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app2": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-1",
					},
					"app3": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app4": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app5": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-2",
					},
					"app6": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-2",
					},
					"app7": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
				},
				VmsToCreate: []models.IaasVm{
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
			},
		},
		{
			name: "case empty VMs to create",
			src: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app1": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app2": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-1",
					},
					"app3": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app4": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app5": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-2",
					},
					"app6": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-2",
					},
					"app7": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
				},
				VmsToCreate: []models.IaasVm{},
			},
		},
		{
			name: "case without VMs to create",
			src: Solution{
				AppsSolution: map[string]SingleAppSolution{
					"app1": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app2": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-1",
					},
					"app3": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app4": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app5": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-2",
					},
					"app6": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-2",
					},
					"app7": SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		dst := SolutionCopy(testCase.src)
		assert.Equal(t, testCase.src, dst)

		t.Log("test add AppsSolution")
		dst.AppsSolution["app9"] = SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "NOKIA9",
			K8sNodeName:     "auto-sched-nokia9-3",
		}
		assert.NotEqual(t, testCase.src, dst)
		_, exist := testCase.src.AppsSolution["app9"]
		assert.False(t, exist)

		t.Log("test append VmsToCreate")
		dst.VmsToCreate = append(dst.VmsToCreate, models.IaasVm{
			Name:    "auto-sched-nokia9-3",
			Cloud:   "NOKIA9",
			VCpu:    12,
			Ram:     2888,
			Storage: 55,
		})
		assert.NotEqual(t, len(testCase.src.VmsToCreate), len(dst.VmsToCreate))

		t.Log("test change AppsSolution")
		dst.AppsSolution["app7"] = SingleAppSolution{
			Accepted: false,
		}
		assert.NotEqual(t, testCase.src, dst)
		assert.NotEqual(t, testCase.src.AppsSolution["app7"], dst.AppsSolution["app7"])

	}
}
