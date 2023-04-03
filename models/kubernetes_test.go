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
	if err := AddNode(IaasVm{Name: "node1", IPs: []string{"192.168.100.43"}}, joinCmd); err != nil {
		t.Errorf("AddNode error: %s", err.Error())
	} else {
		t.Logf("node joined")
	}
}

func TestJoinSeveralNodes(t *testing.T) {
	InitKubernetesClient()

	vms := []IaasVm{
		{Name: "node2", IPs: []string{"192.168.100.97"}},
		{Name: "cnode1", IPs: []string{"10.234.234.38"}},
		{Name: "cnode2", IPs: []string{"10.234.234.93"}},
	}

	errs := AddNodes(vms)
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
