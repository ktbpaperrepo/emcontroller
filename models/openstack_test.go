package models

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	storagequota "github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	computequota "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	networkingquota "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/quotas"
	"testing"
)

func TestGophercloud(t *testing.T) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint:            "https://strato-new.claaudia.aau.dk:5000/v3",
		ApplicationCredentialID:     "b2205e0c72374d918f5b6ee14c4c40c5",
		ApplicationCredentialSecret: "weifan",
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Printf("openstack.AuthenticatedClient error: %s\n", err.Error())
	}

	// compute
	computeClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	result := computequota.GetDetail(computeClient, "833668770f4244e299517d63006c3b46")
	fmt.Printf("%#v\n\n", result)

	extracted, err := result.Extract()
	if err != nil {
		fmt.Printf("compute result.Extract, error: %s\n", err.Error())
	}
	fmt.Println("compute:", extracted)

	// network
	networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Name:   "neutron",
		Region: "RegionOne",
	})
	networkExtracted, err := networkingquota.GetDetail(networkClient, "833668770f4244e299517d63006c3b46").Extract()
	if err != nil {
		fmt.Printf("networking result.Extract, error: %s\n", err.Error())
	}
	fmt.Println("networking:", networkExtracted)

	// storage
	storageClient, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	storageExtracted, err := storagequota.GetUsage(storageClient, "833668770f4244e299517d63006c3b46").Extract()
	if err != nil {
		fmt.Printf("storageExtracted, error: %s\n", err.Error())
	}
	fmt.Println("storage:", storageExtracted)
}

func TestGetComputeQuota(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			fmt.Printf("%+v\n", Clouds[i].(*Openstack).GetComputeQuota())
		}
	}
}

func TestGetNetworkQuota(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			fmt.Printf("%+v\n", Clouds[i].(*Openstack).GetNetworkQuota())
		}
	}
}

func TestGetStorageQuota(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			fmt.Printf("%+v\n", Clouds[i].(*Openstack).GetStorageQuota())
		}
	}
}
