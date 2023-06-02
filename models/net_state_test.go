package models

import (
	"encoding/json"
	"testing"
)

func TestGetNetState(t *testing.T) {
	InitSomeThing()
	netState, err := GetNetState()
	if err != nil {
		t.Errorf("GetNetState() error: %s", err.Error())
	} else {
		t.Logf("netstat: %v", netState)
		netStateJson, err := json.Marshal(netState)
		if err != nil {
			t.Errorf("json.Marshal(netState) error: %s", err.Error())
		}
		t.Logf("netstat json: %s", netStateJson)
	}
}
