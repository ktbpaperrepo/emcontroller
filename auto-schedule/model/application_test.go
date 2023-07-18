package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"emcontroller/models"
)

func TestAppMapCopy(t *testing.T) {
	testCases := []struct {
		name string
		src  map[string]Application
	}{
		{
			name: "case1",
			src: map[string]Application{
				"app1": {
					Name:     "app1",
					Priority: 2,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 2.0,
							Memory:  2048,
							Storage: 10,
						},
					},
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
						{
							AppName: "app3",
						},
					},
				},
				"app2": {
					Name:     "app2",
					Priority: 3,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 3.0,
							Memory:  4096,
							Storage: 20,
						},
					},
				},
				"app3": {
					Name:     "app3",
					Priority: 4,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 4.0,
							Memory:  4096,
							Storage: 30,
						},
					},
				},
			},
		},
		{
			name: "case2",
			src: map[string]Application{
				"app1": {
					Name:     "app1",
					Priority: 2,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 3.0,
							Memory:  2048,
							Storage: 10,
						},
					},
					Dependencies: []models.Dependency{
						{
							AppName: "app2",
						},
						{
							AppName: "app3",
						},
					},
				},
				"app2": {
					Name:     "app2",
					Priority: 3,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 4.0,
							Memory:  4096,
							Storage: 20,
						},
					},
				},
				"app3": {
					Name:     "app3",
					Priority: 4,
					Resources: AppResources{
						GenericResources: GenericResources{
							CpuCore: 5.0,
							Memory:  4096,
							Storage: 30,
						},
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		var dst map[string]Application

		// test extra key
		dst = AppMapCopy(testCase.src)
		extraKey := "app100"
		dst[extraKey] = Application{
			Name:     extraKey,
			Priority: 4,
			Resources: AppResources{
				GenericResources: GenericResources{
					CpuCore: 5.0,
					Memory:  4096,
					Storage: 30,
				},
			},
		}
		if app, exist := testCase.src[extraKey]; exist {
			t.Errorf("key \"%s\" should not exist in testCase.src, the value is: %+v", extraKey, app)
		} else {
			t.Logf("As expected, key \"%s\" does not exist in testCase.src", extraKey)
		}

		// test change
		dst = AppMapCopy(testCase.src)
		changeKey := "app1"
		dst[changeKey] = Application{
			Name:     "app1",
			Priority: 4,
			Resources: AppResources{
				GenericResources: GenericResources{
					CpuCore: 5.0,
					Memory:  4096,
					Storage: 30,
				},
			},
		}

		t.Logf("testCase.src[changeKey]:\n%+v\ndst[changeKey]:\n%+v", testCase.src[changeKey], dst[changeKey])
		assert.NotEqual(t, testCase.src[changeKey], dst[changeKey], fmt.Sprintf("testCase.src[changeKey] [%+v] and dst[changeKey] [%+v], show not be equal.", testCase.src[changeKey], dst[changeKey]))

	}

}
