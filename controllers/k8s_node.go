package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
		beego.Error(fmt.Sprintf("Delete node [%s] in Kubernetes Cluster error: %s", nodeName, err.Error()))
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
	contentType := c.Ctx.Request.Header.Get("Content-Type")
	beego.Info(fmt.Sprintf("The header \"Content-Type\" is [%s]", contentType))

	switch {
	case strings.Contains(strings.ToLower(contentType), JsonContentType):
		beego.Info(fmt.Sprintf("The input body should be json"))
		c.DoAddNodesJson()
	default:
		beego.Info(fmt.Sprintf("The input body should be form"))
		c.DoAddNodesForm()
	}
}

func (c *K8sNodeController) DoAddNodesForm() {
	nodeNum, err := c.GetInt("newNodeNumber")
	if err != nil {
		beego.Error(fmt.Sprintf("Get newNodeNumber error: %s", err.Error()))
		return
	}
	beego.Info(fmt.Sprintf("%d nodes need to join the Kubernetes cluster", nodeNum))

	// prepare the information of the nodes to add
	vms := make([]models.IaasVm, nodeNum, nodeNum)

	for i := 0; i < nodeNum; i++ {
		vms[i].Name = c.GetString(fmt.Sprintf("node%dName", i))
		vms[i].IPs = append(vms[i].IPs, c.GetString(fmt.Sprintf("node%dIP", i)))
	}

	logContent := "Nodes to add:"
	for i := 0; i < nodeNum; i++ {
		logContent += fmt.Sprintf("\n\r%d. Name: %s\tIP: %v", i+1, vms[i].Name, vms[i].IPs)
	}
	beego.Info(logContent)

	vmsJson, err := json.Marshal(vms)
	if err != nil {
		outErr := fmt.Errorf("json Marshal this: %v, error: %w", vms, err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
		c.Data["errorMessage"] = outErr.Error()
		c.TplName = "error.tpl"
		return
	}
	beego.Info(fmt.Sprintf("VMs json is\n%s", string(vmsJson)))

	// add node
	if errs := models.AddNodes(vms); len(errs) != 0 {
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

// test command:
// curl -i -X POST -H Content-Type:application/json -d '[{"name":"hpe1","ips":["192.168.100.124"]},{"name":"cnode1","ips":["10.234.234.99"]},{"name":"cnode2","ips":["10.234.234.99"]},{"name":"nokia7","ips":["192.168.100.69"]}]' http://localhost:20000/k8sNode/doAdd
func (c *K8sNodeController) DoAddNodesJson() {
	var vms []models.IaasVm
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vms); err != nil {
		outErr := fmt.Errorf("json.Unmarshal the vms in RequestBody, error: %w", err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		//c.Ctx.WriteString(outErr.Error())
		if result, err := c.Ctx.ResponseWriter.Write([]byte(outErr.Error())); err != nil {
			beego.Error("Write Error to response, error: %s, result: %d", err.Error(), result)
		}
		return
	}

	beego.Info(fmt.Sprintf("From json input, we successfully parsed vms [%v]", vms))

	// Use the parsed vms to create VMs
	if errs := models.AddNodes(vms); len(errs) != 0 {
		outErr := models.HandleErrSlice(errs)
		beego.Error(fmt.Sprintf("AddNodes Error: %s", outErr.Error()))
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString(outErr.Error())
		return
	}

	//c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Ctx.Output.Status = http.StatusCreated
	c.Data["json"] = vms
	c.ServeJSON()
}
