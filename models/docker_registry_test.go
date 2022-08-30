package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test $GOPATH/src/emcontroller/models/ --run TestGetCatalog -v
func TestGetCatalog(t *testing.T) {
	DockerRegistry = "192.168.100.36:5000"
	testCases := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "case1",
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Logf("test: %s", testCase.name)
		actualResult, actualError := GetCatalog()
		for i, oneImage := range actualResult {
			t.Logf("image %d: %s\n", i, oneImage)
		}
		assert.Equal(t, testCase.expectedError, actualError, fmt.Sprintf("%s: Error is not expected", testCase.name))
	}
}

func TestListTags(t *testing.T) {
	DockerRegistry = "192.168.100.36:5000"
	testCases := []struct {
		name          string
		imageName     string
		expectedError error
	}{
		{
			name:          "case1",
			imageName:     "helloworld12345",
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Logf("test: %s", testCase.name)
		actualResult, actualError := ListTags(testCase.imageName)
		for i, oneTag := range actualResult {
			t.Logf("Tag %d: %s\n", i, oneTag)
		}
		assert.Equal(t, testCase.expectedError, actualError, fmt.Sprintf("%s: Error is not expected", testCase.name))
	}
}

func TestListRepoTags(t *testing.T) {
	DockerRegistry = "192.168.100.36:5000"
	testCases := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "case1",
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Logf("test: %s", testCase.name)
		actualResult, actualError := ListRepoTags()
		for i, repoTag := range actualResult {
			t.Logf("RepoTag %d: %s\n", i, repoTag)
		}
		assert.Equal(t, testCase.expectedError, actualError, fmt.Sprintf("%s: Error is not expected", testCase.name))
	}
}
