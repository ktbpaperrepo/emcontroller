package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// used for the input of creating applications, so we need to define the json
type K8sApp struct {
	Name          string              `json:"name"`
	Replicas      int32               `json:"replicas"`
	HostNetwork   bool                `json:"hostNetwork"`
	NodeName      string              `json:"nodeName,omitempty"`
	NodeSelector  map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations   []corev1.Toleration `json:"tolerations,omitempty"`
	Containers    []K8sContainer      `json:"containers"`
	Priority      int                 `json:"priority"`
	AutoScheduled bool                `json:"autoScheduled"`
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
	NodePortIP    []string  `json:"nodePortIP"`
	SvcPort       []string  `json:"svcPort"`
	NodePort      []string  `json:"nodePort"`
	ContainerPort []string  `json:"containerPort"`
	Hosts         []PodHost `json:"hosts"`
	Status        string    `json:"status"`
	Priority      int       `json:"priority"`
	AutoScheduled bool      `json:"autoScheduled"`
}

type PodHost struct {
	PodIP    string `json:"podIP"`
	HostName string `json:"hostName"`
	HostIP   string `json:"hostIP"`
}

// a method to check whether the application is running
func appRunning(app appsv1.Deployment) bool {
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
func getAllPods(app appsv1.Deployment) []corev1.Pod {
	selector, err := metav1.LabelSelectorAsSelector(app.Spec.Selector)
	if err != nil {
		beego.Error(fmt.Sprintf("Error, get the selector of deployment %s/%s, error: %s", app.Namespace, app.Name, err.Error()))
		return []corev1.Pod{}
	}
	stringSelector := selector.String()
	beego.Info(fmt.Sprintf("List pods belonging to the app %s/%s, with selector: %s", app.Namespace, app.Name, stringSelector))
	pods, err := ListPods(app.Namespace, metav1.ListOptions{LabelSelector: stringSelector})
	if err != nil {
		beego.Error(fmt.Sprintf("Error, List pods in namespace %s, with selector: %s, error: %s", app.Namespace, stringSelector, err.Error()))
		return []corev1.Pod{}
	}
	return pods
}

// get the host Kubernetes Nodes of all pods of this application
func getHosts(app appsv1.Deployment, pods []corev1.Pod) []PodHost {
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
func getNodePortIP(app appsv1.Deployment, pods []corev1.Pod) []string {
	if len(pods) == 0 {
		beego.Info(fmt.Sprintf("No pods belonging to the app %s/%s are got.", app.Namespace, app.Name))
		return []string{}
	}

	var nodePortIPs []string
	for i := 0; i < len(pods); i++ {
		if len(pods[i].Status.HostIP) > 0 {
			nodePortIPs = append(nodePortIPs, pods[i].Status.HostIP)
		}
	}

	// We add kubernetes master IP only when this application is available.
	// And in our environment in AAU 5G Smart Production Lab, Kubernetes master node cannot access the container IPs (I do not have time to figure out why), so we only show Kubernetes Master IP for NodePort for the applications with HostNetwork.
	if len(nodePortIPs) > 0 && pods[0].Spec.HostNetwork {
		nodePortIPs = append(nodePortIPs, beego.AppConfig.String("k8sMasterIP"))
	}

	return nodePortIPs
}

func ListApplications() ([]AppInfo, error) {
	applications, err := ListDeployment(KubernetesNamespace)
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
		go func(d appsv1.Deployment) {
			defer wg.Done()

			thisApp, _ := getAppInfoDeploy(d)

			appListMu.Lock()
			appList = append(appList, thisApp)
			appListMu.Unlock()
		}(app)
	}
	wg.Wait()

	return appList, nil
}

func GetApplication(appName string) (AppInfo, error, int) {
	deployName := appName + DeploymentSuffix

	deploy, err := GetDeployment(KubernetesNamespace, deployName)
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

	outApp, err := getAppInfoDeploy(*deploy)
	if err != nil {
		outErr := fmt.Errorf("getAppInfoDeploy, app [%s], error: %w", appName, err)
		beego.Error(outErr)
		return outApp, outErr, http.StatusInternalServerError
	}

	return outApp, nil, http.StatusOK
}

// get application info from a Kubernetes deployment
func getAppInfoDeploy(d appsv1.Deployment) (AppInfo, error) {
	var thisApp AppInfo
	var appName, svcName string
	appName = strings.TrimSuffix(d.Name, DeploymentSuffix)
	svcName = appName + ServiceSuffix

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

	// read annotation to get the autoschedule value
	autoScheduled, err := strconv.ParseBool(d.Annotations[AutoScheduledAnno])
	if err != nil {
		beego.Info(fmt.Sprintf("Parse %s to bool, error [%s], app [%s] AutoScheduled should be automatically set to \"false\".", d.Annotations[AutoScheduledAnno], err.Error(), appName))
	}
	thisApp.AutoScheduled = autoScheduled

	// read annotation to get the priority value
	if priority, err := strconv.Atoi(d.Annotations[PriorityAnno]); err == nil {
		thisApp.Priority = priority
	} else {
		beego.Info(fmt.Sprintf("Parse %s to int, error [%s], app [%s] Priority should be automatically set to \"0\".", d.Annotations[PriorityAnno], err.Error(), appName))
		thisApp.Priority = 0
	}

	svc, err := GetService(KubernetesNamespace, svcName)
	if err != nil {
		outErr := fmt.Errorf("GetService %s/%s error: %w", KubernetesNamespace, svcName, err)
		beego.Error(outErr)

		thisApp.ClusterIP = ""
		thisApp.NodePortIP = []string{}
		thisApp.SvcPort = []string{}
		thisApp.NodePort = []string{}

		return thisApp, outErr
	}
	if svc != nil {
		thisApp.ClusterIP = svc.Spec.ClusterIP
		if svc.Spec.Type == corev1.ServiceTypeNodePort {
			thisApp.NodePortIP = getNodePortIP(d, pods)
		} else {
			thisApp.NodePortIP = []string{}
		}
		for _, port := range svc.Spec.Ports {
			thisApp.SvcPort = append(thisApp.SvcPort, strconv.FormatInt(int64(port.Port), 10))
			thisApp.NodePort = append(thisApp.NodePort, strconv.FormatInt(int64(port.NodePort), 10))
			thisApp.ContainerPort = append(thisApp.ContainerPort, port.TargetPort.String())
		}
	} else {
		thisApp.ClusterIP = ""
		thisApp.NodePortIP = []string{}
		thisApp.SvcPort = []string{}
		thisApp.NodePort = []string{}
	}
	return thisApp, nil
}

func DeleteApplication(appName string) (error, int) {
	deployName := appName + DeploymentSuffix
	svcName := appName + ServiceSuffix

	deploy, err := GetDeployment(KubernetesNamespace, deployName)
	if err != nil {
		outErr := fmt.Errorf("Get the deployment of app [%s], error: %w", appName, err)
		beego.Error(outErr)
		return outErr, http.StatusInternalServerError
	}

	beego.Info(fmt.Sprintf("Delete deployment [%s/%s]", KubernetesNamespace, deployName))
	if err := DeleteDeployment(KubernetesNamespace, deployName); err != nil {
		outErr := fmt.Errorf("Delete deployment [%s/%s] error: %s", KubernetesNamespace, deployName, err.Error())
		beego.Error(outErr)
		return outErr, http.StatusInternalServerError
	}
	beego.Info(fmt.Sprintf("Successfully sent request to delete deployment [%s/%s]", KubernetesNamespace, deployName))

	beego.Info(fmt.Sprintf("Delete service [%s/%s]", KubernetesNamespace, svcName))
	if err := DeleteService(KubernetesNamespace, svcName); err != nil {
		outErr := fmt.Errorf("Delete deployment [%s/%s] error: %s", KubernetesNamespace, svcName, err.Error())
		beego.Error(outErr)
		return outErr, http.StatusInternalServerError
	}
	beego.Info(fmt.Sprintf("Successful! Delete service [%s/%s]", KubernetesNamespace, svcName))

	beego.Info(fmt.Sprintf("Start to wait for the deployment [%s/%s] deleted.", KubernetesNamespace, deployName))
	if err := WaitForDeployDeleted(WaitForTimeOut, 10, deploy); err != nil {
		outErr := fmt.Errorf("Wait for the deployment [%s/%s] deleted, error: %w", KubernetesNamespace, deployName, err)
		beego.Error(outErr)
		return outErr, http.StatusInternalServerError
	}
	beego.Info(fmt.Sprintf("The deployment [%s/%s] is already deleted.", KubernetesNamespace, deployName))

	beego.Info(fmt.Sprintf("Successful! Deleted deployment [%s/%s]", KubernetesNamespace, deployName))
	return nil, http.StatusOK
}

// for an application, we need to create a Deployment and a Service for it
func CreateApplication(app K8sApp) error {
	if err := ValidateK8sApp(app); err != nil {
		outErr := fmt.Errorf("Validate app [%s] error: %w", app.Name, err)
		beego.Error(outErr)
		return outErr
	}

	// Kubernetes labels of the pods of this application
	labels := map[string]string{
		"app": app.Name,
	}

	// the ports in the service of this application
	var servicePorts []corev1.ServicePort

	var hasNodePort bool

	// make volume configuration in a pod
	beego.Info("make volume configuration in a pod")
	var volumeP2N map[string]string = make(map[string]string) // a map from VM paths to volume names
	var volumes []corev1.Volume                               // put this slice into deployment template
	for i := 0; i < len(app.Containers); i++ {
		beego.Info(fmt.Sprintf("make volume configuration, Container [%d] has [%d] mount items.", i, len(app.Containers[i].Mounts)))

		for j := 0; j < len(app.Containers[i].Mounts); j++ {
			thisVMPath := app.Containers[i].Mounts[j].VmPath
			if _, exist := volumeP2N[thisVMPath]; !exist {
				thisVolumeName := "volume" + strconv.Itoa(len(volumeP2N))
				beego.Info(fmt.Sprintf("add volume name: [%s], VM path: [%s]", thisVolumeName, thisVMPath))
				volumeP2N[thisVMPath] = thisVolumeName

				var hostPathType corev1.HostPathType = corev1.HostPathDirectoryOrCreate
				volumes = append(volumes, corev1.Volume{
					Name: thisVolumeName,
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
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
	var containers []corev1.Container
	for i := 0; i < len(app.Containers); i++ {
		var thisContainer corev1.Container = corev1.Container{
			ImagePullPolicy: corev1.PullIfNotPresent,
		}

		// After rolling update, if we want the rolling update seamless and does not have downtime, we need to give some time to the loadbalancing update before the signal "kill -15" to delete the old pod, so we set this preStop hook to all containers by default.
		var defaultLifecycle *corev1.Lifecycle = &corev1.Lifecycle{
			PreStop: &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"/bin/sh", "-c", "sleep 10"},
				},
			},
		}
		thisContainer.Lifecycle = defaultLifecycle

		beego.Info(fmt.Sprintf("Get the configuration of container %d", i))
		thisContainer.Name = app.Containers[i].Name
		thisContainer.Image = app.Containers[i].Image
		thisContainer.Resources = corev1.ResourceRequirements{
			Requests: make(corev1.ResourceList),
			Limits:   make(corev1.ResourceList),
		}
		if len(app.Containers[i].Resources.Requests.CPU) > 0 {
			thisContainer.Resources.Requests[corev1.ResourceCPU] = resource.MustParse(app.Containers[i].Resources.Requests.CPU)
		}
		if len(app.Containers[i].Resources.Requests.Memory) > 0 {
			thisContainer.Resources.Requests[corev1.ResourceMemory] = resource.MustParse(app.Containers[i].Resources.Requests.Memory)
		}
		if len(app.Containers[i].Resources.Requests.Storage) > 0 {
			thisContainer.Resources.Requests[corev1.ResourceEphemeralStorage] = resource.MustParse(app.Containers[i].Resources.Requests.Storage)
		}
		if len(app.Containers[i].Resources.Limits.CPU) > 0 {
			thisContainer.Resources.Limits[corev1.ResourceCPU] = resource.MustParse(app.Containers[i].Resources.Limits.CPU)
		}
		if len(app.Containers[i].Resources.Limits.Memory) > 0 {
			thisContainer.Resources.Limits[corev1.ResourceMemory] = resource.MustParse(app.Containers[i].Resources.Limits.Memory)
		}
		if len(app.Containers[i].Resources.Limits.Storage) > 0 {
			thisContainer.Resources.Limits[corev1.ResourceEphemeralStorage] = resource.MustParse(app.Containers[i].Resources.Limits.Storage)
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
			thisContainer.Env = append(thisContainer.Env, corev1.EnvVar{
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

			thisContainer.VolumeMounts = append(thisContainer.VolumeMounts, corev1.VolumeMount{
				Name:      volumeName,
				MountPath: thisContainerPath,
			})
		}

		// get ports
		for j := 0; j < len(app.Containers[i].Ports); j++ {
			var onePort PortInfo = app.Containers[i].Ports[j]
			beego.Info(fmt.Sprintf("Container [%d], Port [%d]: [%+v].", i, j, onePort))

			// put port information into the container of the deployment
			thisContainer.Ports = append(thisContainer.Ports, corev1.ContainerPort{
				ContainerPort: int32(onePort.ContainerPort),
				Name:          onePort.Name,
				Protocol:      corev1.Protocol(strings.ToUpper(onePort.Protocol)),
			})

			if len(onePort.ServicePort)+len(onePort.NodePort) > 0 {
				thisServicePort := corev1.ServicePort{
					Name:       onePort.Name,
					Protocol:   corev1.Protocol(strings.ToUpper(onePort.Protocol)),
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

	// When a deployment's Replicas is 1, we should set the maxUnavailable as 0, and MaxSurge as 1. Then minAvailable will be 1 - 0 = 1.
	// In this condition, only when the new pod is available, the availablePodCount will be 2 > 1, then, the old pod will be deleted. This is seamless.
	// If we set maxUnavailable as 1, and MaxSurge as 1, the old pod will be deleted before the new pod is available, so it will not be seamless.
	maxUnavailable := intstr.FromInt(1)
	if app.Replicas == 1 {
		beego.Info(fmt.Sprintf("app [%s], app.Replicas is 1, so we set maxUnavailable as 0, to enable seamless rolling update.", app.Name))
		maxUnavailable = intstr.FromInt(0)
	}
	maxSurge := intstr.FromInt(1)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name + DeploymentSuffix,
			Namespace: KubernetesNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &app.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &maxUnavailable,
					MaxSurge:       &maxSurge,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					HostNetwork: app.HostNetwork,
					DNSPolicy:   corev1.DNSClusterFirstWithHostNet, // without this, pods with HostNetwork cannot access coredns
					Containers:  containers,
					Volumes:     volumes,
					Affinity: &corev1.Affinity{
						PodAntiAffinity: &corev1.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
								corev1.PodAffinityTerm{
									TopologyKey: corev1.LabelHostname,
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

	// If the app in the request body has tolerations, we set them in K8s deployment
	if len(app.Tolerations) > 0 {
		deployment.Spec.Template.Spec.Tolerations = app.Tolerations
	}

	// for auto scheduled applications, we add this annotation, to enable it to be migrated
	if app.AutoScheduled {
		if deployment.Annotations == nil { // avoid panic
			deployment.Annotations = make(map[string]string)
		}
		deployment.Annotations[AutoScheduledAnno] = strconv.FormatBool(app.AutoScheduled)
		deployment.Annotations[PriorityAnno] = strconv.Itoa(app.Priority)
	}

	beego.Info(fmt.Sprintf("Create deployment [%+v]", deployment))
	beego.Info(fmt.Sprintf(""))
	deploymentJson, err := json.Marshal(deployment)
	if err != nil {
		beego.Error(fmt.Sprintf("Json Marshal error: %s", err.Error()))
	}
	beego.Info(fmt.Sprintf("Create deployment (json) [%s]", string(deploymentJson)))

	createdDeployment, err := CreateDeployment(deployment)
	if err != nil {
		outErr := fmt.Errorf("Create deployment [%+v] error: %w", deployment, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("Deployment %s/%s created successful.", createdDeployment.Namespace, createdDeployment.Name))

	// service of this application
	if len(servicePorts) != 0 {
		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Name + ServiceSuffix,
				Namespace: KubernetesNamespace,
			},
			Spec: corev1.ServiceSpec{
				Selector: labels,
				Type:     corev1.ServiceTypeClusterIP,
				Ports:    servicePorts,
			},
		}
		if hasNodePort {
			service.Spec.Type = corev1.ServiceTypeNodePort
		}

		beego.Info(fmt.Sprintf("Create service [%+v]", service))
		beego.Info(fmt.Sprintf(""))
		serviceJson, err := json.Marshal(service)
		if err != nil {
			outErr := fmt.Errorf("Json Marshal error: %w", err)
			beego.Error(outErr)
			return outErr
		}
		beego.Info(fmt.Sprintf("Create service (json) [%s]", string(serviceJson)))

		createdService, err := CreateService(service)
		if err != nil {
			outErr := fmt.Errorf("Create service [%+v] error: %w", service, err)
			beego.Error(outErr)
			return outErr
		}
		beego.Info(fmt.Sprintf("Service %s/%s created successful.", createdService.Namespace, createdService.Name))
	}

	return nil
}

func WaitForAppRunning(timeout int, checkInterval int, appName string) error {
	return MyWaitFor(timeout, checkInterval, func() (bool, error) {
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
