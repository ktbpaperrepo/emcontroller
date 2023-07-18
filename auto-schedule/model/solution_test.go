package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSolutionCopy(t *testing.T) {
	testCases := []struct {
		name string
		src  Solution
	}{
		{
			name: "case1",
			src: Solution{
				"app1": SingleAppSolution{
					Accepted: false,
				},
				"app2": SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "cloud1",
				},
			},
		},
		{
			name: "case2",
			src: Solution{
				"app2": SingleAppSolution{
					Accepted: false,
				},
				"app20": SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "cloud2",
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		var dst Solution

		// test extra key
		dst = SolutionCopy(testCase.src)
		extraKey := "app100"
		dst[extraKey] = SingleAppSolution{
			Accepted: false,
		}
		if app, exist := testCase.src[extraKey]; exist {
			t.Errorf("key \"%s\" should not exist in testCase.src, the value is: %+v", extraKey, app)
		} else {
			t.Logf("As expected, key \"%s\" does not exist in testCase.src", extraKey)
		}

		// test change
		dst = SolutionCopy(testCase.src)
		changeKey := "app2"
		dst[changeKey] = SingleAppSolution{
			Accepted:        true,
			TargetCloudName: "cloud20",
		}

		t.Logf("testCase.src[changeKey]:\n%+v\ndst[changeKey]:\n%+v", testCase.src[changeKey], dst[changeKey])
		assert.NotEqual(t, testCase.src[changeKey], dst[changeKey], fmt.Sprintf("testCase.src[changeKey] [%+v] and dst[changeKey] [%+v], show not be equal.", testCase.src[changeKey], dst[changeKey]))

	}
}
