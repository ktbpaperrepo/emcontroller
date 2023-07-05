package executors

import (
	"fmt"

	"emcontroller/models"
)

func generateAppMap(apps []models.K8sApp) map[string]models.K8sApp {
	var appMap map[string]models.K8sApp = make(map[string]models.K8sApp)

	for _, app := range apps {
		appMap[app.Name] = app
	}

	return appMap
}

// validate the dependencies among the Auto-Schedule applications
func ValidateAutoScheduleDep(apps []models.K8sApp) []error {
	var allErrs []error

	appMap := generateAppMap(apps)

	// The priority of a dependent application should be greater than or equal to that of the one that depends on it.
	// The dependent application should exist in this group of applications.
	for _, app := range appMap {
		for _, dependency := range app.Dependencies {
			if dependentApp, exist := appMap[dependency.AppName]; exist {
				if dependentApp.Priority < app.Priority {
					allErrs = append(allErrs, fmt.Errorf("Auto-schedule application [%s] with priority [%d] depends on application [%s] with priority [%d], but the priority of a dependent application should be greater than or equal to that of the one that depends on it.", app.Name, app.Priority, dependentApp.Name, dependentApp.Priority))
				}
			} else {
				allErrs = append(allErrs, fmt.Errorf("Auto-schedule application [%s] depends on application [%s], but the dependent application does not exist in this group of applications.", app.Name, dependency.AppName))
			}
		}
	}

	if len(allErrs) != 0 {
		return allErrs
	}

	// There should not be "circular dependencies" in these applications.
	allErrs = append(allErrs, validateCircularDep(appMap)...)

	return allErrs
}

// Check the circular dependencies among these applications.
func validateCircularDep(appMap map[string]models.K8sApp) []error {
	var allErrs []error

	if topoOrder, hasCycle := TopoSort(appMap); hasCycle {
		allErrs = append(allErrs, fmt.Errorf("Auto-schedule applications have circular dependencies, the cycles are not in the following applications %+v.", topoOrder))
	}

	return allErrs
}
