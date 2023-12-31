package algorithms

import (
	"github.com/KeepTheBeats/routing-algorithms/random"

	asmodel "emcontroller/auto-schedule/model"
)

// Generate a solution randomly, doing the best to accept more applications.
func RandomAcceptMostSolution(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) asmodel.Solution {

	// initialize an all-reject solution with all applications rejected.
	var solution asmodel.Solution = asmodel.GenEmptySoln()
	for _, app := range apps {
		solution.AppsSolution[app.Name] = asmodel.SasCopy(asmodel.RejSoln)
	}

	// Then, we try to accept applications based on the all-reject solution.

	// avoiding changing the original application map
	untriedApps := asmodel.AppMapCopy(apps)
	// traverse apps in random order
	for len(untriedApps) > 0 {
		// randomly choose an application, trying to deploy it to a cloud.
		pickedAppName, _ := randomAppMapPick(untriedApps)

		// avoiding changing the original cloud map
		untriedClouds := asmodel.CloudMapCopy(clouds)
		// traverse clouds in random order
		for len(untriedClouds) > 0 {
			// randomly choose a cloud, trying to deploy the application to it.
			pickedCloudName, _ := randomCloudMapPick(untriedClouds)
			solution.AppsSolution[pickedAppName] = asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: pickedCloudName,
			}

			// If the randomly chosen cloud and the app constitute an acceptable solution,
			// we map them in the solution, and "break" to look for a cloud for another application.
			refinedSoln, acceptable := RefineSoln(clouds, apps, appsOrder, solution)
			if acceptable {
				// if this solution passes the 3 checks (VM, CPU, acceptable), we set it back to the original solution.
				solution = asmodel.SolutionCopy(refinedSoln)
				break
			}

			// Otherwise,
			// we restore this application as rejected, remove this cloud from the untriedClouds, and try another cloud in the next loop.
			solution.AppsSolution[pickedAppName] = asmodel.SasCopy(asmodel.RejSoln)
			delete(untriedClouds, pickedCloudName)

		}

		// When the code reaches here, we have finished the solution searching for this application, so we remove it from untriedApps and start to look for a solution for another application in the next loop.
		delete(untriedApps, pickedAppName)
	}

	return solution
}

// randomly pick an item from a cloud map
func randomCloudMapPick(m map[string]asmodel.Cloud) (string, asmodel.Cloud) {
	k := random.RandomInt(0, len(m)-1)
	for name, cloud := range m {
		if k == 0 {
			return name, cloud
		}
		k--
	}
	panic("Unexpected condition.")
}

// randomly pick an item from an application map
func randomAppMapPick(m map[string]asmodel.Application) (string, asmodel.Application) {
	k := random.RandomInt(0, len(m)-1)
	for name, app := range m {
		if k == 0 {
			return name, app
		}
		k--
	}
	panic("Unexpected condition.")
}
