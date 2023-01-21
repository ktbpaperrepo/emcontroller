package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gophercloud/gophercloud"
	"github.com/spf13/viper"
)

type Iaas interface {
	CreateVM(name string, vcpu, ram, storage int) (*IaasVm, error)
	DeleteVM(vmID string) error
	CheckResources() (ResourceStatus, error)
}

type ResourceStatus struct {
	Limit ResSet // total amounts of resources
	InUse ResSet // the amounts of resources being used
}

type IaasVm struct {
	ID        string // the id provided by the cloud
	Name      string
	IPs       []string
	Cloud     string
	CloudType string
}

// Resource set
type ResSet struct {
	VCpu    float64 `json:"vcpu"`    // number of logical CPU cores
	Ram     float64 `json:"ram"`     // memory size unit: MB
	Vm      float64 `json:"vm"`      // number of virtual machines
	Volume  float64 `json:"volume"`  // number of volumes
	Storage float64 `json:"storage"` // storage size unit: GB
	Port    float64 `json:"port"`    // number of network ports
}

var Clouds []Iaas
var iaasConfig *viper.Viper

// read config from iaas.json
func readIaasConfig() {
	iaasConfig = viper.New()
	iaasConfig.SetConfigFile("conf/iaas.json")
	iaasConfig.SetConfigType("json")
	if err := iaasConfig.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error perse iaas.json: %w", err))
	}
}

// InitClouds init the slice Clouds
func InitClouds() {
	readIaasConfig()
	var iaasParas []map[string]interface{}
	if err := iaasConfig.UnmarshalKey("iaas", &iaasParas); err != nil {
		panic(fmt.Errorf("UnmarshalKey \"iaas\" of iaas.json error: %w", err))
	}
	// use the configuration parameters to build the elements in the slice Clouds
	for i := 0; i < len(iaasParas); i++ {
		switch iaasParas[i]["type"].(string) {
		case OpenstackIaas:
			Clouds = append(Clouds, InitOpenstack(iaasParas[i]))
		case ProxmoxIaas:
		}
	}
	beego.Info(fmt.Sprintf("All %d clouds are initialized.", len(Clouds)))
}

// after a VM is created, we should wait until the SSH is enabled, then we can to other things.
func WaitForSshPem(user string, pemFilePath string, sshIP string, sshPort int, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		sshClient, err := SshClientWithPem(pemFilePath, user, sshIP, sshPort)
		if err != nil {
			beego.Info(fmt.Sprintf("Waiting for SSH ip %s, this time SshClientWithPem error: %s", sshIP, err.Error()))
			return false, nil // cannot return error, otherwise, gophercloud.WaitFor will stop with error
		}
		defer sshClient.Close()
		output, err := SshOneCommand(sshClient, "pwd")
		if err != nil {
			beego.Info(fmt.Sprintf("Waiting for SSH ip %s, this time SshOneCommand \"pwd\" error: %s", sshIP, err.Error()))
			return false, nil
		}
		beego.Info(fmt.Sprintf("SSH of ip %s is enabled, output: %s", sshIP, output))
		return true, nil
	})
}

// after a VM is created, we should wait until the SSH is enabled, then we can to other things.
func WaitForSshPasswd(user string, passwd string, sshIP string, sshPort int, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		sshClient, err := SshClientWithPasswd(user, passwd, sshIP, sshPort)
		if err != nil {
			beego.Info(fmt.Sprintf("Waiting for SSH ip %s, this time SshClientWithPasswd error: %s", sshIP, err.Error()))
			return false, nil // cannot return error, otherwise, gophercloud.WaitFor will stop with error
		}
		defer sshClient.Close()
		output, err := SshOneCommand(sshClient, "pwd")
		if err != nil {
			beego.Info(fmt.Sprintf("Waiting for SSH ip %s, this time SshOneCommand \"pwd\" error: %s", sshIP, err.Error()))
			return false, nil
		}
		beego.Info(fmt.Sprintf("SSH of ip %s is enabled, output: %s", sshIP, output))
		return true, nil
	})
}
