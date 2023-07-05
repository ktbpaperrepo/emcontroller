package algorithms

import (
	"fmt"

	"github.com/astaxie/beego"

	asmodel "emcontroller/auto-schedule/model"
)

// Multi-Cloud Service Scheduling Genetic Algorithm (MCSSGA)
type Mcssga struct {
}

func (m *Mcssga) Schedule(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application) (asmodel.Solution, error) {
	beego.Info("Clouds:")
	for _, cloud := range clouds {
		beego.Info(fmt.Sprintf("%+v\n", cloud))
	}
	beego.Info("Applications:")
	for _, app := range apps {
		beego.Info(fmt.Sprintf("%+v\n", app))
	}
	return nil, nil
}
