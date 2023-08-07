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

func TestInnerGeneMutate(t *testing.T) {
	testCases := []struct {
		name   string
		m      *Mcssga
		clouds map[string]asmodel.Cloud
		ori    asmodel.SingleAppSolution
	}{
		{
			name:   "case1",
			m:      NewMcssga(100, 5000, 0.7, 0.2, 200),
			clouds: cloudsWithNetForTest()[3],
			ori: asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "NOKIA4",
			},
		},
		{
			name:   "case2",
			m:      NewMcssga(100, 5000, 0.7, 0.2, 200),
			clouds: cloudsWithNetForTest()[3],
			ori: asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "NOKIA5",
			},
		},
		{
			name:   "case3",
			m:      NewMcssga(100, 5000, 0.7, 0.2, 200),
			clouds: cloudsWithNetForTest()[3],
			ori: asmodel.SingleAppSolution{
				Accepted: false,
			},
		},
		{
			name:   "case4",
			m:      NewMcssga(100, 5000, 0.7, 0.2, 200),
			clouds: cloudsWithNetForTest()[3],
			ori: asmodel.SingleAppSolution{
				Accepted:        false,
				TargetCloudName: "NOKIA5",
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		var loopTimes int
		if testCase.ori.Accepted {
			loopTimes = 100
		} else {
			loopTimes = 20
		}

		for j := 0; j < loopTimes; j++ {
			mutated := testCase.m.geneMutate(testCase.clouds, testCase.ori)
			t.Logf("mutated %d: %v", j, mutated)
			if testCase.ori.Accepted {
				assert.NotEqual(t, testCase.ori, mutated)
			}
		}
		fmt.Println()
	}

}

func TestInnerCalcPoint2(t *testing.T) {
	testCases := []struct {
		name           string
		point1         int
		pointWidth     int
		expectedResult int
	}{
		{
			name:           "case1",
			point1:         0,
			pointWidth:     1,
			expectedResult: 0,
		},
		{
			name:           "case2",
			point1:         0,
			pointWidth:     3,
			expectedResult: 2,
		},
		{
			name:           "case3",
			point1:         2,
			pointWidth:     4,
			expectedResult: 5,
		},
		{
			name:           "case5",
			point1:         0,
			pointWidth:     13,
			expectedResult: 12,
		},
		{
			name:           "case6",
			point1:         1,
			pointWidth:     13,
			expectedResult: 13,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		assert.Equal(t, testCase.expectedResult, calcPoint2(testCase.point1, testCase.pointWidth), fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

// This test function can only show the traverse effect by log. It cannot make an error when the code has problems.
func TestAllPossTwoPointCrossover(t *testing.T) {
	testCases := []struct {
		name             string
		firstChromosome  asmodel.Solution
		secondChromosome asmodel.Solution
		appsOrder        []string
		expectedNewCh1   asmodel.Solution
		expectedNewCh2   asmodel.Solution
	}{
		{
			name:             "case1",
			firstChromosome:  solnsForTest()[6],
			secondChromosome: solnsForTest()[7],
			appsOrder:        appOrdersForTest()[1],
			expectedNewCh1:   solnsForTest()[6],
			expectedNewCh2:   solnsForTest()[7],
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualNewCh1, actualNewCh2 := AllPossTwoPointCrossover(testCase.firstChromosome, testCase.secondChromosome, nil, nil, testCase.appsOrder)
		assert.Equal(t, testCase.expectedNewCh1, actualNewCh1, fmt.Sprintf("%s: result is not expected", testCase.name))
		assert.Equal(t, testCase.expectedNewCh2, actualNewCh2, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}

func TestInnerTwoPointCrossover(t *testing.T) {
	testCases := []struct {
		name             string
		firstChromosome  asmodel.Solution
		secondChromosome asmodel.Solution
		appsOrder        []string
		point1           int
		point2           int
		expectedNewCh1   asmodel.Solution
		expectedNewCh2   asmodel.Solution
	}{
		{
			name:             "case1",
			firstChromosome:  solnsForTest()[6],
			secondChromosome: solnsForTest()[7],
			appsOrder:        appOrdersForTest()[1],
			point1:           1,
			point2:           4,
			expectedNewCh1: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        false,
						TargetCloudName: "NOKIA2",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
					},
					"app5": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "HPE2",
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA12",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "HPE1",
					},
				},
			},
			expectedNewCh2: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA1",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA2",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
					},
					"app4": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app5": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA10",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA2",
					},
				},
			},
		},
		{
			name:             "case2",
			firstChromosome:  solnsForTest()[6],
			secondChromosome: solnsForTest()[7],
			appsOrder:        appOrdersForTest()[1],
			point1:           2,
			point2:           3,
			expectedNewCh1: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        false,
						TargetCloudName: "NOKIA2",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA2",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
					},
					"app5": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "HPE2",
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA10",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "HPE1",
					},
				},
			},
			expectedNewCh2: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA1",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
					},
					"app4": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app5": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA12",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA2",
					},
				},
			},
		},
		{
			name:             "case3",
			firstChromosome:  solnsForTest()[6],
			secondChromosome: solnsForTest()[7],
			appsOrder:        appOrdersForTest()[0],
			point1:           4,
			point2:           6,
			expectedNewCh1: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA1",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA2",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app4": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
					},
					"app5": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "HPE2",
					},
					"app6": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA12",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "HPE1",
					},
				},
			},
			expectedNewCh2: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"app1": asmodel.SingleAppSolution{
						Accepted:        false,
						TargetCloudName: "NOKIA2",
					},
					"app2": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA4",
					},
					"app3": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA6",
					},
					"app4": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app5": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"app6": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA7",
					},
					"app7": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA10",
					},
					"app8": asmodel.SingleAppSolution{
						Accepted:        true,
						TargetCloudName: "NOKIA2",
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualNewCh1, actualNewCh2 := twoPointCrossover(testCase.firstChromosome, testCase.secondChromosome, testCase.appsOrder, testCase.point1, testCase.point2)
		assert.Equal(t, testCase.expectedNewCh1, actualNewCh1, fmt.Sprintf("%s: result is not expected", testCase.name))
		assert.Equal(t, testCase.expectedNewCh2, actualNewCh2, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}
