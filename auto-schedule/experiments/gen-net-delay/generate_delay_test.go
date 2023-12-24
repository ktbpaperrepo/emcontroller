package gen_net_delay

import (
	"testing"

	"emcontroller/models"
)

// go test /mnt/c/mine/code/gocode/src/emcontroller/auto-schedule/experiments/gen-net-delay/ -v -count=1 -run TestGenCloudsDelay
func TestGenCloudsDelay(t *testing.T) {
	models.InitClouds()

	delays := map[string]string{
		"NOKIA2":  "9ms",
		"NOKIA3":  "15ms",
		"NOKIA4":  "22ms",
		"NOKIA5":  "28ms",
		"NOKIA6":  "36ms",
		"NOKIA8":  "43ms",
		"NOKIA10": "50ms",
	}

	if err := GenCloudsDelay(delays); err != nil {
		t.Errorf("GenCloudsDelay error: %s", err.Error())
	} else {
		t.Log("GenCloudsDelay finished")
	}
}

// go test /mnt/c/mine/code/gocode/src/emcontroller/auto-schedule/experiments/gen-net-delay/ -v -count=1 -run TestClearAllDelay
func TestClearAllDelay(t *testing.T) {
	models.InitClouds()

	if err := ClearAllDelay(); err != nil {
		t.Errorf("ClearAllDelay error: %s", err.Error())
	} else {
		t.Log("ClearAllDelay finished")
	}
}

// go test /mnt/c/mine/code/gocode/src/emcontroller/auto-schedule/experiments/gen-net-delay/ -v -count=1 -run TestInnerDelayOneCloud
func TestInnerDelayOneCloud(t *testing.T) {
	models.InitClouds()

	if err := delayOneCloud("NOKIA4", "100ms"); err != nil {
		t.Errorf("delayOneCloud error: %s", err.Error())
	} else {
		t.Log("delayOneCloud finished")
	}
}

// go test /mnt/c/mine/code/gocode/src/emcontroller/auto-schedule/experiments/gen-net-delay/ -v -count=1 -run TestInnerClearDelayOneCloud
func TestInnerClearDelayOneCloud(t *testing.T) {
	models.InitClouds()

	if err := clearDelayOneCloud("NOKIA4"); err != nil {
		t.Errorf("clearDelayOneCloud error: %s", err.Error())
	} else {
		t.Log("clearDelayOneCloud finished")
	}
}

// go test /mnt/c/mine/code/gocode/src/emcontroller/auto-schedule/experiments/gen-net-delay/ -v -count=1 -run TestInnerSetDelay
func TestInnerSetDelay(t *testing.T) {
	if err := setDelay("root", "xxxxxxxxxx", "192.168.100.60", "130ms"); err != nil {
		t.Errorf("setDelay error: %s", err.Error())
	} else {
		t.Log("setDelay finished")
	}
}

// go test /mnt/c/mine/code/gocode/src/emcontroller/auto-schedule/experiments/gen-net-delay/ -v -count=1 -run TestInnerClearDelay
func TestInnerClearDelay(t *testing.T) {
	if err := clearDelay("root", "xxxxxxxxxxxx", "192.168.100.60"); err != nil {
		t.Errorf("clearDelay error: %s", err.Error())
	} else {
		t.Log("clearDelay finished")
	}
}
