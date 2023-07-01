package algorithms

import asmodel "emcontroller/auto-schedule/model"

// SchedulingAlgorithm is the interface that all algorithms should implement
type SchedulingAlgorithm interface {
	Schedule(clouds []asmodel.Cloud, apps []asmodel.Application) (asmodel.Solution, error)
}
