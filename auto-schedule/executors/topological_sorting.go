package executors

import "emcontroller/models"

// We should check the circular dependencies among these applications. We use Topological Sorting to check it.
type appForTopoSort struct {
	name         string
	dependencies map[string]struct{}
}

func generateAppMapForTopoSort(appMap map[string]models.K8sApp) map[string]appForTopoSort {
	var appMapCheckCirDep map[string]appForTopoSort = make(map[string]appForTopoSort)
	for name, app := range appMap {

		var thisDeps map[string]struct{} = make(map[string]struct{})
		for _, dependency := range app.Dependencies {
			thisDeps[dependency.AppName] = struct{}{}
		}

		appMapCheckCirDep[name] = appForTopoSort{
			name:         name,
			dependencies: thisDeps,
		}
	}
	return appMapCheckCirDep
}

// Do Topological Sorting for these applications, which can also check the circular dependencies among them.
func TopoSort(appMap map[string]models.K8sApp) ([][]string, bool) {
	appMapCheckCirDep := generateAppMapForTopoSort(appMap)

	var topologicalOrder [][]string
	var hasCycles bool

	for len(appMapCheckCirDep) > 0 {
		// in every loop, we find all applications with no dependencies, which can be put at the head of the rest apps
		var thisGroup []string
		for name, app := range appMapCheckCirDep {
			if len(app.dependencies) == 0 {
				thisGroup = append(thisGroup, name)
			}
		}
		// If the map appMapCheckCirDep still has apps, but no apps in it has 0 dependencies, it means there are cycles in the rest apps.
		if len(thisGroup) == 0 {
			hasCycles = true
			break
		}

		// remove the apps (Node in the Graph) in thisGroup from the map appMapCheckCirDep, and remove the dependencies (Link in the Graph) to them.
		for i := 0; i < len(thisGroup); i++ { // remove the apps
			delete(appMapCheckCirDep, thisGroup[i])
		}
		for i := 0; i < len(thisGroup); i++ { // remove the depencies
			for name, _ := range appMapCheckCirDep {
				delete(appMapCheckCirDep[name].dependencies, thisGroup[i])
			}
		}

		// put this group of apps in the output topological order
		topologicalOrder = append(topologicalOrder, thisGroup)

	}

	return topologicalOrder, hasCycles
}
