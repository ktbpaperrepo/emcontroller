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

	DiskInitCmd string = "fsname=$(df -h / | grep -v Filesystem | awk '{print $1}'); diskname=$(echo ${fsname} | sed 's/2$//'); echo \"d\n2\nn\n2\n\n\nNo\nw\" | fdisk ${diskname}; resize2fs ${fsname}"

	LoopBackMac     = "00:00:00:00:00:00"
	LoopBackIntName = "lo"
	LoopBackIp      = "127.0.0.1"
	IPv4Type        = "ipv4"
	IPv6Type        = "ipv6"
)

var (
	DockerEngine   string = beego.AppConfig.String("dockerEngineIP") + ":" + beego.AppConfig.String("dockerEnginePort")
	DockerRegistry string = beego.AppConfig.String("dockerRegistryIP") + ":" + beego.AppConfig.String("dockerRegistryPort")
	KubeConfigPath string = beego.AppConfig.String("kubeConfigPath")

	// When checking resources, if not successful, return -1
	errRs ResourceStatus = ResourceStatus{
		Limit: ResSet{
			VCpu:    float64(-1),
			Ram:     float64(-1),
			Vm:      float64(-1),
			Volume:  float64(-1),
			Storage: float64(-1),
			Port:    float64(-1),
		},
		InUse: ResSet{
			VCpu:    float64(-1),
			Ram:     float64(-1),
			Vm:      float64(-1),
			Volume:  float64(-1),
			Storage: float64(-1),
			Port:    float64(-1),
		},
	}
)
