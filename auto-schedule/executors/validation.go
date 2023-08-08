package executors

import (
	"fmt"
	"regexp"
	"strings"

	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

func ValidateAutoScheduleApps(apps []models.K8sApp) []error {
	var allErrs []error

	// validate every single applications
	for _, app := range apps {
		allErrs = append(allErrs, ValidateAutoScheduleApp(app)...)
	}

	// validate the dependencies among these applications
	allErrs = append(allErrs, ValidateAutoScheduleDep(apps)...)

	return allErrs
}

// If an application needs to be auto-scheduled, it should meet some requirements.
func ValidateAutoScheduleApp(app models.K8sApp) []error {
	var allErrs []error

	if len(app.Name) <= 0 {
		allErrs = append(allErrs, fmt.Errorf("Auto-schedule application [%s] does not have a name. len(app.Name) is [%d].", app.Name, len(app.Name)))
	}

	if app.Priority > asmodel.MaxPriority || app.Priority < asmodel.MinPriority {
		allErrs = append(allErrs, fmt.Errorf("Auto-schedule application [%s], app.Priority should be in [%d, %d], but the input app.Priority is %d.", app.Name, asmodel.MinPriority, asmodel.MaxPriority, app.Priority))
	}

	if !app.AutoScheduled {
		allErrs = append(allErrs, fmt.Errorf("Auto-schedule application [%s], AutoScheduled should be [%t], but it is [%t].", app.Name, true, app.AutoScheduled))
	}

	var allowedReplicas int32 = 1
	if app.Replicas != allowedReplicas {
		allErrs = append(allErrs, fmt.Errorf("Auto-schedule application [%s], Replicas should be [%d], but it is [%d].", app.Name, allowedReplicas, app.Replicas))
	}

	if len(app.NodeName) != 0 {
		allErrs = append(allErrs, fmt.Errorf("Auto-schedule application [%s] should not be set NodeName, but it is set as [%s].", app.Name, app.NodeName))
	}

	if len(app.NodeSelector) != 0 {
		allErrs = append(allErrs, fmt.Errorf("Auto-schedule application [%s] should not have NodeSelector, but it has [%s].", app.Name, app.NodeSelector))
	}

	if len(app.Containers) != 1 {
		allErrs = append(allErrs, fmt.Errorf("Auto-schedule application [%s] should only have 1 container, but it has [%d].", app.Name, len(app.Containers)))
	}

	for _, container := range app.Containers {
		allErrs = append(allErrs, validateContainer(container)...)
	}

	return allErrs
}

func validateContainer(container models.K8sContainer) []error {
	var allErrs []error

	allErrs = append(allErrs, validateResources(container.Resources)...)

	return allErrs
}

func validateResources(res models.K8sResReq) []error {
	var allErrs []error

	if !res.Requests.Equal(res.Limits) {
		allErrs = append(allErrs, fmt.Errorf("For auto-schedule application container resources, res.Requests [%+v] should be equal to res.Limits [%+v], but they are not equal.", res.Requests, res.Limits))
	}

	// CPU should not end with m
	if match, err := regexp.MatchString(asmodel.CpuResReg, res.Requests.CPU); err != nil {
		allErrs = append(allErrs, fmt.Errorf("For auto-schedule application container resources, res.Requests.CPU [%s] should match the regular expression [%s], error: [%w].", res.Requests.CPU, asmodel.CpuResReg, err))
	} else if !match {
		allErrs = append(allErrs, fmt.Errorf("For auto-schedule application container resources, res.Requests.CPU [%s] should match the regular expression [%s].", res.Requests.CPU, asmodel.CpuResReg))
	}

	// Memory should have the unit Mi
	if len(res.Requests.Memory) != 0 && !strings.HasSuffix(res.Requests.Memory, asmodel.MemUnitSuffix) {
		allErrs = append(allErrs, fmt.Errorf("For auto-schedule application container resources, res.Requests.Memory [%s] should have the unit suffix [%s].", res.Requests.Memory, asmodel.MemUnitSuffix))
	}

	// Storage should have the unit Gi
	if len(res.Requests.Storage) != 0 && !strings.HasSuffix(res.Requests.Storage, asmodel.StorageUnitSuffix) {
		allErrs = append(allErrs, fmt.Errorf("For auto-schedule application container resources, res.Requests.Storage [%s] should have the unit suffix [%s].", res.Requests.Storage, asmodel.StorageUnitSuffix))
	}

	return allErrs
}
