package algorithms

import (
	"testing"

	asmodel "emcontroller/auto-schedule/model"
)

func TestInnerRandomCloudMapPick(t *testing.T) {
	t.Logf("\nStart to test the order.\n")
	testOrder := make(map[string]asmodel.Cloud)
	testOrder["cc1"] = asmodel.Cloud{
		Name: "c1",
	}
	testOrder["cc2"] = asmodel.Cloud{
		Name: "c2",
	}
	testOrder["cc3"] = asmodel.Cloud{
		Name: "c3",
	}
	testOrder["cc4"] = asmodel.Cloud{
		Name: "c4",
	}
	testOrder["cc5"] = asmodel.Cloud{
		Name: "c5",
	}
	testOrder["cc6"] = asmodel.Cloud{
		Name: "c6",
	}

	for len(testOrder) > 0 {
		pickedKey, pickedValue := randomCloudMapPick(testOrder)
		t.Log(pickedKey, pickedValue)
		delete(testOrder, pickedKey)
		t.Logf("The rest of the map: %v\n", testOrder)
	}

}

func TestInnerRandomAppMapPick(t *testing.T) {
	t.Logf("\nStart to test the order.\n")
	testOrder := make(map[string]asmodel.Application)
	testOrder["cc1"] = asmodel.Application{
		Name: "c1",
	}
	testOrder["cc2"] = asmodel.Application{
		Name: "c2",
	}
	testOrder["cc3"] = asmodel.Application{
		Name: "c3",
	}
	testOrder["cc4"] = asmodel.Application{
		Name: "c4",
	}
	testOrder["cc5"] = asmodel.Application{
		Name: "c5",
	}
	testOrder["cc6"] = asmodel.Application{
		Name: "c6",
	}

	for len(testOrder) > 0 {
		func() {
			pickedKey, pickedApp := randomAppMapPick(testOrder)
			defer delete(testOrder, pickedKey)

			t.Log(pickedKey, pickedApp)
			t.Logf("The current map: %v\n", testOrder)
		}()
	}

}
