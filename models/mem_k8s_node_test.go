package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var nodesForTest = []apiv1.Node{
	apiv1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node1",
		},
	},
	apiv1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node2",
		},
	},
	apiv1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node3",
		},
	},
	apiv1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node4",
		},
	},
}

func TestFindIdxNodeInList(t *testing.T) {

	testCases := []struct {
		name           string
		nodeList       []apiv1.Node
		nodeNameToFind string
		expectedResult int
	}{
		{
			name:           "case found 1",
			nodeList:       nodesForTest,
			nodeNameToFind: "node1",
			expectedResult: 0,
		},
		{
			name:           "case found 2",
			nodeList:       nodesForTest,
			nodeNameToFind: "node2",
			expectedResult: 1,
		},
		{
			name:           "case found 3",
			nodeList:       nodesForTest,
			nodeNameToFind: "node3",
			expectedResult: 2,
		},
		{
			name:           "case found 4",
			nodeList:       nodesForTest,
			nodeNameToFind: "node4",
			expectedResult: 3,
		},
		{
			name:           "case not found 1",
			nodeList:       nodesForTest,
			nodeNameToFind: "node5",
			expectedResult: -1,
		},
		{
			name:           "case not found 2",
			nodeList:       nodesForTest,
			nodeNameToFind: "asdfasf",
			expectedResult: -1,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := FindIdxNodeInList(testCase.nodeList, testCase.nodeNameToFind)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}

func TestRemoveNodeFromList(t *testing.T) {

	testCases := []struct {
		name             string
		nodeList         *[]apiv1.Node
		nodeNameToRemove string
		expectedResult   []apiv1.Node
	}{
		{
			name: "case exist 1",
			nodeList: &[]apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node1",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node2",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node3",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node4",
					},
				},
			},
			nodeNameToRemove: "node1",
			expectedResult: []apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node2",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node3",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node4",
					},
				},
			},
		},
		{
			name: "case exist 2",
			nodeList: &[]apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node1",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node2",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node3",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node4",
					},
				},
			},
			nodeNameToRemove: "node2",
			expectedResult: []apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node1",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node3",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node4",
					},
				},
			},
		},
		{
			name: "case exist 3",
			nodeList: &[]apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node1",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node2",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node3",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node4",
					},
				},
			},
			nodeNameToRemove: "node3",
			expectedResult: []apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node1",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node2",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node4",
					},
				},
			},
		},
		{
			name: "case exist 4",
			nodeList: &[]apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node1",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node2",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node3",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node4",
					},
				},
			},
			nodeNameToRemove: "node4",
			expectedResult: []apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node1",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node2",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node3",
					},
				},
			},
		},
		{
			name: "case not exist 1",
			nodeList: &[]apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node1",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node2",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node3",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node4",
					},
				},
			},
			nodeNameToRemove: "node5",
			expectedResult: []apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node1",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node2",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node3",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node4",
					},
				},
			},
		},
		{
			name: "case not exist 2",
			nodeList: &[]apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node1",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node2",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node3",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node4",
					},
				},
			},
			nodeNameToRemove: "asdfasf",
			expectedResult: []apiv1.Node{
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node1",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node2",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node3",
					},
				},
				apiv1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node4",
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		RemoveNodeFromList(testCase.nodeList, testCase.nodeNameToRemove)
		assert.Equal(t, testCase.expectedResult, *(testCase.nodeList), fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}
