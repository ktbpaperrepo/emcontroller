package algorithms

import (
	"testing"

	"github.com/stretchr/testify/assert"

	asmodel "emcontroller/auto-schedule/model"
)

func TestInnerAppOneCloudIter(t *testing.T) {
	cloud, apps, soln := cloudAppsSolnForTest()

	appsThisCloud := findAppsOneCloud(cloud, apps, soln)
	t.Logf("appsThisCloud:\n%+v\n", appsThisCloud)

	appsOrder := GenerateAppsOrder(apps)
	t.Logf("appsOrder:\n%+v\n", appsOrder)

	var appsThisCloudIter *appOneCloudIter
	var curAppName string

	t.Log()
	t.Log("The following 3 orders of curAppName should be the same.")
	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	for {
		curAppName = appsThisCloudIter.nextAppName()
		t.Logf("curAppName: %s\n", curAppName)
		if len(curAppName) == 0 {
			break
		}
	}

	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	for {
		curAppName = appsThisCloudIter.nextAppName()
		t.Logf("curAppName: %s\n", curAppName)
		if len(curAppName) == 0 {
			break
		}
	}

	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	for {
		curAppName = appsThisCloudIter.nextAppName()
		t.Logf("curAppName: %s\n", curAppName)
		if len(curAppName) == 0 {
			break
		}
	}

	t.Log()
	t.Log("Then, we test copying an iterator.")
	appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()
	t.Logf("curAppName: %s\n", curAppName)
	curAppName = appsThisCloudIter.nextAppName()
	t.Logf("curAppName: %s\n", curAppName)

	t.Log("Now, we copy an iterator")
	iterCopy := appsThisCloudIter.Copy()
	curAppNameCopy := curAppName

	t.Log("Now, we run the original iterator.")
	for {
		curAppName = appsThisCloudIter.nextAppName()
		t.Logf("curAppName: %s\n", curAppName)
		if len(curAppName) == 0 {
			break
		}
	}

	t.Log("Now, we run the copied iterator, and it should work from the point where we copied.")
	for {
		curAppNameCopy = iterCopy.nextAppName()
		t.Logf("curAppNameCopy: %s\n", curAppNameCopy)
		if len(curAppNameCopy) == 0 {
			break
		}
	}

	t.Log()
	t.Log("Test the scenario in which no applications are scheduled to this cloud.")
	solnNon := asmodel.Solution{
		AppsSolution: map[string]asmodel.SingleAppSolution{
			"app1": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud2",
			},
			"app2": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud2",
			},
			"app3": asmodel.SingleAppSolution{
				Accepted: false,
			},
			"app4": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud2",
			},
			"app5": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud2",
			},
			"app6": asmodel.SingleAppSolution{
				Accepted: false,
			},
			"app7": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud3",
			},
			"app8": asmodel.SingleAppSolution{
				Accepted:        true,
				TargetCloudName: "cloud3",
			},
		},
	}
	appsThisCloudNon := findAppsOneCloud(cloud, apps, solnNon)
	appsThisCloudIter = newAppOneCloudIter(appsThisCloudNon, appsOrder)
	curAppName = appsThisCloudIter.nextAppName()
	assert.Equal(t, 0, len(curAppName))

	func() {
		t.Log()
		t.Log("Then, we test copying and max priority.")
		appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
		curAppName = appsThisCloudIter.nextMaxPriAppName()
		t.Logf("curAppName: %s\n", curAppName)

		t.Log("Now, we copy an iterator")
		iterCopy := appsThisCloudIter.Copy()
		curAppNameCopy := curAppName

		t.Log("Now, we run the original iterator.")
		for {
			curAppName = appsThisCloudIter.nextMaxPriAppName()
			t.Logf("curAppName: %s\n", curAppName)
			if len(curAppName) == 0 {
				break
			}
		}

		t.Log("Now, we run the copied iterator, and it should work from the point where we copied.")
		for {
			curAppNameCopy = iterCopy.nextMaxPriAppName()
			t.Logf("curAppNameCopy: %s\n", curAppNameCopy)
			if len(curAppNameCopy) == 0 {
				break
			}
		}
	}()

	func() {
		t.Log()
		t.Log("Then, we test copying and not max priority.")
		appsThisCloudIter = newAppOneCloudIter(appsThisCloud, appsOrder)
		curAppName = appsThisCloudIter.nextNotMaxPriAppName()
		t.Logf("curAppName: %s\n", curAppName)

		t.Log("Now, we copy an iterator")
		iterCopy := appsThisCloudIter.Copy()
		curAppNameCopy := curAppName

		t.Log("Now, we run the original iterator.")
		for {
			curAppName = appsThisCloudIter.nextNotMaxPriAppName()
			t.Logf("curAppName: %s\n", curAppName)
			if len(curAppName) == 0 {
				break
			}
		}

		t.Log("Now, we run the copied iterator, and it should work from the point where we copied.")
		for {
			curAppNameCopy = iterCopy.nextNotMaxPriAppName()
			t.Logf("curAppNameCopy: %s\n", curAppNameCopy)
			if len(curAppNameCopy) == 0 {
				break
			}
		}
	}()
}
