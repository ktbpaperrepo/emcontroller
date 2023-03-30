package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"emcontroller/models"
)

type K8sNodeController struct {
	beego.Controller
}

type K8sNodeInfo struct {
	Name   string
	IP     string
	Status string
}

func (c *K8sNodeController) Get() {
	// TODO: This code does not work, I do not know the reason.
	//K8sMasterSelector := labels.NewSelector()
	//K8sMasterReq, err := labels.NewRequirement(models.K8sMasterNodeRole, selection.NotEquals, []string{""})
	//if err != nil {
	//	beego.Error(fmt.Sprintf("Construct Kubernetes Master requirement, error: %s", err.Error()))
	//}
	//K8sMasterSelector.Add(*K8sMasterReq)
	//beego.Info(fmt.Sprintf("List nodes with selector: %v", K8sMasterSelector))
	//beego.Info(fmt.Sprintf("List nodes with selector: %s", K8sMasterSelector.String()))
	//nodes, err := models.ListNodes(metav1.ListOptions{LabelSelector: K8sMasterSelector.String()})

	nodes, err := models.ListNodes(metav1.ListOptions{})
	if err != nil {
		beego.Error(fmt.Sprintf("List Kubernetes nodes, error: %s", err.Error()))
	}

	selectorControlPlane := labels.SelectorFromSet(labels.Set(map[string]string{
		models.K8sMasterNodeRole: "",
	}))

	var k8sNodeList []K8sNodeInfo
	for _, node := range nodes {
		if selectorControlPlane.Matches(labels.Set(node.Labels)) {
			beego.Info(fmt.Sprintf("node %s is a Master node, so we do not show it.", node.Name))
			continue
		}

		k8sNodeList = append(k8sNodeList, K8sNodeInfo{
			Name:   node.Name,
			IP:     models.GetNodeInternalIp(node),
			Status: models.ExtractNodeStatus(node),
		})
	}

	c.Data["k8sNodeList"] = k8sNodeList
	c.TplName = "k8sNode.tpl"
}

// delete a node from the Kubernetes cluster
func (c *K8sNodeController) DeleteNode() {
	nodeName := c.Ctx.Input.Param(":nodeName")
	beego.Info(fmt.Sprintf("Delete node [%s] in Kubernetes Cluster", nodeName))
	if err := models.UninstallNode(nodeName); err != nil {
		beego.Error(fmt.Printf("Delete node [%s] in Kubernetes Cluster error: %s", nodeName, err.Error()))
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	beego.Info(fmt.Sprintf("Successful! Delete node [%s] in Kubernetes Cluster", nodeName))
	c.Ctx.ResponseWriter.WriteHeader(200)
}

func (c *K8sNodeController) AddNodes() {
	c.TplName = "addK8sNodes.tpl"
}

// Add one or several nodes to the Kubernetes cluster
func (c *K8sNodeController) DoAddNodes() {
	nodeNum, err := c.GetInt("newNodeNumber")
	if err != nil {
		beego.Error(fmt.Sprintf("Get newNodeNumber error: %s", err.Error()))
		return
	}
	beego.Info(fmt.Sprintf("%d nodes need to join the Kubernetes cluster", nodeNum))

	// prepare the information of the node to add
	nodeNames := make([]string, nodeNum, nodeNum)
	nodeIPs := make([]string, nodeNum, nodeNum)

	for i := 0; i < nodeNum; i++ {
		nodeNames[i] = c.GetString(fmt.Sprintf("node%dName", i))
		nodeIPs[i] = c.GetString(fmt.Sprintf("node%dIP", i))
	}

	logContent := "Nodes to add:"
	for i := 0; i < nodeNum; i++ {
		logContent += "\n\r" + strconv.Itoa(i+1) + ". Name: " + nodeNames[i] + "\tIP: " + nodeIPs[i]
	}

	beego.Info(logContent)

	// add node
	if errs := models.AddNodes(nodeNames, nodeIPs); len(errs) != 0 {
		sumErr := models.HandleErrSlice(errs)
		beego.Error(fmt.Sprintf("AddNodes Error: %s", sumErr.Error()))
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
		c.Data["errorMessage"] = sumErr.Error()
		c.TplName = "error.tpl"
		return
	} else {
		beego.Info("Successful. Add nodes.")
	}

	c.TplName = "addK8sNodesSuccess.tpl"
}
