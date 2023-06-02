package models

import (
	"testing"
)

func TestInitNetPerfDB(t *testing.T) {
	InitSomeThing()
	err := InitNetPerfDB()
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	} else {
		t.Logf("Success")
	}
}

func TestInnerRunNetTestServer(t *testing.T) {
	InitSomeThing()
	err := runNetTestServer(Clouds["NOKIA8"])
	if err != nil {
		t.Errorf("test error: %s", err.Error())
	} else {
		t.Logf("Success")
	}
}
