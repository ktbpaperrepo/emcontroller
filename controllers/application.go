package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"emcontroller/models"
)

type ApplicationController struct {
	beego.Controller
}

// used for the input of creating applications, so we need to define the json
type K8sApp struct {
	Name         string            `json:"name"`
	Replicas     int32             `json:"replicas"`
	HostNetwork  bool              `json:"hostNetwork"`
	NodeName     string            `json:"nodeName,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	Containers   []K8sContainer    `json:"containers"`
}

type K8sContainer struct {
	Name      string     `json:"name"`
	Image     string     `json:"image"`
	WorkDir   string     `json:"workDir"`
	Resources K8sResReq  `json:"resources"`
	Commands  []string   `json:"commands"`
	Args      []string   `json:"args"`
	Env       []K8sEnv   `json:"env"`
	Mounts    []K8sMount `json:"mounts"`
	Ports     []PortInfo `json:"ports"`
}

type K8sResReq struct {
	Limits   K8sResList `json:"limits"`
	Requests K8sResList `json:"requests"`
}

type K8sResList struct {
	Memory  string `json:"memory"`
	CPU     string `json:"cpu"`
	Storage string `json:"storage"`
}

// Environment Variables
type K8sEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Mount VM paths into the container
type K8sMount struct {
	VmPath        string `json:"vmPath"`
	ContainerPath string `json:"containerPath"`
}

// PortInfo can store the port information from the web form
type PortInfo struct {
	ContainerPort int    `json:"containerPort"`
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	ServicePort   string `json:"servicePort"`
	NodePort      string `json:"nodePort"`
}

type AppInfo struct {
	AppName       string    `json:"appName"`
	SvcName       string    `json:"svcName"`
	DeployName    string    `json:"deployName"`
	ClusterIP     string    `json:"clusterIP"`
	NodePortIP    string    `json:"nodePortIP"`
	SvcPort       []string  `json:"svcPort"`
	NodePort      []string  `json:"nodePort"`
	ContainerPort []string  `json:"containerPort"`
	Hosts         []PodHost `json:"hosts"`
	Status        string    `json:"status"`
}

type PodHost struct {
	PodIP    string `json:"podIP"`
	HostName string `json:"hostName"`
	HostIP   string `json:"hostIP"`
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

// get all Kubernetes pods of this application
func getAllPods(app v1.Deployment) []apiv1.Pod {
	selector, err := metav1.LabelSelectorAsSelector(app.Spec.Selector)
	if err != nil {
		beego.Error(fmt.Sprintf("Error, get the selector of deployment %s/%s, error: %s", app.Namespace, app.Name, err.Error()))
		return []apiv1.Pod{}
	}
	stringSelector := selector.String()
	beego.Info(fmt.Sprintf("List pods belonging to the app %s/%s, with selector: %s", app.Namespace, app.Name, stringSelector))
	pods, err := models.ListPods(app.Namespace, metav1.ListOptions{LabelSelector: stringSelector})
	if err != nil {
		beego.Error(fmt.Sprintf("Error, List pods in namespace %s, with selector: %s, error: %s", app.Namespace, stringSelector, err.Error()))
		return []apiv1.Pod{}
	}
	return pods
}

// get the host Kubernetes Nodes of all pods of this application
func getHosts(app v1.Deployment, pods []apiv1.Pod) []PodHost {
	var hosts []PodHost
	if len(pods) == 0 {
		beego.Info(fmt.Sprintf("No pods belonging to the app %s/%s are got.", app.Namespace, app.Name))
	}
	for _, pod := range pods {
		var podIP string
		if len(pod.Status.PodIPs) > 0 { // pkg/printers/internalversion/printers.go, func printPod, the kubectl code gets pod IP here
			podIP = pod.Status.PodIPs[0].IP
		} else if len(pod.Status.PodIP) > 0 {
			podIP = pod.Status.PodIP
		}

		if len(pod.Spec.NodeName)+len(pod.Status.HostIP)+len(podIP) == 0 {
			continue
		}
		hosts = append(hosts, PodHost{
			PodIP:    podIP,
			HostName: pod.Spec.NodeName,
			HostIP:   pod.Status.HostIP,
		})
	}

	return hosts
}

// We list all pods of this deployment, and take the node IP of a pod as the NodePortIP.
func getNodePortIP(app v1.Deployment, pods []apiv1.Pod) string {
	if len(pods) == 0 {
		beego.Info(fmt.Sprintf("No pods belonging to the app %s/%s are got.", app.Namespace, app.Name))
		return ""
	}
	return pods[0].Status.HostIP
}

func (c *ApplicationController) Get() {
	appList, err := ListApplications()
	if err != nil {
		beego.Error(fmt.Sprintf("ListApplications error: %s", err.Error()))
	}
	c.Data["applicationList"] = appList
	c.TplName = "application.tpl"
}

func ListApplications() ([]AppInfo, error) {
	applications, err := models.ListDeployment(models.KubernetesNamespace)
	if err != nil {
		beego.Error(fmt.Sprintf("ListDeployment error: %s", err.Error()))
		return []AppInfo{}, err
	}

	var appList []AppInfo

	// the slice in golang is not safe for concurrent read/write
	var appListMu sync.Mutex

	// handle every deployment in parallel
	var wg sync.WaitGroup

	for _, app := range applications {
		wg.Add(1)
		go func(d v1.Deployment) {
			defer wg.Done()

			var thisApp AppInfo
			var appName, svcName string
			appName = strings.TrimSuffix(d.Name, models.DeploymentSuffix)
			svcName = appName + models.ServiceSuffix

			pods := getAllPods(d)

			thisApp.AppName = appName
			thisApp.SvcName = svcName
			thisApp.DeployName = d.Name
			thisApp.Hosts = getHosts(d, pods)

			// set the status of this application
			if appRunning(d) {
				thisApp.Status = RunningStatus
			} else {
				thisApp.Status = NotStableStatus
			}

			svc, err := models.GetService(models.KubernetesNamespace, svcName)
			if err != nil {
				beego.Error(fmt.Sprintf("GetService %s/%s error: %s", models.KubernetesNamespace, svcName, err.Error()))
			}
			if svc != nil {
				thisApp.ClusterIP = svc.Spec.ClusterIP
				if svc.Spec.Type == apiv1.ServiceTypeNodePort {
					thisApp.NodePortIP = getNodePortIP(d, pods)
				} else {
					thisApp.NodePortIP = ""
				}
				for _, port := range svc.Spec.Ports {
					thisApp.SvcPort = append(thisApp.SvcPort, strconv.FormatInt(int64(port.Port), 10))
					thisApp.NodePort = append(thisApp.NodePort, strconv.FormatInt(int64(port.NodePort), 10))
					thisApp.ContainerPort = append(thisApp.ContainerPort, port.TargetPort.String())
				}
			} else {
				thisApp.ClusterIP = ""
				thisApp.NodePortIP = ""
				thisApp.SvcPort = []string{}
				thisApp.NodePort = []string{}
			}

			appListMu.Lock()
			appList = append(appList, thisApp)
			appListMu.Unlock()
		}(app)
	}
	wg.Wait()

	return appList, nil
}

// DeleteApp delete the deployment and service of the application
func (c *ApplicationController) DeleteApp() {
	appName := c.Ctx.Input.Param(":appName")

	err, statusCode := DeleteApplication(appName)
	if err != nil {
		beego.Error(err)
		c.Ctx.ResponseWriter.WriteHeader(statusCode)
		c.Ctx.WriteString(err.Error())
		return
	}

	c.Ctx.ResponseWriter.WriteHeader(statusCode)
}

func DeleteApplication(appName string) (error, int) {
	deployName := appName + models.DeploymentSuffix
	svcName := appName + models.ServiceSuffix

	beego.Info(fmt.Sprintf("Delete deployment [%s/%s]", models.KubernetesNamespace, deployName))
	if err := models.DeleteDeployment(models.KubernetesNamespace, deployName); err != nil {
		outErr := fmt.Errorf("Delete deployment [%s/%s] error: %s", models.KubernetesNamespace, deployName, err.Error())
		beego.Error(outErr)
		return outErr, http.StatusInternalServerError
	}
	beego.Info(fmt.Sprintf("Successful! Delete deployment [%s/%s]", models.KubernetesNamespace, deployName))

	beego.Info(fmt.Sprintf("Delete service [%s/%s]", models.KubernetesNamespace, svcName))
	if err := models.DeleteService(models.KubernetesNamespace, svcName); err != nil {
		outErr := fmt.Errorf("Delete deployment [%s/%s] error: %s", models.KubernetesNamespace, svcName, err.Error())
		beego.Error(outErr)
		return outErr, http.StatusInternalServerError
	}
	beego.Info(fmt.Sprintf("Successful! Delete service [%s/%s]", models.KubernetesNamespace, svcName))
	return nil, http.StatusOK
}

// test command:
// curl -i -X GET http://localhost:20000/application/test
func (c *ApplicationController) GetApp() {
	appName := c.Ctx.Input.Param(":appName")

	outApp, err, statusCode := GetApplication(appName)
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

func GetApplication(appName string) (AppInfo, error, int) {
	deployName := appName + models.DeploymentSuffix
	svcName := appName + models.ServiceSuffix

	deploy, err := models.GetDeployment(models.KubernetesNamespace, deployName)
	if err != nil {
		outErr := fmt.Errorf("Get the deployment of app [%s], error: %w", appName, err)
		beego.Error(outErr)
		return AppInfo{}, outErr, http.StatusInternalServerError
	}
	if deploy == nil {
		outErr := fmt.Errorf("The deployment of app [%s] not found", appName)
		beego.Error(outErr)
		return AppInfo{}, outErr, http.StatusNotFound
	}
	svc, err := models.GetService(models.KubernetesNamespace, svcName)
	if err != nil {
		outErr := fmt.Errorf("Get the service of app [%s], error: %w", appName, err)
		beego.Error(outErr)
		return AppInfo{}, outErr, http.StatusInternalServerError
	}

	var outApp AppInfo

	pods := getAllPods(*deploy)

	outApp.AppName = appName
	outApp.SvcName = svcName
	outApp.DeployName = deploy.Name
	outApp.Hosts = getHosts(*deploy, pods)

	// set the status of this application
	if appRunning(*deploy) {
		outApp.Status = RunningStatus
	} else {
		outApp.Status = NotStableStatus
	}

	if svc != nil {
		outApp.ClusterIP = svc.Spec.ClusterIP
		if svc.Spec.Type == apiv1.ServiceTypeNodePort {
			outApp.NodePortIP = getNodePortIP(*deploy, pods)
		} else {
			outApp.NodePortIP = ""
		}
		for _, port := range svc.Spec.Ports {
			outApp.SvcPort = append(outApp.SvcPort, strconv.FormatInt(int64(port.Port), 10))
			outApp.NodePort = append(outApp.NodePort, strconv.FormatInt(int64(port.NodePort), 10))
			outApp.ContainerPort = append(outApp.ContainerPort, port.TargetPort.String())
		}
	} else {
		outApp.ClusterIP = ""
		outApp.NodePortIP = ""
		outApp.SvcPort = []string{}
		outApp.NodePort = []string{}
	}

	return outApp, nil, http.StatusOK
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

// Deprecated: I split this function into 2 functions: DoNewApplication and CreateApplication
// for an application, we need to create a Deployment and a Service for it
func (c *ApplicationController) OldDoNewApplication() {
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

	// networkType have 2 options: "container" and "host"
	var hostNetwork bool = c.GetString("networkType") == "host"

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

		// If the working directory of the container is configured, we transmit it to Kubernetes.
		workdir := c.GetString(fmt.Sprintf("container%dWorkdir", i))
		if len(workdir) > 0 {
			thisContainer.WorkingDir = workdir
		}

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
					HostNetwork: hostNetwork,
					DNSPolicy:   apiv1.DNSClusterFirstWithHostNet, // without this, pods with HostNetwork cannot access coredns
					Containers:  containers,
					Volumes:     volumes,
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
		beego.Error(fmt.Sprintf("Create deployment [%#v] error: %s", deployment, err.Error()))
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
			beego.Error(fmt.Sprintf("Create service [%#v] error: %s", service, err.Error()))
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
	var app K8sApp

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
	app.Containers = make([]K8sContainer, containerNum, containerNum)

	for i := 0; i < containerNum; i++ {
		var thisContainer K8sContainer
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
			thisContainer.Env = append(thisContainer.Env, K8sEnv{
				Name:  thisEnvName,
				Value: thisEnvValue,
			})
		}

		// get mount items
		for j := 0; j < mountNum; j++ {
			thisVMPath := c.GetString(fmt.Sprintf("container%dMount%dVM", i, j))
			thisContainerPath := c.GetString(fmt.Sprintf("container%dMount%dContainer", i, j))
			beego.Info(fmt.Sprintf("Container [%d], mount [%d]: VM Path [%s], Container Path [%s].", i, j, thisVMPath, thisContainerPath))
			thisContainer.Mounts = append(thisContainer.Mounts, K8sMount{
				VmPath:        thisVMPath,
				ContainerPath: thisContainerPath,
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
				outErr := fmt.Errorf("Get ContainerPort error: %w", err)
				beego.Error(outErr)
				c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
				c.Data["errorMessage"] = outErr.Error()
				c.TplName = "error.tpl"
				return
			}
			beego.Info(fmt.Sprintf("Container [%d], Port [%d]: [%#v].", i, j, onePort))
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
	if err := CreateApplication(app); err != nil {
		outErr := fmt.Errorf("Create application %v, error: %w", app, err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain")
		c.Data["errorMessage"] = outErr.Error()
		c.TplName = "error.tpl"
		return
	}

	c.TplName = "newAppSuccess.tpl"
}

// for an application, we need to create a Deployment and a Service for it
func CreateApplication(app K8sApp) error {

	// Kubernetes labels of the pods of this application
	labels := map[string]string{
		"app": app.Name,
	}

	// the ports in the service of this application
	var servicePorts []apiv1.ServicePort

	var hasNodePort bool

	// make volume configuration in a pod
	beego.Info("make volume configuration in a pod")
	var volumeP2N map[string]string = make(map[string]string) // a map from VM paths to volume names
	var volumes []apiv1.Volume                                // put this slice into deployment template
	for i := 0; i < len(app.Containers); i++ {
		beego.Info(fmt.Sprintf("make volume configuration, Container [%d] has [%d] mount items.", i, len(app.Containers[i].Mounts)))

		for j := 0; j < len(app.Containers[i].Mounts); j++ {
			thisVMPath := app.Containers[i].Mounts[j].VmPath
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
	for i := 0; i < len(app.Containers); i++ {
		var thisContainer apiv1.Container = apiv1.Container{
			ImagePullPolicy: apiv1.PullIfNotPresent,
		}
		beego.Info(fmt.Sprintf("Get the configuration of container %d", i))
		thisContainer.Name = app.Containers[i].Name
		thisContainer.Image = app.Containers[i].Image
		thisContainer.Resources = apiv1.ResourceRequirements{
			Requests: make(apiv1.ResourceList),
			Limits:   make(apiv1.ResourceList),
		}
		if len(app.Containers[i].Resources.Requests.CPU) > 0 {
			thisContainer.Resources.Requests[apiv1.ResourceCPU] = resource.MustParse(app.Containers[i].Resources.Requests.CPU)
		}
		if len(app.Containers[i].Resources.Requests.Memory) > 0 {
			thisContainer.Resources.Requests[apiv1.ResourceMemory] = resource.MustParse(app.Containers[i].Resources.Requests.Memory)
		}
		if len(app.Containers[i].Resources.Requests.Storage) > 0 {
			thisContainer.Resources.Requests[apiv1.ResourceEphemeralStorage] = resource.MustParse(app.Containers[i].Resources.Requests.Storage)
		}
		if len(app.Containers[i].Resources.Limits.CPU) > 0 {
			thisContainer.Resources.Limits[apiv1.ResourceCPU] = resource.MustParse(app.Containers[i].Resources.Limits.CPU)
		}
		if len(app.Containers[i].Resources.Limits.Memory) > 0 {
			thisContainer.Resources.Limits[apiv1.ResourceMemory] = resource.MustParse(app.Containers[i].Resources.Limits.Memory)
		}
		if len(app.Containers[i].Resources.Limits.Storage) > 0 {
			thisContainer.Resources.Limits[apiv1.ResourceEphemeralStorage] = resource.MustParse(app.Containers[i].Resources.Limits.Storage)
		}

		// If the working directory of the container is configured, we transmit it to Kubernetes.
		if len(app.Containers[i].WorkDir) > 0 {
			thisContainer.WorkingDir = app.Containers[i].WorkDir
		}

		beego.Info(fmt.Sprintf("Container [%d] has [%d] commands.", i, len(app.Containers[i].Commands)))
		beego.Info(fmt.Sprintf("Container [%d] has [%d] args.", i, len(app.Containers[i].Args)))
		beego.Info(fmt.Sprintf("Container [%d] has [%d] environment variables.", i, len(app.Containers[i].Env)))
		beego.Info(fmt.Sprintf("Container [%d] has [%d] mount items.", i, len(app.Containers[i].Mounts)))
		beego.Info(fmt.Sprintf("Container [%d] has [%d] ports.", i, len(app.Containers[i].Ports)))

		// get commands
		for j := 0; j < len(app.Containers[i].Commands); j++ {
			beego.Info(fmt.Sprintf("Container [%d], Command [%d]: [%s].", i, j, app.Containers[i].Commands[j]))
			thisContainer.Command = append(thisContainer.Command, app.Containers[i].Commands[j])
		}

		// get args
		for j := 0; j < len(app.Containers[i].Args); j++ {
			beego.Info(fmt.Sprintf("Container [%d], Arg [%d]: [%s].", i, j, app.Containers[i].Args[j]))
			thisContainer.Args = append(thisContainer.Args, app.Containers[i].Args[j])
		}

		// get environment variables
		for j := 0; j < len(app.Containers[i].Env); j++ {
			thisEnvName := app.Containers[i].Env[j].Name
			thisEnvValue := app.Containers[i].Env[j].Value
			beego.Info(fmt.Sprintf("Container [%d], Env [%d]: [%s=%s].", i, j, thisEnvName, thisEnvValue))
			thisContainer.Env = append(thisContainer.Env, apiv1.EnvVar{
				Name:  thisEnvName,
				Value: thisEnvValue,
			})
		}

		// get mount items
		for j := 0; j < len(app.Containers[i].Mounts); j++ {
			thisVMPath := app.Containers[i].Mounts[j].VmPath
			thisContainerPath := app.Containers[i].Mounts[j].ContainerPath
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
		for j := 0; j < len(app.Containers[i].Ports); j++ {
			var onePort PortInfo = app.Containers[i].Ports[j]
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
						outErr := fmt.Errorf("Atoi ServicePort error: %w", err)
						beego.Error(outErr)
						return outErr
					}
					thisServicePort.Port = int32(sp)
				}

				// set node port if it exists
				if len(onePort.NodePort) > 0 {
					np, err := strconv.Atoi(onePort.NodePort)
					if err != nil {
						outErr := fmt.Errorf("Atoi NodePort error: %w", err)
						beego.Error(outErr)
						return outErr
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
			Name:      app.Name + models.DeploymentSuffix,
			Namespace: models.KubernetesNamespace,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &app.Replicas,
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
					HostNetwork: app.HostNetwork,
					DNSPolicy:   apiv1.DNSClusterFirstWithHostNet, // without this, pods with HostNetwork cannot access coredns
					Containers:  containers,
					Volumes:     volumes,
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

	// if the app in the request body has node name, we set it in K8s deployment
	if len(app.NodeName) > 0 {
		deployment.Spec.Template.Spec.NodeName = app.NodeName
	}
	// if the app in the request body has node selectors, we set them in K8s deployment
	if len(app.NodeSelector) > 0 {
		deployment.Spec.Template.Spec.NodeSelector = app.NodeSelector
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
		outErr := fmt.Errorf("Create deployment [%#v] error: %w", deployment, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("Deployment %s/%s created successful.", createdDeployment.Namespace, createdDeployment.Name))

	// service of this application
	if len(servicePorts) != 0 {
		service := &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Name + models.ServiceSuffix,
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
			outErr := fmt.Errorf("Json Marshal error: %w", err)
			beego.Error(outErr)
			return outErr
		}
		beego.Info(fmt.Sprintf("Create service (json) [%s]", string(serviceJson)))

		createdService, err := models.CreateService(service)
		if err != nil {
			outErr := fmt.Errorf("Create service [%#v] error: %w", service, err)
			beego.Error(outErr)
			return outErr
		}
		beego.Info(fmt.Sprintf("Service %s/%s created successful.", createdService.Namespace, createdService.Name))
	}

	return nil
}

// Used for json request, input is json
// test command:
// curl -i -X POST -H Content-Type:application/json -d '{"name":"test","replicas":2,"hostNetwork":true,"nodeName":"node1","nodeSelector":{"lnginx":"isnginx","lnginx2":"isnginx2"},"containers":[{"name":"printtime","image":"172.27.15.31:5000/printtime:v1","workDir":"/printtime","resources":{"limits":{"memory":"30Mi","cpu":"200m","storage":"2Gi"},"requests":{"memory":"20Mi","cpu":"100m","storage":"1Gi"}},"commands":["bash"],"args":["-c","python3 -u main.py > $LOGFILE"],"env":[{"name":"PARAMETER1","value":"testRenderenv1"},{"name":"LOGFILE","value":"/tmp/234/printtime.log"}],"mounts":[{"vmPath":"/tmp/asdff","containerPath":"/tmp/234"},{"vmPath":"/tmp/uyyyy","containerPath":"/tmp/2345"}],"ports":null},{"name":"nginx","image":"172.27.15.31:5000/nginx:1.17.0","workDir":"","resources":{"limits":{"memory":"","cpu":"","storage":""},"requests":{"memory":"","cpu":"","storage":""}},"commands":null,"args":null,"env":null,"mounts":null,"ports":[{"containerPort":80,"name":"fsd","protocol":"tcp","servicePort":"80","nodePort":"30001"}]},{"name":"ubuntu","image":"172.27.15.31:5000/ubuntu:latest","workDir":"","resources":{"limits":{"memory":"","cpu":"","storage":""},"requests":{"memory":"","cpu":"","storage":""}},"commands":["bash","-c","while true;do sleep 10;done"],"args":null,"env":[{"name":"asfasf","value":"asfasf"},{"name":"asdfsdf","value":"sfsdf"}],"mounts":[{"vmPath":"/tmp/asdff","containerPath":"/tmp/log"}],"ports":null}]}' http://localhost:20000/doNewApplication
func (c *ApplicationController) DoNewAppJson() {
	var app K8sApp
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

	beego.Info(fmt.Sprintf("From json input, we successfully parsed application [%v]", app))

	// Use the parsed app to create an application
	if err := CreateApplication(app); err != nil {
		outErr := fmt.Errorf("Create application %v, error: %w", app, err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString(outErr.Error())
		return
	}

	// Here, we wait until the app status becomes running.
	// We only do this behavior for json input, because for the form input, users can check the status on the web
	// And we need to put the application information (including the service port, pod IP, or nodePort IP) in the response body, to let the user know the information.
	beego.Info(fmt.Sprintf("Start to wait for the application [%s] running", app.Name))
	if err := WaitForAppRunning(models.WaitForTimeOut, 10, app.Name); err != nil {
		outErr := fmt.Errorf("Wait for application [%s] running, error: %w", app.Name, err)
		beego.Error(outErr)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString(outErr.Error())
		return
	}
	beego.Info(fmt.Sprintf("The application [%s] is already running", app.Name))

	outApp, err, statusCode := GetApplication(app.Name)
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

func WaitForAppRunning(timeout int, checkInterval int, appName string) error {
	return models.MyWaitFor(timeout, checkInterval, func() (bool, error) {
		app, err, statusCode := GetApplication(appName)
		if err != nil {
			if statusCode == http.StatusNotFound {
				return false, err
			}
			return false, nil
		}
		beego.Info(fmt.Sprintf("The status of the application [%s] is [%s]", appName, app.Status))
		if app.Status == RunningStatus {
			return true, nil
		}
		return false, nil
	})
}
