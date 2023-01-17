package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	storagequota "github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	computequota "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	networkquota "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/quotas"
)

type Openstack struct {
	Name          string
	Type          string
	ProjectID     string
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
