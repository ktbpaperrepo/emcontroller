package models

import (
	"time"

	"github.com/astaxie/beego"
)

const (
	ControllerName        string        = "Multi-Cloud Manager"
	UploadDir             string        = "upload/"
	RequestTimeout        time.Duration = 5 * time.Minute
	KubernetesNamespace   string        = "default"
	defaultKubeConfigPath string        = "/root/.kube/config"
	DeploymentSuffix      string        = "-deployment"
	ServiceSuffix         string        = "-service"

	// type of clouds
	OpenstackIaas string = "openstack"
	ProxmoxIaas   string = "proxmox"
)

var (
	DockerEngine   string = beego.AppConfig.String("dockerEngineIP") + ":" + beego.AppConfig.String("dockerEnginePort")
	DockerRegistry string = beego.AppConfig.String("dockerRegistryIP") + ":" + beego.AppConfig.String("dockerRegistryPort")
	KubeConfigPath string = beego.AppConfig.String("kubeConfigPath")
)
