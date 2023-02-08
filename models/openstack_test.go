package models

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	storagequota "github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	computequota "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	networkingquota "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/quotas"
)

const testOsCloudName = "CLAAUDIAweifan"

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

func TestNoConfig(t *testing.T) {
	InitClouds()
	for _, cloud := range Clouds {
		switch cloud.(type) {
		case *Openstack:
			fmt.Printf("Cloud: %s, root password is\n", cloud.(*Openstack).Name)
			fmt.Println(cloud.(*Openstack).RootPasswd)
			fmt.Printf("%t\n", len(cloud.(*Openstack).RootPasswd) > 0)
		}
	}
}

func TestGetComputeQuota(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		computeQuota, _ := cloud.(*Openstack).GetComputeQuota()
		fmt.Printf("%+v\n", computeQuota)
	}
}

func TestGetNetworkQuota(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		networkQuota, _ := cloud.(*Openstack).GetNetworkQuota()
		fmt.Printf("%+v\n", networkQuota)
	}

}

func TestGetStorageQuota(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		storageQuota, _ := cloud.(*Openstack).GetStorageQuota()
		fmt.Printf("%+v\n", storageQuota)
	}

}

func TestListAllVolumes(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		allVolumes, _ := cloud.(*Openstack).ListAllVolumes()
		for j := 0; j < len(allVolumes); j++ {
			vol, _ := cloud.(*Openstack).GetVolume(allVolumes[j].ID)
			fmt.Printf("%+v\n", vol)
		}
	}

}

func TestCreateVolume(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		opts := volumes.CreateOpts{
			Size:    150,
			Name:    "volume-unittest3",
			ImageID: cloud.(*Openstack).ImageID,
		}
		vol, _ := cloud.(*Openstack).CreateVolume(opts)
		fmt.Printf("%+v\n", vol)
	}

}

func TestGetVolume(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		id := "a5afa30f-28d8-46c8-8568-58e068cbde32"
		vol, _ := cloud.(*Openstack).GetVolume(id)
		fmt.Printf("%+v\n", vol)
	}
}

func TestDeleteVolume(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		id := "a5afa30f-28d8-46c8-8568-58e068cbde32"
		_ = cloud.(*Openstack).DeleteVolume(id)
	}
}

func TestListAllFavors(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		allFlavors, _ := cloud.(*Openstack).ListAllFavors()
		for j := 0; j < len(allFlavors); j++ {
			fmt.Printf("%+v\n", allFlavors[j])
		}
	}
}

func TestGetFlavor(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		id := "0238fdc1-2525-4669-be22-a545341c8301"
		vol, _ := cloud.(*Openstack).GetFlavor(id)
		fmt.Printf("%+v\n", vol)
	}
}

func TestListAllServers(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		allServers, _ := cloud.(*Openstack).ListAllServers()
		for j := 0; j < len(allServers); j++ {
			fmt.Printf("%+v\n", allServers[j])
		}
	}

}

func TestCreateServer(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		baseOpts := servers.CreateOpts{
			Name:           "unittest3",
			FlavorRef:      "abb9477b-955c-45fa-bfaf-b53ebc8b2cb7",
			SecurityGroups: []string{cloud.(*Openstack).SecurityGroup},
			Networks: []servers.Network{
				{UUID: cloud.(*Openstack).NetworkID},
			},
		}
		optsWithKeyPair := keypairs.CreateOptsExt{
			CreateOptsBuilder: baseOpts,
			KeyName:           cloud.(*Openstack).KeyName,
		}
		optsBfv := bootfromvolume.CreateOptsExt{
			CreateOptsBuilder: optsWithKeyPair,
			BlockDevice: []bootfromvolume.BlockDevice{
				{
					BootIndex:           0,
					DeleteOnTermination: false,
					UUID:                "a5afa30f-28d8-46c8-8568-58e068cbde32",
					SourceType:          bootfromvolume.SourceVolume,
					DestinationType:     bootfromvolume.DestinationVolume,
				},
			},
		}
		newServer, _ := cloud.(*Openstack).CreateServer(optsBfv)
		fmt.Printf("%+v\n", newServer)
	}

}

func TestGetServerAndExtractIPs(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		id := "0b0d3d61-360d-4cb2-93a1-9f2073e13b50"
		server, err := cloud.(*Openstack).GetServer(id)
		if err != nil {
			t.Fatalf("get server error: %s", err.Error())
		}
		fmt.Printf("%+v\n", server)
		fmt.Println(cloud.(*Openstack).ExtractIPs(server))
	}
}

func TestDeleteServer(t *testing.T) {
	InitClouds()
	cloud := Clouds[testOsCloudName]
	switch cloud.(type) {
	case *Openstack:
		id := "0b0d3d61-360d-4cb2-93a1-9f2073e13b50"
		_ = cloud.(*Openstack).DeleteServer(id)
	}
}
