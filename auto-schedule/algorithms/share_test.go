package algorithms

import (
	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

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
			},
			"app2": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud1",
			},
			"app3": asmodel.SingleAppSolution{
				Accepted: false,
			},
			"app4": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud1",
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

	return clouds
}
