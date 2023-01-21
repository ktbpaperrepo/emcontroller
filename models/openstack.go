package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	storagequota "github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	computequota "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	networkquota "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/quotas"
	"strings"
)

type Openstack struct {
	Name          string
	Type          string
	ProjectID     string
	ImageID       string // we use a fixed image to create VMs
	NetworkID     string // we use a fixed network to create VMs
	SecurityGroup string // we use a fixed security group to create VMs
	KeyName       string // we use a fixed key pair to create VMs
	SshPemPath    string // the SSH identity file private key for the VMs on this cloud
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
		SshPemPath:    paras["sshpempath"].(string),
		Provider:      provider,
		ComputeClient: computeClient,
		NetworkClient: networkClient,
		StorageClient: storageClient,
	}
}

// the unit of vcpu, ram, storage in the input is consistent with ResSet
func (os *Openstack) CreateVM(name string, vcpu, ram, storage int) (*IaasVm, error) {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Create VM: %s", os.Name, os.Type, os.ProjectID, name))
	// 1. create a volume
	volumeOpts := volumes.CreateOpts{
		Size:    storage,
		Name:    "volume-" + name,
		ImageID: os.ImageID,
	}
	beego.Info(fmt.Sprintf("Create volume of VM %s", name))
	vol, err := os.CreateVolume(volumeOpts)
	if err != nil {
		outErr := fmt.Errorf("create volume of VM %s error: %w", name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Wait for finishing creating volume of VM %s", name))
	if err := volumes.WaitForStatus(os.StorageClient, vol.ID, "available", WaitForTimeOut); err != nil {
		outErr := fmt.Errorf("wait for finishing creating volume of VM %s error: %w", name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! Create volume of VM %s", name))

	// 2. create the VM
	// choose a proper flavor according to the input vcpu and ram
	beego.Info(fmt.Sprintf("start choosing flavor for vcpu: %d, ram %d MiB", vcpu, ram))
	allFlavors, err := os.ListAllFavors()
	if err != nil {
		outErr := fmt.Errorf("list flavors error: %w", err)
		beego.Error(outErr)
		return nil, outErr
	}
	chosenFlavor, found := os.ChooseMinFlavor(allFlavors, vcpu, ram)
	if !found {
		outErr := fmt.Errorf("no flavor can meet the vCPU %d and RAM %d", vcpu, ram)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Chosen flavor for vcpu: %d, ram %d MiB is:\n%+v", vcpu, ram, chosenFlavor))

	// create VM
	var metadata = map[string]string{
		McmSign: McmSign, // add this sign to the VM meaning that it is created by multi-cloud manager
	}
	baseVmOpts := servers.CreateOpts{
		Name:           name,
		Metadata:       metadata,
		FlavorRef:      chosenFlavor.ID,
		SecurityGroups: []string{os.SecurityGroup},
		Networks: []servers.Network{
			{UUID: os.NetworkID},
		},
	}
	vmOptsWithKeyPair := keypairs.CreateOptsExt{
		CreateOptsBuilder: baseVmOpts,
		KeyName:           os.KeyName,
	}
	vmOptsBfv := bootfromvolume.CreateOptsExt{
		CreateOptsBuilder: vmOptsWithKeyPair,
		BlockDevice: []bootfromvolume.BlockDevice{
			{
				BootIndex:           0,
				DeleteOnTermination: false,
				UUID:                vol.ID,
				SourceType:          bootfromvolume.SourceVolume,
				DestinationType:     bootfromvolume.DestinationVolume,
			},
		},
	}
	beego.Info(fmt.Sprintf("Create VM %s", name))
	vm, err := os.CreateServer(vmOptsBfv)
	if err != nil {
		outErr := fmt.Errorf("create VM %s error %w", name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Wait for VM %s status ACTIVE", name))
	if err := servers.WaitForStatus(os.ComputeClient, vm.ID, "ACTIVE", WaitForTimeOut); err != nil {
		outErr := fmt.Errorf("wait for VM %s status ACTIVE, error: %w", name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! Wait for VM %s status ACTIVE", name))

	curVM, err := os.GetServer(vm.ID)
	if err != nil {
		outErr := fmt.Errorf("get VM before waiting for SSH %s error: %w", name, err)
		beego.Error(outErr)
		return nil, outErr
	}
	sshIP := os.ExtractIPs(curVM)[0]

	beego.Info(fmt.Sprintf("Wait for VM %s able to be SSHed, ip %s", name, sshIP))

	if err := WaitForSshPem(SshUser, os.SshPemPath, sshIP, SshPort, WaitForTimeOut); err != nil {
		outErr := fmt.Errorf("wait for VM %s able to be SSHed, ip %s, error: %w", name, sshIP, err)
		beego.Error(outErr)
		return nil, outErr
	}
	beego.Info(fmt.Sprintf("Successful! Wait for VM %s able to be SSHed ip %s", name, sshIP))

	finishedVm, err := os.GetServer(vm.ID)
	if err != nil {
		outErr := fmt.Errorf("get finishedVm %s error: %w", name, err)
		beego.Error(outErr)
		return nil, outErr
	}

	var outIaasVm IaasVm = IaasVm{
		Name:      name,
		ID:        finishedVm.ID,
		IPs:       os.ExtractIPs(finishedVm),
		Cloud:     os.Name,
		CloudType: os.Type,
	}
	return &outIaasVm, nil
}

func (os *Openstack) DeleteVM(vmID string) error {
	beego.Info(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Delete VM: %s", os.Name, os.Type, os.ProjectID, vmID))
	// 1. delete server
	oriVM, err := os.GetServer(vmID)
	if err != nil {
		outErr := fmt.Errorf("get oriVM %s error: %w", vmID, err)
		beego.Error(outErr)
		return outErr
	}

	// VMs created by multi-cloud manager has this metadata. This method can only delete VMs created by multi-cloud manager
	if oriVM.Metadata[McmSign] != McmSign {
		outErr := fmt.Errorf("server %s is not created by multi-cloud manager, so we cannot delete it.", vmID)
		beego.Error(outErr)
		return outErr
	}

	beego.Info(fmt.Sprintf("Delete server: %+v\n", oriVM))
	err = os.DeleteServer(vmID)
	if err != nil {
		outErr := fmt.Errorf("Delete server %s error: %w", vmID, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("Successful! Delete server: %+v\n", oriVM))

	beego.Info(fmt.Sprintf("Wait for server %s deleted.", vmID))
	err = gophercloud.WaitFor(WaitForTimeOut, func() (bool, error) {
		beego.Info(fmt.Sprintf("Waiting ... try to get vm %s", vmID))
		if _, err := os.GetServer(vmID); err != nil {
			beego.Info(fmt.Sprintf("This time, get vm %s, error: %s", vmID, err.Error()))
			if IsOs404(err) {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		outErr := fmt.Errorf("Wait for server %s deleted, error: %w", vmID, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("Successful! Wait for server %s deleted.", vmID))

	// 2. delte the attached volumes of the server
	beego.Info(fmt.Sprintf("Start to delete the attached volumes of the server %s", vmID))
	for _, attachedVol := range oriVM.AttachedVolumes {
		beego.Info(fmt.Sprintf("Delete volume %s of server %s", attachedVol.ID, vmID))
		if err = os.DeleteVolume(attachedVol.ID); err != nil {
			outErr := fmt.Errorf("Delete volume %s of server %s, error: %w", attachedVol.ID, vmID, err)
			beego.Error(outErr)
			if !IsOs404(err) { // ignore not found error when deleting
				return outErr
			}
		}
		beego.Info(fmt.Sprintf("Successful! Delete volume %s of server %s", attachedVol.ID, vmID))

		beego.Info(fmt.Sprintf("Wait for volume %s of server %s deleted.", attachedVol.ID, vmID))
		err = gophercloud.WaitFor(WaitForTimeOut, func() (bool, error) {
			beego.Info(fmt.Sprintf("Waiting ... try to get volume %s", attachedVol.ID))
			if _, err := os.GetVolume(attachedVol.ID); err != nil {
				beego.Info(fmt.Sprintf("This time, get volume %s, error: %s", attachedVol.ID, err.Error()))
				if IsOs404(err) {
					return true, nil
				}
			}
			return false, nil
		})
		if err != nil {
			outErr := fmt.Errorf("Wait for volume %s of server %s deleted, error: %w", attachedVol.ID, vmID, err)
			beego.Error(outErr)
			return outErr
		}
		beego.Info(fmt.Sprintf("Successful! Wait for volume %s of server %s deleted.", attachedVol.ID, vmID))
	}
	return nil
}

func (os *Openstack) CheckResources() (ResourceStatus, error) {
	computeQuota, err := os.GetComputeQuota()
	if err != nil {
		return ResourceStatus{}, err
	}
	networkQuota, err := os.GetNetworkQuota()
	if err != nil {
		return ResourceStatus{}, err
	}
	storageQuota, err := os.GetStorageQuota()
	if err != nil {
		return ResourceStatus{}, err
	}

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
	}, nil
}

// check whether this error is the "not found" of Openstack
func IsOs404(err error) bool {
	return strings.Contains(err.Error(), Os404Substr)
}

// Choose the smallest flavor that can meet the requirements of RAM and vCPU
func (os *Openstack) ChooseMinFlavor(allFlavors []flavors.Flavor, reqCpu, reqRam int) (flavors.Flavor, bool) {
	var minFlavor flavors.Flavor
	var found bool = false
	computeQuota, err := os.GetComputeQuota()
	if err != nil {
		beego.Error(fmt.Sprintf("Get comput quota error: %s", err.Error()))
		return minFlavor, found
	}
	beego.Info(fmt.Sprintf("Try to find a flavor to meet vCPU %d, RAM %d MiB. The vCPU Limit quota is %d, in use is %d. The RAM Limit quota is %d MiB, in use is %d MiB", reqCpu, reqRam, computeQuota.Cores.Limit, computeQuota.Cores.InUse, computeQuota.RAM.Limit, computeQuota.RAM.InUse))
	remainingVcpu := computeQuota.Cores.Limit - computeQuota.Cores.InUse
	remainingRam := computeQuota.RAM.Limit - computeQuota.RAM.InUse
	for i := 0; i < len(allFlavors); i++ {
		if allFlavors[i].VCPUs < reqCpu || allFlavors[i].RAM < reqRam || allFlavors[i].VCPUs > remainingVcpu || allFlavors[i].RAM > remainingRam { // count meet the requirements
			continue
		}
		if !found { // the first one that can meet the requirements
			minFlavor = allFlavors[i]
			found = true
			continue
		}
		if os.overflowFlavor(allFlavors[i], reqCpu, reqRam) < os.overflowFlavor(minFlavor, reqCpu, reqRam) {
			minFlavor = allFlavors[i]
		}

	}
	return minFlavor, found
}

// Calculate the amount that a flavor overflows the required RAM and vCPU
func (os *Openstack) overflowFlavor(flavor flavors.Flavor, reqCpu, reqRam int) float64 {
	cpuOverflow := (float64(flavor.VCPUs) - float64(reqCpu)) / float64(reqCpu)
	ramOverflow := (float64(flavor.RAM) - float64(reqRam)) / float64(reqRam)
	return ramOverflow + cpuOverflow
}

// extract IP addresses from a server, imitating github.com\gophercloud\gophercloud@v1.1.1\acceptance\openstack\compute\v2\floatingip_test.go
func (os *Openstack) ExtractIPs(s *servers.Server) []string {
	var IPs []string
	for _, port := range s.Addresses {
		for _, networkAddresses := range port.([]interface{}) {
			address := networkAddresses.(map[string]interface{})
			IPs = append(IPs, address["addr"].(string))
		}
	}
	return IPs
}

func (os *Openstack) GetComputeQuota() (computequota.QuotaDetailSet, error) {
	quotaResult := computequota.GetDetail(os.ComputeClient, os.ProjectID)
	extracted, err := quotaResult.Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Check compute quota error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	return extracted, err
}

func (os *Openstack) GetNetworkQuota() (*networkquota.QuotaDetailSet, error) {
	quotaResult := networkquota.GetDetail(os.NetworkClient, os.ProjectID)
	extracted, err := quotaResult.Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Check network quota error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	return extracted, err
}

func (os *Openstack) GetStorageQuota() (storagequota.QuotaUsageSet, error) {
	quotaResult := storagequota.GetUsage(os.StorageClient, os.ProjectID)
	extracted, err := quotaResult.Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Check storage quota error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	return extracted, err
}

// the listed volumes only have the attribute ID, other attributes are empty, so we should use GetVolume to see the details of each one. See the unit test.
func (os *Openstack) ListAllVolumes() ([]volumes.Volume, error) {
	opts := volumes.ListOpts{}
	allPages, err := volumes.List(os.StorageClient, opts).AllPages()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], list volumes error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
		return []volumes.Volume{}, err
	}
	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], volumes.ExtractVolumes error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	return allVolumes, err
}

func (os *Openstack) CreateVolume(opts volumes.CreateOpts) (*volumes.Volume, error) {
	beego.Info(fmt.Sprintf("Create volume opts: [%+v]", opts))
	vol, err := volumes.Create(os.StorageClient, opts).Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Create Volume opts [%+v] error: %s", os.Name, os.Type, os.ProjectID, opts, err.Error()))
	} else {
		beego.Info(fmt.Sprintf("Successful! Create volume opts: [%+v]", opts))
	}
	return vol, err
}

func (os *Openstack) GetVolume(id string) (*volumes.Volume, error) {
	vol, err := volumes.Get(os.StorageClient, id).Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Get Volume id [%s] error: %s", os.Name, os.Type, os.ProjectID, id, err.Error()))
	}
	return vol, err
}

func (os *Openstack) DeleteVolume(id string) error {
	beego.Info(fmt.Sprintf("Delete volume id: [%s]", id))
	res := volumes.Delete(os.StorageClient, id)
	beego.Info(fmt.Sprintf("Delete volume id: [%s] response:\n%+v\n", id, res))
	if res.Err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Delete Volume id [%s] error: %s", os.Name, os.Type, os.ProjectID, id, res.Err.Error()))
		beego.Error(fmt.Sprintf("res: %+v\n", res))
	} else {
		beego.Info(fmt.Sprintf("Successful! Delete volume id: [%s]", id))
	}
	return res.Err
}

func (os *Openstack) ListAllFavors() ([]flavors.Flavor, error) {
	opts := flavors.ListOpts{}
	allPages, err := flavors.ListDetail(os.ComputeClient, opts).AllPages()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], list flavors error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
		return []flavors.Flavor{}, err
	}
	allFlavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], flavors.ExtractFlavors error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
		return allFlavors, err
	}
	return allFlavors, err
}

func (os *Openstack) ListAllServers() ([]servers.Server, error) {
	opts := servers.ListOpts{}
	allPages, err := servers.List(os.ComputeClient, opts).AllPages()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], list servers error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
		return []servers.Server{}, err
	}
	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], servers.ExtractServers error: %s", os.Name, os.Type, os.ProjectID, err.Error()))
	}
	return allServers, err
}

func (os *Openstack) CreateServer(opts servers.CreateOptsBuilder) (*servers.Server, error) {
	beego.Info(fmt.Sprintf("Create server opts: [%+v]", opts))
	newServer, err := servers.Create(os.ComputeClient, opts).Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Create Server opts [%+v] error: %s", os.Name, os.Type, os.ProjectID, opts, err.Error()))
	} else {
		beego.Info(fmt.Sprintf("Successful! Create server opts: [%+v]", opts))
	}
	return newServer, err
}

func (os *Openstack) GetServer(id string) (*servers.Server, error) {
	server, err := servers.Get(os.ComputeClient, id).Extract()
	if err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Get Server id [%s] error: %s", os.Name, os.Type, os.ProjectID, id, err.Error()))
	}
	return server, err
}

func (os *Openstack) DeleteServer(id string) error {
	beego.Info(fmt.Sprintf("Delete server id: [%s]", id))
	res := servers.Delete(os.ComputeClient, id)
	beego.Info(fmt.Sprintf("Delete server id: [%s] response:\n%+v\n", id, res))
	if res.Err != nil {
		beego.Error(fmt.Sprintf("Cloud name [%s], type [%s], project id [%s], Delete Server id [%s] error: %s", os.Name, os.Type, os.ProjectID, id, res.Err.Error()))
	} else {
		beego.Info(fmt.Sprintf("Successful! Delete server id: [%s]", id))
	}
	return res.Err
}
