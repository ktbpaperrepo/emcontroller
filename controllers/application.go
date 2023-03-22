package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"emcontroller/models"
)

// PortInfo can store the port information from the web form
type PortInfo struct {
	ContainerPort int
	Name          string
	Protocol      string
	ServicePort   string
	NodePort      string
}

type ApplicationController struct {
	beego.Controller
}

type AppInfo struct {
	AppName    string
	SvcName    string
	DeployName string
	ClusterIP  string
	NodePortIP string
	SvcPort    []string
	NodePort   []string
	Status     string
}

// a method to check whether the application is running
func appRunning(app v1.Deployment) bool {
	if *app.Spec.Replicas != app.Status.Replicas {
		return false
	}
	if *app.Spec.Replicas != app.Status.UpdatedReplicas {
		return false
	}
	if *app.Spec.Replicas != app.Status.ReadyReplicas {
		return false
	}
	if *app.Spec.Replicas != app.Status.AvailableReplicas {
		return false
	}
	return true
}

// We list all pods of this deployment, and take the node IP of a pod as the NodePortIP.
func getNodePortIP(app v1.Deployment) string {
	selector, err := metav1.LabelSelectorAsSelector(app.Spec.Selector)
	if err != nil {
		beego.Error(fmt.Sprintf("Error, get the selector of deployment %s/%s, error: %s", app.Namespace, app.Name, err.Error()))
		return ""
	}
	stringSelector := selector.String()
	beego.Info(fmt.Sprintf("List pods in namespace %s, with selector: %s", app.Namespace, stringSelector))
	pods, err := models.ListPods(app.Namespace, metav1.ListOptions{LabelSelector: stringSelector})
	if err != nil {
		beego.Error(fmt.Sprintf("Error, List pods in namespace %s, with selector: %s, error: %s", app.Namespace, stringSelector, err.Error()))
		return ""
	}
	if len(pods) == 0 {
		beego.Info(fmt.Sprintf("There are no pods in namespace %s that can mapping: %s", app.Namespace, stringSelector))
	}
	return pods[0].Status.HostIP
}

func (c *ApplicationController) Get() {
	applications, err := models.ListDeployment(models.KubernetesNamespace)
	if err != nil {
		beego.Error(fmt.Sprintf("error: %s", err.Error()))
		c.Data["applicationList"] = []AppInfo{}
	}
	var appList []AppInfo
	for _, app := range applications {
		var thisApp AppInfo
		var appName, svcName string
		appName = strings.TrimSuffix(app.Name, models.DeploymentSuffix)
		svcName = appName + models.ServiceSuffix

		thisApp.AppName = appName
		thisApp.SvcName = svcName
		thisApp.DeployName = app.Name

		// set the status of this app
		if appRunning(app) {
			thisApp.Status = "Stable Running"
		} else {
			thisApp.Status = "Not Yet Stable"
		}

		svc, _ := models.GetService(models.KubernetesNamespace, svcName)
		if svc != nil {
			thisApp.ClusterIP = svc.Spec.ClusterIP
			if svc.Spec.Type == apiv1.ServiceTypeNodePort {
				thisApp.NodePortIP = getNodePortIP(app)
			} else {
				thisApp.NodePortIP = ""
			}
			for _, port := range svc.Spec.Ports {
				thisApp.SvcPort = append(thisApp.SvcPort, strconv.FormatInt(int64(port.Port), 10))
				thisApp.NodePort = append(thisApp.NodePort, strconv.FormatInt(int64(port.NodePort), 10))
			}
		} else {
			thisApp.ClusterIP = ""
			thisApp.NodePortIP = ""
			thisApp.SvcPort = []string{}
			thisApp.NodePort = []string{}
		}

		appList = append(appList, thisApp)
	}

	c.Data["applicationList"] = appList
	c.TplName = "application.tpl"
}

// DeleteApp delete the deployment and service of the application
func (c *ApplicationController) DeleteApp() {
	appName := c.Ctx.Input.Param(":appName")
	deployName := appName + models.DeploymentSuffix
	svcName := appName + models.ServiceSuffix

	beego.Info(fmt.Sprintf("Delete deployment [%s/%s]", models.KubernetesNamespace, deployName))
	if err := models.DeleteDeployment(models.KubernetesNamespace, deployName); err != nil {
		beego.Error(fmt.Printf("Delete deployment [%s/%s] error: %s", models.KubernetesNamespace, deployName, err.Error()))
		return
	}
	beego.Info(fmt.Sprintf("Successful! Delete deployment [%s/%s]", models.KubernetesNamespace, deployName))

	beego.Info(fmt.Sprintf("Delete service [%s/%s]", models.KubernetesNamespace, svcName))
	if err := models.DeleteService(models.KubernetesNamespace, svcName); err != nil {
		beego.Error(fmt.Printf("Delete deployment [%s/%s] error: %s", models.KubernetesNamespace, svcName, err.Error()))
		return
	}
	beego.Info(fmt.Sprintf("Successful! Delete service [%s/%s]", models.KubernetesNamespace, svcName))

	c.Ctx.ResponseWriter.WriteHeader(200)
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

// for an application, we need to create a Deployment and a Service for it
func (c *ApplicationController) DoNewApplication() {
	appName := c.GetString("name")
	if appName == "" { // in basic mode, appName is containerName
		beego.Info("basic new application mode, set app name as container name")
		appName = c.GetString("container0Name")
	}
	replicas, err := c.GetInt32("replicas")
	if err != nil {
		beego.Error(fmt.Sprintf("Get replicas error: %s", err.Error()))
		return
	}
	containerNum, err := c.GetInt("containerNumber")
	if err != nil {
		beego.Error(fmt.Sprintf("Get containerNumber error: %s", err.Error()))
		return
	}
	beego.Info(fmt.Sprintf("Application [%s] has [%d] pods. Each pod has [%d] containers.", appName, replicas, containerNum))

	// Kubernetes labels of the pods of this application
	labels := map[string]string{
		"app": appName,
	}

	// the ports in the service of this application
	var servicePorts []apiv1.ServicePort

	var hasNodePort bool

	// make volume configuration in a pod
	beego.Info("make volume configuration in a pod")
	var volumeP2N map[string]string = make(map[string]string) // a map from VM paths to volume names
	var volumes []apiv1.Volume                                // put this slice into deployment template
	for i := 0; i < containerNum; i++ {
		mountNum, err := c.GetInt(fmt.Sprintf("container%dMountNumber", i))
		if err != nil {
			beego.Error(fmt.Sprintf("make volume configuration, Get mount Number error: %s", err.Error()))
			return
		}
		beego.Info(fmt.Sprintf("make volume configuration, Container [%d] has [%d] mount items.", i, mountNum))

		for j := 0; j < mountNum; j++ {
			thisVMPath := c.GetString(fmt.Sprintf("container%dMount%dVM", i, j))
			if _, exist := volumeP2N[thisVMPath]; !exist {
				thisVolumeName := "volume" + strconv.Itoa(len(volumeP2N))
				beego.Info(fmt.Sprintf("add volume name: [%s], VM path: [%s]", thisVolumeName, thisVMPath))
				volumeP2N[thisVMPath] = thisVolumeName

				var hostPathType apiv1.HostPathType = apiv1.HostPathDirectoryOrCreate
				volumes = append(volumes, apiv1.Volume{
					Name: thisVolumeName,
					VolumeSource: apiv1.VolumeSource{
						HostPath: &apiv1.HostPathVolumeSource{
							Path: thisVMPath,
							Type: &hostPathType,
						},
					},
				})
			}
		}
	}

	// get the configuration of every container
	beego.Info("make containers configuration")
	var containers []apiv1.Container
	for i := 0; i < containerNum; i++ {
		var thisContainer apiv1.Container = apiv1.Container{
			ImagePullPolicy: apiv1.PullIfNotPresent,
		}
		beego.Info(fmt.Sprintf("Get the configuration of container %d", i))
		thisContainer.Name = c.GetString(fmt.Sprintf("container%dName", i))
		thisContainer.Image = c.GetString(fmt.Sprintf("container%dImage", i))
		thisContainer.Resources = getContainerResources(c, i)

		CommandNum, err := c.GetInt(fmt.Sprintf("container%dCommandNumber", i))
		if err != nil {
			beego.Error(fmt.Sprintf("Get command Number error: %s", err.Error()))
			return
		}
		beego.Info(fmt.Sprintf("Container [%d] has [%d] commands.", i, CommandNum))

		ArgNum, err := c.GetInt(fmt.Sprintf("container%dArgNumber", i))
		if err != nil {
			beego.Error(fmt.Sprintf("Get Arg Number error: %s", err.Error()))
			return
		}
		beego.Info(fmt.Sprintf("Container [%d] has [%d] args.", i, ArgNum))

		envNum, err := c.GetInt(fmt.Sprintf("container%dEnvNumber", i))
		if err != nil {
			beego.Error(fmt.Sprintf("Get environment variables Number error: %s", err.Error()))
			return
		}
		beego.Info(fmt.Sprintf("Container [%d] has [%d] environment variables.", i, envNum))

		mountNum, err := c.GetInt(fmt.Sprintf("container%dMountNumber", i))
		if err != nil {
			beego.Error(fmt.Sprintf("Get mount Number error: %s", err.Error()))
			return
		}
		beego.Info(fmt.Sprintf("Container [%d] has [%d] mount items.", i, mountNum))

		PortNum, err := c.GetInt(fmt.Sprintf("container%dPortNumber", i))
		if err != nil {
			beego.Error(fmt.Sprintf("Get Port Number error: %s", err.Error()))
			return
		}
		beego.Info(fmt.Sprintf("Container [%d] has [%d] ports.", i, PortNum))

		// get commands
		for j := 0; j < CommandNum; j++ {
			thisCommand := c.GetString(fmt.Sprintf("container%dCommand%d", i, j))
			beego.Info(fmt.Sprintf("Container [%d], Command [%d]: [%s].", i, j, thisCommand))
			thisContainer.Command = append(thisContainer.Command, thisCommand)
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
			thisContainer.Env = append(thisContainer.Env, apiv1.EnvVar{
				Name:  thisEnvName,
				Value: thisEnvValue,
			})
		}

		// get mount items
		for j := 0; j < mountNum; j++ {
			thisVMPath := c.GetString(fmt.Sprintf("container%dMount%dVM", i, j))
			thisContainerPath := c.GetString(fmt.Sprintf("container%dMount%dContainer", i, j))
			volumeName, found := volumeP2N[thisVMPath]
			if !found {
				beego.Error(fmt.Sprintf("Container [%d], mount [%d]: VM Path [%s], Container Path [%s], cannot found volume name.", i, j, thisVMPath, thisContainerPath))
			} else {
				beego.Info(fmt.Sprintf("Container [%d], mount [%d]: VM Path [%s], Container Path [%s], volume name [%s].", i, j, thisVMPath, thisContainerPath, volumeName))
			}

			thisContainer.VolumeMounts = append(thisContainer.VolumeMounts, apiv1.VolumeMount{
				Name:      volumeName,
				MountPath: thisContainerPath,
			})
		}

		// get ports
		for j := 0; j < PortNum; j++ {
			var onePort PortInfo = PortInfo{
				Name:        c.GetString(fmt.Sprintf("container%dPort%dName", i, j)),
				Protocol:    c.GetString(fmt.Sprintf("container%dPort%dProtocol", i, j)),
				ServicePort: c.GetString(fmt.Sprintf("container%dPort%dServicePort", i, j)),
				NodePort:    c.GetString(fmt.Sprintf("container%dPort%dNodePort", i, j)),
			}
			onePort.ContainerPort, err = c.GetInt(fmt.Sprintf("container%dPort%dContainerPort", i, j))
			if err != nil {
				beego.Error(fmt.Sprintf("Get ContainerPort error: %s", err.Error()))
				return
			}

			beego.Info(fmt.Sprintf("Container [%d], Port [%d]: [%#v].", i, j, onePort))

			// put port information into the container of the deployment
			thisContainer.Ports = append(thisContainer.Ports, apiv1.ContainerPort{
				ContainerPort: int32(onePort.ContainerPort),
				Name:          onePort.Name,
				Protocol:      apiv1.Protocol(strings.ToUpper(onePort.Protocol)),
			})

			if len(onePort.ServicePort)+len(onePort.NodePort) > 0 {
				thisServicePort := apiv1.ServicePort{
					Name:       onePort.Name,
					Protocol:   apiv1.Protocol(strings.ToUpper(onePort.Protocol)),
					TargetPort: intstr.FromInt(onePort.ContainerPort),
				}

				// set service port if it exists
				if len(onePort.ServicePort) > 0 {
					sp, err := strconv.Atoi(onePort.ServicePort)
					if err != nil {
						beego.Error(fmt.Sprintf("Atoi ServicePort error: %s", err.Error()))
						return
					}
					thisServicePort.Port = int32(sp)
				}

				// set node port if it exists
				if len(onePort.NodePort) > 0 {
					np, err := strconv.Atoi(onePort.NodePort)
					if err != nil {
						beego.Error(fmt.Sprintf("Atoi NodePort error: %s", err.Error()))
						return
					}
					hasNodePort = true
					thisServicePort.NodePort = int32(np)
				}

				// put the port information into the service
				servicePorts = append(servicePorts, thisServicePort)
			}
		}

		// add this container
		containers = append(containers, thisContainer)
	}

	maxUnavailable := intstr.FromInt(1)
	maxSurge := intstr.FromInt(1)

	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName + models.DeploymentSuffix,
			Namespace: models.KubernetesNamespace,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Strategy: v1.DeploymentStrategy{
				Type: v1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &v1.RollingUpdateDeployment{
					MaxUnavailable: &maxUnavailable,
					MaxSurge:       &maxSurge,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: apiv1.PodSpec{
					Containers: containers,
					Volumes:    volumes,
					Affinity: &apiv1.Affinity{
						PodAntiAffinity: &apiv1.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []apiv1.PodAffinityTerm{
								apiv1.PodAffinityTerm{
									TopologyKey: apiv1.LabelHostname,
									LabelSelector: &metav1.LabelSelector{
										MatchLabels: labels,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	beego.Info(fmt.Sprintf("Create deployment [%#v]", deployment))
	beego.Info(fmt.Sprintf(""))
	deploymentJson, err := json.Marshal(deployment)
	if err != nil {
		beego.Error(fmt.Sprintf("Json Marshal error: %s", err.Error()))
	}
	beego.Info(fmt.Sprintf("Create deployment (json) [%s]", string(deploymentJson)))

	createdDeployment, err := models.CreateDeployment(deployment)
	if err != nil {
		beego.Error(fmt.Printf("Create deployment [%#v] error: %s", deployment, err.Error()))
		return
	}
	beego.Info(fmt.Sprintf("Deployment %s/%s created successful.", createdDeployment.Namespace, createdDeployment.Name))

	// service of this application
	if len(servicePorts) != 0 {
		service := &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      appName + models.ServiceSuffix,
				Namespace: models.KubernetesNamespace,
			},
			Spec: apiv1.ServiceSpec{
				Selector: labels,
				Type:     apiv1.ServiceTypeClusterIP,
				Ports:    servicePorts,
			},
		}
		if hasNodePort {
			service.Spec.Type = apiv1.ServiceTypeNodePort
		}

		beego.Info(fmt.Sprintf("Create service [%#v]", service))
		beego.Info(fmt.Sprintf(""))
		serviceJson, err := json.Marshal(service)
		if err != nil {
			beego.Error(fmt.Sprintf("Json Marshal error: %s", err.Error()))
		}
		beego.Info(fmt.Sprintf("Create service (json) [%s]", string(serviceJson)))

		createdService, err := models.CreateService(service)
		if err != nil {
			beego.Error(fmt.Printf("Create service [%#v] error: %s", service, err.Error()))
			return
		}
		beego.Info(fmt.Sprintf("Service %s/%s created successful.", createdService.Namespace, createdService.Name))
	}
	c.TplName = "newAppSuccess.tpl"
}

func getContainerResources(c *ApplicationController, containerIndex int) apiv1.ResourceRequirements {
	var resources apiv1.ResourceRequirements = apiv1.ResourceRequirements{
		Requests: make(apiv1.ResourceList),
		Limits:   make(apiv1.ResourceList),
	}

	var resourcesToGet map[apiv1.ResourceName]string = map[apiv1.ResourceName]string{
		apiv1.ResourceCPU:              "container%d%sCPU",
		apiv1.ResourceMemory:           "container%d%sMemory",
		apiv1.ResourceEphemeralStorage: "container%d%sEphemeralStorage",
	}

	var getResFunc func(resType string, resList *apiv1.ResourceList) = func(resType string, resList *apiv1.ResourceList) {
		for k, v := range resourcesToGet {
			var htmlInputName string = fmt.Sprintf(v, containerIndex, resType)
			var amount string = c.GetString(htmlInputName)
			beego.Info(fmt.Sprintf("The value of html input [%s] is [%s]", htmlInputName, amount))
			if len(amount) > 0 {
				(*resList)[k] = resource.MustParse(amount)
			}
		}
	}

	// get request resources
	getResFunc("Request", &(resources.Requests))
	// get limit resources
	getResFunc("Limit", &(resources.Limits))

	return resources
}
