package algorithms

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

func TestDfsGroup(t *testing.T) {
	testCases := []struct {
		name           string
		apps           map[string]asmodel.Application
		expectedResult [][]string
	}{
		{
			name: "case 1",
			apps: map[string]asmodel.Application{
				"app1": asmodel.Application{
					Name: "app1",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app2",
						},
					},
				},
				"app2": asmodel.Application{
					Name:         "app2",
					Dependencies: []models.Dependency{},
				},
				"app3": asmodel.Application{
					Name: "app3",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app2",
						},
						models.Dependency{
							AppName: "app4",
						},
					},
				},
				"app4": asmodel.Application{
					Name:         "app4",
					Dependencies: []models.Dependency{},
				},
				"app5": asmodel.Application{
					Name: "app5",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app6",
						},
					},
				},
				"app6": asmodel.Application{
					Name:         "app6",
					Dependencies: []models.Dependency{},
				},
			},
			expectedResult: [][]string{
				[]string{"app1", "app2", "app3", "app4"},
				[]string{"app5", "app6"},
			},
		},
		{
			name: "case 2",
			apps: map[string]asmodel.Application{
				"app1": asmodel.Application{
					Name:         "app1",
					Dependencies: []models.Dependency{},
				},
				"app2": asmodel.Application{
					Name: "app2",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app1",
						},
					},
				},
				"app3": asmodel.Application{
					Name: "app3",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app1",
						},
					},
				},
				"app4": asmodel.Application{
					Name: "app4",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app3",
						},
					},
				},
				"app5": asmodel.Application{
					Name: "app5",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app2",
						},
						models.Dependency{
							AppName: "app3",
						},
					},
				},
				"app6": asmodel.Application{
					Name: "app6",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app7",
						},
					},
				},
				"app7": asmodel.Application{
					Name:         "app7",
					Dependencies: []models.Dependency{},
				},
				"app8": asmodel.Application{
					Name: "app8",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app7",
						},
						models.Dependency{
							AppName: "app9",
						},
					},
				},
				"app9": asmodel.Application{
					Name:         "app9",
					Dependencies: []models.Dependency{},
				},
				"app10": asmodel.Application{
					Name: "app10",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app8",
						},
					},
				},
			},
			expectedResult: [][]string{
				[]string{"app1", "app2", "app3", "app4", "app5"},
				[]string{"app6", "app7", "app8", "app9", "app10"},
			},
		},
		{
			name: "case 3",
			apps: map[string]asmodel.Application{
				"app1": asmodel.Application{
					Name: "app1",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app2",
						},
						models.Dependency{
							AppName: "app3",
						},
					},
				},
				"app2": asmodel.Application{
					Name: "app2",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app3",
						},
						models.Dependency{
							AppName: "app4",
						},
					},
				},
				"app3": asmodel.Application{
					Name:         "app3",
					Dependencies: []models.Dependency{},
				},
				"app4": asmodel.Application{
					Name:         "app4",
					Dependencies: []models.Dependency{},
				},
				"app5": asmodel.Application{
					Name: "app5",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app10",
						},
					},
				},
				"app6": asmodel.Application{
					Name:         "app6",
					Dependencies: []models.Dependency{},
				},
				"app7": asmodel.Application{
					Name: "app7",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app5",
						},
						models.Dependency{
							AppName: "app6",
						},
					},
				},
				"app8": asmodel.Application{
					Name:         "app8",
					Dependencies: []models.Dependency{},
				},
				"app9": asmodel.Application{
					Name: "app9",
					Dependencies: []models.Dependency{
						models.Dependency{
							AppName: "app8",
						},
					},
				},
				"app10": asmodel.Application{
					Name:         "app10",
					Dependencies: []models.Dependency{},
				},
				"app11": asmodel.Application{
					Name:         "app11",
					Dependencies: []models.Dependency{},
				},
			},
			expectedResult: [][]string{
				[]string{"app1", "app2", "app3", "app4"},
				[]string{"app5", "app6", "app7", "app10"},
				[]string{"app9", "app8"},
				[]string{"app11"},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		biApps := genBiDir(testCase.apps)
		t.Log("biApps:")
		for _, biApp := range biApps {
			t.Log(biApp)
		}

		depGroups := groupByDep(testCase.apps)
		t.Log("depGroups:")
		for _, depGroup := range depGroups {
			t.Log(depGroup)
		}

		actualResult := depGroups
		expectedResult := testCase.expectedResult
		for _, s := range actualResult {
			sort.Strings(s)
		}
		for _, s := range expectedResult {
			sort.Strings(s)
		}

		assert.ElementsMatch(t, expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}
