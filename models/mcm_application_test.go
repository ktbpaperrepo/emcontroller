package models

import (
	"testing"
)

var (
	appToCreate1 = K8sApp{
		Priority:      4,
		AutoScheduled: true,
		Name:          "group-ubuntu",
		Replicas:      1,
		HostNetwork:   true,
		Containers: []K8sContainer{
			K8sContainer{
				Name:    "ubuntu",
				Image:   "172.27.15.31:5000/ubuntu:latest",
				WorkDir: "",
				Resources: K8sResReq{
					Limits: K8sResList{
						Memory:  "512Mi",
						CPU:     "2.1",
						Storage: "20Gi",
					},
					Requests: K8sResList{
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
				Env: []K8sEnv{
					{
						Name:  "asfasf",
						Value: "asfasf",
					},
					{
						Name:  "asdfsdf",
						Value: "sfsdf",
					},
				},
				Mounts: []K8sMount{
					{
						VmPath:        "/tmp/asdff",
						ContainerPath: "/tmp/log",
					},
				},
				Ports: nil,
			},
		},
		Dependencies: []Dependency{},
	}

	appToCreate2 = K8sApp{
		Priority:      4,
		AutoScheduled: true,
		Name:          "group-nginx",
		Replicas:      1,
		HostNetwork:   true,
		Containers: []K8sContainer{
			K8sContainer{
				Name:    "nginx",
				Image:   "172.27.15.31:5000/nginx:1.17.1",
				WorkDir: "",
				Resources: K8sResReq{
					Limits: K8sResList{
						Memory:  "1024Mi",
						CPU:     "2.1",
						Storage: "20Gi",
					},
					Requests: K8sResList{
						Memory:  "1024Mi",
						CPU:     "2.1",
						Storage: "20Gi",
					},
				},
				Ports: []PortInfo{
					{
						ContainerPort: 80,
						Name:          "fsd",
						Protocol:      "tcp",
						ServicePort:   "80",
					},
				},
			},
		},
		Dependencies: []Dependency{
			{
				AppName: "group-ubuntu",
			},
		},
	}

	appToCreate3 = K8sApp{
		Priority:      2,
		AutoScheduled: true,
		Name:          "group-printtime",
		Replicas:      1,
		HostNetwork:   false,
		Containers: []K8sContainer{
			K8sContainer{
				Name:    "printtime",
				Image:   "172.27.15.31:5000/printtime:v1",
				WorkDir: "/printtime",
				Resources: K8sResReq{
					Limits: K8sResList{
						Memory:  "30Mi",
						CPU:     "1.2",
						Storage: "2Gi",
					},
					Requests: K8sResList{
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
				Env: []K8sEnv{
					{
						Name:  "PARAMETER1",
						Value: "testRenderenv1",
					},
					{
						Name:  "LOGFILE",
						Value: "/tmp/234/printtime.log",
					},
				},
				Mounts: []K8sMount{
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
		Dependencies: []Dependency{
			{
				AppName: "group-nginx",
			},
			{
				AppName: "group-ubuntu",
			},
		},
	}
)

func TestCreateAppAndWait(t *testing.T) {
	InitSomeThing()

	appToCreate := appToCreate2

	outAppInfo, err := CreateAppAndWait(appToCreate)
	if err != nil {
		t.Errorf("create application [%s], error: [%s]", appToCreate.Name, err.Error())
	} else {
		t.Logf("The created application is [%s]", JsonString(outAppInfo))
	}

}

func TestCreateAppsWait(t *testing.T) {
	InitSomeThing()

	appsToCreate := []K8sApp{appToCreate1, appToCreate2, appToCreate3}

	outAppsInfo, err := CreateAppsWait(appsToCreate)
	if err != nil {
		t.Errorf("create applications, error: [%s]", err.Error())
	} else {
		t.Logf("The created applications are [%s]", JsonString(outAppsInfo))
	}

}

func TestDeleteBatchApps(t *testing.T) {
	InitSomeThing()

	appsNamesToDelete := []string{
		"test-app-20",
		"test-app-29",
		"test-app-10",
		"test-app-5",
		"test-app-18",
	}

	errs := DeleteBatchApps(appsNamesToDelete)
	if errs != nil {
		t.Errorf("Delete applications, error: [%s]", HandleErrSlice(errs).Error())
	} else {
		t.Logf("Delete applications successfully")
	}

}
