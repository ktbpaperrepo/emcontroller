package model

import (
	"fmt"
	"testing"

	"github.com/KeepTheBeats/routing-algorithms/random"
	"github.com/stretchr/testify/assert"

	"emcontroller/models"
)

func TestCloudMapCopy(t *testing.T) {
	testCases := []struct {
		name string
		src  map[string]Cloud
	}{
		{
			name: "case1",
			src: map[string]Cloud{
				"cloud1": Cloud{
					Name: "cloud1",
					K8sNodes: []K8sNode{
						K8sNode{
							Name: "node1",
						},
						K8sNode{
							Name: "node2",
						},
					},
				},
				"cloud2": Cloud{
					Name: "cloud2",
					K8sNodes: []K8sNode{
						K8sNode{
							Name: "node3",
						},
						K8sNode{
							Name: "node4",
						},
					},
				},
			},
		},
		{
			name: "case2",
			src: map[string]Cloud{
				"cloud1": Cloud{
					Name: "cloud1",
					K8sNodes: []K8sNode{
						K8sNode{
							Name: "node10",
						},
						K8sNode{
							Name: "node20",
						},
					},
				},
				"cloud2": Cloud{
					Name: "cloud2",
					K8sNodes: []K8sNode{
						K8sNode{
							Name: "node30",
						},
						K8sNode{
							Name: "node40",
						},
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)

		var dst map[string]Cloud

		// test extra key
		dst = CloudMapCopy(testCase.src)
		extraKey := "cloud100"
		dst[extraKey] = Cloud{
			Name: extraKey,
			K8sNodes: []K8sNode{
				K8sNode{
					Name: "node300",
				},
				K8sNode{
					Name: "node400",
				},
			},
		}
		if app, exist := testCase.src[extraKey]; exist {
			t.Errorf("key \"%s\" should not exist in testCase.src, the value is: %+v", extraKey, app)
		} else {
			t.Logf("As expected, key \"%s\" does not exist in testCase.src", extraKey)
		}

		// test change
		dst = CloudMapCopy(testCase.src)
		changeKey := "cloud1"
		dst[changeKey] = Cloud{
			Name: "cloud1",
			K8sNodes: []K8sNode{
				K8sNode{
					Name: "node150",
				},
				K8sNode{
					Name: "node250",
				},
			},
		}

		t.Logf("testCase.src[changeKey]:\n%+v\ndst[changeKey]:\n%+v", testCase.src[changeKey], dst[changeKey])
		assert.NotEqual(t, testCase.src[changeKey], dst[changeKey], fmt.Sprintf("testCase.src[changeKey] [%+v] and dst[changeKey] [%+v], show not be equal.", testCase.src[changeKey], dst[changeKey]))

	}

}

func TestGetInfoVmToCreate(t *testing.T) {
	// test 1
	func() {
		cloud := Cloud{
			Name: "NOKIA4",
			Resources: models.ResourceStatus{
				Limit: models.ResSet{
					VCpu:    56,
					Ram:     128796.75390625,
					Storage: 1396.5185890197754,
					Vm:      -1,
					Port:    -1,
					Volume:  -1,
				},
				InUse: models.ResSet{
					VCpu:    56,
					Ram:     128796.75390625,
					Storage: 1396.5185890197754,
					Vm:      -1,
					Port:    -1,
					Volume:  -1,
				},
			},
			K8sNodes: []K8sNode{},
		}

		for i := 0; i < 5; i++ {
			resPct := random.RandomFloat64(0, 1)
			t.Logf("resource percent: %f\n", resPct)
			newNode := cloud.GetInfoVmToCreate(resPct)
			t.Logf("%+v\n", newNode)
			cloud.K8sNodes = append(cloud.K8sNodes, newNode)
		}

		resPct := 1.0
		t.Logf("resource percent: %f\n", resPct)
		newNode := cloud.GetInfoVmToCreate(resPct)
		t.Logf("%+v\n", newNode)

		fmt.Println()
		t.Logf("%+v\n", cloud)
	}()

	// test 2
	func() {
		cloud := Cloud{
			Name: "NOKIA4",
			Resources: models.ResourceStatus{
				Limit: models.ResSet{
					VCpu:    56,
					Ram:     128796.75390625,
					Storage: 1396.5185890197754,
					Vm:      -1,
					Port:    -1,
					Volume:  -1,
				},
				InUse: models.ResSet{
					VCpu:    56,
					Ram:     128796.75390625,
					Storage: 1396.5185890197754,
					Vm:      -1,
					Port:    -1,
					Volume:  -1,
				},
			},
			K8sNodes: []K8sNode{
				K8sNode{Name: "auto-sched-nokia4-2"},
				K8sNode{Name: "auto-sched-nokia4-4"},
			},
		}

		for i := 0; i < 5; i++ {
			resPct := random.RandomFloat64(0, 1)
			t.Logf("resource percent: %f\n", resPct)
			newNode := cloud.GetInfoVmToCreate(resPct)
			t.Logf("%+v\n", newNode)
			cloud.K8sNodes = append(cloud.K8sNodes, newNode)
		}

		resPct := 1.0
		t.Logf("resource percent: %f\n", resPct)
		newNode := cloud.GetInfoVmToCreate(resPct)
		t.Logf("%+v\n", newNode)

		fmt.Println()
		t.Logf("%+v\n", cloud)
	}()
}

func TestSupportCreateNewVM(t *testing.T) {
	testCases := []struct {
		name           string
		cloud          Cloud
		expectedResult bool
	}{
		{
			name: "case-proxmox",
			cloud: Cloud{
				Name: "NOKIA4",
				Type: models.ProxmoxIaas,
			},
			expectedResult: true,
		},
		{
			name: "case-openstack",
			cloud: Cloud{
				Name: "CLAAUDIAweifan",
				Type: models.OpenstackIaas,
			},
			expectedResult: false,
		},
		{
			name: "case-other",
			cloud: Cloud{
				Name: "HPE2",
				Type: "proxmox (powered off)",
			},
			expectedResult: false,
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := testCase.cloud.SupportCreateNewVM()
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}
