package models

const (
	RunningStatus   = "Stable Running"
	NotStableStatus = "Not Yet Stable"

	// Kubernetes Annotation keys for the auto-schedule functionality
	AutoScheduledAnno    string = "auto-schedule"
	PriorityAnno         string = "priority"
	AutoScheduleInfoAnno string = "auto-schedule/info"
)
