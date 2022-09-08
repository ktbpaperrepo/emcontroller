package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test $GOPATH/src/emcontroller/models/ --run TestListDeployment -v
func TestListDeployment(t *testing.T) {
	testCases := []struct {
		name          string
		namespace     string
		expectedError error
	}{
		{
			name:          "case1",
			namespace:     KubernetesNamespace,
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Logf("test: %s", testCase.name)
		actualResult, actualError := ListDeployment(testCase.namespace)
		for i, oneDeployment := range actualResult {
			t.Logf("deployment %d: %#v\n", i, oneDeployment)
		}
		assert.Equal(t, testCase.expectedError, actualError, fmt.Sprintf("%s: Error is not expected", testCase.name))
	}
}
