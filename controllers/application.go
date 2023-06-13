package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/astaxie/beego"

	"emcontroller/models"
)

type ApplicationController struct {
	beego.Controller
}

func (c *ApplicationController) Get() {
	appList, err := models.ListApplications()
	if err != nil {
		beego.Error(fmt.Sprintf("ListApplications error: %s", err.Error()))
	}
	c.Data["applicationList"] = appList
	c.TplName = "application.tpl"
}

// DeleteApp delete the deployment and service of the application
// test command:
// curl -i -X DELETE http://localhost:20000/application/test
func (c *ApplicationController) DeleteApp() {
	appName := c.Ctx.Input.Param(":appName")

	err, statusCode := models.DeleteApplication(appName)
	if err != nil {
		beego.Error(err)
		c.Ctx.ResponseWriter.WriteHeader(statusCode)
		c.Ctx.WriteString(err.Error())
		return
	}

	c.Ctx.ResponseWriter.WriteHeader(statusCode)
}

// test command:
// curl -i -X GET http://localhost:20000/application/test
func (c *ApplicationController) GetApp() {
	appName := c.Ctx.Input.Param(":appName")

	outApp, err, statusCode := models.GetApplication(appName)
	if err != nil {
		beego.Error(err)
		c.Ctx.ResponseWriter.WriteHeader(statusCode)
		c.Ctx.WriteString(err.Error())
		return
	}

	c.Ctx.Output.Status = http.StatusOK
	c.Data["json"] = outApp
	c.ServeJSON()
}

func (c *ApplicationController) NewApplication() {
	mode := c.GetString("mode")
	beego.Info("New application mode:", mode)
	// Basic mode can cover most scenarios. Advanced mode support more configurations.
	var tplname string
	switch mode {
	case "basic":
		tplname = "newApplicationBasic.tpl"
	case "advanced":
		tplname = "newApplicationAdvanced.tpl"
	default:
		tplname = "newApplicationBasic.tpl"
	}
	c.TplName = tplname
}

func (c *ApplicationController) DoNewApplication() {
	contentType := c.Ctx.Request.Header.Get("Content-Type")
	beego.Info(fmt.Sprintf("The header \"Content-Type\" is [%s]", contentType))

	switch {
	case strings.Contains(strings.ToLower(contentType), JsonContentType):
		beego.Info(fmt.Sprintf("The input body should be json"))
		c.DoNewAppJson()
	default:
		beego.Info(fmt.Sprintf("The input body should be form"))
		c.DoNewAppForm()
	}
}

// Used for front end request, input is form
func (c *ApplicationController) DoNewAppForm() {
	var app models.K8sApp

	appName := c.GetString("name")
	if appName == "" { // in basic mode, appName is containerName
		beego.Info("basic new application mode, set app name as container name")
		appName = c.GetString("container0Name")
	}
	replicas, err := c.GetInt32("replicas")
	if err != nil {
		outErr := fmt.Errorf("Get replicas error: %w", err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
		c.Data["errorMessage"] = outErr.Error()
		c.TplName = "error.tpl"
		return
	}

	nodeName := c.GetString("nodeName")

	// read node selectors
	var nodeSelector map[string]string = make(map[string]string)
	nodeSelectorNum, err := c.GetInt("nodeSelectorNumber")
	if err != nil {
		outErr := fmt.Errorf("Get nodeSelectorNumber error: %w", err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
		c.Data["errorMessage"] = outErr.Error()
		c.TplName = "error.tpl"
		return
	}
	for i := 0; i < nodeSelectorNum; i++ {
		thisKey := c.GetString(fmt.Sprintf("nodeSelector%dKey", i))
		thisValue := c.GetString(fmt.Sprintf("nodeSelector%dValue", i))
		nodeSelector[thisKey] = thisValue
	}

	// networkType have 2 options: "container" and "host"
	var hostNetwork bool = c.GetString("networkType") == "host"

	containerNum, err := c.GetInt("containerNumber")
	if err != nil {
		outErr := fmt.Errorf("Get containerNumber error: %w", err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
		c.Data["errorMessage"] = outErr.Error()
		c.TplName = "error.tpl"
		return
	}
	beego.Info(fmt.Sprintf("Application [%s] has [%d] pods. Each pod has [%d] containers.", appName, replicas, containerNum))

	app.Name = appName
	app.Replicas = replicas
	app.NodeName = nodeName
	app.NodeSelector = nodeSelector
	app.HostNetwork = hostNetwork
	app.Containers = make([]models.K8sContainer, containerNum, containerNum)

	for i := 0; i < containerNum; i++ {
		var thisContainer models.K8sContainer
		thisContainer.Name = c.GetString(fmt.Sprintf("container%dName", i))
		thisContainer.Image = c.GetString(fmt.Sprintf("container%dImage", i))
		thisContainer.Resources.Requests.Memory = c.GetString(fmt.Sprintf("container%dRequestMemory", i))
		thisContainer.Resources.Requests.CPU = c.GetString(fmt.Sprintf("container%dRequestCPU", i))
		thisContainer.Resources.Requests.Storage = c.GetString(fmt.Sprintf("container%dRequestEphemeralStorage", i))
		thisContainer.Resources.Limits.Memory = c.GetString(fmt.Sprintf("container%dLimitMemory", i))
		thisContainer.Resources.Limits.CPU = c.GetString(fmt.Sprintf("container%dLimitCPU", i))
		thisContainer.Resources.Limits.Storage = c.GetString(fmt.Sprintf("container%dLimitEphemeralStorage", i))
		thisContainer.WorkDir = c.GetString(fmt.Sprintf("container%dWorkdir", i))

		CommandNum, err := c.GetInt(fmt.Sprintf("container%dCommandNumber", i))
		if err != nil {
			outErr := fmt.Errorf("Get command Number error: %w", err)
			beego.Error(outErr)
			c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
			c.Data["errorMessage"] = outErr.Error()
			c.TplName = "error.tpl"
			return
		}
		beego.Info(fmt.Sprintf("Container [%d] has [%d] commands.", i, CommandNum))

		ArgNum, err := c.GetInt(fmt.Sprintf("container%dArgNumber", i))
		if err != nil {
			outErr := fmt.Errorf("Get Arg Number error: %w", err)
			beego.Error(outErr)
			c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
			c.Data["errorMessage"] = outErr.Error()
			c.TplName = "error.tpl"
			return
		}
		beego.Info(fmt.Sprintf("Container [%d] has [%d] args.", i, ArgNum))

		envNum, err := c.GetInt(fmt.Sprintf("container%dEnvNumber", i))
		if err != nil {
			outErr := fmt.Errorf("Get environment variables Number error: %w", err)
			beego.Error(outErr)
			c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
			c.Data["errorMessage"] = outErr.Error()
			c.TplName = "error.tpl"
			return
		}
		beego.Info(fmt.Sprintf("Container [%d] has [%d] environment variables.", i, envNum))

		mountNum, err := c.GetInt(fmt.Sprintf("container%dMountNumber", i))
		if err != nil {
			outErr := fmt.Errorf("Get mount Number error: %w", err)
			beego.Error(outErr)
			c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
			c.Data["errorMessage"] = outErr.Error()
			c.TplName = "error.tpl"
			return
		}
		beego.Info(fmt.Sprintf("Container [%d] has [%d] mount items.", i, mountNum))

		PortNum, err := c.GetInt(fmt.Sprintf("container%dPortNumber", i))
		if err != nil {
			outErr := fmt.Errorf("Get Port Number error: %w", err)
			beego.Error(outErr)
			c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
			c.Data["errorMessage"] = outErr.Error()
			c.TplName = "error.tpl"
			return
		}
		beego.Info(fmt.Sprintf("Container [%d] has [%d] ports.", i, PortNum))

		// get commands
		for j := 0; j < CommandNum; j++ {
			thisCommand := c.GetString(fmt.Sprintf("container%dCommand%d", i, j))
			beego.Info(fmt.Sprintf("Container [%d], Command [%d]: [%s].", i, j, thisCommand))
			thisContainer.Commands = append(thisContainer.Commands, thisCommand)
		}

		// get args
		for j := 0; j < ArgNum; j++ {
			thisArg := c.GetString(fmt.Sprintf("container%dArg%d", i, j))
			beego.Info(fmt.Sprintf("Container [%d], Arg [%d]: [%s].", i, j, thisArg))
			thisContainer.Args = append(thisContainer.Args, thisArg)
		}

		// get environment variables
		for j := 0; j < envNum; j++ {
			thisEnvName := c.GetString(fmt.Sprintf("container%dEnv%dName", i, j))
			thisEnvValue := c.GetString(fmt.Sprintf("container%dEnv%dValue", i, j))
			beego.Info(fmt.Sprintf("Container [%d], Env [%d]: [%s=%s].", i, j, thisEnvName, thisEnvValue))
			thisContainer.Env = append(thisContainer.Env, models.K8sEnv{
				Name:  thisEnvName,
				Value: thisEnvValue,
			})
		}

		// get mount items
		for j := 0; j < mountNum; j++ {
			thisVMPath := c.GetString(fmt.Sprintf("container%dMount%dVM", i, j))
			thisContainerPath := c.GetString(fmt.Sprintf("container%dMount%dContainer", i, j))
			beego.Info(fmt.Sprintf("Container [%d], mount [%d]: VM Path [%s], Container Path [%s].", i, j, thisVMPath, thisContainerPath))
			thisContainer.Mounts = append(thisContainer.Mounts, models.K8sMount{
				VmPath:        thisVMPath,
				ContainerPath: thisContainerPath,
			})
		}

		// get ports
		for j := 0; j < PortNum; j++ {
			var onePort models.PortInfo = models.PortInfo{
				Name:        c.GetString(fmt.Sprintf("container%dPort%dName", i, j)),
				Protocol:    c.GetString(fmt.Sprintf("container%dPort%dProtocol", i, j)),
				ServicePort: c.GetString(fmt.Sprintf("container%dPort%dServicePort", i, j)),
				NodePort:    c.GetString(fmt.Sprintf("container%dPort%dNodePort", i, j)),
			}
			onePort.ContainerPort, err = c.GetInt(fmt.Sprintf("container%dPort%dContainerPort", i, j))
			if err != nil {
				outErr := fmt.Errorf("Get ContainerPort error: %w", err)
				beego.Error(outErr)
				c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
				c.Data["errorMessage"] = outErr.Error()
				c.TplName = "error.tpl"
				return
			}
			beego.Info(fmt.Sprintf("Container [%d], Port [%d]: [%+v].", i, j, onePort))
			thisContainer.Ports = append(thisContainer.Ports, onePort)
		}

		app.Containers[i] = thisContainer
	}

	appJson, err := json.Marshal(app)
	if err != nil {
		outErr := fmt.Errorf("json Marshal this: %v, error: %w", app, err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
		c.Data["errorMessage"] = outErr.Error()
		c.TplName = "error.tpl"
		return
	}
	beego.Info(fmt.Sprintf("App json is\n%s", string(appJson)))

	// Use the parsed app to create an application
	if err := models.CreateApplication(app); err != nil {
		outErr := fmt.Errorf("Create application %+v, error: %w", app, err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
		c.Data["errorMessage"] = outErr.Error()
		c.TplName = "error.tpl"
		return
	}

	c.TplName = "newAppSuccess.tpl"
}

// Used for json request, input is json
// test command:
// curl -i -X POST -H Content-Type:application/json -d '{"priority":-1,"autoScheduled":true,"name":"test","replicas":2,"hostNetwork":true,"tolerations":[{"key":"mcm","operator":"Equal","effect": "NoSchedule","value": "net-test"}],"nodeName":"node1","nodeSelector":{"lnginx":"isnginx","lnginx2":"isnginx2"},"containers":[{"name":"printtime","image":"172.27.15.31:5000/printtime:v1","workDir":"/printtime","resources":{"limits":{"memory":"30Mi","cpu":"200m","storage":"2Gi"},"requests":{"memory":"20Mi","cpu":"100m","storage":"1Gi"}},"commands":["bash"],"args":["-c","python3 -u main.py > $LOGFILE"],"env":[{"name":"PARAMETER1","value":"testRenderenv1"},{"name":"LOGFILE","value":"/tmp/234/printtime.log"}],"mounts":[{"vmPath":"/tmp/asdff","containerPath":"/tmp/234"},{"vmPath":"/tmp/uyyyy","containerPath":"/tmp/2345"}],"ports":null},{"name":"nginx","image":"172.27.15.31:5000/nginx:1.17.1","workDir":"","resources":{"limits":{"memory":"","cpu":"","storage":""},"requests":{"memory":"","cpu":"","storage":""}},"commands":null,"args":null,"env":null,"mounts":null,"ports":[{"containerPort":80,"name":"fsd","protocol":"tcp","servicePort":"80","nodePort":"30001"}]},{"name":"ubuntu","image":"172.27.15.31:5000/ubuntu:latest","workDir":"","resources":{"limits":{"memory":"","cpu":"","storage":""},"requests":{"memory":"","cpu":"","storage":""}},"commands":["bash","-c","while true;do sleep 10;done"],"args":null,"env":[{"name":"asfasf","value":"asfasf"},{"name":"asdfsdf","value":"sfsdf"}],"mounts":[{"vmPath":"/tmp/asdff","containerPath":"/tmp/log"}],"ports":null}]}' http://localhost:20000/doNewApplication
func (c *ApplicationController) DoNewAppJson() {
	var app models.K8sApp
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &app); err != nil {
		outErr := fmt.Errorf("json.Unmarshal the application in RequestBody, error: %w", err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		//c.Ctx.WriteString(outErr.Error())
		if result, err := c.Ctx.ResponseWriter.Write([]byte(outErr.Error())); err != nil {
			beego.Error("Write Error to response, error: %s, result: %d", err.Error(), result)
		}
		return
	}

	beego.Info(fmt.Sprintf("From json input, we successfully parsed application [%+v]", app))

	// Use the parsed app to create an application
	if err := models.CreateApplication(app); err != nil {
		outErr := fmt.Errorf("Create application %+v, error: %w", app, err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString(outErr.Error())
		return
	}

	// Here, we wait until the app status becomes running.
	// We only do this behavior for json input, because for the form input, users can check the status on the web
	// And we need to put the application information (including the service port, pod IP, or nodePort IP) in the response body, to let the user know the information.
	beego.Info(fmt.Sprintf("Start to wait for the application [%s] running", app.Name))
	if err := models.WaitForAppRunning(models.WaitForTimeOut, 10, app.Name); err != nil {
		outErr := fmt.Errorf("Wait for application [%s] running, error: %w", app.Name, err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString(outErr.Error())
		return
	}
	beego.Info(fmt.Sprintf("The application [%s] is already running", app.Name))

	outApp, err, statusCode := models.GetApplication(app.Name)
	if err != nil {
		beego.Error(err)
		c.Ctx.ResponseWriter.WriteHeader(statusCode)
		c.Ctx.WriteString(err.Error())
		return
	}

	//c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
	c.Ctx.Output.Status = http.StatusCreated
	c.Data["json"] = outApp
	c.ServeJSON()
}
