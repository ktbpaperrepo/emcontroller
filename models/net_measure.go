package models

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NetTestFuncOffMsg string = "Network State Measurement function is turned off."

	DefaultNetTestPeriodSec = 300

	// The names of McmNetTestVMs are "CloudName+suffix"
	mcmNetTestServerSuffix string = "-mcm-net-test-server"
	mcmNetTestClientSuffix string = "-mcm-net-test-client"

	// The names of mcm net test apps are "CloudName+suffix"
	mcmNetTestServerAppSuffix string = "-mcm-nt-app-server"
	mcmNetTestClientAppSuffix string = "-mcm-nt-app-client"

	mcmNetTestCpu  int = 2    // number of logical CPU cores
	mcmNetTestMem  int = 2048 // memory size unit: MB
	mcmNetTestDisk int = 20   // storage size unit: GB (I tried 10 GB, but it is too small.)

	taintValue   string                   = "net-test"
	taintEffect  apiv1.TaintEffect        = apiv1.TaintEffectNoSchedule
	tolerationOp apiv1.TolerationOperator = apiv1.TolerationOpEqual

	NetPerfDbName    string = "multi_cloud"
	DbFieldCloudName string = "target_cloud_name" // field in the tables of network performance database
	DbFieldRtt       string = "rtt_ms"            // field in the tables of network performance database
)

var (
	NtContainerImage string = DockerRegistry + "/mcnettest:latest"
	NetTestPeriodSec int
	NetTestFuncOn    bool = false
)

var NetTestTaint *apiv1.Taint = &apiv1.Taint{
	Key:    McmKey,
	Effect: taintEffect,
	Value:  taintValue,
}

var NetTestToleration apiv1.Toleration = apiv1.Toleration{
	Operator: tolerationOp,
	Key:      McmKey,
	Effect:   taintEffect,
	Value:    taintValue,
}

// The function to measure network performance between every two clouds
// This function should be executed every time period
func MeasNetPerf() {
	defer func() {
		// Delete server Deployments
		if err := deleteNetTestServers(); err != nil {
			outErr := fmt.Errorf("Cannot delete network performance test servers, Error: %w", err)
			beego.Error(outErr)
			return
		}
	}()

	defer func() {
		// Delete client Jobs
		if err := deleteNetTestClients(); err != nil {
			outErr := fmt.Errorf("Cannot delete network performance test clients, Error: %w", err)
			beego.Error(outErr)
			return
		}
	}()

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

	// Execute network test between every two clouds
	beego.Info("Start to measure network performance between every two clouds.")

	// Run server Deployments
	if err := runNetTestServers(); err != nil {
		outErr := fmt.Errorf("Cannot run network performance test servers, Error: %w", err)
		beego.Error(outErr)
		return
	}

	// Execute client Jobs
	if err := executeNetTestClients(); err != nil {
		outErr := fmt.Errorf("Cannot run network performance test servers, Error: %w", err)
		beego.Error(outErr)
		return
	}

	beego.Info("Finish measuring network performance between every two clouds. Then we will clean up the environment.")
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

func getNetTestServerAppName(cloud Iaas) string {
	return strings.ToLower(cloud.ShowName()) + mcmNetTestServerAppSuffix
}

func getNetTestClientAppName(srcCloud, targetCloud Iaas) string {
	return strings.ToLower(srcCloud.ShowName()) + "-" + strings.ToLower(targetCloud.ShowName()) + mcmNetTestClientAppSuffix
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

// Initialize the tables by specifying the name of the database in every command.
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

// Initialize the tables after executing the command `use <database name>;` without specifying the name of the database in every command.
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

// run all network performance test servers
func runNetTestServers() error {
	beego.Info("Start to run the network performance test server Deployment on every net test server VM.")

	// Do it in parallel
	var wg sync.WaitGroup
	var errsMu sync.Mutex // the slice in golang is not safe for concurrent read/write
	var errs []error
	for name, cloud := range Clouds {
		beego.Info(fmt.Sprintf("Run the network test server for cloud %s", name))
		wg.Add(1)
		go func(c Iaas) {
			defer wg.Done()
			if err := runNetTestServer(c); err != nil {
				outErr := fmt.Errorf("Cannot run the network test server for cloud [%s], error: [%w]", c.ShowName(), err)
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
		outErr := fmt.Errorf("Failed to run the network performance test server Deployments, Error: %w", sumErr)
		beego.Error(outErr)
		return outErr
	}

	beego.Info("Finish running the network performance test server Deployment on every net test server VM.")
	return nil
}

// run the network performance test server on a specified cloud
func runNetTestServer(cloud Iaas) error {
	serverVmName := getNetTestServerName(cloud)
	serverAppName := getNetTestServerAppName(cloud)

	// In some conditions, for example, restricted by AAU CLAAUDIA network security rules, we can only use Host Network to cross different clouds
	var hostNetwork bool
	if hostNetTest, err := beego.AppConfig.Bool("HostNetTest"); err == nil {
		hostNetwork = hostNetTest
	}

	var app K8sApp = K8sApp{
		Name:        serverAppName,
		Replicas:    1,
		HostNetwork: hostNetwork,
		NodeName:    serverVmName,
		Tolerations: []apiv1.Toleration{NetTestToleration},
		Containers: []K8sContainer{
			K8sContainer{
				Name:  serverAppName,
				Image: NtContainerImage,
			},
		},
	}

	if err := CreateApplication(app); err != nil {
		outErr := fmt.Errorf("Cloud [%s] type [%s], Create network test server application [%s], error: [%w].", cloud.ShowName(), cloud.ShowType(), app.Name, err)
		beego.Error(outErr)
		return outErr
	}

	beego.Info(fmt.Sprintf("Start to wait for the application [%s] running", app.Name))
	if err := WaitForAppRunning(WaitForTimeOut, 10, app.Name); err != nil {
		outErr := fmt.Errorf("Wait for application [%s] running, error: %w", app.Name, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("The application [%s] is already running", app.Name))

	return nil
}

// delete all network performance test servers
func deleteNetTestServers() error {
	beego.Info("Start to delete the network performance test server Deployment on every net test server VM.")

	// Do it in parallel
	var wg sync.WaitGroup
	var errsMu sync.Mutex // the slice in golang is not safe for concurrent read/write
	var errs []error
	for name, cloud := range Clouds {
		beego.Info(fmt.Sprintf("Delete the network test server for cloud %s", name))
		wg.Add(1)
		go func(c Iaas) {
			defer wg.Done()
			if err := deleteNetTestServer(c); err != nil {
				outErr := fmt.Errorf("Cannot delete the network test server for cloud [%s], error: [%w]", c.ShowName(), err)
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
		outErr := fmt.Errorf("Failed to delete the network performance test server Deployments, Error: %w", sumErr)
		beego.Error(outErr)
		return outErr
	}

	beego.Info("Finish deleting the network performance test server Deployment on every net test server VM.")
	return nil
}

// delete the network performance test server on a specified cloud
func deleteNetTestServer(cloud Iaas) error {
	serverAppName := getNetTestServerAppName(cloud)

	if err, _ := DeleteApplication(serverAppName); err != nil {
		outErr := fmt.Errorf("Cloud [%s] type [%s], Delete network test server application [%s], error: [%w].", cloud.ShowName(), cloud.ShowType(), serverAppName, err)
		beego.Error(outErr)
		return outErr
	}

	beego.Info(fmt.Sprintf("The application [%s] is already deleted", serverAppName))

	return nil
}

// From each cloud to each cloud, we run a Kubernetes Job to measure the RTT and write the RTT in the database.
func executeNetTestClients() error {
	beego.Info("Start to execute the network performance test client Job on every net test client VM.")

	// Do it in parallel
	var wg sync.WaitGroup
	var errsMu sync.Mutex // the slice in golang is not safe for concurrent read/write
	var errs []error
	for nameFrom, cloudFrom := range Clouds {
		for nameTo, cloudTo := range Clouds {
			beego.Info(fmt.Sprintf("Execute the network performance test client Job from cloud [%s] to cloud [%s]", nameFrom, nameTo))
			wg.Add(1)
			go func(cF Iaas, cT Iaas) {
				defer wg.Done()
				if err := executeNetTestClient(cF, cT); err != nil {
					outErr := fmt.Errorf("Cannot execute the network performance test client Job from cloud [%s] to cloud [%s], error: [%w]", cF.ShowName(), cT.ShowName(), err)
					beego.Error(outErr)
					errsMu.Lock()
					errs = append(errs, outErr)
					errsMu.Unlock()
					return
				}
			}(cloudFrom, cloudTo)
		}
	}
	wg.Wait()

	if len(errs) != 0 {
		sumErr := HandleErrSlice(errs)
		outErr := fmt.Errorf("Failed to execute the network performance test client Jobs, Error: %w", sumErr)
		beego.Error(outErr)
		return outErr
	}

	beego.Info("Finish executing the network performance test client Job on every net test client VM.")
	return nil
}

// Run a network performance test Job on cloudFrom to measure the RTT from cloudFrom to cloudTo, and write the RTT value in the MySQL database.
func executeNetTestClient(cloudFrom, cloudTo Iaas) error {
	srcK8sNodeName := getNetTestClientName(cloudFrom)
	dstK8sAppName := getNetTestServerAppName(cloudTo)
	cliK8sJobName := getNetTestClientAppName(cloudFrom, cloudTo)

	dstK8sApp, err, _ := GetApplication(dstK8sAppName)
	if err != nil {
		outErr := fmt.Errorf("Get dstK8sApp [%s], error: %w", dstK8sAppName, err)
		beego.Error(outErr)
		return outErr
	}
	if len(dstK8sApp.Hosts) == 0 {
		outErr := fmt.Errorf("len(dstK8sApp.Hosts) is [%d], so we cannot get the IP of the target pod", len(dstK8sApp.Hosts))
		beego.Error(outErr)
		return outErr
	}
	if len(dstK8sApp.Hosts[0].PodIP) == 0 {
		outErr := fmt.Errorf("len(dstK8sApp.Hosts[0].PodIP) is [%d], so we cannot get the IP of the target pod", len(dstK8sApp.Hosts))
		beego.Error(outErr)
		return outErr
	}
	dstPodIp := dstK8sApp.Hosts[0].PodIP

	// In some conditions, for example, restricted by AAU CLAAUDIA network security rules, we can only use Host Network to cross different clouds
	var hostNetwork bool
	if hostNetTest, err := beego.AppConfig.Bool("HostNetTest"); err == nil {
		hostNetwork = hostNetTest
	}

	// For a job I do not need to set the labels and selectors. If I want to do it, I can set `.spec.manualSelector: true` in the job's spec
	var job *batchv1.Job = &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cliK8sJobName,
			Namespace: KubernetesNamespace,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					HostNetwork:   hostNetwork,
					RestartPolicy: apiv1.RestartPolicyOnFailure,
					NodeName:      srcK8sNodeName, // execute this job on the network test client VM of the cloudFrom
					// Tolerations:   []apiv1.Toleration{NetTestToleration}, // toleration is not necessary, I think because our taint is NoSchedule, not NoExecute, and we specify the NodeName.
					Containers: []apiv1.Container{
						apiv1.Container{
							Name:            cliK8sJobName,
							Image:           NtContainerImage,
							ImagePullPolicy: apiv1.PullIfNotPresent,
							// My test docker run command is:
							// docker run -d --entrypoint "bash"  mcnettest:latest client.sh "192.168.100.136" "NOKIA7" "192.168.100.136" "3306" "multicloud" 'AAUproxmox1234!@#' "multi_cloud" "NOKIA8"
							// The client.sh is `net-perf-container-image/client.sh`
							Command: []string{"bash", "client.sh"},
							Args:    []string{dstPodIp, cloudTo.ShowName(), MySqlIp, MySqlPort, MySqlUser, MySqlPasswd, NetPerfDbName, cloudFrom.ShowName()},
						},
					},
				},
			},
		},
	}

	// create the job
	createdJob, err := CreateJob(job)
	if err != nil {
		outErr := fmt.Errorf("CreateJob [%v], error: [%w]", job, err)
		beego.Error(outErr)
		return outErr
	}

	beego.Info(fmt.Sprintf("Job [%s/%s] is created.", createdJob.Namespace, createdJob.Name))

	beego.Info(fmt.Sprintf("Start to wait for the Job [%s/%s] completed.", createdJob.Namespace, createdJob.Name))
	if err := WaitForJobCompleted(WaitForTimeOut, 10, createdJob.Namespace, createdJob.Name); err != nil {
		outErr := fmt.Errorf("Wait for the Job [%s/%s] completed, error: %w", createdJob.Namespace, createdJob.Name, err)
		beego.Error(outErr)
		return outErr
	}
	beego.Info(fmt.Sprintf("The Job [%s/%s] is already completed.", createdJob.Namespace, createdJob.Name))

	return nil
}

// delete all network performance test clients
func deleteNetTestClients() error {
	beego.Info("Start to delete the network performance test client Job on every net test client VM.")

	// Do it in parallel
	var wg sync.WaitGroup
	var errsMu sync.Mutex // the slice in golang is not safe for concurrent read/write
	var errs []error
	for _, cloudFrom := range Clouds {
		for _, cloudTo := range Clouds {
			wg.Add(1)
			go func(cF Iaas, cT Iaas) {
				defer wg.Done()
				if err := deleteNetTestClient(cF, cT); err != nil {
					outErr := fmt.Errorf("Cannot delete the network test client from cloud [%s] to cloud [%s], error: [%w]", cF.ShowName(), cT.ShowName(), err)
					beego.Error(outErr)
					errsMu.Lock()
					errs = append(errs, outErr)
					errsMu.Unlock()
					return
				}
			}(cloudFrom, cloudTo)
		}
	}
	wg.Wait()

	if len(errs) != 0 {
		sumErr := HandleErrSlice(errs)
		outErr := fmt.Errorf("Failed to delete the network performance test client Jobs, Error: %w", sumErr)
		beego.Error(outErr)
		return outErr
	}

	beego.Info("Finish deleting the network performance test client Job on every net test client VM.")
	return nil
}

// delete one client job
func deleteNetTestClient(cloudFrom, cloudTo Iaas) error {
	beego.Info(fmt.Sprintf("Delete the network test client from cloud [%s] to cloud [%s]", cloudFrom.ShowName(), cloudTo.ShowName()))
	cliK8sJobName := getNetTestClientAppName(cloudFrom, cloudTo)
	if err := DeleteJob(KubernetesNamespace, cliK8sJobName); err != nil {
		outErr := fmt.Errorf("Delete Job from cloud [%s] to cloud [%s], Error: %w", cloudFrom.ShowName(), cloudTo.ShowName(), err)
		beego.Error(outErr)
		return outErr
	}
	return nil
}
