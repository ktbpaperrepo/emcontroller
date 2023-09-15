package algorithms

import (
	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

const testDelta float64 = floatDelta

func cloudAppsSolnForTest() (asmodel.Cloud, map[string]asmodel.Application, asmodel.Solution) {
	cloud := asmodel.Cloud{
		Name: "cloud1",
	}
	apps := map[string]asmodel.Application{
		"app1": asmodel.Application{
			Name:     "app1",
			Priority: 5,
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 2.4,
					Memory:  1024,
					Storage: 10,
				},
			},
		},
		"app2": asmodel.Application{
			Name:     "app2",
			Priority: 10,
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 3.9,
					Memory:  990,
					Storage: 15,
				},
			},
		},
		"app3": asmodel.Application{
			Name:     "app3",
			Priority: 1,
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 1.4,
					Memory:  990,
					Storage: 15,
				},
			},
		},
		"app4": asmodel.Application{
			Name:     "app4",
			Priority: 3,
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 5.0,
					Memory:  660,
					Storage: 6,
				},
			},
		},
		"app5": asmodel.Application{
			Name:     "app5",
			Priority: 2,
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 5.0,
					Memory:  990,
					Storage: 15,
				},
			},
		},
		"app6": asmodel.Application{
			Name:     "app6",
			Priority: 7,
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 4.0,
					Memory:  990,
					Storage: 15,
				},
			},
		},
		"app7": asmodel.Application{
			Name:     "app7",
			Priority: 10,
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 3.0,
					Memory:  540,
					Storage: 35,
				},
			},
		},
		"app8": asmodel.Application{
			Name:     "app8",
			Priority: 8,
			Resources: asmodel.AppResources{
				GenericResources: asmodel.GenericResources{
					CpuCore: 2.0,
					Memory:  540,
					Storage: 15,
				},
			},
		},
	}
	soln := asmodel.Solution{
		AppsSolution: map[string]asmodel.SingleAppSolution{
			"app1": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud1",
				K8sNodeName:     "auto-sched-cloud1-1",
			},
			"app2": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud1",
				K8sNodeName:     "auto-sched-cloud1-1",
			},
			"app3": asmodel.SingleAppSolution{
				Accepted: false,
			},
			"app4": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud1",
				K8sNodeName:     "auto-sched-cloud1-2",
			},
			"app5": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud2",
			},
			"app6": asmodel.SingleAppSolution{
				Accepted: false,
			},
			"app7": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud1",
				K8sNodeName:     "auto-sched-cloud1-3",
			},
			"app8": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud3",
			},
		},
	}
	return cloud, apps, soln
}

func cloudsForTest() map[string]asmodel.Cloud {
	var clouds map[string]asmodel.Cloud = make(map[string]asmodel.Cloud)

	clouds["nokia4"] = asmodel.Cloud{
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
				VCpu:    26,
				Ram:     59392,
				Storage: 629,
				Vm:      -1,
				Port:    -1,
				Volume:  -1,
			},
		},
		K8sNodes: []asmodel.K8sNode{},
	}

	clouds["nokia4WithOneNode"] = asmodel.Cloud{
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
				VCpu:    26,
				Ram:     59392,
				Storage: 629,
				Vm:      -1,
				Port:    -1,
				Volume:  -1,
			},
		},
		K8sNodes: []asmodel.K8sNode{{Name: "auto-sched-nokia4-0"}},
	}

	clouds["nokia4With2OriNode"] = asmodel.Cloud{
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
				VCpu:    26,
				Ram:     59392,
				Storage: 629,
				Vm:      -1,
				Port:    -1,
				Volume:  -1,
			},
		},
		K8sNodes: []asmodel.K8sNode{
			{Name: "nokia4-ori-node1"},
			{Name: "nokia4-ori-node2"},
		},
	}

	clouds["nokia4With2SchedNode"] = asmodel.Cloud{
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
				VCpu:    26,
				Ram:     59392,
				Storage: 629,
				Vm:      -1,
				Port:    -1,
				Volume:  -1,
			},
		},
		K8sNodes: []asmodel.K8sNode{
			{Name: "auto-sched-nokia4-4"},
			{Name: "auto-sched-nokia4-5"},
		},
	}

	return clouds
}

func cloudsWithNetForTest() []map[string]asmodel.Cloud {
	return []map[string]asmodel.Cloud{
		map[string]asmodel.Cloud{ // index 0
			"NOKIA4": asmodel.Cloud{
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: 0.753,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA7": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
				},
			},
			"NOKIA5": asmodel.Cloud{
				Name: "NOKIA5",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: 0.753,
					},
					"NOKIA6": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA7": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
				},
			},
			"NOKIA6": asmodel.Cloud{
				Name: "NOKIA6",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: 0.753,
					},
					"NOKIA7": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
				},
			},
			"NOKIA7": asmodel.Cloud{
				Name: "NOKIA7",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA7": models.NetworkState{
						Rtt: 0.753,
					},
				},
			},
		},
		map[string]asmodel.Cloud{ // index 1
			"NOKIA4": asmodel.Cloud{
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA7": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
				},
			},
			"NOKIA5": asmodel.Cloud{
				Name: "NOKIA5",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA7": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
				},
			},
			"NOKIA6": asmodel.Cloud{
				Name: "NOKIA6",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA7": models.NetworkState{
						Rtt: 0.753,
					},
				},
			},
			"NOKIA7": asmodel.Cloud{
				Name: "NOKIA7",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: 0.5,
					},
					"NOKIA7": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
				},
			},
		},
		map[string]asmodel.Cloud{ // index 2
			"NOKIA4": asmodel.Cloud{
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA7": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
				},
			},
			"NOKIA5": asmodel.Cloud{
				Name: "NOKIA5",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA7": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
				},
			},
			"NOKIA6": asmodel.Cloud{
				Name: "NOKIA6",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: 0.753,
					},
					"NOKIA7": models.NetworkState{
						Rtt: 0.753,
					},
				},
			},
			"NOKIA7": asmodel.Cloud{
				Name: "NOKIA7",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: 0.5,
					},
					"NOKIA7": models.NetworkState{
						Rtt: 0.753,
					},
				},
			},
		},
		map[string]asmodel.Cloud{ // index 3
			"NOKIA4": asmodel.Cloud{
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: 2.124,
					},
					"NOKIA5": models.NetworkState{
						Rtt: 3.236,
					},
					"NOKIA6": models.NetworkState{
						Rtt: 0.578,
					},
					"NOKIA7": models.NetworkState{
						Rtt: 5.2,
					},
				},
			},
			"NOKIA5": asmodel.Cloud{
				Name: "NOKIA5",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA7": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
				},
			},
			"NOKIA6": asmodel.Cloud{
				Name: "NOKIA6",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: 0.753,
					},
					"NOKIA7": models.NetworkState{
						Rtt: 0.753,
					},
				},
			},
			"NOKIA7": asmodel.Cloud{
				Name: "NOKIA7",
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
						VCpu:    26,
						Ram:     59392,
						Storage: 629,
						Vm:      -1,
						Port:    -1,
						Volume:  -1,
					},
				},
				K8sNodes: []asmodel.K8sNode{},
				NetState: map[string]models.NetworkState{
					"NOKIA4": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA5": models.NetworkState{
						Rtt: models.UnreachableRttMs,
					},
					"NOKIA6": models.NetworkState{
						Rtt: 0.5,
					},
					"NOKIA7": models.NetworkState{
						Rtt: 0.753,
					},
				},
			},
		},
	}
}

func appOrdersForTest() [][]string {
	return [][]string{
		[]string{"app1", "app2", "app3", "app4", "app5", "app6", "app7", "app8"}, // index 0
		[]string{"app3", "app2", "app5", "app1", "app7", "app6", "app8", "app4"}, // index 1
	}
}

func appsForTest() []map[string]asmodel.Application {
	return []map[string]asmodel.Application{
		map[string]asmodel.Application{ // index 0
			"app1": asmodel.Application{
				Name:     "app1",
				Priority: 5,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 3.3,
						Memory:  1024,
						Storage: 10,
					},
				},
			},
			"app2": asmodel.Application{
				Name:     "app2",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 1.4,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app3": asmodel.Application{
				Name:     "app3",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 2.2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app4": asmodel.Application{
				Name:     "app4",
				Priority: 3,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.1,
						Memory:  660,
						Storage: 6,
					},
				},
			},
			"app5": asmodel.Application{
				Name:     "app5",
				Priority: 2,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 6.3,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app6": asmodel.Application{
				Name:     "app6",
				Priority: 7,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 3.4,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app7": asmodel.Application{
				Name:     "app7",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.5,
						Memory:  540,
						Storage: 35,
					},
				},
			},
			"app8": asmodel.Application{
				Name:     "app8",
				Priority: 8,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.6,
						Memory:  540,
						Storage: 15,
					},
				},
			},
		},
		map[string]asmodel.Application{ // index 1
			"app1": asmodel.Application{
				Name:     "app1",
				Priority: 5,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.2,
						Memory:  1024,
						Storage: 10,
					},
				},
			},
			"app2": asmodel.Application{
				Name:     "app2",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 1.3,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app3": asmodel.Application{
				Name:     "app3",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 2.2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app4": asmodel.Application{
				Name:     "app4",
				Priority: 3,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.6,
						Memory:  660,
						Storage: 6,
					},
				},
			},
			"app5": asmodel.Application{
				Name:     "app5",
				Priority: 2,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 6.8,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app6": asmodel.Application{
				Name:     "app6",
				Priority: 7,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 2.3,
						Memory:  990,
						Storage: 15,
					},
				},
			},
		},
		map[string]asmodel.Application{ // index 2
			"app1": asmodel.Application{
				Name:     "app1",
				Priority: 5,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.2,
						Memory:  1024,
						Storage: 10,
					},
				},
			},
			"app2": asmodel.Application{
				Name:     "app2",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 1.3,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app4": asmodel.Application{
				Name:     "app4",
				Priority: 3,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.6,
						Memory:  660,
						Storage: 6,
					},
				},
			},
		},
		map[string]asmodel.Application{ // index 3
			"app2": asmodel.Application{
				Name:     "app2",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 1.3,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app3": asmodel.Application{
				Name:     "app3",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 2.2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app7": asmodel.Application{
				Name:     "app7",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.5,
						Memory:  540,
						Storage: 35,
					},
				},
			},
		},
		map[string]asmodel.Application{ // index 4
			"app1": asmodel.Application{
				Name:     "app1",
				Priority: 5,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.2,
						Memory:  1024,
						Storage: 10,
					},
				},
			},
			"app4": asmodel.Application{
				Name:     "app4",
				Priority: 3,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.6,
						Memory:  660,
						Storage: 6,
					},
				},
			},
			"app6": asmodel.Application{
				Name:     "app6",
				Priority: 7,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 2.3,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app7": asmodel.Application{
				Name:     "app7",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.5,
						Memory:  540,
						Storage: 35,
					},
				},
			},
		},
		map[string]asmodel.Application{ // index 5
			"app1": asmodel.Application{
				Name:     "app1",
				Priority: 5,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 3.3,
						Memory:  1024,
						Storage: 10,
					},
				},
			},
			"app2": asmodel.Application{
				Name:     "app2",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 1.4,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app3": asmodel.Application{
				Name:     "app3",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 2.2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app4": asmodel.Application{
				Name:     "app4",
				Priority: 3,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5.1,
						Memory:  660,
						Storage: 6,
					},
				},
			},
			"app6": asmodel.Application{
				Name:     "app6",
				Priority: 2,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 3.4,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app7": asmodel.Application{
				Name:     "app7",
				Priority: 2,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.5,
						Memory:  540,
						Storage: 35,
					},
				},
			},
			"app8": asmodel.Application{
				Name:     "app8",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.6,
						Memory:  540,
						Storage: 15,
					},
				},
			},
		},
		map[string]asmodel.Application{ // index 6
			"app1": asmodel.Application{
				Name:     "app1",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.3,
						Memory:  1024,
						Storage: 10,
					},
				},
			},
			"app2": asmodel.Application{
				Name:     "app2",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app3": asmodel.Application{
				Name:     "app3",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app4": asmodel.Application{
				Name:     "app4",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.1,
						Memory:  660,
						Storage: 6,
					},
				},
			},
			"app6": asmodel.Application{
				Name:     "app6",
				Priority: 5,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.3,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app7": asmodel.Application{
				Name:     "app7",
				Priority: 2,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.4,
						Memory:  540,
						Storage: 35,
					},
				},
			},
			"app8": asmodel.Application{
				Name:     "app8",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.1,
						Memory:  540,
						Storage: 15,
					},
				},
			},
		},
		map[string]asmodel.Application{ // index 7
			"app1": asmodel.Application{
				Name:     "app1",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 3,
						Memory:  1024,
						Storage: 10,
					},
				},
			},
			"app2": asmodel.Application{
				Name:     "app2",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app3": asmodel.Application{
				Name:     "app3",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app4": asmodel.Application{
				Name:     "app4",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 1,
						Memory:  660,
						Storage: 6,
					},
				},
			},
			"app6": asmodel.Application{
				Name:     "app6",
				Priority: 5,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 5,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app7": asmodel.Application{
				Name:     "app7",
				Priority: 2,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 4,
						Memory:  540,
						Storage: 35,
					},
				},
			},
			"app8": asmodel.Application{
				Name:     "app8",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 1,
						Memory:  540,
						Storage: 15,
					},
				},
			},
		},
		map[string]asmodel.Application{ // index 8
			"app1": asmodel.Application{
				Name:     "app1",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 3,
						Memory:  1024,
						Storage: 10,
					},
				},
			},
			"app2": asmodel.Application{
				Name:     "app2",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app3": asmodel.Application{
				Name:     "app3",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app4": asmodel.Application{
				Name:     "app4",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 1,
						Memory:  660,
						Storage: 6,
					},
				},
			},
			"app5": asmodel.Application{
				Name:     "app5",
				Priority: 2,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 12,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app6": asmodel.Application{
				Name:     "app6",
				Priority: 5,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 20,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app7": asmodel.Application{
				Name:     "app7",
				Priority: 2,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 4,
						Memory:  540,
						Storage: 35,
					},
				},
			},
			"app8": asmodel.Application{
				Name:     "app8",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 1,
						Memory:  540,
						Storage: 15,
					},
				},
			},
		},
		map[string]asmodel.Application{ // index 9
			"app1": asmodel.Application{
				Name:     "app1",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.3,
						Memory:  1024,
						Storage: 10,
					},
				},
				Dependencies: []models.Dependency{
					models.Dependency{
						AppName: "app2",
					},
					models.Dependency{
						AppName: "app3",
					},
				},
			},
			"app2": asmodel.Application{
				Name:     "app2",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.2,
						Memory:  990,
						Storage: 15,
					},
				},
				Dependencies: []models.Dependency{
					models.Dependency{
						AppName: "app4",
					},
					models.Dependency{
						AppName: "app3",
					},
				},
			},
			"app3": asmodel.Application{
				Name:     "app3",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app4": asmodel.Application{
				Name:     "app4",
				Priority: 10,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.1,
						Memory:  660,
						Storage: 6,
					},
				},
			},
			"app5": asmodel.Application{
				Name:     "app5",
				Priority: 2,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 6.3,
						Memory:  990,
						Storage: 15,
					},
				},
				Dependencies: []models.Dependency{
					models.Dependency{
						AppName: "app6",
					},
					models.Dependency{
						AppName: "app7",
					},
				},
			},
			"app6": asmodel.Application{
				Name:     "app6",
				Priority: 5,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 2,
						Memory:  990,
						Storage: 15,
					},
				},
			},
			"app7": asmodel.Application{
				Name:     "app7",
				Priority: 2,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.4,
						Memory:  540,
						Storage: 35,
					},
				},
				Dependencies: []models.Dependency{
					models.Dependency{
						AppName: "app8",
					},
				},
			},
			"app8": asmodel.Application{
				Name:     "app8",
				Priority: 1,
				Resources: asmodel.AppResources{
					GenericResources: asmodel.GenericResources{
						CpuCore: 0.1,
						Memory:  540,
						Storage: 15,
					},
				},
			},
		},
	}
}

func solnsForTest() []asmodel.Solution {
	return []asmodel.Solution{
		asmodel.Solution{ // index 0
			AppsSolution: map[string]asmodel.SingleAppSolution{
				"app1": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "cloud1",
					K8sNodeName:     "auto-sched-cloud1-1",
				},
				"app2": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "cloud1",
					K8sNodeName:     "auto-sched-cloud1-1",
				},
				"app3": asmodel.SingleAppSolution{
					Accepted: false,
				},
				"app4": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "cloud1",
					K8sNodeName:     "auto-sched-cloud1-2",
				},
				"app5": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "cloud2",
				},
				"app6": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "cloud1",
					K8sNodeName:     "auto-sched-cloud1-2",
				},
				"app7": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "cloud1",
					K8sNodeName:     "auto-sched-cloud1-3",
				},
				"app8": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "cloud3",
				},
			},
			VmsToCreate: []models.IaasVm{},
		},
		asmodel.Solution{ // index 1
			AppsSolution: map[string]asmodel.SingleAppSolution{
				"app1": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA4",
					K8sNodeName:     "nokia4-ori-node1",
				},
				"app2": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA4",
					K8sNodeName:     "nokia4-ori-node2",
				},
				"app3": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA4",
					K8sNodeName:     "auto-sched-nokia4-4",
				},
				"app4": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "nokia6-ori-node1",
				},
				"app5": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA8",
					K8sNodeName:     "auto-sched-nokia8-0",
				},
				"app6": asmodel.SingleAppSolution{
					Accepted: false,
				},
				"app7": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
			},
			VmsToCreate: []models.IaasVm{
				models.IaasVm{
					Name:    "auto-sched-nokia4-4",
					Cloud:   "NOKIA4",
					VCpu:    8,
					Ram:     3122,
					Storage: 44,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia6-3",
					Cloud:   "NOKIA6",
					VCpu:    10,
					Ram:     2567,
					Storage: 80,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia8-0",
					Cloud:   "NOKIA8",
					VCpu:    7,
					Ram:     4000,
					Storage: 100,
				},
			},
		},
		asmodel.Solution{ // index 2
			AppsSolution: map[string]asmodel.SingleAppSolution{
				"app1": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA4",
					K8sNodeName:     "nokia4-ori-node1",
				},
				"app2": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA4",
					K8sNodeName:     "nokia4-ori-node1",
				},
				"app3": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA4",
					K8sNodeName:     "auto-sched-nokia4-4",
				},
				"app4": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA4",
					K8sNodeName:     "nokia4-ori-node1",
				},
				"app5": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA8",
					K8sNodeName:     "auto-sched-nokia8-0",
				},
				"app6": asmodel.SingleAppSolution{
					Accepted: false,
				},
				"app7": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
			},
			VmsToCreate: []models.IaasVm{
				models.IaasVm{
					Name:    "auto-sched-nokia4-4",
					Cloud:   "NOKIA4",
					VCpu:    8,
					Ram:     3122,
					Storage: 44,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia6-3",
					Cloud:   "NOKIA6",
					VCpu:    10,
					Ram:     2567,
					Storage: 80,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia8-0",
					Cloud:   "NOKIA8",
					VCpu:    7,
					Ram:     4000,
					Storage: 100,
				},
			},
		},
		asmodel.Solution{ // index 3
			AppsSolution: map[string]asmodel.SingleAppSolution{
				"app1": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA4",
					K8sNodeName:     "nokia4-ori-node1",
				},
				"app2": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app3": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app4": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA4",
					K8sNodeName:     "nokia4-ori-node1",
				},
				"app5": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA8",
					K8sNodeName:     "auto-sched-nokia8-0",
				},
				"app6": asmodel.SingleAppSolution{
					Accepted: false,
				},
				"app7": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
			},
			VmsToCreate: []models.IaasVm{
				models.IaasVm{
					Name:    "auto-sched-nokia4-4",
					Cloud:   "NOKIA4",
					VCpu:    8,
					Ram:     3122,
					Storage: 44,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia6-3",
					Cloud:   "NOKIA6",
					VCpu:    10,
					Ram:     2567,
					Storage: 80,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia8-0",
					Cloud:   "NOKIA8",
					VCpu:    7,
					Ram:     4000,
					Storage: 100,
				},
			},
		},
		asmodel.Solution{ // index 4
			AppsSolution: map[string]asmodel.SingleAppSolution{
				"app1": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app2": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA4",
					K8sNodeName:     "auto-sched-nokia4-3",
				},
				"app3": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA8",
					K8sNodeName:     "auto-sched-nokia8-2",
				},
				"app4": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app5": asmodel.SingleAppSolution{
					Accepted: false,
				},
				"app6": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app7": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app8": asmodel.SingleAppSolution{
					Accepted: false,
				},
			},
			VmsToCreate: []models.IaasVm{
				models.IaasVm{
					Name:    "auto-sched-nokia4-4",
					Cloud:   "NOKIA4",
					VCpu:    8,
					Ram:     3122,
					Storage: 44,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia6-3",
					Cloud:   "NOKIA6",
					VCpu:    10,
					Ram:     2567,
					Storage: 80,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia8-0",
					Cloud:   "NOKIA8",
					VCpu:    7,
					Ram:     4000,
					Storage: 100,
				},
			},
		},
		asmodel.Solution{ // index 5
			AppsSolution: map[string]asmodel.SingleAppSolution{
				"app1": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app2": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app3": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app4": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app5": asmodel.SingleAppSolution{
					Accepted: false,
				},
				"app6": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app7": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
				"app8": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
					K8sNodeName:     "auto-sched-nokia6-3",
				},
			},
			VmsToCreate: []models.IaasVm{
				models.IaasVm{
					Name:    "auto-sched-nokia4-4",
					Cloud:   "NOKIA4",
					VCpu:    8,
					Ram:     3122,
					Storage: 44,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia6-3",
					Cloud:   "NOKIA6",
					VCpu:    10,
					Ram:     2567,
					Storage: 80,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia8-0",
					Cloud:   "NOKIA8",
					VCpu:    7,
					Ram:     4000,
					Storage: 100,
				},
			},
		},
		asmodel.Solution{ // index 6
			AppsSolution: map[string]asmodel.SingleAppSolution{
				"app1": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA1",
				},
				"app2": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA2",
				},
				"app3": asmodel.SingleAppSolution{
					Accepted: false,
				},
				"app4": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
				},
				"app5": asmodel.SingleAppSolution{
					Accepted: false,
				},
				"app6": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA7",
				},
				"app7": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA10",
				},
				"app8": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "HPE1",
				},
			},
			VmsToCreate: []models.IaasVm{
				models.IaasVm{
					Name:    "auto-sched-nokia4-4",
					Cloud:   "NOKIA4",
					VCpu:    8,
					Ram:     3122,
					Storage: 44,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia6-3",
					Cloud:   "NOKIA6",
					VCpu:    10,
					Ram:     2567,
					Storage: 80,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia8-0",
					Cloud:   "NOKIA8",
					VCpu:    7,
					Ram:     4000,
					Storage: 100,
				},
			},
		},
		asmodel.Solution{ // index 7
			AppsSolution: map[string]asmodel.SingleAppSolution{
				"app1": asmodel.SingleAppSolution{
					Accepted:        false,
					TargetCloudName: "NOKIA2",
				},
				"app2": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA4",
				},
				"app3": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA6",
				},
				"app4": asmodel.SingleAppSolution{
					Accepted: false,
				},
				"app5": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "HPE2",
				},
				"app6": asmodel.SingleAppSolution{
					Accepted: false,
				},
				"app7": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA12",
				},
				"app8": asmodel.SingleAppSolution{
					Accepted:        true,
					TargetCloudName: "NOKIA2",
				},
			},
			VmsToCreate: []models.IaasVm{
				models.IaasVm{
					Name:    "auto-sched-nokia4-4",
					Cloud:   "NOKIA4",
					VCpu:    8,
					Ram:     3122,
					Storage: 44,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia6-3",
					Cloud:   "NOKIA6",
					VCpu:    10,
					Ram:     2567,
					Storage: 80,
				},
				models.IaasVm{
					Name:    "auto-sched-nokia8-0",
					Cloud:   "NOKIA8",
					VCpu:    7,
					Ram:     4000,
					Storage: 100,
				},
			},
		},
	}
}
