package models

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
)

type Proxmox struct {
	Name            string
	Type            string
	IP              string      // IP of Proxmox server
	Port            string      // Port of Proxmox service
	Endpoint        string      // IP:Port of Proxmox service
	ProxmoxUser     string      // The user used in Proxmox Web and Proxmox server SSH
	ProxmoxPassword string      // The password used in Proxmox Web and Proxmox server SSH
	TokenName       string      // The name of API Token for HTTP request
	AuthHeader      string      // The header "Authorization" used in HTTP request
	RootPasswd      string      // root password for SSH of VMs.
	HTTPClient      http.Client // used to call the API of proxmox
}

func InitProxmox(paras map[string]interface{}) *Proxmox {
	beego.Info(fmt.Sprintf("Start to initialize cloud name [%s] type [%s]", paras["name"].(string), paras["type"].(string)))

	ip := paras["ip"].(string)
	port := paras["port"].(string)
	proxmoxUser := paras["proxmox_user"].(string)
	tokenName := paras["token_name"].(string)
	tokenSecret := paras["token_secret"].(string)

	// initialize the http client to call the API of proxmox
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // curl -k for https
		},
	}
	client := http.Client{
		Transport: tr,
	}

	return &Proxmox{
		Name:            paras["name"].(string),
		Type:            paras["type"].(string),
		IP:              ip,
		Port:            port,
		Endpoint:        ip + ":" + port,
		ProxmoxUser:     proxmoxUser,
		ProxmoxPassword: paras["proxmox_password"].(string),
		AuthHeader:      fmt.Sprintf("PVEAPIToken=%s@pam!%s=%s", proxmoxUser, tokenName, tokenSecret),
		RootPasswd:      paras["root_password"].(string),
		HTTPClient:      client,
	}
}

func (p *Proxmox) ShowName() string {
	return p.Name
}

func (p *Proxmox) ShowType() string {
	return p.Type
}

func (p *Proxmox) ShowWebUrl() string {
	return fmt.Sprintf("https://%s:%s/", p.IP, p.Port)
}

// in Proxmox, a node is a cloud, this function is to get the cloud status
func (p *Proxmox) NodeStatus() ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], get node status.", p.Name, p.Type))

	// send HTTP request to get node status firstly, in Proxmox, a node is a cloud
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/status", p.Endpoint, p.Name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get node status, construct request, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get node status, do HTTP request, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get node status, read response body, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], get node status.", p.Name, p.Type))
	return body, nil
}

// in Proxmox, a qemu is a VM, this function is to list all qemus
func (p *Proxmox) ListQemus() ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], list all qemus.", p.Name, p.Type))

	// send HTTP request to list all qemus firstly, in Proxmox, a qemu is a VM
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu", p.Endpoint, p.Name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], list all qemus, construct request, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], list all qemus, do HTTP request, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], list all qemus, read response body, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], list all qemus.", p.Name, p.Type))
	return body, nil
}

// in Proxmox, a qemu is a VM, this function is get the current status of a qemu
func (p *Proxmox) GetQemu(vmid string) ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], get qemu ID [%s].", p.Name, p.Type, vmid))

	// send HTTP request to get this qemu status, in Proxmox, a qemu is a VM
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu/%s/status/current", p.Endpoint, p.Name, vmid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu id [%s], construct request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu id [%s], do HTTP request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu id [%s], read response body, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], get qemu id [%s].", p.Name, p.Type, vmid))
	return body, nil
}

// To get all net interfaces of a qemu
func (p *Proxmox) GetNetInterfaces(vmid string) ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], get vm id [%s] network interfaces .", p.Name, p.Type, vmid))
	// send HTTP request to get vm network interfaces, in Proxmox, a qemu is a VM
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu/%s/agent/network-get-interfaces", p.Endpoint, p.Name, vmid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get vm id [%s], construct request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get vm id [%s], do HTTP request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get vm id [%s], read response body, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], get vm id [%s].", p.Name, p.Type, vmid))
	return body, nil
}

// get all IPs of a Proxmox VM. If error, return an empty string slice, which means we cannot get IPs, but it does not affect other information.
func (p *Proxmox) getVmIps(vmid string) []string {
	netIntsBytes, err := p.GetNetInterfaces(vmid)
	if err != nil {
		beego.Error(fmt.Errorf("Cloud name [%s], type [%s], GetNetInterfaces vmid %s, error: %w", p.Name, p.Type, vmid, err))
		return []string{}
	}
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], GetNetInterfaces vmid %s, response: %s", p.Name, p.Type, vmid, string(netIntsBytes)))

	var netInts map[string]interface{}
	if err := json.Unmarshal(netIntsBytes, &netInts); err != nil {
		beego.Error(fmt.Errorf("Cloud name [%s], type [%s], GetNetInterfaces vmid %s, Unmarshal netIntsBytes, error: %w", p.Name, p.Type, vmid, err))
		return []string{}
	}

	var netIntSlice []interface{}
	switch netInts["data"].(type) {
	case map[string]interface{}:
		netIntSlice = netInts["data"].(map[string]interface{})["result"].([]interface{})
	default:
		beego.Info(fmt.Errorf("netInts[\"data\"] is not a map[string]interface{}"))
		return []string{}
	}

	var vmIps []string
	for _, netInt := range netIntSlice {
		// we do not need loopback IP
		if netInt.(map[string]interface{})["hardware-address"].(string) == LoopBackMac || netInt.(map[string]interface{})["name"].(string) == LoopBackIntName {
			continue
		}

		// We only have requirements about interface name
		if !IsIfNeeded(netInt.(map[string]interface{})["name"].(string)) {
			continue
		}

		ipAddrs := netInt.(map[string]interface{})["ip-addresses"].([]interface{})
		for _, ipAddr := range ipAddrs {
			if ipAddr.(map[string]interface{})["ip-address-type"].(string) == IPv4Type {
				vmIps = append(vmIps, ipAddr.(map[string]interface{})["ip-address"].(string))
			}
		}
	}
	return vmIps
}

func (p *Proxmox) CheckResources() (ResourceStatus, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], check resources.", p.Name, p.Type))

	// if we cannot get a resource, return -1
	var errResult ResourceStatus = errRs

	// From node status get total CPU, total Storage, total Memory.
	nodeStatusBytes, err := p.NodeStatus()
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], check resources, get node status, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return errResult, outErr
	}

	beego.Info(fmt.Sprintf("Node status of Cloud [%s] is [%v]", p.Name, string(nodeStatusBytes)))

	var nodeStatus map[string]interface{}
	if err := json.Unmarshal(nodeStatusBytes, &nodeStatus); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], check resources, unmarshal nodeStatus, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return errResult, outErr
	}
	totalCPU := nodeStatus["data"].(map[string]interface{})["cpuinfo"].(map[string]interface{})["cpus"].(float64)
	totalMemoryUnitB := nodeStatus["data"].(map[string]interface{})["memory"].(map[string]interface{})["total"].(float64)
	totalStorageUnitB := nodeStatus["data"].(map[string]interface{})["rootfs"].(map[string]interface{})["total"].(float64)
	beego.Info(fmt.Sprintf("Cloud [%s], type [%s], totalCPU: %f, totalMemoryUnitB: %f, totalStorageUnitB: %f.", p.Name, p.Type, totalCPU, totalMemoryUnitB, totalStorageUnitB))

	// From qemus get inuse CPU, inuse Storage, inuse Memory.
	qemusBytes, err := p.ListQemus()
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], check resources, list qemus, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return errResult, outErr
	}
	beego.Info(fmt.Sprintf("Qemus of Cloud [%s] is [%v]", p.Name, string(qemusBytes)))

	var qemus map[string]interface{}
	if err := json.Unmarshal(qemusBytes, &qemus); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], check resources, unmarshal qemus, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return errResult, outErr
	}
	qemuSlice := qemus["data"].([]interface{})
	var usedCPU, usedMemoryUnitB, usedStorageUnitB float64
	for _, qemu := range qemuSlice {
		usedCPU += qemu.(map[string]interface{})["cpus"].(float64)
		usedMemoryUnitB += qemu.(map[string]interface{})["maxmem"].(float64)
		usedStorageUnitB += qemu.(map[string]interface{})["maxdisk"].(float64)
	}
	beego.Info(fmt.Sprintf("Cloud [%s], type [%s], usedCPU: %f, usedMemoryUnitB: %f, usedStorageUnitB: %f.", p.Name, p.Type, usedCPU, usedMemoryUnitB, usedStorageUnitB))

	return ResourceStatus{
		Limit: ResSet{
			VCpu:    totalCPU,
			Ram:     totalMemoryUnitB / 1024 / 1024, // unit MB
			Vm:      -1,
			Volume:  -1,
			Storage: totalStorageUnitB / 1024 / 1024 / 1024, // unit GB
			Port:    -1,
		},
		InUse: ResSet{
			VCpu:    usedCPU,
			Ram:     usedMemoryUnitB / 1024 / 1024, // unit MB
			Vm:      -1,
			Volume:  -1,
			Storage: usedStorageUnitB / 1024 / 1024 / 1024, // unit GB
			Port:    -1,
		},
	}, nil
}

func (p *Proxmox) GetVM(vmid string) (*IaasVm, error) {
	qemuBytes, err := p.GetQemu(vmid)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get VM id [%s], get qemu, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], get qemu id %s, response: %s", p.Name, p.Type, vmid, string(qemuBytes)))

	var qemu map[string]interface{}
	if err := json.Unmarshal(qemuBytes, &qemu); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get VM id %s, Unmarshal qemuBytes, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}

	return &IaasVm{
		ID:        strconv.FormatFloat(qemu["data"].(map[string]interface{})["vmid"].(float64), 'f', -1, 64),
		Name:      qemu["data"].(map[string]interface{})["name"].(string),
		IPs:       p.getVmIps(vmid),
		VCpu:      qemu["data"].(map[string]interface{})["cpus"].(float64),
		Ram:       qemu["data"].(map[string]interface{})["maxmem"].(float64) / 1024 / 1024,         // unit MB
		Storage:   qemu["data"].(map[string]interface{})["maxdisk"].(float64) / 1024 / 1024 / 1024, // unit GB
		Status:    qemu["data"].(map[string]interface{})["status"].(string),
		Cloud:     p.Name,
		CloudType: p.Type,
	}, nil
}
func (p *Proxmox) ListAllVMs() ([]IaasVm, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], list all VMs.", p.Name, p.Type))

	// From qemus get inuse CPU, inuse Storage, inuse Memory.
	qemusBytes, err := p.ListQemus()
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], list all VMs, list qemus, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return []IaasVm{}, outErr
	}
	beego.Info(fmt.Sprintf("Qemus of Cloud [%s] is [%v]", p.Name, string(qemusBytes)))

	var qemus map[string]interface{}
	if err := json.Unmarshal(qemusBytes, &qemus); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], list all VMs, unmarshal qemus, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return []IaasVm{}, outErr
	}
	qemuSlice := qemus["data"].([]interface{})
	var vms []IaasVm
	for _, qemu := range qemuSlice {
		var vmid string = strconv.FormatFloat(qemu.(map[string]interface{})["vmid"].(float64), 'f', -1, 64)

		// get the ip address of this VM.
		var ips []string = p.getVmIps(vmid)

		thisVM := IaasVm{
			ID:        vmid,
			Name:      qemu.(map[string]interface{})["name"].(string),
			IPs:       ips,
			VCpu:      qemu.(map[string]interface{})["cpus"].(float64),
			Ram:       qemu.(map[string]interface{})["maxmem"].(float64) / 1024 / 1024,         // unit MB
			Storage:   qemu.(map[string]interface{})["maxdisk"].(float64) / 1024 / 1024 / 1024, // unit GB
			Status:    qemu.(map[string]interface{})["status"].(string),
			Cloud:     p.Name,
			CloudType: p.Type,
		}
		vms = append(vms, thisVM)
	}

	return vms, nil
}

func (p *Proxmox) CreateVM(name string, vcpu, ram, storage int) (*IaasVm, error) {
	return nil, nil
}

func (p *Proxmox) DeleteVM(vmID string) error {
	return nil
}
