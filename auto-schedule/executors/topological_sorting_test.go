package executors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"emcontroller/models"
)

func TestTopoSort(t *testing.T) {
	testCases := []struct {
		name                     string
		apps                     map[string]models.K8sApp
		expectedTopologicalOrder [][]string
		expectedHasCycles        bool
	}{
		{
			name: "noCycle",
			apps: map[string]models.K8sApp{
				"app1": models.K8sApp{
					Name: "app1",
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				"app2": models.K8sApp{
					Name: "app2",
					Dependencies: []models.Dependency{
						{
							AppName: "app3",
						},
						{
							AppName: "app5",
						},
					},
				},
				"app3": models.K8sApp{
					Name:         "app3",
					Dependencies: []models.Dependency{},
				},
				"app4": models.K8sApp{
					Name: "app4",
					Dependencies: []models.Dependency{
						{
							AppName: "app5",
						},
					},
				},
				"app5": models.K8sApp{
					Name: "app5",
					Dependencies: []models.Dependency{
						{
							AppName: "app6",
						},
					},
				},
				"app6": models.K8sApp{
					Name:         "app6",
					Dependencies: []models.Dependency{},
				},
			},
			expectedTopologicalOrder: [][]string{
				[]string{"app3", "app6"},
				[]string{"app5"},
				[]string{"app2", "app4"},
				[]string{"app1"},
			},
			expectedHasCycles: false,
		},
		{
			name: "noCycle",
			apps: map[string]models.K8sApp{
				"app1": models.K8sApp{
					Name: "app1",
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				"app2": models.K8sApp{
					Name: "app2",
					Dependencies: []models.Dependency{
						{
							AppName: "app3",
						},
						{
							AppName: "app5",
						},
					},
				},
				"app3": models.K8sApp{
					Name:         "app3",
					Dependencies: []models.Dependency{},
				},
				"app4": models.K8sApp{
					Name: "app4",
					Dependencies: []models.Dependency{
						{
							AppName: "app5",
						},
					},
				},
				"app5": models.K8sApp{
					Name: "app5",
					Dependencies: []models.Dependency{
						{
							AppName: "app6",
						},
					},
				},
				"app6": models.K8sApp{
					Name: "app6",
					Dependencies: []models.Dependency{
						{
							AppName: "app3",
						},
					},
				},
			},
			expectedTopologicalOrder: [][]string{
				[]string{"app3"},
				[]string{"app6"},
				[]string{"app5"},
				[]string{"app2", "app4"},
				[]string{"app1"},
			},
			expectedHasCycles: false,
		},
		{
			name: "hasCycle",
			apps: map[string]models.K8sApp{
				"app1": models.K8sApp{
					Name: "app1",
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				"app2": models.K8sApp{
					Name: "app2",
					Dependencies: []models.Dependency{
						{
							AppName: "app3",
						},
						{
							AppName: "app5",
						},
					},
				},
				"app3": models.K8sApp{
					Name:         "app3",
					Dependencies: []models.Dependency{},
				},
				"app4": models.K8sApp{
					Name: "app4",
					Dependencies: []models.Dependency{
						{
							AppName: "app5",
						},
					},
				},
				"app5": models.K8sApp{
					Name: "app5",
					Dependencies: []models.Dependency{
						{
							AppName: "app6",
						},
					},
				},
				"app6": models.K8sApp{
					Name: "app6",
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
			},
			expectedTopologicalOrder: [][]string{
				[]string{"app3"},
			},
			expectedHasCycles: true,
		},
		{
			name: "hasCycle",
			apps: map[string]models.K8sApp{
				"app1": models.K8sApp{
					Name: "app1",
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				"app2": models.K8sApp{
					Name: "app2",
					Dependencies: []models.Dependency{
						{
							AppName: "app3",
						},
						{
							AppName: "app5",
						},
					},
				},
				"app3": models.K8sApp{
					Name:         "app3",
					Dependencies: []models.Dependency{},
				},
				"app4": models.K8sApp{
					Name: "app4",
					Dependencies: []models.Dependency{
						{
							AppName: "app5",
						},
					},
				},
				"app5": models.K8sApp{
					Name: "app5",
					Dependencies: []models.Dependency{
						{
							AppName: "app6",
						},
						{
							AppName: "app1",
						},
					},
				},
				"app6": models.K8sApp{
					Name:         "app6",
					Dependencies: []models.Dependency{},
				},
			},
			expectedTopologicalOrder: [][]string{
				[]string{"app3", "app6"},
			},
			expectedHasCycles: true,
		},
		{
			name: "noCycle",
			apps: map[string]models.K8sApp{
				"app1": models.K8sApp{
					Name: "app1",
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				"app2": models.K8sApp{
					Name: "app2",
					Dependencies: []models.Dependency{
						{
							AppName: "app3",
						},
					},
				},
				"app3": models.K8sApp{
					Name:         "app3",
					Dependencies: []models.Dependency{},
				},
			},
			expectedTopologicalOrder: [][]string{
				[]string{"app3"},
				[]string{"app2"},
				[]string{"app1"},
			},
			expectedHasCycles: false,
		},
		{
			name: "hasCycle",
			apps: map[string]models.K8sApp{
				"app1": models.K8sApp{
					Name: "app1",
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				"app2": models.K8sApp{
					Name: "app2",
					Dependencies: []models.Dependency{
						{
							AppName: "app3",
						},
					},
				},
				"app3": models.K8sApp{
					Name: "app3",
					Dependencies: []models.Dependency{
						{
							AppName: "app1",
						},
					},
				},
			},
			expectedTopologicalOrder: [][]string{},
			expectedHasCycles:        true,
		},
		{
			name: "hasCycle",
			apps: map[string]models.K8sApp{
				"app1": models.K8sApp{
					Name: "app1",
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
					},
				},
				"app2": models.K8sApp{
					Name: "app2",
					Dependencies: []models.Dependency{
						{
							AppName: "app1",
						},
					},
				},
				"app3": models.K8sApp{
					Name:         "app3",
					Dependencies: []models.Dependency{},
				},
			},
			expectedTopologicalOrder: [][]string{
				[]string{"app3"},
			},
			expectedHasCycles: true,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		topoOrder, hasCycles := TopoSort(testCase.apps)
		fmt.Printf("hasCycle is %t. topo order is:\n", hasCycles)
		for i := 0; i < len(topoOrder); i++ {
			fmt.Println(topoOrder[i])
		}

		assert.Equal(t, len(testCase.expectedTopologicalOrder), len(topoOrder), fmt.Sprintf("%s: len(topoOrder) is not expected", testCase.name))
		for i := 0; i < len(topoOrder); i++ {
			assert.ElementsMatch(t, testCase.expectedTopologicalOrder[i], topoOrder[i], fmt.Sprintf("%s: topoOrder is not expected", testCase.name))
		}
		assert.Equal(t, testCase.expectedHasCycles, hasCycles, fmt.Sprintf("%s: hasCycles is not expected", testCase.name))
	}
}
