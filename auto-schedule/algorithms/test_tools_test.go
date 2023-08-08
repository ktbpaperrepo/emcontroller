package algorithms

import (
	"fmt"
	"sort"
	"testing"

	"github.com/KeepTheBeats/routing-algorithms/random"

	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

type depGenHelper struct {
	appName  string
	priority int
	oriIdx   int // original index
}

type depGenHelperSlice []depGenHelper

func (dghs depGenHelperSlice) Len() int {
	return len(dghs)
}

func (dghs depGenHelperSlice) Swap(i, j int) {
	dghs[i], dghs[j] = dghs[j], dghs[i]
}

func (dghs depGenHelperSlice) Less(i, j int) bool {
	return dghs[i].priority < dghs[j].priority
}

func newDepGenHelper(app models.K8sApp, originalIndex int) depGenHelper {
	return depGenHelper{
		appName:  app.Name,
		priority: app.Priority,
		oriIdx:   originalIndex,
	}
}

// the variables of applications.
type appVars struct {
	image    string
	commands []string
	ports    []models.PortInfo
}

func MakeAppsForTest(namePrefix string, count int, possibleVars []appVars) []models.K8sApp {
	outApps := make([]models.K8sApp, count)

	for i := 0; i < len(outApps); i++ {
		randomVar := possibleVars[random.RandomInt(0, len(possibleVars)-1)]

		outApps[i].Name = fmt.Sprintf("%s-%d", namePrefix, i)
		outApps[i].AutoScheduled = true
		outApps[i].Replicas = 1
		outApps[i].Priority = random.RandomInt(asmodel.MinPriority, asmodel.MaxPriority)
		//outApps[i].HostNetwork = random.RandomInt(0, 1) == 0
		outApps[i].HostNetwork = false
		outApps[i].Containers = []models.K8sContainer{
			models.K8sContainer{
				Name:     "container",
				Image:    randomVar.image,
				Commands: randomVar.commands,
				Ports:    randomVar.ports,
				WorkDir:  "",
				Resources: models.K8sResReq{
					Limits: models.K8sResList{
						CPU:    fmt.Sprintf("%.1f", random.NormalRandomBM(1.0, 32.0, 6.0, 6.0)),
						Memory: fmt.Sprintf("%.0fMi", random.NormalRandomBM(200.0, 16384.0, 2048.0, 8192.0)),
					},
				},
			},
		}

		resStorage := random.RandomInt(0, 200)
		if resStorage > 0 {
			outApps[i].Containers[0].Resources.Limits.Storage = fmt.Sprintf("%dGi", resStorage)
		}
		outApps[i].Containers[0].Resources.Requests = outApps[i].Containers[0].Resources.Limits

	}

	// generate dependencies
	var depHelpers []depGenHelper = make([]depGenHelper, count)
	for i := 0; i < len(outApps); i++ {
		depHelpers[i] = newDepGenHelper(outApps[i], i)
	}
	sort.Sort(depGenHelperSlice(depHelpers)) // sort apps, and an app can only depend on those after it.

	for i := 0; i < len(depHelpers); i++ {
		for j := i + 1; j < len(depHelpers); j++ {
			if random.RandomInt(0, 3) == 0 {
				outApps[depHelpers[i].oriIdx].Dependencies = append(outApps[depHelpers[i].oriIdx].Dependencies, models.Dependency{AppName: depHelpers[j].appName})
			}
		}
	}

	fmt.Println("Generated Applications Json:")
	fmt.Println(models.JsonString(outApps))

	return outApps
}

func TestMakeAppsForTest(t *testing.T) {
	var namePrefix string = "test-app"
	var count int = 40

	var possibleVars []appVars = []appVars{
		appVars{
			image: "172.27.15.31:5000/nginx:1.17.1",
			ports: []models.PortInfo{
				models.PortInfo{
					ContainerPort: 80,
					Name:          "tcp",
					Protocol:      "tcp",
					ServicePort:   "100",
				},
			},
		},
		appVars{
			image: "172.27.15.31:5000/ubuntu:latest",
			commands: []string{
				"bash",
				"-c",
				"while true;do sleep 10;done",
			},
		},
	}

	_ = MakeAppsForTest(namePrefix, count, possibleVars)
}
