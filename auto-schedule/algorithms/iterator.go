package algorithms

import asmodel "emcontroller/auto-schedule/model"

// Iterator to iterate the applications in a fixed order
type iterForApps struct {
	appsToIterate map[string]asmodel.Application
	appsOrder     []string // a slice of app names
	nextAppIdx    int
}

func newIterForApps(appsToIterate map[string]asmodel.Application, appsOrder []string) *iterForApps {
	return &iterForApps{
		appsToIterate: appsToIterate,
		appsOrder:     appsOrder,
		nextAppIdx:    0,
	}
}

func (it *iterForApps) nextAppName() string {
	for it.nextAppIdx < len(it.appsOrder) {
		curAppName := it.appsOrder[it.nextAppIdx]
		if _, exist := it.appsToIterate[curAppName]; exist {
			it.nextAppIdx++
			return curAppName
		}
		it.nextAppIdx++
	}
	return ""
}

// get the name of the next application with the priority asmodel.MaxPriority
func (it *iterForApps) nextMaxPriAppName() string {
	for it.nextAppIdx < len(it.appsOrder) {
		curAppName := it.appsOrder[it.nextAppIdx]
		if app, exist := it.appsToIterate[curAppName]; exist && app.Priority == asmodel.MaxPriority {
			it.nextAppIdx++
			return curAppName
		}
		it.nextAppIdx++
	}
	return ""
}

// get the name of the next application with the priority not asmodel.MaxPriority
func (it *iterForApps) nextNotMaxPriAppName() string {
	for it.nextAppIdx < len(it.appsOrder) {
		curAppName := it.appsOrder[it.nextAppIdx]
		if app, exist := it.appsToIterate[curAppName]; exist && app.Priority != asmodel.MaxPriority {
			it.nextAppIdx++
			return curAppName
		}
		it.nextAppIdx++
	}
	return ""
}

func (it *iterForApps) Copy() *iterForApps {
	// copy the map (not completely deep, but enough for our scenarios)
	appsToIterateCopy := asmodel.AppMapCopy(it.appsToIterate)

	// deep copy the slice
	appsOrderCopy := make([]string, len(it.appsOrder))
	copy(appsOrderCopy, it.appsOrder)

	return &iterForApps{
		appsToIterate: appsToIterateCopy,
		appsOrder:     appsOrderCopy,
		nextAppIdx:    it.nextAppIdx,
	}
}
