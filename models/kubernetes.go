package models

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubernetesClient *kubernetes.Clientset

func InitKubernetesClient() {
	// default kubeconfig path
	if KubeConfigPath == "" {
		KubeConfigPath = defaultKubeConfigPath
	}
	kubernetesClient = initKubernetesClient()
	beego.Info("Kubernetes client initialized.")
}

func initKubernetesClient() *kubernetes.Clientset {
	// Use the configuration of "kubeconfig" to access kubernetes
	config, err := clientcmd.BuildConfigFromFlags("", KubeConfigPath)
	if err != nil {
		beego.Error(fmt.Printf("Build kubernetes config error: %s", err.Error()))
		panic(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		beego.Error(fmt.Printf("Create kubernetes client error: %s", err.Error()))
		panic(err)
	}
	return client
}

func ListDeployment(namespace string) ([]v1.Deployment, error) {
	ctx := context.Background()
	deployments, err := kubernetesClient.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		beego.Error(fmt.Printf("List deployments error: %s", err.Error()))
		return []v1.Deployment{}, err
	}
	return deployments.Items, nil
}

func CreateDeployment(d *v1.Deployment) (*v1.Deployment, error) {
	ctx := context.Background()
	createdDeployment, err := kubernetesClient.AppsV1().Deployments(d.Namespace).Create(ctx, d, metav1.CreateOptions{})
	if err != nil {
		beego.Error(fmt.Printf("Create deployment %s/%s error: %s", d.Namespace, d.Name, err.Error()))
	}
	return createdDeployment, err
}

func DeleteDeployment(namespace, name string) error {
	ctx := context.Background()
	//deletePolicy := metav1.DeletePropagationForeground
	err := kubernetesClient.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{
		//PropagationPolicy: &deletePolicy,
	})
	if err != nil && errors.IsNotFound(err) {
		beego.Info(fmt.Printf("Deployment %s/%s not found: %s, do nothing", namespace, name, err.Error()))
		return nil
	}
	if err != nil {
		beego.Error(fmt.Printf("Delete deployment %s/%s error: %s", namespace, name, err.Error()))
		return err
	}
	return nil
}

func CreateService(s *apiv1.Service) (*apiv1.Service, error) {
	ctx := context.Background()
	createdService, err := kubernetesClient.CoreV1().Services(s.Namespace).Create(ctx, s, metav1.CreateOptions{})
	if err != nil {
		beego.Error(fmt.Printf("Create service %s/%s error: %s", s.Namespace, s.Name, err.Error()))
	}
	return createdService, err
}

func DeleteService(namespace, name string) error {
	ctx := context.Background()
	//deletePolicy := metav1.DeletePropagationForeground
	err := kubernetesClient.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{
		//PropagationPolicy: &deletePolicy,
	})
	if err != nil && errors.IsNotFound(err) {
		beego.Info(fmt.Printf("Service %s/%s not found: %s, do nothing", namespace, name, err.Error()))
		return nil
	}
	if err != nil {
		beego.Error(fmt.Printf("Delete service %s/%s error: %s", namespace, name, err.Error()))
		return err
	}
	return nil
}

func GetService(namespace, name string) (*apiv1.Service, error) {
	ctx := context.Background()
	service, err := kubernetesClient.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		beego.Info(fmt.Printf("Service %s/%s not found: %s", namespace, name, err.Error()))
		return nil, nil
	}
	if err != nil {
		beego.Error(fmt.Printf("Get service %s/%s error: %s", namespace, name, err.Error()))
		return nil, err
	}
	return service, nil
}

func ListPods(namespace string, listOptions metav1.ListOptions) ([]apiv1.Pod, error) {
	ctx := context.Background()
	pods, err := kubernetesClient.CoreV1().Pods(namespace).List(ctx, listOptions)
	if err != nil {
		beego.Error(fmt.Printf("List pods error: %s", err.Error()))
		return []apiv1.Pod{}, err
	}
	return pods.Items, nil
}

func ListNodes(listOptions metav1.ListOptions) ([]apiv1.Node, error) {
	ctx := context.Background()
	nodes, err := kubernetesClient.CoreV1().Nodes().List(ctx, listOptions)
	if err != nil {
		beego.Error(fmt.Printf("List nodes error: %s", err.Error()))
		return []apiv1.Node{}, err
	}
	return nodes.Items, nil
}

func GetNode(name string, getOptions metav1.GetOptions) (*apiv1.Node, error) {
	ctx := context.Background()
	node, err := kubernetesClient.CoreV1().Nodes().Get(ctx, name, getOptions)
	if err != nil {
		beego.Error(fmt.Printf("Get node %s with options %v error: %s", name, getOptions, err.Error()))
		return nil, err
	}
	return node, nil
}

func DeleteNode(name string, deleteOptions metav1.DeleteOptions) error {
	ctx := context.Background()
	err := kubernetesClient.CoreV1().Nodes().Delete(ctx, name, deleteOptions)
	if err != nil {
		beego.Error(fmt.Printf("Delete node %s with options %v error: %s", name, deleteOptions, err.Error()))
		return err
	}
	return nil
}

func GetNodeInternalIp(node apiv1.Node) string {
	var ip string
	for _, addr := range node.Status.Addresses {
		// Kubernetes gets the VM IP from the network interface used as default gateway, and put it in the internal IP
		// pkg/kubelet/nodestatus/setters.go, func NodeAddress
		if addr.Type == apiv1.NodeInternalIP {
			ip = addr.Address
			break
		}
	}
	return ip
}

// Imitate kubectl get node xxx -o wide.
// pkg/printers/internalversion/printers.go, func printNode
func ExtractNodeStatus(node apiv1.Node) string {
	var status string = "Unknown"
	for _, cond := range node.Status.Conditions {
		if cond.Type == apiv1.NodeReady {
			if cond.Status == apiv1.ConditionTrue {
				status = string(cond.Type)
			} else {
				status = "Not" + string(cond.Type)
			}
		}
	}
	return status
}

// SSH to Kubernetes Master node to generate a kubeadm join command. The command will expire by default after 24 hours
func GetJoinCmd() (string, error) {
	K8sMasterIP := beego.AppConfig.String("k8sMasterIP")
	sshPrivateKey := beego.AppConfig.String("k8sVmSshPrivateKey")
	sshPort := SshPort
	sshUser := SshRootUser

	sshClient, err := SshClientWithPem(sshPrivateKey, sshUser, K8sMasterIP, sshPort)
	if err != nil {
		outErr := fmt.Errorf("Create ssh client with SSH key fail: error: %w", err)
		beego.Error(outErr)
		return "", outErr
	}
	defer sshClient.Close()

	joinCmdBytes, err := SshOneCommand(sshClient, "kubeadm token create --print-join-command")
	if err != nil {
		outErr := fmt.Errorf("ssh error: %w", err)
		beego.Error(outErr)
		return "", outErr
	}

	joinCmd := strings.TrimRight(string(joinCmdBytes), "\r\n ")

	return joinCmd, nil
}

// kubectl drain <nodeName> --ignore-daemonsets
func DrainNode(nodeName string) error {
	K8sMasterIP := beego.AppConfig.String("k8sMasterIP")
	sshPrivateKey := beego.AppConfig.String("k8sVmSshPrivateKey")
	sshPort := SshPort
	sshUser := SshRootUser

	sshClient, err := SshClientWithPem(sshPrivateKey, sshUser, K8sMasterIP, sshPort)
	if err != nil {
		outErr := fmt.Errorf("Create ssh client with SSH key fail: error: %w", err)
		beego.Error(outErr)
		return outErr
	}
	defer sshClient.Close()

	if _, err := SshOneCommand(sshClient, fmt.Sprintf("kubectl drain %s --ignore-daemonsets", nodeName)); err != nil {
		outErr := fmt.Errorf("ssh error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	return nil
}

// add a new node into Kubernetes cluster
func AddNode(vm IaasVm, joinCmd string) error {
	if len(vm.IPs) == 0 {
		outErr := fmt.Errorf("the input vm [%s] has no ip address", vm.Name)
		beego.Error(outErr)
		return outErr
	}

	// get the name and IP of the input VM
	name, ip := vm.Name, vm.IPs[0]

	sshPrivateKey := beego.AppConfig.String("k8sVmSshPrivateKey")
	sshPort := SshPort
	sshUser := SshRootUser

	sshClient, err := SshClientWithPem(sshPrivateKey, sshUser, ip, sshPort)
	if err != nil {
		outErr := fmt.Errorf("Create ssh client with SSH key fail: error: %w", err)
		beego.Error(outErr)
		return outErr
	}
	defer sshClient.Close()

	// replace the placeholder in containerd configuration file
	K8sMasterIP := beego.AppConfig.String("k8sMasterIP")
	if _, err := SshOneCommand(sshClient, fmt.Sprintf("sed -i 's/<IP>/%s/g' /etc/containerd/config.toml", K8sMasterIP)); err != nil {
		outErr := fmt.Errorf("ssh error: %w", err)
		beego.Error(outErr)
		return outErr
	}
	// restart containerd
	if _, err := SshOneCommand(sshClient, "systemctl restart containerd"); err != nil {
		outErr := fmt.Errorf("ssh error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	// set that kubelet start when VM start
	if _, err := SshOneCommand(sshClient, "systemctl enable kubelet"); err != nil {
		outErr := fmt.Errorf("ssh error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	// execute kubeadm join command to add this VM into Kubernetes cluster
	if _, err := SshOneCommand(sshClient, fmt.Sprintf("%s --node-name=%s", joinCmd, name)); err != nil {
		outErr := fmt.Errorf("ssh error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	if err := WaitForNodeJoin(WaitForTimeOut, 5, name); err != nil {
		outErr := fmt.Errorf("Wait for node %s join, error: %w", name, err)
		beego.Error(outErr)
		return outErr
	}

	return nil
}

func WaitForNodeJoin(timeout int, checkInterval int, nodeName string) error {
	return MyWaitFor(timeout, checkInterval, func() (bool, error) {
		node, err := GetNode(nodeName, metav1.GetOptions{})
		if err != nil {
			outErr := fmt.Errorf("Get Kubernetes node %s, error: %w", nodeName, err)
			beego.Error(outErr)
			return false, nil
		}
		beego.Info(fmt.Sprintf("node: [%v]", node))

		if ExtractNodeStatus(*node) == string(apiv1.NodeReady) {
			return true, nil
		}

		return false, nil
	})
}

// add some new nodes into Kubernetes cluster
func AddNodes(vms []IaasVm) []error {
	if len(vms) == 0 {
		outErr := fmt.Errorf("no vms in the input")
		beego.Error(outErr)
		return []error{outErr}
	}

	// use one joinCmd to add all nodes
	joinCmd, err := GetJoinCmd()
	if err != nil {
		outErr := fmt.Errorf("AddNodes, GetJoinCmd error: %w", err)
		beego.Error(outErr)
		return []error{outErr}
	}

	var errs []error
	// add all nodes in parallel
	// use one goroutine to add one node
	var wg sync.WaitGroup
	var errsMu sync.Mutex // the slice (errs) in golang is not safe for concurrent read/write
	for i := 0; i < len(vms); i++ {
		wg.Add(1)
		go func(vm IaasVm) {
			defer wg.Done()
			if err := AddNode(vm, joinCmd); err != nil {
				outErr := fmt.Errorf("AddNode %s, error: %w", vm.Name, err)
				beego.Error(outErr)
				errsMu.Lock()
				errs = append(errs, outErr)
				errsMu.Unlock()
			}
		}(vms[i])
	}
	wg.Wait() // wait for all nodes joined

	return errs
}

// delete a node from the Kubernetes cluster
func UninstallNode(name string) error {
	// // kubectl drain xxxxxx --ignore-daemonsets
	if err := DrainNode(name); err != nil {
		outErr := fmt.Errorf("DrainNode %s: error: %w", name, err)
		beego.Error(outErr)
		return outErr
	}

	// // put kubeconfig in /root/.kube
	sshPrivateKey := beego.AppConfig.String("k8sVmSshPrivateKey")
	sshPort := SshPort
	sshUser := SshRootUser

	node, err := GetNode(name, metav1.GetOptions{})
	if err != nil {
		outErr := fmt.Errorf("Kubernetes get node: error: %w", err)
		beego.Error(outErr)
		return outErr
	}
	nodeIp := GetNodeInternalIp(*node)
	sshClient, err := SshClientWithPem(sshPrivateKey, sshUser, nodeIp, sshPort)
	if err != nil {
		outErr := fmt.Errorf("Create ssh client with SSH key fail: error: %w", err)
		beego.Error(outErr)
		return outErr
	}
	defer sshClient.Close()

	// create the folder at the destination VM
	if _, err := SshOneCommand(sshClient, fmt.Sprintf("mkdir -p /root/.kube")); err != nil {
		outErr := fmt.Errorf("ssh error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	// put the kubeConfig into the folder of the destination VM
	var kubeConfigPath string
	if KubeConfigPath == "" {
		kubeConfigPath = defaultKubeConfigPath
	} else {
		kubeConfigPath = KubeConfigPath
	}
	if err := SftpCopyFile(kubeConfigPath, defaultKubeConfigPath, sshClient); err != nil {
		outErr := fmt.Errorf("SFTP error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	// // kubeadm reset the destination VM
	if _, err := SshOneCommand(sshClient, fmt.Sprintf("kubeadm reset -f")); err != nil {
		outErr := fmt.Errorf("ssh error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	// // Delete the folder /root/.kube
	if _, err := SshOneCommand(sshClient, fmt.Sprintf("rm -rf /root/.kube")); err != nil {
		outErr := fmt.Errorf("ssh error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	// // delete the node from the Kubernetes Cluster
	beego.Info(fmt.Sprintf("Use client-go API to delete node %s in Kubernetes cluster", name))
	if err := DeleteNode(name, metav1.DeleteOptions{}); err != nil {
		outErr := fmt.Errorf("use client-go API to delete node error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	return nil
}
