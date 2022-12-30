package models

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubernetesClient *kubernetes.Clientset

func init() {
	// default kubeconfig path
	if KubeConfigPath == "" {
		KubeConfigPath = defaultKubeConfigPath
	}
	kubernetesClient = initKubernetesClient()
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
