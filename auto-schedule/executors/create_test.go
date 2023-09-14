package executors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

func TestInnerAddScheInfoToApps(t *testing.T) {
	testCases := []struct {
		name           string
		apps           []models.K8sApp
		scheSoln       asmodel.Solution
		expectedResult []models.K8sApp
	}{
		{
			name: "case1",
			apps: []models.K8sApp{
				models.K8sApp{
					Priority:      2,
					AutoScheduled: true,
					Name:          "group-printtime",
					Replicas:      1,
					HostNetwork:   false,
					Containers: []models.K8sContainer{
						models.K8sContainer{
							Name:    "printtime",
							Image:   "172.27.15.31:5000/printtime:v1",
							WorkDir: "/printtime",
							Resources: models.K8sResReq{
								Limits: models.K8sResList{
									Memory:  "30Mi",
									CPU:     "1.2",
									Storage: "2Gi",
								},
								Requests: models.K8sResList{
									Memory:  "30Mi",
									CPU:     "1.2",
									Storage: "2Gi",
								},
							},
							Commands: []string{
								"bash",
							},
							Args: []string{
								"-c",
								"python3 -u main.py > $LOGFILE",
							},
							Env: []models.K8sEnv{
								{
									Name:  "PARAMETER1",
									Value: "testRenderenv1",
								},
								{
									Name:  "LOGFILE",
									Value: "/tmp/234/printtime.log",
								},
							},
							Mounts: []models.K8sMount{
								{
									VmPath:        "/tmp/asdff",
									ContainerPath: "/tmp/234",
								},
								{
									VmPath:        "/tmp/uyyyy",
									ContainerPath: "/tmp/2345",
								},
							},
						},
					},
					Dependencies: []models.Dependency{
						{
							AppName: "group-nginx",
						},
						{
							AppName: "group-ubuntu",
						},
					},
				},
				models.K8sApp{
					Priority:      4,
					AutoScheduled: true,
					Name:          "group-nginx",
					Replicas:      1,
					HostNetwork:   true,
					Containers: []models.K8sContainer{
						models.K8sContainer{
							Name:    "nginx",
							Image:   "172.27.15.31:5000/nginx:1.17.1",
							WorkDir: "",
							Resources: models.K8sResReq{
								Limits: models.K8sResList{
									Memory:  "1024Mi",
									CPU:     "2.1",
									Storage: "20Gi",
								},
								Requests: models.K8sResList{
									Memory:  "1024Mi",
									CPU:     "2.1",
									Storage: "20Gi",
								},
							},
							Ports: []models.PortInfo{
								{
									ContainerPort: 80,
									Name:          "fsd",
									Protocol:      "tcp",
									ServicePort:   "80",
									NodePort:      "30001",
								},
							},
						},
					},
					Dependencies: []models.Dependency{
						{
							AppName: "group-ubuntu",
						},
					},
				},
				models.K8sApp{
					Priority:      4,
					AutoScheduled: true,
					Name:          "group-ubuntu",
					Replicas:      1,
					HostNetwork:   true,
					Containers: []models.K8sContainer{
						models.K8sContainer{
							Name:    "ubuntu",
							Image:   "172.27.15.31:5000/ubuntu:latest",
							WorkDir: "",
							Resources: models.K8sResReq{
								Limits: models.K8sResList{
									Memory:  "512Mi",
									CPU:     "2.1",
									Storage: "20Gi",
								},
								Requests: models.K8sResList{
									Memory:  "512Mi",
									CPU:     "2.1",
									Storage: "20Gi",
								},
							},
							Commands: []string{
								"bash",
								"-c",
								"while true;do sleep 10;done",
							},
							Args: nil,
							Env: []models.K8sEnv{
								{
									Name:  "asfasf",
									Value: "asfasf",
								},
								{
									Name:  "asdfsdf",
									Value: "sfsdf",
								},
							},
							Mounts: []models.K8sMount{
								{
									VmPath:        "/tmp/asdff",
									ContainerPath: "/tmp/log",
								},
							},
							Ports: nil,
						},
					},
					Dependencies: []models.Dependency{},
				},
			},
			scheSoln: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"group-nginx": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA7",
						K8sNodeName:      "auto-sched-nokia7-0",
						AllocatedCpuCore: 2.1,
					},
					"group-printtime": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA10",
						K8sNodeName:      "auto-sched-nokia10-0",
						AllocatedCpuCore: 1,
					},
					"group-ubuntu": asmodel.SingleAppSolution{
						Accepted:         false,
						TargetCloudName:  "",
						K8sNodeName:      "",
						AllocatedCpuCore: 0,
					},
				},
				VmsToCreate: []models.IaasVm{
					models.IaasVm{
						ID:        "",
						Name:      "auto-sched-nokia7-0",
						IPs:       nil,
						VCpu:      16,
						Ram:       38639,
						Storage:   279,
						Status:    "",
						Cloud:     "NOKIA7",
						CloudType: "",
						McmCreate: false,
					},
					models.IaasVm{
						ID:        "",
						Name:      "auto-sched-nokia10-0",
						IPs:       nil,
						VCpu:      2,
						Ram:       4896,
						Storage:   44,
						Status:    "",
						Cloud:     "NOKIA10",
						CloudType: "",
						McmCreate: false,
					},
				},
			},
			expectedResult: []models.K8sApp{
				models.K8sApp{
					Priority:      2,
					AutoScheduled: true,
					Name:          "group-printtime",
					Replicas:      1,
					HostNetwork:   false,
					NodeName:      "auto-sched-nokia10-0",
					Containers: []models.K8sContainer{
						models.K8sContainer{
							Name:    "printtime",
							Image:   "172.27.15.31:5000/printtime:v1",
							WorkDir: "/printtime",
							Resources: models.K8sResReq{
								Limits: models.K8sResList{
									Memory:  "30Mi",
									CPU:     "1.0",
									Storage: "2Gi",
								},
								Requests: models.K8sResList{
									Memory:  "30Mi",
									CPU:     "1.0",
									Storage: "2Gi",
								},
							},
							Commands: []string{
								"bash",
							},
							Args: []string{
								"-c",
								"python3 -u main.py > $LOGFILE",
							},
							Env: []models.K8sEnv{
								{
									Name:  "PARAMETER1",
									Value: "testRenderenv1",
								},
								{
									Name:  "LOGFILE",
									Value: "/tmp/234/printtime.log",
								},
							},
							Mounts: []models.K8sMount{
								{
									VmPath:        "/tmp/asdff",
									ContainerPath: "/tmp/234",
								},
								{
									VmPath:        "/tmp/uyyyy",
									ContainerPath: "/tmp/2345",
								},
							},
						},
					},
					Dependencies: []models.Dependency{
						{
							AppName: "group-nginx",
						},
						{
							AppName: "group-ubuntu",
						},
					},
				},
				models.K8sApp{
					Priority:      4,
					AutoScheduled: true,
					Name:          "group-nginx",
					Replicas:      1,
					HostNetwork:   true,
					NodeName:      "auto-sched-nokia7-0",
					Containers: []models.K8sContainer{
						models.K8sContainer{
							Name:    "nginx",
							Image:   "172.27.15.31:5000/nginx:1.17.1",
							WorkDir: "",
							Resources: models.K8sResReq{
								Limits: models.K8sResList{
									Memory:  "1024Mi",
									CPU:     "2.1",
									Storage: "20Gi",
								},
								Requests: models.K8sResList{
									Memory:  "1024Mi",
									CPU:     "2.1",
									Storage: "20Gi",
								},
							},
							Ports: []models.PortInfo{
								{
									ContainerPort: 80,
									Name:          "fsd",
									Protocol:      "tcp",
									ServicePort:   "80",
									NodePort:      "30001",
								},
							},
						},
					},
					Dependencies: []models.Dependency{
						{
							AppName: "group-ubuntu",
						},
					},
				},
			},
		},
		{
			name: "case2",
			apps: []models.K8sApp{
				models.K8sApp{
					Priority:      2,
					AutoScheduled: true,
					Name:          "group-printtime",
					Replicas:      1,
					HostNetwork:   false,
					Containers: []models.K8sContainer{
						models.K8sContainer{
							Name:    "printtime",
							Image:   "172.27.15.31:5000/printtime:v1",
							WorkDir: "/printtime",
							Resources: models.K8sResReq{
								Limits: models.K8sResList{
									Memory:  "30Mi",
									CPU:     "1.2",
									Storage: "2Gi",
								},
								Requests: models.K8sResList{
									Memory:  "30Mi",
									CPU:     "1.2",
									Storage: "2Gi",
								},
							},
							Commands: []string{
								"bash",
							},
							Args: []string{
								"-c",
								"python3 -u main.py > $LOGFILE",
							},
							Env: []models.K8sEnv{
								{
									Name:  "PARAMETER1",
									Value: "testRenderenv1",
								},
								{
									Name:  "LOGFILE",
									Value: "/tmp/234/printtime.log",
								},
							},
							Mounts: []models.K8sMount{
								{
									VmPath:        "/tmp/asdff",
									ContainerPath: "/tmp/234",
								},
								{
									VmPath:        "/tmp/uyyyy",
									ContainerPath: "/tmp/2345",
								},
							},
						},
					},
					Dependencies: []models.Dependency{
						{
							AppName: "group-nginx",
						},
						{
							AppName: "group-ubuntu",
						},
					},
				},
				models.K8sApp{
					Priority:      4,
					AutoScheduled: true,
					Name:          "group-nginx",
					Replicas:      1,
					HostNetwork:   true,
					Containers: []models.K8sContainer{
						models.K8sContainer{
							Name:    "nginx",
							Image:   "172.27.15.31:5000/nginx:1.17.1",
							WorkDir: "",
							Resources: models.K8sResReq{
								Limits: models.K8sResList{
									Memory:  "1024Mi",
									CPU:     "2.1",
									Storage: "20Gi",
								},
								Requests: models.K8sResList{
									Memory:  "1024Mi",
									CPU:     "2.1",
									Storage: "20Gi",
								},
							},
							Ports: []models.PortInfo{
								{
									ContainerPort: 80,
									Name:          "fsd",
									Protocol:      "tcp",
									ServicePort:   "80",
									NodePort:      "30001",
								},
							},
						},
					},
					Dependencies: []models.Dependency{
						{
							AppName: "group-ubuntu",
						},
					},
				},
				models.K8sApp{
					Priority:      4,
					AutoScheduled: true,
					Name:          "group-ubuntu",
					Replicas:      1,
					HostNetwork:   true,
					Containers: []models.K8sContainer{
						models.K8sContainer{
							Name:    "ubuntu",
							Image:   "172.27.15.31:5000/ubuntu:latest",
							WorkDir: "",
							Resources: models.K8sResReq{
								Limits: models.K8sResList{
									Memory:  "512Mi",
									CPU:     "2.1",
									Storage: "20Gi",
								},
								Requests: models.K8sResList{
									Memory:  "512Mi",
									CPU:     "2.1",
									Storage: "20Gi",
								},
							},
							Commands: []string{
								"bash",
								"-c",
								"while true;do sleep 10;done",
							},
							Args: nil,
							Env: []models.K8sEnv{
								{
									Name:  "asfasf",
									Value: "asfasf",
								},
								{
									Name:  "asdfsdf",
									Value: "sfsdf",
								},
							},
							Mounts: []models.K8sMount{
								{
									VmPath:        "/tmp/asdff",
									ContainerPath: "/tmp/log",
								},
							},
							Ports: nil,
						},
					},
					Dependencies: []models.Dependency{},
				},
			},
			scheSoln: asmodel.Solution{
				AppsSolution: map[string]asmodel.SingleAppSolution{
					"group-nginx": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA7",
						K8sNodeName:      "auto-sched-nokia7-0",
						AllocatedCpuCore: 1.100000001,
					},
					"group-printtime": asmodel.SingleAppSolution{
						Accepted: false,
					},
					"group-ubuntu": asmodel.SingleAppSolution{
						Accepted:         true,
						TargetCloudName:  "NOKIA8",
						K8sNodeName:      "auto-sched-nokia8-0",
						AllocatedCpuCore: 0.7999999,
					},
				},
				VmsToCreate: []models.IaasVm{
					models.IaasVm{
						ID:        "",
						Name:      "auto-sched-nokia7-0",
						IPs:       nil,
						VCpu:      16,
						Ram:       38639,
						Storage:   279,
						Status:    "",
						Cloud:     "NOKIA7",
						CloudType: "",
						McmCreate: false,
					},
					models.IaasVm{
						ID:        "",
						Name:      "auto-sched-nokia10-0",
						IPs:       nil,
						VCpu:      2,
						Ram:       4896,
						Storage:   44,
						Status:    "",
						Cloud:     "NOKIA10",
						CloudType: "",
						McmCreate: false,
					},
				},
			},
			expectedResult: []models.K8sApp{
				models.K8sApp{
					Priority:      4,
					AutoScheduled: true,
					Name:          "group-nginx",
					Replicas:      1,
					HostNetwork:   true,
					NodeName:      "auto-sched-nokia7-0",
					Containers: []models.K8sContainer{
						models.K8sContainer{
							Name:    "nginx",
							Image:   "172.27.15.31:5000/nginx:1.17.1",
							WorkDir: "",
							Resources: models.K8sResReq{
								Limits: models.K8sResList{
									Memory:  "1024Mi",
									CPU:     "1.1",
									Storage: "20Gi",
								},
								Requests: models.K8sResList{
									Memory:  "1024Mi",
									CPU:     "1.1",
									Storage: "20Gi",
								},
							},
							Ports: []models.PortInfo{
								{
									ContainerPort: 80,
									Name:          "fsd",
									Protocol:      "tcp",
									ServicePort:   "80",
									NodePort:      "30001",
								},
							},
						},
					},
					Dependencies: []models.Dependency{
						{
							AppName: "group-ubuntu",
						},
					},
				},
				models.K8sApp{
					Priority:      4,
					AutoScheduled: true,
					Name:          "group-ubuntu",
					Replicas:      1,
					HostNetwork:   true,
					NodeName:      "auto-sched-nokia8-0",
					Containers: []models.K8sContainer{
						models.K8sContainer{
							Name:    "ubuntu",
							Image:   "172.27.15.31:5000/ubuntu:latest",
							WorkDir: "",
							Resources: models.K8sResReq{
								Limits: models.K8sResList{
									Memory:  "512Mi",
									CPU:     "0.8",
									Storage: "20Gi",
								},
								Requests: models.K8sResList{
									Memory:  "512Mi",
									CPU:     "0.8",
									Storage: "20Gi",
								},
							},
							Commands: []string{
								"bash",
								"-c",
								"while true;do sleep 10;done",
							},
							Args: nil,
							Env: []models.K8sEnv{
								{
									Name:  "asfasf",
									Value: "asfasf",
								},
								{
									Name:  "asdfsdf",
									Value: "sfsdf",
								},
							},
							Mounts: []models.K8sMount{
								{
									VmPath:        "/tmp/asdff",
									ContainerPath: "/tmp/log",
								},
							},
							Ports: nil,
						},
					},
					Dependencies: []models.Dependency{},
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Logf("test: %d, %s", i, testCase.name)
		actualResult := addScheInfoToApps(testCase.apps, testCase.scheSoln)
		assert.Equal(t, testCase.expectedResult, actualResult, fmt.Sprintf("%s: result is not expected", testCase.name))
	}
}
