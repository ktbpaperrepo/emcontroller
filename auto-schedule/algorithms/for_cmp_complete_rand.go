package algorithms

import (
	"fmt"

	"github.com/KeepTheBeats/routing-algorithms/random"
	"github.com/astaxie/beego"

	asmodel "emcontroller/auto-schedule/model"
)

/**
An algorithm for comparison, with name "completely random algorithm".
In the experiment, we will compare MCSSGA with this algorithm.
*/

// completely random algorithm
type CompRand struct {
}

func NewCompRand() *CompRand {
	return &CompRand{}
}

func (m *CompRand) Schedule(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) (asmodel.Solution, error) {
	beego.Info("Using scheduling algorithm:", CompRandName)

	// copy clouds, avoiding changing the original ones.
	cloudsCopy := asmodel.CloudMapCopy(clouds)

	var solution asmodel.Solution = asmodel.GenEmptySoln()
	for appName, _ := range apps {
		var thisAppSoln asmodel.SingleAppSolution

		thisAppSoln.Accepted = random.RandomInt(0, 1) == 0 // randomly set accepted

		if thisAppSoln.Accepted { // randomly set cloud
			pickedCloudName, _ := randomCloudMapPick(cloudsCopy)
			thisAppSoln.TargetCloudName = pickedCloudName
		}

		solution.AppsSolution[appName] = thisAppSoln
	}

	refinedSoln, acceptable := CmpRefineSoln(clouds, apps, appsOrder, solution)
	if !acceptable {
		return asmodel.Solution{}, fmt.Errorf("This time, \"completely random algorithm\" get an unusable solution.")
	}

	return refinedSoln, nil
}
