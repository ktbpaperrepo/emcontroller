package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

func NewMySqlCli() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", MySqlUser, MySqlPasswd, MySqlIp, MySqlPort))
	if err != nil {
		outErr := fmt.Errorf("sql.Open, error [%w].", err)
		beego.Error(outErr)
		return nil, outErr
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}

func ListDbs() ([]string, error) {
	db, err := NewMySqlCli()
	if err != nil {
		outErr := fmt.Errorf("Create MySQL client, error [%w].", err)
		beego.Error(outErr)
		return []string{}, outErr
	}
	defer db.Close()

	result, err := db.Query("show databases")
	if err != nil {
		outErr := fmt.Errorf("Query \"show databases\", error [%w].", err)
		beego.Error(outErr)
		return []string{}, outErr
	}
	defer result.Close()

	var dbs []string
	for result.Next() {
		var thisDbName string
		if err := result.Scan(&thisDbName); err != nil {
			outErr := fmt.Errorf("Query \"show databases\", result.Scan, error [%w].", err)
			beego.Error(outErr)
			beego.Error(fmt.Sprintf("Current dbs: %v", dbs))
			return []string{}, outErr
		}
		dbs = append(dbs, thisDbName)
	}

	return dbs, nil
}

func DeleteDb(dbName string) error {
	db, err := NewMySqlCli()
	if err != nil {
		outErr := fmt.Errorf("Delete MySQL client, error [%w].", err)
		beego.Error(outErr)
		return outErr
	}
	defer db.Close()

	query := fmt.Sprintf("drop database %s", dbName)

	// I set a timeout for this delete database request,
	// because I find a problem:
	// when the VM of the MySQL server has not space left on the disk, the "delete database request" will be stuck forever without a timeout at the client side.
	ctx, cancel := context.WithTimeout(context.Background(), ReqShortTimeout)
	defer cancel()

	result, err := db.QueryContext(ctx, query)
	if err != nil {
		outErr := fmt.Errorf("Query [%s], error [%w].", query, err)
		beego.Error(outErr)
		return outErr
	}
	defer result.Close()

	beego.Info(fmt.Sprintf("Query [%s] successfully.", query))
	return nil
}

func CreateDb(dbName string) error {
	db, err := NewMySqlCli()
	if err != nil {
		outErr := fmt.Errorf("Create MySQL client, error [%w].", err)
		beego.Error(outErr)
		return outErr
	}
	defer db.Close()

	query := fmt.Sprintf("create database %s", dbName)

	result, err := db.Query(query)
	if err != nil {
		outErr := fmt.Errorf("Query [%s], error [%w].", query, err)
		beego.Error(outErr)
		return outErr
	}
	defer result.Close()

	beego.Info(fmt.Sprintf("Query [%s] successfully.", query))
	return nil
}

func UseDb(db *sql.DB, dbName string) error {
	query := fmt.Sprintf("use %s", dbName)

	result, err := db.Query(query)
	if err != nil {
		outErr := fmt.Errorf("Query [%s], error [%w].", query, err)
		beego.Error(outErr)
		return outErr
	}
	defer result.Close()

	beego.Info(fmt.Sprintf("Query [%s] successfully.", query))
	return nil
}

// show which databse is used currently
func ShowCurUsedDb(db *sql.DB) (string, error) {
	query := "select database() from dual"
	result, err := db.Query(query)
	if err != nil {
		outErr := fmt.Errorf("Query [%s], error [%w].", query, err)
		beego.Error(outErr)
		return "", outErr
	}
	defer result.Close()

	beego.Info(fmt.Sprintf("Query [%s] successfully.", query))

	var curUsedDbs []string
	for result.Next() {
		// the value in a column of the result of this query may be null or string, so we need to use this type and check valid.
		var thisDbName sql.NullString
		if err := result.Scan(&thisDbName); err != nil {
			outErr := fmt.Errorf("Query [%s], result.Scan, error [%w].", query, err)
			beego.Error(outErr)
			beego.Error(fmt.Sprintf("Current curUsedDbs: %v", curUsedDbs))
			return "", outErr
		}
		// the value in a column of the result of this query may be null or string, so we need to use this type and check valid.
		if thisDbName.Valid {
			curUsedDbs = append(curUsedDbs, thisDbName.String)
		}
	}
	beego.Info(fmt.Sprintf("Query [%s], result.Scan finished, curUsedDbs: %v", query, curUsedDbs))

	if len(curUsedDbs) > 1 {
		outErr := fmt.Errorf("len(curUsedDbs) is %d", len(curUsedDbs))
		beego.Error(outErr)
		return "", outErr
	}

	if len(curUsedDbs) == 0 {
		beego.Info("No database is used.")
		return "", nil
	}

	return curUsedDbs[0], nil
}
