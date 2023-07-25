package algorithms

import asmodel "emcontroller/auto-schedule/model"

// Iterator to iterate the applications on one cloud in a fixed order
type appOneCloudIter struct {
	appsThisCloud map[string]asmodel.Application
	appsOrder     []string // a slice of app names
	nextAppIdx    int
}

func newAppOneCloudIter(appsThisCloud map[string]asmodel.Application, appsOrder []string) *appOneCloudIter {
	return &appOneCloudIter{
		appsThisCloud: appsThisCloud,
		appsOrder:     appsOrder,
		nextAppIdx:    0,
	}
}

func (it *appOneCloudIter) nextAppName() string {
	for it.nextAppIdx < len(it.appsOrder) {
		curAppName := it.appsOrder[it.nextAppIdx]
		if _, exist := it.appsThisCloud[curAppName]; exist {
			it.nextAppIdx++
			return curAppName
		}
		it.nextAppIdx++
	}
	return ""
}

// get the name of the next application with the priority asmodel.MaxPriority
func (it *appOneCloudIter) nextMaxPriAppName() string {
	for it.nextAppIdx < len(it.appsOrder) {
		curAppName := it.appsOrder[it.nextAppIdx]
		if app, exist := it.appsThisCloud[curAppName]; exist && app.Priority == asmodel.MaxPriority {
			it.nextAppIdx++
			return curAppName
		}
		it.nextAppIdx++
	}
	return ""
}

// get the name of the next application with the priority not asmodel.MaxPriority
func (it *appOneCloudIter) nextNotMaxPriAppName() string {
	for it.nextAppIdx < len(it.appsOrder) {
		curAppName := it.appsOrder[it.nextAppIdx]
		if app, exist := it.appsThisCloud[curAppName]; exist && app.Priority != asmodel.MaxPriority {
			it.nextAppIdx++
			return curAppName
		}
		it.nextAppIdx++
	}
	return ""
}

func (it *appOneCloudIter) Copy() *appOneCloudIter {
	// copy the map (not completely deep, but enough for our scenarios)
	appsThisCloudCopy := asmodel.AppMapCopy(it.appsThisCloud)

	// deep copy the slice
	appsOrderCopy := make([]string, len(it.appsOrder))
	copy(appsOrderCopy, it.appsOrder)

	return &appOneCloudIter{
		appsThisCloud: appsThisCloudCopy,
		appsOrder:     appsOrderCopy,
		nextAppIdx:    it.nextAppIdx,
	}
}
