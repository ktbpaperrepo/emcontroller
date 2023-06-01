package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"sync"
)

type NetworkState struct {
	Rtt float64 `json:"rtt"` // Round-Trip Time, unit millisecond (ms)
}

// Check network state from MySQL and return the result as a matrix.
func GetNetState() (map[string]map[string]NetworkState, error) {
	var allNetSt map[string]map[string]NetworkState = make(map[string]map[string]NetworkState)

	beego.Info("Start to check network state.")

	// Do it in parallel
	var wg sync.WaitGroup
	var allNetStMu sync.Mutex // the map in golang is not safe for concurrent read/write
	var errsMu sync.Mutex     // the slice in golang is not safe for concurrent read/write
	var errs []error
	for name, _ := range Clouds {
		beego.Info(fmt.Sprintf("check network state from cloud %s", name))
		wg.Add(1)
		go func(cloudName string) {
			defer wg.Done()
			thisCloudNetSt, err := GetNetStateOneCloud(cloudName)
			if err != nil {
				outErr := fmt.Errorf("check network state from cloud [%s], error: [%w]", cloudName, err)
				beego.Error(outErr)
				errsMu.Lock()
				errs = append(errs, outErr)
				errsMu.Unlock()
				return
			}
			allNetStMu.Lock()
			allNetSt[cloudName] = thisCloudNetSt
			allNetStMu.Unlock()
		}(name)
	}
	wg.Wait()

	if len(errs) != 0 {
		sumErr := HandleErrSlice(errs)
		outErr := fmt.Errorf("Failed to check network state, Error: %w", sumErr)
		beego.Error(outErr)
		return nil, outErr
	}

	beego.Info("Finish running the network performance test server Deployment on every net test server VM.")
	return allNetSt, nil
}

// Check network state in the MySQL database from the table of one cloud.
func GetNetStateOneCloud(cloudName string) (map[string]NetworkState, error) {
	db, err := NewMySqlCli()
	if err != nil {
		outErr := fmt.Errorf("Create MySQL client, error [%w].", err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer db.Close()

	query := fmt.Sprintf("select * from %s.%s", NetPerfDbName, cloudName)

	result, err := db.Query(query)
	if err != nil {
		outErr := fmt.Errorf("Query [%s], error [%w].", query, err)
		beego.Error(outErr)
		return nil, outErr
	}
	defer result.Close()

	var netStates map[string]NetworkState = make(map[string]NetworkState)
	for result.Next() {
		var targetCloudName string
		var rtt float64
		if err := result.Scan(&targetCloudName, &rtt); err != nil {
			outErr := fmt.Errorf("Query [%s], result.Scan, error [%w].", query, err)
			beego.Error(outErr)
			beego.Error(fmt.Sprintf("Current netStates: %v", netStates))
			return nil, outErr
		}
		netStates[targetCloudName] = NetworkState{Rtt: rtt}
	}

	beego.Info(fmt.Sprintf("Query [%s] successfully.", query))
	return netStates, nil
}
