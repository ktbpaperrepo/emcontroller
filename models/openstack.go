package models

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	storagequota "github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
	computequota "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	networkquota "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/quotas"
)

type Openstack struct {
	Name          string
	Type          string
	ProjectID     string
	ImageID       string // we use a fixed image to create VMs
	NetworkID     string // we use a fixed network to create VMs
	SecurityGroup string // we use a fixed security group to create VMs
	KeyName       string // we use a fixed key pair to create VMs
	Provider      *gophercloud.ProviderClient
	ComputeClient *gophercloud.ServiceClient
	NetworkClient *gophercloud.ServiceClient
	StorageClient *gophercloud.ServiceClient
}

func InitOpenstack(paras map[string]interface{}) *Openstack {
	beego.Info(fmt.Sprintf("Start to initialize cloud name [%s] type [%s]", paras["name"].(string), paras["type"].(string)))
	opts := gophercloud.AuthOptions{
		IdentityEndpoint:            paras["authurl"].(string),
		ApplicationCredentialID:     paras["applicationcredentialid"].(string),
		ApplicationCredentialSecret: paras["applicationcredentialsecret"].(string),
		AllowReauth:                 true,
	}

	// init provider
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		panic(fmt.Errorf("cloud name [%s] type [%s], openstack.AuthenticatedClient error: %w", paras["name"].(string), paras["type"].(string), err))
	}

	// init compute client
	computeClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: paras["region"].(string),
	})
	if err != nil {
		panic(fmt.Errorf("cloud name [%s] type [%s], openstack.NewComputeV2 error: %w", paras["name"].(string), paras["type"].(string), err))
	}

	// init network client
	networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Name:   "neutron",
		Region: "RegionOne",
	})
	if err != nil {
		panic(fmt.Errorf("cloud name [%s] type [%s], openstack.NewNetworkV2 error: %w", paras["name"].(string), paras["type"].(string), err))
	}

	// init storage client
	storageClient, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		panic(fmt.Errorf("cloud name [%s] type [%s], openstack.NewBlockStorageV3 error: %w", paras["name"].(string), paras["type"].(string), err))
	}

	return &Openstack{
		Name:          paras["name"].(string),
		Type:          paras["type"].(string),
		ProjectID:     paras["project_id"].(string),
		ImageID:       paras["imageid"].(string),
		NetworkID:     paras["networkid"].(string),
		SecurityGroup: paras["securitygroup"].(string),
		KeyName:       paras["keyname"].(string),
		Provider:      provider,
		ComputeClient: computeClient,
		NetworkClient: networkClient,
		StorageClient: storageClient,
	}
}

func (os *Openstack) CreateVM() {

}

func (os *Openstack) DeleteVM() {

}

func (os *Openstack) CheckResources() ResourceStatus {
	var computeQuota computequota.QuotaDetailSet = os.GetComputeQuota()
	var networkQuota *networkquota.QuotaDetailSet = os.GetNetworkQuota()
	var storageQuota storagequota.QuotaUsageSet = os.GetStorageQuota()

	return ResourceStatus{
		Limit: ResSet{
			VCpu:    float64(computeQuota.Cores.Limit),
			Ram:     float64(computeQuota.RAM.Limit),
			Vm:      float64(computeQuota.Instances.Limit),
			Volume:  float64(storageQuota.Volumes.Limit),
			Storage: float64(storageQuota.Gigabytes.Limit),
			Port:    float64(networkQuota.Port.Limit),
		},
		InUse: ResSet{
			VCpu:    float64(computeQuota.Cores.InUse),
			Ram:     float64(computeQuota.RAM.InUse),
			Vm:      float64(computeQuota.Instances.InUse),
			Volume:  float64(storageQuota.Volumes.InUse),
			Storage: float64(storageQuota.Gigabytes.InUse),
			Port:    float64(networkQuota.Port.Used),
		},
	}

}

func (os *Openstack) GetComputeQuota() computequota.QuotaDetailSet {
	quotaResult := computequota.GetDetail(os.ComputeClient, os.ProjectID)
	extracted, err := quotaResult.Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Check compute quota error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	return extracted
}

func (os *Openstack) GetNetworkQuota() *networkquota.QuotaDetailSet {
	quotaResult := networkquota.GetDetail(os.NetworkClient, os.ProjectID)
	extracted, err := quotaResult.Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Check network quota error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	return extracted
}

func (os *Openstack) GetStorageQuota() storagequota.QuotaUsageSet {
	quotaResult := storagequota.GetUsage(os.StorageClient, os.ProjectID)
	extracted, err := quotaResult.Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Check storage quota error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	return extracted
}

// the listed volumes only have the attribute ID, other attributes are empty, so we should use GetVolume to see the details of each one. See the unit test.
func (os *Openstack) ListAllVolumes() []volumes.Volume {
	opts := volumes.ListOpts{}
	allPages, err := volumes.List(os.StorageClient, opts).AllPages()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], list volumes error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], volumes.ExtractVolumes error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	return allVolumes
}

func (os *Openstack) CreateVolume(opts volumes.CreateOpts) *volumes.Volume {
	beego.Info(fmt.Sprintf("Create volume opts: [%+v]", opts))
	vol, err := volumes.Create(os.StorageClient, opts).Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Create Volume opts [%+v] error: %s", os.Name, os.Type, os.ProjectID, opts, err.Error()))
	} else {
		beego.Info(fmt.Sprintf("Successful! Create volume opts: [%+v]", opts))
	}
	return vol
}

func (os *Openstack) GetVolume(id string) *volumes.Volume {
	vol, err := volumes.Get(os.StorageClient, id).Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Get Volume id [%s] error: %s", os.Name, os.Type, os.ProjectID, id, err.Error()))
	}
	return vol
}

func (os *Openstack) DeleteVolume(id string) {
	beego.Info(fmt.Sprintf("Delete volume id: [%s]", id))
	res := volumes.Delete(os.StorageClient, id)
	beego.Info(fmt.Sprintf("Delete volume id: [%s] response:\n%+v\n", id, res))
	if res.Err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Delete Volume id [%s] error: %s", os.Name, os.Type, os.ProjectID, id, res.Err.Error()))
	} else {
		beego.Info(fmt.Sprintf("Successful! Delete volume id: [%s]", id))
	}
}

func (os *Openstack) ListAllFavors() []flavors.Flavor {
	opts := flavors.ListOpts{}
	allPages, err := flavors.ListDetail(os.ComputeClient, opts).AllPages()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], list flavors error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	allFlavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], flavors.ExtractFlavors error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	return allFlavors
}

func (os *Openstack) ListAllServers() []servers.Server {
	opts := servers.ListOpts{}
	allPages, err := servers.List(os.ComputeClient, opts).AllPages()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], list servers error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], servers.ExtractServers error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	return allServers
}

func (os *Openstack) CreateServer(opts servers.CreateOptsBuilder) *servers.Server {
	beego.Info(fmt.Sprintf("Create server opts: [%+v]", opts))
	newServer, err := servers.Create(os.ComputeClient, opts).Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Create Server opts [%+v] error: %s", os.Name, os.Type, os.ProjectID, opts, err.Error()))
	} else {
		beego.Info(fmt.Sprintf("Successful! Create server opts: [%+v]", opts))
	}
	return newServer
}

func (os *Openstack) GetServer(id string) *servers.Server {
	server, err := servers.Get(os.ComputeClient, id).Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Get Server id [%s] error: %s", os.Name, os.Type, os.ProjectID, id, err.Error()))
	}
	return server
}

func (os *Openstack) DeleteServer(id string) {
	beego.Info(fmt.Sprintf("Delete server id: [%s]", id))
	res := servers.Delete(os.ComputeClient, id)
	beego.Info(fmt.Sprintf("Delete server id: [%s] response:\n%+v\n", id, res))
	if res.Err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Delete Server id [%s] error: %s", os.Name, os.Type, os.ProjectID, id, res.Err.Error()))
	} else {
		beego.Info(fmt.Sprintf("Successful! Delete server id: [%s]", id))
	}
}
