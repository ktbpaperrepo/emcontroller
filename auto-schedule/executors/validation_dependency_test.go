package executors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"emcontroller/models"
)

func TestValidateAutoScheduleDep(t *testing.T) {
	testCases := []struct {
		name           string
		apps           []models.K8sApp
		expectedErrNum int
	}{
		{
			name: "noErr",
			apps: []models.K8sApp{
				models.K8sApp{
					Name:     "app1",
					Priority: 2,
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				models.K8sApp{
					Name:     "app2",
					Priority: 3,
					Dependencies: []models.Dependency{
						{
							AppName: "app3",
						},
					},
				},
				models.K8sApp{
					Name:         "app3",
					Priority:     10,
					Dependencies: []models.Dependency{},
				},
			},
			expectedErrNum: 0,
		},
		{
			name: "noErr",
			apps: []models.K8sApp{
				models.K8sApp{
					Name:     "app1",
					Priority: 2,
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				models.K8sApp{
					Name:     "app2",
					Priority: 10,
					Dependencies: []models.Dependency{
						{
							AppName: "app3",
						},
					},
				},
				models.K8sApp{
					Name:         "app3",
					Priority:     10,
					Dependencies: []models.Dependency{},
				},
			},
			expectedErrNum: 0,
		},
		{
			name: "PriorityErr",
			apps: []models.K8sApp{
				models.K8sApp{
					Name:     "app1",
					Priority: 2,
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				models.K8sApp{
					Name:     "app2",
					Priority: 1,
					Dependencies: []models.Dependency{
						{
							AppName: "app3",
						},
					},
				},
				models.K8sApp{
					Name:         "app3",
					Priority:     10,
					Dependencies: []models.Dependency{},
				},
			},
			expectedErrNum: 1,
		},
		{
			name: "ExistPriorityErr",
			apps: []models.K8sApp{
				models.K8sApp{
					Name:     "app1",
					Priority: 2,
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				models.K8sApp{
					Name:     "app2",
					Priority: 1,
					Dependencies: []models.Dependency{
						{
							AppName: "app5",
						},
					},
				},
				models.K8sApp{
					Name:         "app3",
					Priority:     10,
					Dependencies: []models.Dependency{},
				},
			},
			expectedErrNum: 2,
		},
		{
			name: "ExistErr",
			apps: []models.K8sApp{
				models.K8sApp{
					Name:     "app1",
					Priority: 2,
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				models.K8sApp{
					Name:     "app2",
					Priority: 4,
					Dependencies: []models.Dependency{
						{
							AppName: "app5",
						},
					},
				},
				models.K8sApp{
					Name:         "app3",
					Priority:     10,
					Dependencies: []models.Dependency{},
				},
			},
			expectedErrNum: 1,
		},
		{
			name: "ExistPriorityErr",
			apps: []models.K8sApp{
				models.K8sApp{
					Name:     "app1",
					Priority: 2,
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				models.K8sApp{
					Name:     "app2",
					Priority: 4,
					Dependencies: []models.Dependency{
						{
							AppName: "app1",
						},
					},
				},
				models.K8sApp{
					Name:     "app3",
					Priority: 10,
					Dependencies: []models.Dependency{
						{
							AppName: "app5",
						},
					},
				},
			},
			expectedErrNum: 2,
		},
		{
			name: "PriorityErr",
			apps: []models.K8sApp{
				models.K8sApp{
					Name:     "app1",
					Priority: 2,
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				models.K8sApp{
					Name:     "app2",
					Priority: 4,
					Dependencies: []models.Dependency{
						{
							AppName: "app1",
						},
					},
				},
				models.K8sApp{
					Name:         "app3",
					Priority:     10,
					Dependencies: []models.Dependency{},
				},
			},
			expectedErrNum: 1,
		},
		{
			name: "CircularDepErr",
			apps: []models.K8sApp{
				models.K8sApp{
					Name:     "app1",
					Priority: 4,
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				models.K8sApp{
					Name:     "app2",
					Priority: 4,
					Dependencies: []models.Dependency{
						{
							AppName: "app1",
						},
					},
				},
				models.K8sApp{
					Name:         "app3",
					Priority:     10,
					Dependencies: []models.Dependency{},
				},
			},
			expectedErrNum: 1,
		},
		{
			name: "CircularDepErr",
			apps: []models.K8sApp{
				models.K8sApp{
					Name:     "app1",
					Priority: 10,
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				models.K8sApp{
					Name:     "app2",
					Priority: 10,
					Dependencies: []models.Dependency{
						{
							AppName: "app3",
						},
					},
				},
				models.K8sApp{
					Name:     "app3",
					Priority: 10,
					Dependencies: []models.Dependency{
						{
							AppName: "app1",
						},
					},
				},
			},
			expectedErrNum: 1,
		},
	}
	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		errs := ValidateAutoScheduleDep(testCase.apps)
		fmt.Println("errors:", models.HandleErrSlice(errs))
		assert.Equal(t, testCase.expectedErrNum, len(errs), fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}
