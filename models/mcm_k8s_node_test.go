package models

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
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

func TestGetResOccupiedByPod(t *testing.T) {

	testCases := []struct {
		name           string
		pod            apiv1.Pod
		expectedResult K8sNodeRes
	}{
		{
			name: "case1",
			pod: apiv1.Pod{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Resources: apiv1.ResourceRequirements{},
						},
					},
				},
			},
			expectedResult: K8sNodeRes{
				CpuCore: 0,
				Memory:  0,
				Storage: 0,
			},
		},
		{
			name: "case2",
			pod: apiv1.Pod{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Resources: apiv1.ResourceRequirements{
								Limits: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("100m"),
									apiv1.ResourceMemory:           resource.MustParse("500Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("10Gi"),
								},
								Requests: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("100m"),
									apiv1.ResourceMemory:           resource.MustParse("500Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("10Gi"),
								},
							},
						},
					},
				},
			},
			expectedResult: K8sNodeRes{
				CpuCore: 0.1,
				Memory:  500,
				Storage: 10,
			},
		},
		{
			name: "case3",
			pod: apiv1.Pod{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Resources: apiv1.ResourceRequirements{
								Limits: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("1.2"),
									apiv1.ResourceMemory:           resource.MustParse("666Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("7Gi"),
								},
								Requests: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("1.2"),
									apiv1.ResourceMemory:           resource.MustParse("666Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("7Gi"),
								},
							},
						},
					},
				},
			},
			expectedResult: K8sNodeRes{
				CpuCore: 1.2,
				Memory:  666,
				Storage: 7,
			},
		},
		{
			name: "case4",
			pod: apiv1.Pod{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Resources: apiv1.ResourceRequirements{
								Limits: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("1.2"),
									apiv1.ResourceMemory:           resource.MustParse("666Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("7Gi"),
								},
								Requests: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("1.2"),
									apiv1.ResourceMemory:           resource.MustParse("666Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("7Gi"),
								},
							},
						},
						{
							Resources: apiv1.ResourceRequirements{
								Limits: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("100m"),
									apiv1.ResourceMemory:           resource.MustParse("500Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("10Gi"),
								},
								Requests: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("100m"),
									apiv1.ResourceMemory:           resource.MustParse("500Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("10Gi"),
								},
							},
						},
					},
				},
			},
			expectedResult: K8sNodeRes{
				CpuCore: 1.3,
				Memory:  1166,
				Storage: 17,
			},
		},
		{
			name: "case5",
			pod: apiv1.Pod{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Resources: apiv1.ResourceRequirements{
								Limits: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("5"),
									apiv1.ResourceMemory:           resource.MustParse("2000Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("107Gi"),
								},
								Requests: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("3.2"),
									apiv1.ResourceMemory:           resource.MustParse("366Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("57Gi"),
								},
							},
						},
						{
							Resources: apiv1.ResourceRequirements{
								Limits: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("10"),
									apiv1.ResourceMemory:           resource.MustParse("500Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("10Gi"),
								},
								Requests: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("2"),
									apiv1.ResourceMemory:           resource.MustParse("100Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("6Gi"),
								},
							},
						},
					},
				},
			},
			expectedResult: K8sNodeRes{
				CpuCore: 5.2,
				Memory:  466,
				Storage: 63,
			},
		},
		{
			name: "case6",
			pod: apiv1.Pod{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Resources: apiv1.ResourceRequirements{
								Limits: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("5"),
									apiv1.ResourceMemory:           resource.MustParse("2000Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("107Gi"),
								},
								Requests: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("3.2"),
									apiv1.ResourceMemory:           resource.MustParse("366Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("57Gi"),
								},
							},
						},
						{
							Resources: apiv1.ResourceRequirements{
								Limits: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("10"),
									apiv1.ResourceMemory:           resource.MustParse("500Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("10Gi"),
								},
								Requests: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("2"),
									apiv1.ResourceMemory:           resource.MustParse("100Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("6Gi"),
								},
							},
						},
						{
							Resources: apiv1.ResourceRequirements{
								Limits: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("11"),
									apiv1.ResourceMemory:           resource.MustParse("400Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("22Gi"),
								},
								Requests: map[apiv1.ResourceName]resource.Quantity{
									apiv1.ResourceCPU:              resource.MustParse("3"),
									apiv1.ResourceMemory:           resource.MustParse("200Mi"),
									apiv1.ResourceEphemeralStorage: resource.MustParse("11Gi"),
								},
							},
						},
					},
				},
			},
			expectedResult: K8sNodeRes{
				CpuCore: 8.2,
				Memory:  666,
				Storage: 74,
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := GetResOccupiedByPod(testCase.pod)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}

}
