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
			computeQuota, _ := Clouds[i].(*Openstack).GetComputeQuota()
			fmt.Printf("%+v\n", computeQuota)
		}
	}
}

func TestGetNetworkQuota(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			networkQuota, _ := Clouds[i].(*Openstack).GetNetworkQuota()
			fmt.Printf("%+v\n", networkQuota)
		}
	}
}

func TestGetStorageQuota(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			storageQuota, _ := Clouds[i].(*Openstack).GetStorageQuota()
			fmt.Printf("%+v\n", storageQuota)
		}
	}
}

func TestListAllVolumes(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			allVolumes := Clouds[i].(*Openstack).ListAllVolumes()
			for j := 0; j < len(allVolumes); j++ {
				fmt.Printf("%+v\n", Clouds[i].(*Openstack).GetVolume(allVolumes[j].ID))
			}
		}
	}
}

func TestCreateVolume(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			opts := volumes.CreateOpts{
				Size:    150,
				Name:    "volume-unittest2",
				ImageID: Clouds[i].(*Openstack).ImageID,
			}
			vol, _ := Clouds[i].(*Openstack).CreateVolume(opts)
			fmt.Printf("%+v\n", vol)
		}
	}
}

func TestGetVolume(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			id := "3dcc23ec-ecf1-460c-9c90-890acd077367"
			vol := Clouds[i].(*Openstack).GetVolume(id)
			fmt.Printf("%+v\n", vol)
		}
	}
}

func TestDeleteVolume(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			id := "ec494289-2f8a-46c3-aa3e-3378ad2a3233"
			Clouds[i].(*Openstack).DeleteVolume(id)
		}
	}
}

func TestListAllFavors(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			allFlavors, _ := Clouds[i].(*Openstack).ListAllFavors()
			for j := 0; j < len(allFlavors); j++ {
				fmt.Printf("%+v\n", allFlavors[j])
			}
		}
	}
}

func TestListAllServers(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			allServers := Clouds[i].(*Openstack).ListAllServers()
			for j := 0; j < len(allServers); j++ {
				fmt.Printf("%+v\n", allServers[j])
			}
		}
	}
}

func TestCreateServer(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			baseOpts := servers.CreateOpts{
				Name:           "unittest2",
				FlavorRef:      "abb9477b-955c-45fa-bfaf-b53ebc8b2cb7",
				SecurityGroups: []string{Clouds[i].(*Openstack).SecurityGroup},
				Networks: []servers.Network{
					{UUID: Clouds[i].(*Openstack).NetworkID},
				},
			}
			optsWithKeyPair := keypairs.CreateOptsExt{
				CreateOptsBuilder: baseOpts,
				KeyName:           Clouds[i].(*Openstack).KeyName,
			}
			optsBfv := bootfromvolume.CreateOptsExt{
				CreateOptsBuilder: optsWithKeyPair,
				BlockDevice: []bootfromvolume.BlockDevice{
					{
						BootIndex:           0,
						DeleteOnTermination: false,
						UUID:                "7f4d550f-08e2-44f5-8aaf-2ed3a639e5a3",
						SourceType:          bootfromvolume.SourceVolume,
						DestinationType:     bootfromvolume.DestinationVolume,
					},
				},
			}
			newServer, _ := Clouds[i].(*Openstack).CreateServer(optsBfv)
			fmt.Printf("%+v\n", newServer)
		}
	}
}

func TestGetServerAndExtractIPs(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			id := "c6f7e22b-3e55-4772-9efa-8d072d4cfc31"
			server, _ := Clouds[i].(*Openstack).GetServer(id)
			fmt.Printf("%+v\n", server)
			fmt.Println(Clouds[i].(*Openstack).ExtractIPs(server))
		}
	}
}

func TestDeleteServer(t *testing.T) {
	InitClouds()
	for i := 0; i < len(Clouds); i++ {
		switch Clouds[i].(type) {
		case *Openstack:
			id := "f128b1f2-3574-44ef-b107-f664b6ce6967"
			Clouds[i].(*Openstack).DeleteServer(id)
		}
	}
}
