package models

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

var generateFakeDockerClient client.Opt = client.WithHost("http://192.168.100.36:19998")

// go test $GOPATH/src/emcontroller/models/ --run TestListImages -v
func TestListImages(t *testing.T) {
	err := generateFakeDockerClient(cli)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
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
		actualResult, actualError := ListImages()
		for i, oneImage := range actualResult {
			t.Logf("image %d: %#v\n", i, oneImage)
		}
		assert.Equal(t, testCase.expectedError, actualError, fmt.Sprintf("%s: Error is not expected", testCase.name))
	}
}

func TestLoadImage(t *testing.T) {
	err := generateFakeDockerClient(cli)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	gopath, found := os.LookupEnv("GOPATH")
	if !found {
		t.Fatalf("GOPATH env not found\n")
	}
	// use this file to test
	var imageFilePath1 string
	if runtime.GOOS == "windows" {
		imageFilePath1 = gopath + "\\src\\emcontroller\\upload\\helloworld.tar"
	} else {
		imageFilePath1 = gopath + "/src/emcontroller/upload/helloworld.tar"
	}
	imageFile1, err := os.Open(imageFilePath1)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer imageFile1.Close()

	var imageFilePath2 string
	if runtime.GOOS == "windows" {
		imageFilePath2 = gopath + "\\src\\emcontroller\\upload\\helloworld2.tar"
	} else {
		imageFilePath2 = gopath + "/src/emcontroller/upload/helloworld2.tar"
	}
	imageFile2, err := os.Open(imageFilePath2)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}
	defer imageFile2.Close()

	testCases := []struct {
		name           string
		imageFile      *os.File
		expectedResult string
	}{
		{
			name:           "case1",
			imageFile:      imageFile1,
			expectedResult: "feb5d9fea6a5e9606aa995e879d862b825965ba48de054caab5ef356dc6b3412",
		},
		{
			name:           "case2",
			imageFile:      imageFile2,
			expectedResult: "192.168.100.36:5000/helloworld:latest",
		},
	}
	for _, testCase := range testCases {
		t.Logf("test: %s", testCase.name)
		actualResult, actualError := LoadImage(testCase.imageFile)
		if actualError != nil {
			t.Errorf("Error: %s\n", actualError.Error())
		}
		t.Logf("actualResult: %#v\n", actualResult)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: Result is not expected", testCase.name))
	}
}
