package models

import (
	"testing"
)

func TestDeleteDb(t *testing.T) {
	if err := DeleteDb(NetPerfDbName); err != nil {
		t.Errorf("Delete database [%s], error %s.", NetPerfDbName, err.Error())
	} else {
		t.Logf("Success")
	}
}

func TestCreateDb(t *testing.T) {
	if err := CreateDb(NetPerfDbName); err != nil {
		t.Errorf("Create database [%s], error %s.", NetPerfDbName, err.Error())
	} else {
		t.Logf("Success")
	}
}

func TestUseDbShowCurUsedDb(t *testing.T) {
	db, err := NewMySqlCli()
	if err != nil {
		t.Errorf("NewMySqlCli(), error [%s].", err.Error())
	}
	defer db.Close()

	curDb, err := ShowCurUsedDb(db)
	if err != nil {
		t.Errorf("ShowCurUsedDb(db), error [%s].", err.Error())
	} else {
		t.Logf("curDb is: %s", curDb)
	}

	if err := UseDb(db, NetPerfDbName); err != nil {
		t.Errorf("Use database [%s], error %s.", NetPerfDbName, err.Error())
	} else {
		t.Logf("UseDb successful")
	}

	curDb, err = ShowCurUsedDb(db)
	if err != nil {
		t.Errorf("ShowCurUsedDb(db), error [%s].", err.Error())
	} else {
		t.Logf("curDb is: %s", curDb)
	}
}
