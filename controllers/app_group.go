package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/astaxie/beego"

	"emcontroller/auto-schedule/algorithms"
	"emcontroller/auto-schedule/executors"
	"emcontroller/models"
)

type AppGroupController struct {
	beego.Controller
}

func (c *AppGroupController) DoNewAppGroup() {
	// scheduling, migration, and cleanup cannot be done at the same time
	if !algorithms.ScheMu.TryLock() {
		outErr := fmt.Errorf("Another task of Scheduling, Migration or Cleanup is running. Please try later.")
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusLocked)
		if result, err := c.Ctx.ResponseWriter.Write([]byte(outErr.Error())); err != nil {
			beego.Error(fmt.Sprintf("Write Error to response, error: %s, result: %d", err.Error(), result))
		}
		return
	}
	defer algorithms.ScheMu.Unlock()

	contentType := c.Ctx.Request.Header.Get("Content-Type")
	beego.Info(fmt.Sprintf("The header \"Content-Type\" is [%s]", contentType))

	switch {
	case strings.Contains(strings.ToLower(contentType), JsonContentType):
		beego.Info(fmt.Sprintf("The input body should be json"))
		c.DoNewAppGroupJson()
	default:
		beego.Info(fmt.Sprintf("The input body should be form"))
		c.DoNewAppGroupForm()
	}
}

// Used for json request, input is json
// test command:
// curl -i -X POST -H Content-Type:application/json -d '[ { "priority": 2, "autoScheduled": true, "name": "group-printtime", "replicas": 1, "hostNetwork": false, "containers": [ { "name": "printtime", "image": "172.27.15.31:5000/printtime:v1", "workDir": "/printtime", "resources": { "limits": { "memory": "30Mi", "cpu": "1.2", "storage": "2Gi" }, "requests": { "memory": "30Mi", "cpu": "1.2", "storage": "2Gi" } }, "commands": [ "bash" ], "args": [ "-c", "python3 -u main.py > $LOGFILE" ], "env": [ { "name": "PARAMETER1", "value": "testRenderenv1" }, { "name": "LOGFILE", "value": "/tmp/234/printtime.log" } ], "mounts": [ { "vmPath": "/tmp/asdff", "containerPath": "/tmp/234" }, { "vmPath": "/tmp/uyyyy", "containerPath": "/tmp/2345" } ] } ], "dependencies": [ { "appName": "group-nginx" }, { "appName": "group-ubuntu" } ] }, { "priority": 4, "autoScheduled": true, "name": "group-nginx", "replicas": 1, "hostNetwork": true, "containers": [ { "name": "nginx", "image": "172.27.15.31:5000/nginx:1.17.1", "workDir": "", "resources": { "limits": { "memory": "1024Mi", "cpu": "2.1", "storage": "20Gi" }, "requests": { "memory": "1024Mi", "cpu": "2.1", "storage": "20Gi" } }, "ports": [ { "containerPort": 80, "name": "fsd", "protocol": "tcp", "servicePort": "80", "nodePort": "30001" } ] } ], "dependencies": [ { "appName": "group-ubuntu" } ] }, { "priority": 4, "autoScheduled": true, "name": "group-ubuntu", "replicas": 1, "hostNetwork": true, "containers": [ { "name": "ubuntu", "image": "172.27.15.31:5000/ubuntu:latest", "workDir": "", "resources": { "limits": { "memory": "512Mi", "cpu": "2.1", "storage": "20Gi" }, "requests": { "memory": "512Mi", "cpu": "2.1", "storage": "20Gi" } }, "commands": [ "bash", "-c", "while true;do sleep 10;done" ], "args": null, "env": [ { "name": "asfasf", "value": "asfasf" }, { "name": "asdfsdf", "value": "sfsdf" } ], "mounts": [ { "vmPath": "/tmp/asdff", "containerPath": "/tmp/log" } ], "ports": null } ], "dependencies": [] } ]' http://localhost:20000/doNewAppGroup
func (c *AppGroupController) DoNewAppGroupJson() {
	var apps []models.K8sApp
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &apps); err != nil {
		outErr := fmt.Errorf("json.Unmarshal the applications in RequestBody, error: %w", err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		if result, err := c.Ctx.ResponseWriter.Write([]byte(outErr.Error())); err != nil {
			beego.Error(fmt.Sprintf("Write Error to response, error: %s, result: %d", err.Error(), result))
		}
		return
	}

	beego.Info(fmt.Sprintf("From json input, we successfully parsed applications [%+v]", apps))

	outApps, err, statusCode := executors.CreateAutoScheduleApps(apps)
	if err != nil {
		outErr := fmt.Errorf("executors.CreateAutoScheduleApps(apps), error: %w", err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(statusCode)
		if result, err := c.Ctx.ResponseWriter.Write([]byte(outErr.Error())); err != nil {
			beego.Error(fmt.Sprintf("Write Error to response, error: %s, result: %d", err.Error(), result))
		}
		return
	}

	c.Ctx.Output.Status = http.StatusCreated
	c.Data["json"] = outApps
	c.ServeJSON()
}

func (c *AppGroupController) DoNewAppGroupForm() {
	outErr := fmt.Errorf("Please set the \"Content-Type\" as \"%s\", because the functions to handle other content types have not been implemented.", JsonContentType)
	beego.Error(outErr)
	c.Ctx.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)
	if result, err := c.Ctx.ResponseWriter.Write([]byte(outErr.Error())); err != nil {
		beego.Error(fmt.Sprintf("Write Error to response, error: %s, result: %d", err.Error(), result))
	}
	return
}
