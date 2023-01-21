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

	McmSign        string = "mcmcreated" // add this sign something, meaning that it is created by multi-cloud manager
	WaitForTimeOut int    = 600          // unit second. wait for 10 minutes when creating or deleting something

	SshUser     string        = "ubuntu"
	SshRootUser string        = "root"
	SshPort     int           = 22
	SshTimeout  time.Duration = 10 * time.Second

	Os404Substr string = "itemNotFound" // this string exists in the "not found" error of Openstack
)

var (
	DockerEngine   string = beego.AppConfig.String("dockerEngineIP") + ":" + beego.AppConfig.String("dockerEnginePort")
	DockerRegistry string = beego.AppConfig.String("dockerRegistryIP") + ":" + beego.AppConfig.String("dockerRegistryPort")
	KubeConfigPath string = beego.AppConfig.String("kubeConfigPath")
)
