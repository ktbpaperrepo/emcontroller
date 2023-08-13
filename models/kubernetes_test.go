package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// go test $GOPATH/src/emcontroller/models/ --run TestListDeployment -v -count=1
func TestListDeployment(t *testing.T) {
	InitSomeThing()
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

func TestGetDeployment(t *testing.T) {
	InitSomeThing()
	testCases := []struct {
		name          string
		namespace     string
		deployName    string
		expectedError error
	}{
		{
			name:          "case1",
			namespace:     KubernetesNamespace,
			deployName:    "test-deployment",
			expectedError: nil,
		},
		{
			name:          "case2",
			namespace:     "kube-system",
			deployName:    "coredns",
			expectedError: nil,
		},
		{
			name:          "case3",
			namespace:     "kube-system",
			deployName:    "coredns1",
			expectedError: nil,
		},
		{
			name:          "case4",
			namespace:     "kube-system1",
			deployName:    "coredns1",
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Logf("test: %s", testCase.name)
		actualResult, actualError := GetDeployment(testCase.namespace, testCase.deployName)
		if testCase.expectedError == nil {
			assert.NoError(t, actualError, fmt.Sprintf("%s: Error is not expected", testCase.name))
		} else {
			assert.Error(t, actualError, fmt.Sprintf("%s: Error is not expected", testCase.name))
		}
		t.Logf("Deploy: %s/%s, %v", testCase.namespace, testCase.deployName, actualResult)
	}
}

func TestListPodsOnNode(t *testing.T) {
	InitSomeThing()
	pods, err := ListPodsOnNode("", "n4test")
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	}
	for _, pod := range pods {
		t.Logf("%s/%s", pod.Namespace, pod.Name)
	}
}

func TestListNodes(t *testing.T) {
	InitSomeThing()
	nodes, err := ListNodes(metav1.ListOptions{})
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	} else {
		t.Logf("nodes: %v", nodes)
	}
}

func TestListNodesNamePrefix(t *testing.T) {
	InitSomeThing()
	nodes, err := ListNodesNamePrefix("auto-sched-")
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	} else {
		for _, node := range nodes {
			t.Log(node.Name, GetNodeInternalIp(node))
		}
	}
}

func TestGetNode(t *testing.T) {
	InitSomeThing()
	node, err := GetNode("node1", metav1.GetOptions{})
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	}
	t.Logf("node: %v", node)
}

func TestGetJoinCmd(t *testing.T) {
	InitSomeThing()
	joinCmd, err := GetJoinCmd()
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	}
	t.Logf("joinCmd: %s --node-name=test", joinCmd)
}

func TestJoinOneNode(t *testing.T) {
	InitSomeThing()
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
	InitSomeThing()

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
	InitSomeThing()
	err := UninstallNode("node1")
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	}
}

func TestUninstallBatchNodes(t *testing.T) {
	InitSomeThing()

	nodesToDelete, err := ListNodesNamePrefix("auto-sched-")
	if err != nil {
		t.Errorf("get nodes to delete error: %s", err.Error())
	}

	var nodeNamesToDelete []string
	for _, node := range nodesToDelete {
		nodeNamesToDelete = append(nodeNamesToDelete, node.Name)
	}

	t.Log("nodes to uninstall:", nodeNamesToDelete)

	errs := UninstallBatchNodes(nodeNamesToDelete)
	if errs != nil {
		t.Errorf("uninstall nodes error: %s", HandleErrSlice(errs).Error())
	} else {
		t.Logf("Uninstall Kubernetes nodes successfully")
	}
}
