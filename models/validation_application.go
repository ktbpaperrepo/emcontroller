package models

import (
	"fmt"
)

func ValidateK8sApp(app K8sApp) error {
	var maxPriority int = 10
	var minPriority int = 0
	if app.Priority > maxPriority || app.Priority < minPriority {
		return fmt.Errorf("app.Priority should be in [%d, %d], but the input app.Priority is %d", minPriority, maxPriority, app.Priority)
	}
	return nil
}
