package models

import (
	"fmt"
	"sync"

	"github.com/astaxie/beego"
	"github.com/gophercloud/gophercloud"
	"github.com/spf13/viper"
)

type Iaas interface {
	ShowName() string
	ShowType() string
	ShowWebUrl() string
	GetVM(vmID string) (*IaasVm, error)
	ListAllVMs() ([]IaasVm, error)
	CreateVM(name string, vcpu, ram, storage int) (*IaasVm, error)
	DeleteVM(vmID string) error
	CheckResources() (ResourceStatus, error)
	IsCreatedByMcm(vmID string) (bool, error)
}

type ResourceStatus struct {
	Limit ResSet `json:"limit"` // total amounts of resources
	InUse ResSet `json:"inUse"` // the amounts of resources being used
}

// The backend APIs (Create VMs and Add K8s Nodes) request uses this struct, so we need to define the json of it.
type IaasVm struct {
	ID   string `json:"id"` // the id provided by the cloud
	Name string `json:"name"`

	// all IPs of this VM.
	// Although we can show multiple IPs, the VMs created by multi-cloud manager should only have 1 IP.
	// So when we need to get the IP of a VM, we can directly get its 1st IP.
	IPs []string `json:"ips"`

	VCpu      float64 `json:"vcpu"`    // number of logical CPU cores
	Ram       float64 `json:"ram"`     // memory size unit: MB
	Storage   float64 `json:"storage"` // storage size unit: GB
	Status    string  `json:"status"`
	Cloud     string  `json:"cloud"` // the name of the cloud that this VM belongs to
	CloudType string  `json:"cloudType"`
	McmCreate bool    `json:"mcmCreate"` // whether this VM is created by Multi-cloud manager
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

// the global variable to record all clouds
var Clouds map[string]Iaas = make(map[string]Iaas)
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
			osCloud := InitOpenstack(iaasParas[i])
			Clouds[osCloud.Name] = osCloud
		case ProxmoxIaas:
			pCloud := InitProxmox(iaasParas[i])
			Clouds[pCloud.Name] = pCloud
		default:
			beego.Info(fmt.Sprintf("Multi-cloud manager does not support cloud type [%s] of cloud [%s]", iaasParas[i]["type"].(string), iaasParas[i]["name"].(string)))
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
		output, err := SshOneCommand(sshClient, DiskInitCmd)
		if err != nil {
			beego.Info(fmt.Sprintf("Waiting for SSH ip %s, this time SshOneCommand \"\n%s\n\" error: %s", sshIP, DiskInitCmd, err.Error()))
			return false, nil
		}
		beego.Info(fmt.Sprintf("SSH of ip %s is enabled, output: %s", sshIP, output))
		return true, nil
	})
}

// after a VM is created, we should wait until the SSH is enabled, then we can to other things.
// We should also extend the disk to use the increased space
func WaitForSshPasswdAndInit(user string, passwd string, sshIP string, sshPort int, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		sshClient, err := SshClientWithPasswd(user, passwd, sshIP, sshPort)
		if err != nil {
			beego.Info(fmt.Sprintf("Waiting for SSH ip %s, this time SshClientWithPasswd error: %s", sshIP, err.Error()))
			return false, nil // cannot return error, otherwise, gophercloud.WaitFor will stop with error
		}
		defer sshClient.Close()
		output, err := SshOneCommand(sshClient, DiskInitCmd)
		if err != nil {
			beego.Info(fmt.Sprintf("Waiting for SSH ip %s, this time SshOneCommand \"\n%s\n\" error: %s", sshIP, DiskInitCmd, err.Error()))
			return false, nil
		}
		beego.Info(fmt.Sprintf("SSH of ip %s is enabled, output: %s", sshIP, output))
		return true, nil
	})
}

func CreateVms(vms []IaasVm) ([]IaasVm, error) {
	// create the VMs concurrently
	// use one goroutine to create one VM
	var errs []error
	var createdVms []IaasVm
	var errsMu sync.Mutex // the slice in golang is not safe for concurrent read/write
	var createdVmsMu sync.Mutex
	var wg sync.WaitGroup
	for _, vm := range vms {
		wg.Add(1)
		go func(v IaasVm) {
			defer wg.Done()
			beego.Info(fmt.Sprintf("Start to create vm Name [%s] CLoud [%s], vcpu cores [%f], ram [%f MB], storage [%f GB].", v.Name, v.Cloud, v.VCpu, v.Ram, v.Storage))
			createdVM, err := Clouds[v.Cloud].CreateVM(v.Name, int(v.VCpu), int(v.Ram), int(v.Storage))
			if err != nil {
				outErr := fmt.Errorf("Create vm %s error %w.", v.Name, err)
				beego.Error(outErr)
				errsMu.Lock()
				errs = append(errs, outErr)
				errsMu.Unlock()
			} else {
				beego.Info(fmt.Sprintf("Successful! Create vm:\n%+v\n", createdVM))
				createdVmsMu.Lock()
				createdVms = append(createdVms, *createdVM)
				createdVmsMu.Unlock()
			}
		}(vm)
	}
	wg.Wait()

	if len(errs) != 0 {
		sumErr := HandleErrSlice(errs)
		outErr := fmt.Errorf("CreateVms, Error: %w", sumErr)
		beego.Error(outErr)
		return createdVms, outErr
	}

	return createdVms, nil
}
