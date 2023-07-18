package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"

	"emcontroller/models"
)

// For migration, we can put the Json of this structure into the Annotation with key variable "AutoScheduleInfoAnno" of the Kubernetes Deployment, so that multi-cloud manager can get the information needed for auto-scheduling from the Annotation with key variable "AutoScheduleInfoAnno".
type Application struct {
	Name         string              `json:"name"`
	Priority     int                 `json:"priority"`
	Resources    AppResources        `json:"resources"`    // The resources information of this application
	Dependencies []models.Dependency `json:"dependencies"` // The information of all applications that this application depends on.
}

func AppMapCopy(src map[string]Application) map[string]Application {
	var dst map[string]Application = make(map[string]Application)
	for name, app := range src {
		dst[name] = app
	}
	return dst
}

func GenerateApplications(inputApps []models.K8sApp) (map[string]Application, error) {
	var outApps map[string]Application = make(map[string]Application)

	for _, inApp := range inputApps {

		// traverse containers to calculate the resources requested by this applications
		var resources AppResources
		for _, container := range inApp.Containers {
			floatCpu, err := strconv.ParseFloat(container.Resources.Requests.CPU, 64)
			if err != nil {
				outErr := fmt.Errorf("Container [%s] container.Resources.Requests.CPU [%s] parse to float64, Error: [%w]", container.Name, container.Resources.Requests.CPU, err)
				beego.Error(outErr)
				return nil, outErr
			}
			floatRamMi, err := strconv.ParseFloat(strings.TrimSuffix(container.Resources.Requests.Memory, MemUnitSuffix), 64)
			if err != nil {
				outErr := fmt.Errorf("Container [%s] container.Resources.Requests.Memory [%s] parse to float64, Error: [%w]", container.Name, container.Resources.Requests.Memory, err)
				beego.Error(outErr)
				return nil, outErr
			}
			floatStorGi, err := strconv.ParseFloat(strings.TrimSuffix(container.Resources.Requests.Storage, StorageUnitSuffix), 64)
			if err != nil {
				outErr := fmt.Errorf("Container [%s] container.Resources.Requests.Storage [%s] parse to float64, Error: [%w]", container.Name, container.Resources.Requests.Storage, err)
				beego.Error(outErr)
				return nil, outErr
			}

			resources.CpuCore += floatCpu
			resources.Memory += floatRamMi
			resources.Storage += floatStorGi
		}

		// put the needed information in the output structure
		var thisOutApp Application
		thisOutApp.Name = inApp.Name
		thisOutApp.Priority = inApp.Priority
		thisOutApp.Resources = resources
		thisOutApp.Dependencies = inApp.Dependencies
		outApps[thisOutApp.Name] = thisOutApp
	}

	return outApps, nil
}
