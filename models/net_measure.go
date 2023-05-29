package models

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// The names of McmNetTestVMs are "CloudName+suffix"
	mcmNetTestServerSuffix string = "-mcm-net-test-server"
	mcmNetTestClientSuffix string = "-mcm-net-test-client"

	mcmNetTestCpu  int = 2    // number of logical CPU cores
	mcmNetTestMem  int = 2048 // memory size unit: MB
	mcmNetTestDisk int = 20   // storage size unit: GB (I tried 10 GB, but it is too small.)

	taintValue   string                   = "net-test"
	taintEffect  apiv1.TaintEffect        = apiv1.TaintEffectNoSchedule
	tolerationOp apiv1.TolerationOperator = apiv1.TolerationOpEqual
)

var NetTestTaint *apiv1.Taint = &apiv1.Taint{
	Key:    McmKey,
	Effect: taintEffect,
	Value:  taintValue,
}

// The function to measure network performance between every two clouds
// This function should be executed every time period
func MeasNetPerf() {
	beego.Info("Start to ensure the network test preconditions for each cloud.")

	// Ensure the network test preconditions for each cloud in parallel
	var wg sync.WaitGroup
	var errsMu sync.Mutex // the slice in golang is not safe for concurrent read/write
	var errs []error
	for name, cloud := range Clouds {
		beego.Info(fmt.Sprintf("Ensure the network test preconditions for cloud %s", name))
		wg.Add(1)
		go func(c Iaas) {
			defer wg.Done()
			if err := ensureTestPreC(c); err != nil {
				outErr := fmt.Errorf("Cannot ensure the network test preconditions for cloud [%s], error: [%w]", c.ShowName(), err)
				beego.Error(outErr)
				errsMu.Lock()
				errs = append(errs, outErr)
				errsMu.Unlock()
				return
			}
		}(cloud)
	}
	wg.Wait()

	if len(errs) != 0 {
		sumErr := HandleErrSlice(errs)
		outErr := fmt.Errorf("Cannot ensure the network test preconditions, Error: %w", sumErr)
		beego.Error(outErr)
		return
	}

	beego.Info("Finish ensuring the network test preconditions for each cloud.")

	// TODO
	// Execute network test between every two clouds
	beego.Info("Start to measure network performance between every two clouds.")

	beego.Info("Finish measuring network performance between every two clouds.")
}

// Ensure the preconditions for network test, including VMs, K8s nodes, and K8s taints.
func ensureTestPreC(cloud Iaas) error {
	// 1. network test VMs exist
	netTestVmMap, err := ensureVms(cloud)
	if err != nil {
		outErr := fmt.Errorf("Cannot ensure that network test VMs exist, Error: %w", err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("On cloud [%s], network test VMs exist.", cloud.ShowName()))

	// 2. network test VMs are managed by Kubernetes
	err = ensureK8sMg(netTestVmMap)
	if err != nil {
		outErr := fmt.Errorf("Cannot ensure that network test VMs are managed by Kubernetes, Error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	// 3. network test VMs have the NetTestTaint in Kubernetes
	err = ensureTaints(netTestVmMap)
	if err != nil {
		outErr := fmt.Errorf("Cannot ensure that VMs have the NetTestTaint in Kubernetes, Error: %w", err)
		beego.Error(outErr)
		return outErr
	}

	beego.Info(fmt.Sprintf("On cloud [%s], The precondition is OK for network test.", cloud.ShowName()))
	return nil
}

// Ensure that the VMs for network test exists
func ensureVms(cloud Iaas) (map[string]*IaasVm, error) {
	var netTestVmNames []string = getNetTestVMNames(cloud)
	vms, err := cloud.ListAllVMs()
	if err != nil {
		outErr := fmt.Errorf("List vms in cloud [%s] type [%s], error %s.", cloud.ShowName(), cloud.ShowType(), err)
		beego.Error(outErr)
		return nil, outErr
	}

	// We should use the information of the network performance test VMs in the following steps, so we use a map to return the VMs information
	var netTestVmMap map[string]*IaasVm = make(map[string]*IaasVm)

	// Here we cannot check each VM in parallel, because if we create more than one VM in proxmox, there will be problems.
	// It seems that proxmox does not support creating VMs in parallel.

	// Check each VM in serial
	for _, netTestVmName := range netTestVmNames {
		if vm, found := FindVm(netTestVmName, vms); found { // If a test vm exists, we return it.
			netTestVmMap[netTestVmName] = vm
		} else { // If a test vm does not exist, we create and return a new one.
			newVm, err := cloud.CreateVM(netTestVmName, mcmNetTestCpu, mcmNetTestMem, mcmNetTestDisk)
			if err != nil {
				outErr := fmt.Errorf("Create vm [%s] in cloud [%s] type [%s], error %w.", netTestVmName, cloud.ShowName(), cloud.ShowType(), err)
				beego.Error(outErr)
				return nil, outErr
			}
			netTestVmMap[netTestVmName] = newVm
		}
	}

	//// Check each VM in parallel
	//var wg sync.WaitGroup
	//var errsMu sync.Mutex // the slice in golang is not safe for concurrent read/write
	//var errs []error
	//for _, netTestVmName := range netTestVmNames {
	//	wg.Add(1)
	//	go func(name string) {
	//		defer wg.Done()
	//		if !FindVm(name, vms) {
	//			if _, err := cloud.CreateVM(name, mcmNetTestCpu, mcmNetTestMem, mcmNetTestDisk); err != nil {
	//				outErr := fmt.Errorf("Create vm [%s] in cloud [%s] type [%s], error %w.", name, cloud.ShowName(), cloud.ShowType(), err)
	//				beego.Error(outErr)
	//				errsMu.Lock()
	//				errs = append(errs, outErr)
	//				errsMu.Unlock()
	//			}
	//		}
	//	}(netTestVmName)
	//}
	//wg.Wait()
	//
	//if len(errs) != 0 {
	//	sumErr := HandleErrSlice(errs)
	//	outErr := fmt.Errorf("Ensure network test VMs in cloud [%s], Error: %w", cloud.ShowName(), sumErr)
	//	beego.Error(outErr)
	//	return outErr
	//}

	return netTestVmMap, nil
}

func getNetTestVMNames(cloud Iaas) []string {
	return []string{
		getNetTestServerName(cloud),
		getNetTestClientName(cloud),
	}
}

func getNetTestServerName(cloud Iaas) string {
	return strings.ToLower(cloud.ShowName()) + mcmNetTestServerSuffix
}

func getNetTestClientName(cloud Iaas) string {
	return strings.ToLower(cloud.ShowName()) + mcmNetTestClientSuffix
}

// Ensure that the VMs for network test are managed by Kubernetes
func ensureK8sMg(vmMap map[string]*IaasVm) error {
	for vmName, vm := range vmMap {
		_, err := GetNode(vmName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				beego.Info(fmt.Sprintf("Node [%s] is not managed by Kubernetes. Then we add it.", vmName))

				joinCmd, err := GetJoinCmd()
				if err != nil {
					outErr := fmt.Errorf("Cannot add node [%s] to Kubernetes cluster, GetJoinCmd error: %w", vmName, err)
					beego.Error(outErr)
					return outErr
				}
				if err := AddNode(*vm, joinCmd); err != nil {
					outErr := fmt.Errorf("Cannot add node [%s] to Kubernetes cluster, error: %w.", vmName, err)
					beego.Error(outErr)
					return outErr
				}

				beego.Info(fmt.Sprintf("Node [%s] is already added to Kubernetes cluster.", vmName))
			} else {
				outErr := fmt.Errorf("Get node [%s] in Kubernetes, error %w.", vmName, err)
				beego.Error(outErr)
				return outErr
			}
		} else {
			beego.Info(fmt.Sprintf("Node [%s] is already managed by Kubernetes.", vmName))
		}
	}
	return nil
}

// Ensure that the nodes of VMs in Kubernetes have the network test taints.
func ensureTaints(vmMap map[string]*IaasVm) error {
	for nodeName, _ := range vmMap {
		if err := TaintNode(nodeName, NetTestTaint); err != nil {
			outErr := fmt.Errorf("ensureTaints, node [%s], error %w.", nodeName, err)
			beego.Error(outErr)
			return outErr
		}
	}

	return nil
}

// Create the database, tables, and columns in MySQL.
func InitNetPerfDB() error {
	beego.Info("Create the database, tables, and columns in MySQL.")

	beego.Info(fmt.Sprintf("Check existing database"))

	dbs, err := ListDbs()
	if err != nil {
		outErr := fmt.Errorf("Check existing database, error %w.", err)
		beego.Error(outErr)
		return outErr
	}

	beego.Info(fmt.Sprintf("Existing databases: %v", dbs))

	var oldDbExist bool
	for _, db := range dbs {
		if db == NetPerfDbName {
			oldDbExist = true
			break
		}
	}

	if oldDbExist {
		beego.Info(fmt.Sprintf("Database [%s] exists, so we delete it.", NetPerfDbName))
		if err := DeleteDb(NetPerfDbName); err != nil {
			outErr := fmt.Errorf("Delete database [%s], error %w.", NetPerfDbName, err)
			beego.Error(outErr)
			return outErr
		}
	} else {
		beego.Info(fmt.Sprintf("Database [%s] does not exist, so we do not need to delete it.", NetPerfDbName))
	}

	beego.Info(fmt.Sprintf("Create new database: %s", NetPerfDbName))
	if err := CreateDb(NetPerfDbName); err != nil {
		outErr := fmt.Errorf("Create database [%s], error %w.", NetPerfDbName, err)
		beego.Error(outErr)
		return outErr
	}

	// In MySQL, we can use the command "use" to select a database.
	// We can also specify the database name in every of our commands.
	// For the purpose of learning, I implement both of the 2 ways, either of which is OK to use.

	//if err := initDbTables(); err != nil {
	if err := initDbTablesWithUse(); err != nil {
		outErr := fmt.Errorf("Create and initialize tables for network performance in database [%s], error %w.", NetPerfDbName, err)
		beego.Error(outErr)
		return outErr
	}

	return nil
}

func initDbTables() error {
	beego.Info(fmt.Sprintf("Create and initialize tables for network performance in database [%s].", NetPerfDbName))
	db, err := NewMySqlCli()
	if err != nil {
		outErr := fmt.Errorf("Create MySQL client, error [%w].", err)
		beego.Error(outErr)
		return outErr
	}
	defer db.Close()

	for cloudName, _ := range Clouds {
		tableName := cloudName
		query := fmt.Sprintf("create table %s.%s(%s varchar(768) not null,%s double not null, primary key(%s))", NetPerfDbName, tableName, DbFieldCloudName, DbFieldRtt, DbFieldCloudName)
		result, err := db.Query(query)
		if err != nil {
			outErr := fmt.Errorf("Query [%s], error [%w].", query, err)
			beego.Error(outErr)
			return outErr
		}
		result.Close()
		beego.Info(fmt.Sprintf("Query [%s] successfully.", query))

		for targetCloudName, _ := range Clouds {
			query := fmt.Sprintf("insert into %s.%s (%s, %s) values (?, ?)", NetPerfDbName, tableName, DbFieldCloudName, DbFieldRtt)
			result, err := db.Query(query, targetCloudName, math.MaxFloat64)
			if err != nil {
				outErr := fmt.Errorf("Query [%s], args: [%s, %g], error [%w].", query, targetCloudName, math.MaxFloat64, err)
				beego.Error(outErr)
				return outErr
			}
			result.Close()
			beego.Info(fmt.Sprintf("Query [%s], args: [%s, %g] successfully.", query, targetCloudName, math.MaxFloat64))
		}
	}
	beego.Info(fmt.Sprintf("Successful! Create and initialize tables for network performance in database [%s].", NetPerfDbName))

	return nil
}

func initDbTablesWithUse() error {
	beego.Info(fmt.Sprintf("Create and initialize tables for network performance in database [%s].", NetPerfDbName))
	db, err := NewMySqlCli()
	if err != nil {
		outErr := fmt.Errorf("Create MySQL client, error [%w].", err)
		beego.Error(outErr)
		return outErr
	}
	defer db.Close()

	if err := UseDb(db, NetPerfDbName); err != nil {
		outErr := fmt.Errorf("UseDb [%s], error %w.", NetPerfDbName, err)
		beego.Error(outErr)
		return outErr
	}

	curDb, err := ShowCurUsedDb(db)
	if err != nil {
		outErr := fmt.Errorf("ShowCurUsedDb(db), error %w.", err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("Currently database [%s] is used.", curDb))

	for cloudName, _ := range Clouds {
		tableName := cloudName
		query := fmt.Sprintf("create table %s(%s varchar(768) not null,%s double not null, primary key(%s))", tableName, DbFieldCloudName, DbFieldRtt, DbFieldCloudName)
		result, err := db.Query(query)
		if err != nil {
			outErr := fmt.Errorf("Query [%s], error [%w].", query, err)
			beego.Error(outErr)
			return outErr
		}
		result.Close()
		beego.Info(fmt.Sprintf("Query [%s] successfully.", query))

		for targetCloudName, _ := range Clouds {
			query := fmt.Sprintf("insert into %s (%s, %s) values (?, ?)", tableName, DbFieldCloudName, DbFieldRtt)
			result, err := db.Query(query, targetCloudName, math.MaxFloat64)
			if err != nil {
				outErr := fmt.Errorf("Query [%s], args: [%s, %g], error [%w].", query, targetCloudName, math.MaxFloat64, err)
				beego.Error(outErr)
				return outErr
			}
			result.Close()
			beego.Info(fmt.Sprintf("Query [%s], args: [%s, %g] successfully.", query, targetCloudName, math.MaxFloat64))
		}
	}
	beego.Info(fmt.Sprintf("Successful! Create and initialize tables for network performance in database [%s].", NetPerfDbName))

	return nil
}
