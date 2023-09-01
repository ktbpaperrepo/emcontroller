package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcVmAvailVcpu(t *testing.T) {
	testCases := []struct {
		name           string
		totalVcpu      float64
		expectedResult float64
	}{
		{
			name:           "case1",
			totalVcpu:      20,
			expectedResult: 19,
		},
		{
			name:           "case2",
			totalVcpu:      0,
			expectedResult: 0,
		},
		{
			name:           "case decimal",
			totalVcpu:      17.5,
			expectedResult: 16,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := CalcVmAvailVcpu(testCase.totalVcpu)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestCalcVmAvailRamMiB(t *testing.T) {
	testCases := []struct {
		name           string
		totalRamMiB    float64
		expectedResult float64
	}{
		{
			name:           "case1",
			totalRamMiB:    16384,
			expectedResult: 13721,
		},
		{
			name:           "case2",
			totalRamMiB:    1000,
			expectedResult: 0,
		},
		{
			name:           "case decimal 1",
			totalRamMiB:    16384.7,
			expectedResult: 13722,
		},
		{
			name:           "case decimal 2",
			totalRamMiB:    16384.2,
			expectedResult: 13721,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := CalcVmAvailRamMiB(testCase.totalRamMiB)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestCalcVmAvailStorGiB(t *testing.T) {
	testCases := []struct {
		name           string
		totalStorGiB   float64
		expectedResult float64
	}{
		{
			name:           "case1",
			totalStorGiB:   200,
			expectedResult: 140,
		},
		{
			name:           "case2",
			totalStorGiB:   11,
			expectedResult: 0,
		},
		{
			name:           "case3",
			totalStorGiB:   9,
			expectedResult: 0,
		},
		{
			name:           "case decimal",
			totalStorGiB:   128.6,
			expectedResult: 86,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := CalcVmAvailStorGiB(testCase.totalStorGiB)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestCalcVmTotalVcpu(t *testing.T) {
	testCases := []struct {
		name           string
		availVcpu      float64
		expectedResult float64
	}{
		{
			name:           "case1",
			availVcpu:      20,
			expectedResult: 21,
		},
		{
			name:           "case2",
			availVcpu:      0,
			expectedResult: 1,
		},
		{
			name:           "case decimal",
			availVcpu:      17.5,
			expectedResult: 19,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := CalcVmTotalVcpu(testCase.availVcpu)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestCalcVmTotalRamMiB(t *testing.T) {
	testCases := []struct {
		name           string
		availRamMiB    float64
		expectedResult float64
	}{
		{
			name:           "case1",
			availRamMiB:    16384,
			expectedResult: 19343,
		},
		{
			name:           "case2",
			availRamMiB:    0,
			expectedResult: 1138,
		},
		{
			name:           "case decimal 1",
			availRamMiB:    16384.7,
			expectedResult: 19343,
		},
		{
			name:           "case decimal 2",
			availRamMiB:    16384.2,
			expectedResult: 19343,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := CalcVmTotalRamMiB(testCase.availRamMiB)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestCalcVmTotalStorGiB(t *testing.T) {
	testCases := []struct {
		name           string
		availStorGiB   float64
		expectedResult float64
	}{
		{
			name:           "case1",
			availStorGiB:   160,
			expectedResult: 227,
		},
		{
			name:           "case2",
			availStorGiB:   0,
			expectedResult: 14,
		},
		{
			name:           "case3",
			availStorGiB:   9,
			expectedResult: 26,
		},
		{
			name:           "case decimal",
			availStorGiB:   128.6,
			expectedResult: 185,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := CalcVmTotalStorGiB(testCase.availStorGiB)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestBackForth(t *testing.T) {
	var availStorGiBs []float64 = []float64{160, 128.6, 25.3, 33.5, 40.8}
	for _, availStorGiB := range availStorGiBs {
		t.Log("original availStorGiB:", availStorGiB)
		totalStorGiB := CalcVmTotalStorGiB(availStorGiB)
		t.Log("first totalStorGiB:", totalStorGiB)
		for i := 0; i < 5; i++ {
			availStorGiB = CalcVmAvailStorGiB(totalStorGiB)
			t.Log("availStorGiB", availStorGiB)
			totalStorGiB = CalcVmTotalStorGiB(availStorGiB)
			t.Log("totalStorGiB", totalStorGiB)
		}
	}
}
