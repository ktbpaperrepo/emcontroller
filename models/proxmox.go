package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type Proxmox struct {
	Name            string
	Type            string
	IP              string // IP of Proxmox server
	Port            string // Port of Proxmox service
	Endpoint        string // IP:Port of Proxmox service
	ProxmoxUser     string // The user used in Proxmox Web and Proxmox server SSH
	ProxmoxPassword string // The password used in Proxmox Web and Proxmox server SSH
	TokenName       string // The name of API Token for HTTP request
	AuthHeader      string // The header "Authorization" used in HTTP request
	RootPasswd      string // root password for SSH of VMs.
	TemplateId      string // the ID of the VM template used to create new VMs

	HTTPClient http.Client // used to call the API of proxmox
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
		TemplateId:      paras["template_id"].(string),
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
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get node status, read response body, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get node status, HTTP response status code is [%d]", p.Name, p.Type, resp.StatusCode)
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
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], list all qemus, read response body, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], list all qemus, HTTP response status code is [%d]", p.Name, p.Type, resp.StatusCode)
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
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu id [%s], read response body, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu id [%s], HTTP response status code is [%d]", p.Name, p.Type, vmid, resp.StatusCode)
		beego.Error(outErr)
		return nil, outErr
	}

	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], get qemu id [%s].", p.Name, p.Type, vmid))
	return body, nil
}

// Clone a new VM from a template
func (p *Proxmox) CloneQemu(newVmId, newVmName string) ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], clone new VM [%s/%s] from template ID [%s].", p.Name, p.Type, newVmId, newVmName, p.TemplateId))

	// send HTTP request to clone a new VM from a template
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu/%s/clone", p.Endpoint, p.Name, p.TemplateId)

	// We add the sign McmSign to the VMs created by multi-cloud manager. When deleting a VM, multi-cloud manager is only allowed to delete the VMs created by itself with this sign.
	bodyJson := []byte(fmt.Sprintf(`{ "newid": %s, "full": true, "name": "%s", "description": "%s"}`, newVmId, newVmName, McmSign))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], clone new VM [%s/%s] from template ID [%s], construct request, error: %w", p.Name, p.Type, newVmId, newVmName, p.TemplateId, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	req.Header.Add("Content-Type", "application/json")
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], clone new VM [%s/%s] from template ID [%s], do HTTP request, error: %w", p.Name, p.Type, newVmId, newVmName, p.TemplateId, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], clone new VM [%s/%s] from template ID [%s], read response body, error: %w", p.Name, p.Type, newVmId, newVmName, p.TemplateId, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], clone new VM [%s/%s] from template ID [%s], HTTP response status code is [%d]", p.Name, p.Type, newVmId, newVmName, p.TemplateId, resp.StatusCode)
		beego.Error(outErr)
		return nil, outErr
	}

	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], clone new VM [%s/%s] from template ID [%s].", p.Name, p.Type, newVmId, newVmName, p.TemplateId))
	return body, nil
}

// Shutdown a VM
func (p *Proxmox) ShutdownQemu(vmid string) ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], shutdown VM [%s].", p.Name, p.Type, vmid))

	// send HTTP request to shut down a VM
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu/%s/status/shutdown", p.Endpoint, p.Name, vmid)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], shutdown VM [%s], construct request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], shutdown VM [%s], do HTTP request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], shutdown VM [%s], read response body, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], shutdown VM [%s], HTTP response status code is [%d]", p.Name, p.Type, vmid, resp.StatusCode)
		beego.Error(outErr)
		return nil, outErr
	}

	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], shutdown VM [%s].", p.Name, p.Type, vmid))
	return body, nil
}

// Delete a VM
func (p *Proxmox) DeleteQemu(vmid string) ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], delete Qemu [%s].", p.Name, p.Type, vmid))

	// send HTTP request to shut down a VM
	// The following two lines are copied from the Proxmox API documentation.
	// purge: boolean. Remove VMID from configurations, like backup & replication jobs and HA.
	// destroy-unreferenced-disks: boolean. If set, destroy additionally all disks not referenced in the config but with a matching VMID from all enabled storages.
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu/%s?purge=1&destroy-unreferenced-disks=1", p.Endpoint, p.Name, vmid)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], delete Qemu [%s], construct request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], delete Qemu [%s], do HTTP request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], delete Qemu [%s], read response body, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], delete Qemu [%s], HTTP response status code is [%d]", p.Name, p.Type, vmid, resp.StatusCode)
		beego.Error(outErr)
		return nil, outErr
	}

	beego.Info(fmt.Sprintf("Request Successful! Cloud name [%s], type [%s], delete Qemu [%s].", p.Name, p.Type, vmid))
	return body, nil
}

// Get the status of a task
func (p *Proxmox) GetTaskStatus(upid string) ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], get task upid ID [%s].", p.Name, p.Type, upid))

	// send HTTP request to get a task
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/tasks/%s/status", p.Endpoint, p.Name, upid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get task upid ID [%s], construct request, error: %w", p.Name, p.Type, upid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get task upid ID [%s], do HTTP request, error: %w", p.Name, p.Type, upid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get task upid ID [%s], read response body, error: %w", p.Name, p.Type, upid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get task upid ID [%s], HTTP response status code is [%d]", p.Name, p.Type, upid, resp.StatusCode)
		beego.Error(outErr)
		return nil, outErr
	}

	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], get task upid ID [%s].", p.Name, p.Type, upid))
	return body, nil
}

// Config CPU and Memory of a VM, the unit of the input RAM is MB, and they must be integers, which means we cannot set the 15555.5555MB Ram, or 5.55 cores
func (p *Proxmox) ConfigCoreRam(vmid, ramMB, cores int) ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], Config VM [%d], Memory [%d]MB, CPU cores [%d].", p.Name, p.Type, vmid, ramMB, cores))

	// send HTTP request to config the Memory and CPU of a VM
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu/%d/config", p.Endpoint, p.Name, vmid)

	var reqBody map[string]interface{} = map[string]interface{}{
		"memory": ramMB,
		"cores":  cores,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		outErr := fmt.Errorf("json.Marshal: %+v, error: %w", reqBody, err)
		beego.Error(outErr)
		return nil, outErr
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], Config VM [%d], Memory [%d]MB, CPU cores [%d], construct request, error: %w", p.Name, p.Type, vmid, ramMB, cores, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	req.Header.Add("Content-Type", "application/json")
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], Config VM [%d], Memory [%d]MB, CPU cores [%d], do HTTP request, error: %w", p.Name, p.Type, vmid, ramMB, cores, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], Config VM [%d], Memory [%d]MB, CPU cores [%d], read response body, error: %w", p.Name, p.Type, vmid, ramMB, cores, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], Config VM [%d], Memory [%d]MB, CPU cores [%d], HTTP response status code is [%d]", p.Name, p.Type, vmid, ramMB, cores, resp.StatusCode)
		beego.Error(outErr)
		return nil, outErr
	}

	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], Config VM [%d], Memory [%d]MB, CPU cores [%d].", p.Name, p.Type, vmid, ramMB, cores))
	return body, nil
}

// Get the config of a Qemu
func (p *Proxmox) GetQemuConfig(vmid string) ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], get qemu config ID [%s].", p.Name, p.Type, vmid))

	// send HTTP request to get this qemu config, in Proxmox, a qemu is a VM
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu/%s/config", p.Endpoint, p.Name, vmid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu config id [%s], construct request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu config id [%s], do HTTP request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu config id [%s], read response body, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu config id [%s], HTTP response status code is [%d]", p.Name, p.Type, vmid, resp.StatusCode)
		beego.Error(outErr)
		return nil, outErr
	}

	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], get qemu config id [%s].", p.Name, p.Type, vmid))
	return body, nil
}

// Check whether a VM is created by multi-cloud manager
func (p *Proxmox) IsCreatedByMcm(vmid string) (bool, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], check whether VM [%s] is created by multi-cloud manager", p.Name, p.Type, vmid))

	qemuConfigBytes, err := p.GetQemuConfig(vmid)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], check whether VM [%s] is created by multi-cloud manager, get qemu config id, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return false, outErr
	}
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], check whether VM [%s] is created by multi-cloud manager, get qemu config id, response: %s", p.Name, p.Type, vmid, string(qemuConfigBytes)))
	if err := p.CheckErrInResp(qemuConfigBytes); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], check whether VM [%s] is created by multi-cloud manager, get qemu config id, error in resp: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return false, outErr
	}

	var qemuConfig map[string]interface{}
	if err := json.Unmarshal(qemuConfigBytes, &qemuConfig); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], check whether VM [%s] is created by multi-cloud manager, get qemu config id, Unmarshal qemuBytes, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return false, outErr
	}

	if description, exist := qemuConfig["data"].(map[string]interface{})["description"]; exist && description.(string) == McmSign {
		beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], VM [%s] is created by multi-cloud manager", p.Name, p.Type, vmid))
		return true, nil
	}

	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], VM [%s] is not created by multi-cloud manager", p.Name, p.Type, vmid))
	return false, nil
}

// Get the Name of the disk of a qemu
func (p *Proxmox) getDiskName(vmid string) (string, error) {
	qemuConfigBytes, err := p.GetQemuConfig(vmid)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], getDiskName, get qemu config id [%s], error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return "", outErr
	}
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], getDiskName, get qemu config id %s, response: %s", p.Name, p.Type, vmid, string(qemuConfigBytes)))
	if err := p.CheckErrInResp(qemuConfigBytes); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], getDiskName, get qemu config id %s, error in resp: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return "", outErr
	}

	var qemuConfig map[string]interface{}
	if err := json.Unmarshal(qemuConfigBytes, &qemuConfig); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], getDiskName, get qemu config id %s, Unmarshal qemuBytes, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return "", outErr
	}

	var bootValue string = qemuConfig["data"].(map[string]interface{})["boot"].(string)
	var bootOrder string = strings.TrimPrefix(bootValue, "order=")
	var diskName string = strings.Split(bootOrder, ";")[0] // we should guarantee that the first boot source is the disk
	return diskName, nil
}

// Resize the storage of a VM, the input is a string
func (p *Proxmox) ResizeDisk(vmid int, disk, size string) ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], Resize VM [%d], disk [%s] size [%s].", p.Name, p.Type, vmid, disk, size))

	// send HTTP request to resize the disk a VM
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu/%d/resize", p.Endpoint, p.Name, vmid)

	var reqBody map[string]interface{} = map[string]interface{}{
		"disk": disk,
		"size": size,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		outErr := fmt.Errorf("json.Marshal: %+v, error: %w", reqBody, err)
		beego.Error(outErr)
		return nil, outErr
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], Resize VM [%d], disk [%s] size [%s], construct request, error: %w", p.Name, p.Type, vmid, disk, size, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	req.Header.Add("Content-Type", "application/json")
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], Resize VM [%d], disk [%s] size [%s], do HTTP request, error: %w", p.Name, p.Type, vmid, disk, size, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], Resize VM [%d], disk [%s] size [%s], read response body, error: %w", p.Name, p.Type, vmid, disk, size, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], Resize VM [%d], disk [%s] size [%s], HTTP response status code is [%d]", p.Name, p.Type, vmid, disk, size, resp.StatusCode)
		beego.Error(outErr)
		return nil, outErr
	}

	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], Resize VM [%d], disk [%s] size [%s].", p.Name, p.Type, vmid, disk, size))
	return body, nil
}

// Resize the storage of a VM, the input is a string
func (p *Proxmox) StartQemu(vmid int) ([]byte, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], start (turn on) VM [%d].", p.Name, p.Type, vmid))

	// send HTTP request to resize the disk a VM
	url := fmt.Sprintf("https://%s/api2/json/nodes/%s/qemu/%d/status/start", p.Endpoint, p.Name, vmid)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], start (turn on) VM [%d], construct request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	req.Header.Add("Authorization", p.AuthHeader)
	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], start (turn on) VM [%d], do HTTP request, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer resp.Body.Close()
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], start (turn on) VM [%d], read response body, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], start (turn on) VM [%d], HTTP response status code is [%d]", p.Name, p.Type, vmid, resp.StatusCode)
		beego.Error(outErr)
		return nil, outErr
	}

	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], start (turn on) VM [%d].", p.Name, p.Type, vmid))
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
	beego.Info(fmt.Sprintf("HTTP Status is [%s], HTTP Status Code is [%d]", resp.Status, resp.StatusCode))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get vm id [%s], read response body, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("HTTP response body is [%s].", string(body)))

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get vm id [%s], HTTP response status code is [%d]", p.Name, p.Type, vmid, resp.StatusCode)
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

// After a VM is created, I suspect that we may not be able to get its IPs in a short time, so I make this function to get IPs with retry.
func (p *Proxmox) getVmIpsWithRetry(vmid string, retryTimes int, retryInterval time.Duration) []string {
	var ips []string

	for i := 0; i < retryTimes; i++ {
		beego.Info(fmt.Sprintf("The %d time to try to get the VM IPs.", i+1))
		if i != 0 {
			beego.Info(fmt.Sprintf("Before try again to get the VM IPs, we wait for %v.", retryInterval))
			time.Sleep(retryInterval)
		}

		ips = p.getVmIps(vmid)

		if len(ips) > 0 {
			break
		}
	}

	return ips
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

	if err := p.CheckErrInResp(qemuBytes); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], get VM id [%s], get qemu, error in resp: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return nil, outErr
	}

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

// the unit of vcpu, ram, storage in the input is consistent with ResSet
func (p *Proxmox) CreateVM(name string, vcpu, ram, storage int) (*IaasVm, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], Create VM: %s", p.Name, p.Type, name))

	// 1. Find the first available VM ID, and use it to create this VM
	// get the IDs of all existing qemus
	qemusBytes, err := p.ListQemus()
	time.Sleep(ProxmoxAPIInterval)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], list qemus, error: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Qemus of Cloud [%s] is [%v]", p.Name, string(qemusBytes)))
	if err := p.CheckErrInResp(qemusBytes); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], list qemus, error in resp: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}

	var qemus map[string]interface{}
	if err := json.Unmarshal(qemusBytes, &qemus); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], check resources, unmarshal qemus, error: %w", p.Name, p.Type, err)
		beego.Error(outErr)
		return nil, outErr
	}
	qemuSlice := qemus["data"].([]interface{})
	existingIDs := make(map[int]struct{})
	for _, qemu := range qemuSlice {
		existingIDs[int(qemu.(map[string]interface{})["vmid"].(float64))] = struct{}{}
	}
	beego.Info(fmt.Sprintf("existing VM IDs are %v", existingIDs))
	// the qemuIDs in proxmox are integers starting from 100
	// find the first unused ID, and use it as the ID of the VM that we are creating
	var vmid int = 100
	for ; ; vmid++ {
		if _, exist := existingIDs[vmid]; !exist {
			break
		}
	}
	beego.Info(fmt.Sprintf("We use %d as this VM ID", vmid))

	// 2. Call Clone Qemu API to create a new VM from a template
	vmidStr := strconv.Itoa(vmid)
	cloneRespBytes, err := p.CloneQemu(vmidStr, name)
	time.Sleep(ProxmoxAPIInterval)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], CloneQemu, error: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Cloud [%s], CloneQemu response of is [%v]", p.Name, string(cloneRespBytes)))
	if err := p.CheckErrInResp(cloneRespBytes); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], CloneQemu, error in resp: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Clone Qemu [%s] request is sent successfully.", name))

	var cloneResp map[string]interface{}
	if err := json.Unmarshal(cloneRespBytes, &cloneResp); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], unmarshal cloneResp, error: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	cloneUpid := cloneResp["data"].(string) // used for check task status

	// 3-1. After we create the clone task, we can use this polling function to wait for it to be finished
	beego.Info(fmt.Sprintf("Wait for Task [%s] is finished.", cloneUpid))
	if err := p.waitForTaskFinished(WaitForTimeOut, 5, cloneUpid); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], waitForTaskFinished, CloneQemu, error: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Task [%s] is already finished successfully.", cloneUpid))

	// 3-2. After a VM is created, the config is locked in a short time. We need to wait until it is unlocked.
	// We use a polling function to check the lock of the newly created VM
	beego.Info(fmt.Sprintf("Wait for Qemu [%s] is unlocked after clone.", name))
	if err := p.waitForUnlock(WaitForTimeOut, 5, vmidStr); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], waitForUnlock, error: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Qemu [%s] is already unlocked after clone.", name))

	// After the VM is unlocked, we can set the CPU, Memory, and Storage of it

	// 4. configure CPU and Ram
	configCpuMemRespBytes, err := p.ConfigCoreRam(vmid, ram, vcpu)
	time.Sleep(ProxmoxAPIInterval)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], config ram [%d]MB, cpu [%d], error: %w", p.Name, p.Type, name, ram, vcpu, err)
		beego.Error(outErr)
		return nil, outErr
	}
	if err := p.CheckErrInResp(configCpuMemRespBytes); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], config ram [%d]MB, cpu [%d], in resp error: %w", p.Name, p.Type, name, ram, vcpu, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], CreateVM [%s], config ram [%d]MB, cpu [%d].", p.Name, p.Type, name, ram, vcpu))

	// 5. configure storage
	// We should get the disk name firstly
	diskName, err := p.getDiskName(vmidStr)
	time.Sleep(ProxmoxAPIInterval)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], getDiskName, error: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("diskName is %s", diskName))

	// call API to resize disk
	// There is a bug of proxmox described in https://forum.proxmox.com/threads/bug-a-bug-about-the-api-to-resize-disk-in-version-7-3-3.123400/.
	// After a VM is created, the disk cannot be resized in a short time (I know it in my practice). We need to retry until the resize is successful.
	// Moreover, only the retry is not enough, because if the resize fails one time, even if we succeed afterward, the config in Proxmox will still have error, i.e., the disk size shown in Proxmox will be the old value rather than the resized value.
	// So we need to SSH to the Proxmox Node to execute a command like qm rescan to refresh the status.

	// resize the disk, the input size of this API should be a string with value and unit.
	diskSize := fmt.Sprintf("%dG", storage)
	beego.Info(fmt.Sprintf("set the disk [%s] size %s, which may trigger a Proxmox bug.", diskName, diskSize))
	if err := p.waitForResizeDisk(WaitForTimeOut, 5, vmid, diskName, diskSize); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], waitForResizeDisk, error: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], CreateVM [%s], config disk [%s] size [%s].", p.Name, p.Type, name, diskName, diskSize))

	// SSH to the Proxmox node to execute the command to fix the problem of the proxmox bug.
	qmRescanCmd := fmt.Sprintf("qm rescan --vmid %d", vmid)
	beego.Info(fmt.Sprintf("SSH to the Proxmox node [%s] to run command [%s] to refresh the state to fix the problem of the Proxmox bug.", p.IP, qmRescanCmd))
	sshClient, err := SshClientWithPasswd(p.ProxmoxUser, p.ProxmoxPassword, p.IP, SshPort)
	if err != nil {
		outErr := fmt.Errorf("Create SshClientWithPasswd for ip %s, this time SshClientWithPasswd error: %w", p.IP, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer sshClient.Close()
	output, err := SshOneCommand(sshClient, qmRescanCmd)
	if err != nil {
		outErr := fmt.Errorf("Execute command %s on Proxmox node %s error: %w", qmRescanCmd, p.IP, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! SSH to the Proxmox node [%s] to run command [%s] to refresh the state to fix the Proxmox bug. output: %s", p.IP, qmRescanCmd, output))

	// 5. Start the VM
	startResqBytes, err := p.StartQemu(vmid)
	time.Sleep(ProxmoxAPIInterval)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], start qemu, error: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	if err := p.CheckErrInResp(startResqBytes); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], start qemu, in resp error: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], CreateVM [%s], start qemu.", p.Name, p.Type, name))

	// 6. Get the 1st IP as the SSH IP. Then, wait for SSH enabled. Then, SSH to the VM and execute commands to extend the disk partition.

	// get the 1st IP as the SSH IP
	vmIPs := p.getVmIpsWithRetry(vmidStr, 4, 30*time.Second)
	if len(vmIPs) == 0 {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], no IPs are got", p.Name, p.Type, name)
		beego.Error(outErr)
		return nil, outErr
	}
	sshIP := vmIPs[0]
	beego.Info(fmt.Sprintf("found IPs [%v], we use ip [%s] to ssh", vmIPs, sshIP))

	// Then, wait for SSH enabled. Then, SSH to the VM and execute commands to extend the disk partition.
	if err := WaitForSshPasswdAndInit(SshRootUser, p.RootPasswd, sshIP, SshPort, WaitForTimeOut); err != nil {
		outErr := fmt.Errorf("wait for VM %s able to be SSHed, ip %s, error: %w", name, sshIP, err)
		beego.Error(outErr)
		return nil, outErr
	}

	// 7. Get the VM and return it.
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], CreateVM [%s]. The creation and initialization are finished, then we get the VM and return it.", p.Name, p.Type, name))
	iaasVm, err := p.GetVM(vmidStr)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%s], get VM, error: %w", p.Name, p.Type, name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], Create VM: %s", p.Name, p.Type, name))
	return iaasVm, nil
}

func (p *Proxmox) DeleteVM(vmid string) error {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], Delete VM: %s", p.Name, p.Type, vmid))

	// Multi-cloud manager is only allowed to delete VMs created by itself
	createdByMcm, err := p.IsCreatedByMcm(vmid)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], DeleteVM [%s], check whether the VM is created by multi-cloud manager, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return outErr
	}
	if !createdByMcm {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], DeleteVM [%s], the VM is not created by multi-cloud manager, so we cannot delete it.", p.Name, p.Type, vmid)
		beego.Error(outErr)
		return outErr
	}

	// 1. Shut down the VM
	sdRespBytes, err := p.ShutdownQemu(vmid)
	time.Sleep(ProxmoxAPIInterval)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], DeleteVM [%s], ShutdownQemu, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("Cloud [%s], ShutdownQemu [%s] response of is [%s]", p.Name, vmid, string(sdRespBytes)))
	if err := p.CheckErrInResp(sdRespBytes); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], DeleteVM [%s], ShutdownQemu, error in resp: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("ShutdownQemu [%s] request is sent successfully.", vmid))

	var sdResp map[string]interface{}
	if err := json.Unmarshal(sdRespBytes, &sdResp); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], DeleteVM [%s], unmarshal sdResp, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return outErr
	}
	sdUpid := sdResp["data"].(string) // used for check task status

	// 2-1. After we create the shutdown task, we can use this polling function to wait for it to be finished
	beego.Info(fmt.Sprintf("Wait for Task [%s] is finished.", sdUpid))
	if err := p.waitForTaskFinished(WaitForTimeOut, 5, sdUpid); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], DeleteVM [%s], waitForTaskFinished, ShutdownQemu, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("Task [%s] is already finished successfully.", sdUpid))

	// 2-2. Before we delete the VM, to be safe, we make sure that the VM status is "stopped"
	beego.Info(fmt.Sprintf("Wait for VM [%s] status to be [%s].", vmid, ProxQSStopped))
	if err := p.waitForQemuStopped(WaitForTimeOut, 5, vmid); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], DeleteVM [%s], waitForQemuStopped, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("VM [%s] is already stopped.", vmid))

	// 3. Send request to delete this Qemu
	dqRespBytes, err := p.DeleteQemu(vmid)
	time.Sleep(ProxmoxAPIInterval)
	if err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], DeleteVM [%s], DeleteQemu, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("Cloud [%s], DeleteQemu [%s] response of is [%s]", p.Name, vmid, string(dqRespBytes)))
	if err := p.CheckErrInResp(dqRespBytes); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], DeleteVM [%s], DeleteQemu, error in resp: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("DeleteQemu [%s] request is sent successfully.", vmid))

	var dqResp map[string]interface{}
	if err := json.Unmarshal(dqRespBytes, &dqResp); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], DeleteVM [%s], unmarshal dqResp, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return outErr
	}
	dqUpid := dqResp["data"].(string) // used for check task status
	beego.Info(fmt.Sprintf("Wait for Task [%s] is finished.", dqUpid))
	if err := p.waitForTaskFinished(WaitForTimeOut, 5, dqUpid); err != nil {
		outErr := fmt.Errorf("Cloud name [%s], type [%s], DeleteVM [%s], waitForTaskFinished, DeleteQemu, error: %w", p.Name, p.Type, vmid, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("Task [%s] is already finished successfully.", dqUpid))

	beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], Delete VM: %s", p.Name, p.Type, vmid))

	return nil
}

// After we create a task, we can use this polling function to wait for it to be finished
func (p *Proxmox) waitForTaskFinished(timeout int, checkInterval int, upid string) error {
	return MyWaitFor(timeout, checkInterval, func() (bool, error) {
		taskStatusBytes, err := p.GetTaskStatus(upid)
		if err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get task status upid [%s] error %w:", p.Name, p.Type, upid, err)
			beego.Error(outErr)
			return false, nil
		}
		beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], get task status upid %s, response: %s", p.Name, p.Type, upid, string(taskStatusBytes)))
		if err := p.CheckErrInResp(taskStatusBytes); err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get task status upid %s, error in resp: %w", p.Name, p.Type, upid, err)
			beego.Error(outErr)
			return false, nil
		}

		var taskStatus map[string]interface{}
		if err := json.Unmarshal(taskStatusBytes, &taskStatus); err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get task status upid %s, Unmarshal qemuBytes, error: %w", p.Name, p.Type, upid, err)
			beego.Error(outErr)
			return false, nil
		}

		if taskStatus["data"].(map[string]interface{})["status"].(string) == ProxTSRunning {
			beego.Info(fmt.Sprintf("task status is still %s", ProxTSRunning))
			return false, nil
		}

		if taskStatus["data"].(map[string]interface{})["status"].(string) != ProxTSStopped {
			outErr := fmt.Errorf("Abnormal! Task status is neither %s nor %s", ProxTSRunning, ProxTSStopped)
			beego.Error(outErr)
			return false, outErr
		}

		// this is what we expect
		if taskStatus["data"].(map[string]interface{})["exitstatus"].(string) == "OK" {
			beego.Info(fmt.Sprintf("Successful! Task status is %s, and exitstatus is %s", ProxTSStopped, taskStatus["data"].(map[string]interface{})["exitstatus"].(string)))
			return true, nil
		}

		outErr := fmt.Errorf("Abnormal! Task status is %s, but exitstatus is %s", ProxTSStopped, taskStatus["data"].(map[string]interface{})["exitstatus"].(string))
		beego.Error(outErr)
		return false, outErr
	})
}

// After a VM is created, the config is locked in a short time. We need to wait until it is unlocked.
func (p *Proxmox) waitForUnlock(timeout int, checkInterval int, vmid string) error {
	return MyWaitFor(timeout, checkInterval, func() (bool, error) {

		qemuBytes, err := p.GetQemu(vmid)
		if err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu id [%s], error: %w", p.Name, p.Type, vmid, err)
			beego.Error(outErr)
			return false, nil
		}
		beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], get qemu id %s, response: %s", p.Name, p.Type, vmid, string(qemuBytes)))
		if err := p.CheckErrInResp(qemuBytes); err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu id %s, error in resp: %w", p.Name, p.Type, vmid, err)
			beego.Error(outErr)
			return false, nil
		}
		var qemu map[string]interface{}
		if err := json.Unmarshal(qemuBytes, &qemu); err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get VM id %s, Unmarshal qemuBytes, error: %w", p.Name, p.Type, vmid, err)
			beego.Error(outErr)
			return false, nil
		}

		// see whether it is locked
		if lockState, exist := qemu["data"].(map[string]interface{})["lock"]; exist {
			beego.Info(fmt.Sprintf("VM with ID [%s] is still locked. The lock state is [%v]", vmid, lockState))
			return false, nil
		}

		// If the code reaches here, it means the VM is unlocked
		return true, nil
	})
}

// After a VM is created, the disk cannot be resized in a short time (I know it in my practice). We need to retry until the resize is successful.
func (p *Proxmox) waitForResizeDisk(timeout int, checkInterval int, vmid int, diskName string, diskSize string) error {
	return MyWaitFor(timeout, checkInterval, func() (bool, error) {
		resizeDiskRespBytes, err := p.ResizeDisk(vmid, diskName, diskSize)
		time.Sleep(ProxmoxAPIInterval)
		if err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%d], config disk [%s] size [%s], error: %w", p.Name, p.Type, vmid, diskName, diskSize, err)
			beego.Error(outErr)
			return false, nil
		}
		if err := p.CheckErrInResp(resizeDiskRespBytes); err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%d], config disk [%s] size [%s], in resp error: %w", p.Name, p.Type, vmid, diskName, diskSize, err)
			beego.Error(outErr)
			return false, nil
		}
		beego.Info(fmt.Sprintf("Request successful! Cloud name [%s], type [%s], CreateVM [%d], config disk [%s] size [%s].", p.Name, p.Type, vmid, diskName, diskSize))

		return true, nil
	})
}

// To avoid a Proxmox bug, we need to resize the disk size 2 times.
// Deprecated: From https://forum.proxmox.com/threads/bug-a-bug-about-the-api-to-resize-disk-in-version-7-3-3.123400/post-536856, I know that the better method to avoid the Proxmox bug is to use qm rescan command.
func (p *Proxmox) waitForAddDiskSize(timeout int, checkInterval int, vmid int, diskName string, addedSize string, expectedSize int) error {
	return MyWaitFor(timeout, checkInterval, func() (bool, error) {
		// send request to resize the disk
		addedDiskRespBytes, err := p.ResizeDisk(vmid, diskName, addedSize)
		time.Sleep(ProxmoxAPIInterval)
		if err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%d], config disk [%s] size [%s], error: %w", p.Name, p.Type, vmid, diskName, addedSize, err)
			beego.Error(outErr)
			return false, nil
		}
		if err := p.CheckErrInResp(addedDiskRespBytes); err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], CreateVM [%d], config disk [%s] size [%s], in resp error: %w", p.Name, p.Type, vmid, diskName, addedSize, err)
			beego.Error(outErr)
			return false, nil
		}
		beego.Info(fmt.Sprintf("Successful! Cloud name [%s], type [%s], CreateVM [%d], config disk [%s] size [%s].", p.Name, p.Type, vmid, diskName, addedSize))

		// Then we get the VM disk size to see its value, only when we can get the resized value, it is successful.
		beego.Info("Then we need to get the disk size to see whether we can get the resized value.")
		vmidStr := strconv.Itoa(vmid)
		qemuBytes, err := p.GetQemu(vmidStr)
		if err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get VM id [%s], get qemu, error: %w", p.Name, p.Type, vmidStr, err)
			beego.Error(outErr)
			return false, nil
		}
		beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], get qemu id %s, response: %s", p.Name, p.Type, vmidStr, string(qemuBytes)))

		if err := p.CheckErrInResp(qemuBytes); err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get VM id [%s], get qemu, error in resp: %w", p.Name, p.Type, vmidStr, err)
			beego.Error(outErr)
			return false, nil
		}

		var qemu map[string]interface{}
		if err := json.Unmarshal(qemuBytes, &qemu); err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get VM id %s, Unmarshal qemuBytes, error: %w", p.Name, p.Type, vmidStr, err)
			beego.Error(outErr)
			return false, nil
		}

		// Only when we can get the resized value, it is successful
		if qemu["data"].(map[string]interface{})["maxdisk"].(float64)/1024/1024/1024 < float64(expectedSize) {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], VM id %s, disk size %fG, smaller than expected %dG", p.Name, p.Type, vmidStr, qemu["data"].(map[string]interface{})["maxdisk"].(float64)/1024/1024/1024, expectedSize)
			beego.Error(outErr)
			return false, nil
		}

		return true, nil
	})
}

// After we shut down a VM, we can use this polling function to make sure it is stopped.
func (p *Proxmox) waitForQemuStopped(timeout int, checkInterval int, vmid string) error {
	return MyWaitFor(timeout, checkInterval, func() (bool, error) {
		qemuBytes, err := p.GetQemu(vmid)
		if err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu id [%s], error: %w", p.Name, p.Type, vmid, err)
			beego.Error(outErr)
			return false, nil
		}
		beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], get qemu id %s, response: %s", p.Name, p.Type, vmid, string(qemuBytes)))

		if err := p.CheckErrInResp(qemuBytes); err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu id [%s], error in resp: %w", p.Name, p.Type, vmid, err)
			beego.Error(outErr)
			return false, nil
		}

		var qemu map[string]interface{}
		if err := json.Unmarshal(qemuBytes, &qemu); err != nil {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], get qemu id [%s], Unmarshal qemuBytes, error: %w", p.Name, p.Type, vmid, err)
			beego.Error(outErr)
			return false, nil
		}

		if qemu["data"].(map[string]interface{})["status"].(string) != ProxQSStopped {
			outErr := fmt.Errorf("Cloud name [%s], type [%s], qemu id [%s], status: %s, not %s", p.Name, p.Type, vmid, qemu["data"].(map[string]interface{})["status"].(string), ProxQSStopped)
			beego.Error(outErr)
			return false, nil
		}

		beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], qemu id [%s], status is %s, it's already %s.", p.Name, p.Type, vmid, qemu["data"].(map[string]interface{})["status"].(string), ProxQSStopped))
		return true, nil
	})
}

func (p *Proxmox) CheckErrInResp(respByte []byte) error {
	beego.Info(fmt.Sprintf("check resp: [%s]", string(respByte)))
	var resp map[string]interface{}
	if err := json.Unmarshal(respByte, &resp); err != nil {
		outErr := fmt.Errorf("CheckErrInResp, json.Unmarshal, error: %w", err)
		beego.Error(outErr)
		return outErr
	}
	if err, exist := resp["errors"]; exist {
		outErr := fmt.Errorf("In resp, error: %v", err)
		beego.Error(outErr)
		return outErr
	}
	return nil
}
