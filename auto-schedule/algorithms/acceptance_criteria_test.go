package algorithms

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

func TestInnerDepAcc(t *testing.T) {
	testCases := []struct {
		name           string
		clouds         map[string]asmodel.Cloud
		apps           map[string]asmodel.Application
		soln           asmodel.Solution
		expectedResult bool
	}{
		{
			name:   "case should be same clouds",
			clouds: cloudsWithNetForTest()[0],
			apps:   appsForTest()[9],
			soln: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app5": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
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
			expectedResult: true,
		},
		{
			name:   "case should be same clouds but not",
			clouds: cloudsWithNetForTest()[0],
			apps:   appsForTest()[9],
			soln: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "auto-sched-nokia6-3",
					},
					"app5": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
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
			expectedResult: false,
		},
		{
			name:   "case inside cloud unreachable",
			clouds: cloudsWithNetForTest()[1],
			apps:   appsForTest()[9],
			soln: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-1",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app5": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-2",
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
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
			expectedResult: false,
		},
		{
			name:   "case inside cloud unreachable but schedule to same VM",
			clouds: cloudsWithNetForTest()[1],
			apps:   appsForTest()[9],
			soln: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
						K8sNodeName:     "auto-sched-nokia4-3",
					},
					"app5": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
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
			expectedResult: true,
		},
		{
			name:   "case schedule to 2 reachable clouds",
			clouds: cloudsWithNetForTest()[2],
			apps:   appsForTest()[9],
			soln: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "auto-sched-nokia6-3",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "auto-sched-nokia6-2",
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-2",
					},
					"app5": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-1",
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "auto-sched-nokia6-2",
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "auto-sched-nokia6-4",
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
			expectedResult: true,
		},
		{
			name:   "case dependency not accepted",
			clouds: cloudsWithNetForTest()[2],
			apps:   appsForTest()[9],
			soln: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "auto-sched-nokia6-3",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "auto-sched-nokia6-2",
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-2",
					},
					"app5": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-1",
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
						K8sNodeName:     "auto-sched-nokia6-2",
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
						K8sNodeName:     "auto-sched-nokia7-3",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted: false,
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
			expectedResult: false,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		actualResult := depAcc(testCase.clouds, testCase.apps, testCase.soln)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}
