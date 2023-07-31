/*
The functions group applications using DFS (Depth First Search) according to their dependencies.
*/

package algorithms

import asmodel "emcontroller/auto-schedule/model"

// application with double directions of dependencies recorded
type biDirDepApp struct {
	name      string
	fatherDep map[string]struct{} // the name of its dependent applications
	childDep  map[string]struct{} // the name of the applications depending on it
	visited   bool                // in dfs we need to record whether an item has been visited or not
}

func newBiDirDepApp(name string) biDirDepApp {
	return biDirDepApp{
		name:      name,
		fatherDep: make(map[string]struct{}),
		childDep:  make(map[string]struct{}),
		visited:   false,
	}
}

// generate map[string]biDirDepApp from map[string]asmodel.Application
func genBiDir(apps map[string]asmodel.Application) map[string]biDirDepApp {
	var biApps map[string]biDirDepApp = make(map[string]biDirDepApp)
	// initialize
	for appName, _ := range apps {
		biApps[appName] = newBiDirDepApp(appName)
	}

	// build bidirectional dependency
	for appName, app := range apps {
		for _, dep := range app.Dependencies {
			depAppName := dep.AppName

			// This means the dependent application is not in the input app group, so we do not need to consider it when grouping applications.
			// For example, in our algorithm, if a 10-pri App A depends on a 10-pre App B, but A and B are scheduled to 2 different clouds, the code will go into this if.
			if _, exist := biApps[depAppName]; !exist {
				continue
			}

			// The value of the member variable of a struct in a map is not allowed to be changed, so we use tmp variables (thisBiApp and depBiApp) to change it.

			// App1's father: the applications that App1 depends on
			thisBiApp := biApps[appName]
			thisBiApp.fatherDep[depAppName] = struct{}{}
			biApps[appName] = thisBiApp

			// App1's child: the applications depending on App1
			depBiApp := biApps[depAppName]
			depBiApp.childDep[appName] = struct{}{}
			biApps[depAppName] = depBiApp
		}

	}

	return biApps
}

// group applications according to their dependency using DFS
func groupByDep(apps map[string]asmodel.Application) [][]string {
	var groups [][]string

	biApps := genBiDir(apps)
	for appName, biApp := range biApps {
		if !biApp.visited {
			var group []string
			dfsDep(&biApps, appName, &group)
			groups = append(groups, group)
		}
	}

	return groups
}

// do DFS for dependencies
func dfsDep(biApps *map[string]biDirDepApp, appName string, group *[]string) {
	// mark this application as visited
	biApp := (*biApps)[appName]
	biApp.visited = true
	(*biApps)[appName] = biApp

	// add this application into this group
	*group = append(*group, appName)

	// add all applications which have dependency with this one into the same group.
	for fatherAppName, _ := range (*biApps)[appName].fatherDep {
		if !(*biApps)[fatherAppName].visited {
			dfsDep(biApps, fatherAppName, group)
		}
	}
	for childAppName, _ := range (*biApps)[appName].childDep {
		if !(*biApps)[childAppName].visited {
			dfsDep(biApps, childAppName, group)
		}
	}
}
