package executors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"emcontroller/models"
)

func TestValidateAutoScheduleApp(t *testing.T) {
	type oneTestCase struct {
		name           string
		app            models.K8sApp
		expectedErrNum int
	}

	var testCases []oneTestCase

	// test cases about priority
	testCasesPriority := []oneTestCase{
		{
			name: "priority-1",
			app: models.K8sApp{
				Name:          "priority-1",
				Priority:      -1,
				Replicas:      1,
				AutoScheduled: true,
			},
			expectedErrNum: 2,
		},
		{
			name: "priority11",
			app: models.K8sApp{
				Name:          "priority11",
				Priority:      11,
				Replicas:      1,
				AutoScheduled: true,
			},
			expectedErrNum: 2,
		},
		{
			name: "priority10",
			app: models.K8sApp{
				Name:          "priority10",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
			},
			expectedErrNum: 1,
		},
		{
			name: "priority1",
			app: models.K8sApp{
				Name:          "priority1",
				Priority:      1,
				Replicas:      1,
				AutoScheduled: true,
			},
			expectedErrNum: 1,
		},
		{
			name: "priority5",
			app: models.K8sApp{
				Name:          "priority5",
				Priority:      5,
				Replicas:      1,
				AutoScheduled: true,
			},
			expectedErrNum: 1,
		},
	}
	testCases = append(testCases, testCasesPriority...)

	// test cases about AutoScheduled
	testCasesAutoScheduled := []oneTestCase{
		{
			name: "autoscheduledTrue",
			app: models.K8sApp{
				Name:          "autoscheduledTrue",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
			},
			expectedErrNum: 1,
		},
		{
			name: "autoscheduledFalse",
			app: models.K8sApp{
				Name:          "autoscheduledFalse",
				Priority:      3,
				Replicas:      1,
				AutoScheduled: false,
			},
			expectedErrNum: 2,
		},
		{
			name: "autoscheduledFalsePri-3",
			app: models.K8sApp{
				Name:          "autoscheduledFalsePri-3",
				Priority:      -3,
				Replicas:      1,
				AutoScheduled: false,
			},
			expectedErrNum: 3,
		},
	}
	testCases = append(testCases, testCasesAutoScheduled...)

	// test cases about Replicas
	testCasesReplicas := []oneTestCase{
		{
			name: "Replicas1",
			app: models.K8sApp{
				Name:          "Replicas1",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
			},
			expectedErrNum: 1,
		},
		{
			name: "Replicas2",
			app: models.K8sApp{
				Name:          "Replicas2",
				Priority:      9,
				Replicas:      2,
				AutoScheduled: true,
			},
			expectedErrNum: 2,
		},
		{
			name: "Replicas-1AutoFalse",
			app: models.K8sApp{
				Name:          "Replicas-1AutoFalse",
				Priority:      8,
				Replicas:      -1,
				AutoScheduled: false,
			},
			expectedErrNum: 3,
		},
	}
	testCases = append(testCases, testCasesReplicas...)

	// test cases about NodeName
	testCasesNodeName := []oneTestCase{
		{
			name: "WithNodeName",
			app: models.K8sApp{
				Name:          "WithNodeName",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				NodeName:      "ads",
			},
			expectedErrNum: 2,
		},
	}
	testCases = append(testCases, testCasesNodeName...)

	// test cases about NodeSelector
	testCasesNodeSelector := []oneTestCase{
		{
			name: "WithNodeSelector",
			app: models.K8sApp{
				Name:          "WithNodeSelector",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				NodeSelector: map[string]string{
					"aaa": "aaa",
					"bbb": "aaa",
				},
			},
			expectedErrNum: 2,
		},
	}
	testCases = append(testCases, testCasesNodeSelector...)

	// test cases about Container
	testCasesContainer := []oneTestCase{
		{
			name: "resLiReCpuNotEqual",
			app: models.K8sApp{
				Name:          "resLiReCpuNotEqual",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "1",
								Storage: "10Gi",
							},
							Requests: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "2",
								Storage: "10Gi",
							},
						},
					},
				},
			},
			expectedErrNum: 1,
		},
		{
			name: "resLiReMemNotEqual",
			app: models.K8sApp{
				Name:          "resLiReMemNotEqual",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "11Mi",
								CPU:     "2",
								Storage: "10Gi",
							},
							Requests: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "2",
								Storage: "10Gi",
							},
						},
					},
				},
			},
			expectedErrNum: 1,
		},
		{
			name: "resLiReStoNotEqual",
			app: models.K8sApp{
				Name:          "resLiReStoNotEqual",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "2",
								Storage: "12Gi",
							},
							Requests: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "2",
								Storage: "10Gi",
							},
						},
					},
				},
			},
			expectedErrNum: 1,
		},
		{
			name: "resLiReStoMemNotEqual",
			app: models.K8sApp{
				Name:          "resLiReStoMemNotEqual",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "11Mi",
								CPU:     "2",
								Storage: "12Gi",
							},
							Requests: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "2",
								Storage: "10Gi",
							},
						},
					},
				},
			},
			expectedErrNum: 1,
		},
		{
			name: "resLiReEqual",
			app: models.K8sApp{
				Name:          "resLiReEqual",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "11Mi",
								CPU:     "2",
								Storage: "12Gi",
							},
							Requests: models.K8sResList{
								Memory:  "11Mi",
								CPU:     "2",
								Storage: "12Gi",
							},
						},
					},
				},
			},
			expectedErrNum: 0,
		},
		{
			name: "resLiReCpuDecimal",
			app: models.K8sApp{
				Name:          "resLiReCpuDecimal",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "11Mi",
								CPU:     "0.5",
								Storage: "12Gi",
							},
							Requests: models.K8sResList{
								Memory:  "11Mi",
								CPU:     "0.5",
								Storage: "12Gi",
							},
						},
					},
				},
			},
			expectedErrNum: 0,
		},
		{
			name: "resLiReCpuWrong",
			app: models.K8sApp{
				Name:          "resLiReCpuWrong",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "11Mi",
								CPU:     "200m",
								Storage: "12Gi",
							},
							Requests: models.K8sResList{
								Memory:  "11Mi",
								CPU:     "200m",
								Storage: "12Gi",
							},
						},
					},
				},
			},
			expectedErrNum: 1,
		},
		{
			name: "resLiReMemWrongGi",
			app: models.K8sApp{
				Name:          "resLiReMemWrongGi",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "3Gi",
								CPU:     "2",
								Storage: "12Gi",
							},
							Requests: models.K8sResList{
								Memory:  "3Gi",
								CPU:     "2",
								Storage: "12Gi",
							},
						},
					},
				},
			},
			expectedErrNum: 1,
		},
		{
			name: "resLiReMemWrongNoUnit",
			app: models.K8sApp{
				Name:          "resLiReMemWrongNoUnit",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "12132165485",
								CPU:     "2",
								Storage: "12Gi",
							},
							Requests: models.K8sResList{
								Memory:  "12132165485",
								CPU:     "2",
								Storage: "12Gi",
							},
						},
					},
				},
			},
			expectedErrNum: 1,
		},
		{
			name: "resLiReStoWrong",
			app: models.K8sApp{
				Name:          "resLiReStoWrong",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "11Mi",
								CPU:     "0.5",
								Storage: "12000Mi",
							},
							Requests: models.K8sResList{
								Memory:  "11Mi",
								CPU:     "0.5",
								Storage: "12000Mi",
							},
						},
					},
				},
			},
			expectedErrNum: 1,
		},
	}
	testCases = append(testCases, testCasesContainer...)

	// test cases about Container
	testCasesMultiContainer := []oneTestCase{
		{
			name: "multiContainer",
			app: models.K8sApp{
				Name:          "multiContainer",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "1",
								Storage: "10Gi",
							},
							Requests: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "1",
								Storage: "10Gi",
							},
						},
					},
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "1",
								Storage: "10Gi",
							},
							Requests: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "1",
								Storage: "10Gi",
							},
						},
					},
				},
			},
			expectedErrNum: 1,
		},
	}
	testCases = append(testCases, testCasesMultiContainer...)

	// test no app name
	testCasesNoAppName := []oneTestCase{
		{
			name: "no name",
			app: models.K8sApp{
				Name:          "",
				Priority:      10,
				Replicas:      1,
				AutoScheduled: true,
				Containers: []models.K8sContainer{
					{
						Resources: models.K8sResReq{
							Limits: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "1",
								Storage: "10Gi",
							},
							Requests: models.K8sResList{
								Memory:  "10Mi",
								CPU:     "1",
								Storage: "10Gi",
							},
						},
					},
				},
			},
			expectedErrNum: 1,
		},
	}
	testCases = append(testCases, testCasesNoAppName...)

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		errs := ValidateAutoScheduleApp(testCase.app)
		t.Log("errors:", models.HandleErrSlice(errs))
		assert.Equal(t, testCase.expectedErrNum, len(errs), fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}
