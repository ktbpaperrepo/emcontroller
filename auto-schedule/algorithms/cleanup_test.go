package algorithms

import (
	"testing"

	"emcontroller/models"
)

func TestGcASVms(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	models.InitSomeThing()
	GcASVms()
}
