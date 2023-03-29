package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// go test $GOPATH/src/emcontroller/models/ --run TestListDeployment -v
func TestListDeployment(t *testing.T) {
	InitKubernetesClient()
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

func TestListNodes(t *testing.T) {
	InitKubernetesClient()
	nodes, err := ListNodes(metav1.ListOptions{})
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	}
	t.Logf("nodes: %v", nodes)
}

func TestGetNode(t *testing.T) {
	InitKubernetesClient()
	node, err := GetNode("node1", metav1.GetOptions{})
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	}
	t.Logf("node: %v", node)
}

func TestGetJoinCmd(t *testing.T) {
	InitKubernetesClient()
	joinCmd, err := GetJoinCmd()
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	}
	t.Logf("joinCmd: %s --node-name=test", joinCmd)
}

func TestJoinOneNode(t *testing.T) {
	InitKubernetesClient()
	joinCmd, err := GetJoinCmd()
	if err != nil {
		t.Errorf("GetJoinCmd error: %s", err.Error())
	}
	if err := AddNode("cnode2", "10.234.234.234", joinCmd); err != nil {
		t.Errorf("AddNode error: %s", err.Error())
	} else {
		t.Logf("node joined")
	}
}

func TestJoinSeveralNodes(t *testing.T) {
	InitKubernetesClient()

	type nodeNameIP struct {
		name string
		ip   string
	}
	nodes := []nodeNameIP{
		{"cnode1", "10.234.234.100"},
		{"cnode2", "10.234.234.197"},
		{"cnode4", "10.234.234.133"},
		{"node1", "192.168.100.145"},
		{"node2", "192.168.100.25"},
	}

	var nodeNames, nodeIPs []string

	for _, node := range nodes {
		nodeNames = append(nodeNames, node.name)
		nodeIPs = append(nodeIPs, node.ip)
	}

	errs := AddNodes(nodeNames, nodeIPs)
	if len(errs) != 0 {
		sumErr := "AddNodes Error:"
		for _, err := range errs {
			sumErr += "\n" + err.Error()
		}
		t.Error(sumErr)
	} else {
		t.Logf("nodes joined")
	}
}

func TestUninstallNode(t *testing.T) {
	InitKubernetesClient()
	err := UninstallNode("node1")
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	}
}
