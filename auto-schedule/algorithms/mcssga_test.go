package algorithms

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	asmodel "emcontroller/auto-schedule/model"
)

func TestInnerSetMaxReaRtt(t *testing.T) {
	testCases := []struct {
		name           string
		m              *Mcssga
		clouds         map[string]asmodel.Cloud
		expectedResult float64
	}{
		{
			name:           "case1",
			m:              NewMcssga(100, 5000, 0.7, 0.2, 200),
			clouds:         cloudsWithNetForTest()[3],
			expectedResult: 5.2,
		},
		{
			name:           "case2",
			m:              NewMcssga(100, 5000, 0.7, 0.2, 200),
			clouds:         cloudsWithNetForTest()[2],
			expectedResult: 0.753,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		assert.InDelta(t, 0, testCase.m.MaxReachableRtt, testDelta, fmt.Sprintf("%s: result is not expected", testCase.name))
		testCase.m.setMaxReaRtt(testCase.clouds)
		assert.InDelta(t, testCase.expectedResult, testCase.m.MaxReachableRtt, testDelta, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}
