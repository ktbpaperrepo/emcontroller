package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"emcontroller/auto-schedule/algorithms"
	asmodel "emcontroller/auto-schedule/model"
)

const (
	cloudsJson string = `{"CLAAUDIAweifan":{"name":"CLAAUDIAweifan","type":"openstack","resources":{"limit":{"vcpu":20,"ram":512000,"vm":5,"volume":10,"storage":10000,"port":500},"inUse":{"vcpu":20,"ram":90112,"vm":5,"volume":5,"storage":5560,"port":11}},"netState":{"CLAAUDIAweifan":{"rtt":0.521},"NOKIA10":{"rtt":2.485},"NOKIA4":{"rtt":2.546},"NOKIA7":{"rtt":2.492},"NOKIA8":{"rtt":2.807}},"k8sNodes":[{"name":"claaudia-large-disk","residualResources":{"cpuCore":7,"memory":13721,"storage":3830}}]},"NOKIA10":{"name":"NOKIA10","type":"proxmox","resources":{"limit":{"vcpu":40,"ram":64288.671875,"vm":-1,"volume":-1,"storage":793.7522621154785,"port":-1},"inUse":{"vcpu":38,"ram":59392,"vm":-1,"volume":-1,"storage":849,"port":-1}},"netState":{"CLAAUDIAweifan":{"rtt":3.127},"NOKIA10":{"rtt":1.064},"NOKIA4":{"rtt":1.193},"NOKIA7":{"rtt":0.915},"NOKIA8":{"rtt":0.884}},"k8sNodes":null},"NOKIA4":{"name":"NOKIA4","type":"proxmox","resources":{"limit":{"vcpu":56,"ram":128796.75390625,"vm":-1,"volume":-1,"storage":1296.5185890197754,"port":-1},"inUse":{"vcpu":26,"ram":67584,"vm":-1,"volume":-1,"storage":729,"port":-1}},"netState":{"CLAAUDIAweifan":{"rtt":2.857},"NOKIA10":{"rtt":0.868},"NOKIA4":{"rtt":0.654},"NOKIA7":{"rtt":0.864},"NOKIA8":{"rtt":1.112}},"k8sNodes":[{"name":"testmem","residualResources":{"cpuCore":3,"memory":10035,"storage":140}},{"name":"n4test","residualResources":{"cpuCore":2.9,"memory":5848,"storage":130}}]},"NOKIA7":{"name":"NOKIA7","type":"proxmox","resources":{"limit":{"vcpu":56,"ram":128796.75390625,"vm":-1,"volume":-1,"storage":831.012393951416,"port":-1},"inUse":{"vcpu":38,"ram":67180,"vm":-1,"volume":-1,"storage":177,"port":-1}},"netState":{"CLAAUDIAweifan":{"rtt":2.765},"NOKIA10":{"rtt":0.783},"NOKIA4":{"rtt":0.987},"NOKIA7":{"rtt":0.718},"NOKIA8":{"rtt":0.856}},"k8sNodes":null},"NOKIA8":{"name":"NOKIA8","type":"proxmox","resources":{"limit":{"vcpu":56,"ram":128796.75390625,"vm":-1,"volume":-1,"storage":831.012393951416,"port":-1},"inUse":{"vcpu":18,"ram":30720,"vm":-1,"volume":-1,"storage":449,"port":-1}},"netState":{"CLAAUDIAweifan":{"rtt":3.273},"NOKIA10":{"rtt":1.054},"NOKIA4":{"rtt":0.834},"NOKIA7":{"rtt":0.814},"NOKIA8":{"rtt":0.689}},"k8sNodes":[{"name":"n8test","residualResources":{"cpuCore":2.9,"memory":5848,"storage":130}}]}}`

	appsJson string = `{"expt-app-0":{"name":"expt-app-0","priority":10,"resources":{"cpuCore":4,"memory":15360,"storage":30},"dependencies":null},"expt-app-1":{"name":"expt-app-1","priority":3,"resources":{"cpuCore":4,"memory":16384,"storage":100},"dependencies":null},"expt-app-10":{"name":"expt-app-10","priority":7,"resources":{"cpuCore":4,"memory":15360,"storage":30},"dependencies":[{"appName":"expt-app-37"}]},"expt-app-11":{"name":"expt-app-11","priority":5,"resources":{"cpuCore":2,"memory":1024,"storage":2},"dependencies":[{"appName":"expt-app-57"},{"appName":"expt-app-34"}]},"expt-app-12":{"name":"expt-app-12","priority":8,"resources":{"cpuCore":2,"memory":1024,"storage":8},"dependencies":[{"appName":"expt-app-59"}]},"expt-app-13":{"name":"expt-app-13","priority":2,"resources":{"cpuCore":4,"memory":16384,"storage":100},"dependencies":[{"appName":"expt-app-25"}]},"expt-app-14":{"name":"expt-app-14","priority":6,"resources":{"cpuCore":2,"memory":1024,"storage":2},"dependencies":[{"appName":"expt-app-12"}]},"expt-app-15":{"name":"expt-app-15","priority":4,"resources":{"cpuCore":2,"memory":8192,"storage":128},"dependencies":null},"expt-app-16":{"name":"expt-app-16","priority":10,"resources":{"cpuCore":4,"memory":2048,"storage":20},"dependencies":null},"expt-app-17":{"name":"expt-app-17","priority":2,"resources":{"cpuCore":4,"memory":15360,"storage":30},"dependencies":[{"appName":"expt-app-28"}]},"expt-app-18":{"name":"expt-app-18","priority":9,"resources":{"cpuCore":1,"memory":500,"storage":0},"dependencies":[{"appName":"expt-app-24"}]},"expt-app-19":{"name":"expt-app-19","priority":9,"resources":{"cpuCore":1,"memory":256,"storage":6},"dependencies":null},"expt-app-2":{"name":"expt-app-2","priority":4,"resources":{"cpuCore":4,"memory":15360,"storage":30},"dependencies":[{"appName":"expt-app-30"},{"appName":"expt-app-53"}]},"expt-app-20":{"name":"expt-app-20","priority":1,"resources":{"cpuCore":2,"memory":8192,"storage":128},"dependencies":[{"appName":"expt-app-22"},{"appName":"expt-app-29"}]},"expt-app-21":{"name":"expt-app-21","priority":7,"resources":{"cpuCore":2,"memory":1024,"storage":2},"dependencies":[{"appName":"expt-app-45"}]},"expt-app-22":{"name":"expt-app-22","priority":6,"resources":{"cpuCore":2,"memory":2048,"storage":1},"dependencies":[{"appName":"expt-app-6"},{"appName":"expt-app-19"}]},"expt-app-23":{"name":"expt-app-23","priority":3,"resources":{"cpuCore":2,"memory":1024,"storage":8},"dependencies":[{"appName":"expt-app-36"},{"appName":"expt-app-53"}]},"expt-app-24":{"name":"expt-app-24","priority":9,"resources":{"cpuCore":4,"memory":2048,"storage":3},"dependencies":null},"expt-app-25":{"name":"expt-app-25","priority":3,"resources":{"cpuCore":2,"memory":8192,"storage":128},"dependencies":[{"appName":"expt-app-23"},{"appName":"expt-app-48"}]},"expt-app-26":{"name":"expt-app-26","priority":2,"resources":{"cpuCore":4,"memory":2048,"storage":20},"dependencies":[{"appName":"expt-app-54"},{"appName":"expt-app-4"}]},"expt-app-27":{"name":"expt-app-27","priority":5,"resources":{"cpuCore":2,"memory":2048,"storage":1},"dependencies":[{"appName":"expt-app-10"}]},"expt-app-28":{"name":"expt-app-28","priority":4,"resources":{"cpuCore":2,"memory":2048,"storage":1},"dependencies":[{"appName":"expt-app-22"},{"appName":"expt-app-53"}]},"expt-app-29":{"name":"expt-app-29","priority":7,"resources":{"cpuCore":2,"memory":1024,"storage":4},"dependencies":null},"expt-app-3":{"name":"expt-app-3","priority":1,"resources":{"cpuCore":8,"memory":8192,"storage":255},"dependencies":[{"appName":"expt-app-13"},{"appName":"expt-app-33"}]},"expt-app-30":{"name":"expt-app-30","priority":5,"resources":{"cpuCore":4,"memory":15360,"storage":30},"dependencies":[{"appName":"expt-app-38"}]},"expt-app-31":{"name":"expt-app-31","priority":4,"resources":{"cpuCore":4,"memory":2048,"storage":3},"dependencies":[{"appName":"expt-app-11"},{"appName":"expt-app-22"}]},"expt-app-32":{"name":"expt-app-32","priority":7,"resources":{"cpuCore":2,"memory":1024,"storage":4},"dependencies":[{"appName":"expt-app-9"}]},"expt-app-33":{"name":"expt-app-33","priority":6,"resources":{"cpuCore":2,"memory":1024,"storage":4},"dependencies":[{"appName":"expt-app-10"}]},"expt-app-34":{"name":"expt-app-34","priority":8,"resources":{"cpuCore":1,"memory":256,"storage":6},"dependencies":null},"expt-app-35":{"name":"expt-app-35","priority":7,"resources":{"cpuCore":2,"memory":1024,"storage":8},"dependencies":null},"expt-app-36":{"name":"expt-app-36","priority":5,"resources":{"cpuCore":4,"memory":15360,"storage":30},"dependencies":[{"appName":"expt-app-27"}]},"expt-app-37":{"name":"expt-app-37","priority":8,"resources":{"cpuCore":4,"memory":2048,"storage":3},"dependencies":[{"appName":"expt-app-49"}]},"expt-app-38":{"name":"expt-app-38","priority":6,"resources":{"cpuCore":2,"memory":8192,"storage":128},"dependencies":[{"appName":"expt-app-9"}]},"expt-app-39":{"name":"expt-app-39","priority":8,"resources":{"cpuCore":2,"memory":8192,"storage":128},"dependencies":[{"appName":"expt-app-43"}]},"expt-app-4":{"name":"expt-app-4","priority":6,"resources":{"cpuCore":4,"memory":15360,"storage":30},"dependencies":[{"appName":"expt-app-56"}]},"expt-app-40":{"name":"expt-app-40","priority":5,"resources":{"cpuCore":1,"memory":500,"storage":0},"dependencies":[{"appName":"expt-app-4"}]},"expt-app-41":{"name":"expt-app-41","priority":6,"resources":{"cpuCore":8,"memory":8192,"storage":255},"dependencies":[{"appName":"expt-app-39"}]},"expt-app-42":{"name":"expt-app-42","priority":3,"resources":{"cpuCore":8,"memory":8192,"storage":255},"dependencies":[{"appName":"expt-app-25"},{"appName":"expt-app-34"}]},"expt-app-43":{"name":"expt-app-43","priority":10,"resources":{"cpuCore":1,"memory":256,"storage":6},"dependencies":null},"expt-app-44":{"name":"expt-app-44","priority":3,"resources":{"cpuCore":1,"memory":256,"storage":6},"dependencies":[{"appName":"expt-app-6"}]},"expt-app-45":{"name":"expt-app-45","priority":7,"resources":{"cpuCore":1,"memory":500,"storage":0},"dependencies":[{"appName":"expt-app-52"}]},"expt-app-46":{"name":"expt-app-46","priority":2,"resources":{"cpuCore":4,"memory":16384,"storage":100},"dependencies":[{"appName":"expt-app-54"},{"appName":"expt-app-7"}]},"expt-app-47":{"name":"expt-app-47","priority":3,"resources":{"cpuCore":1,"memory":500,"storage":0},"dependencies":[{"appName":"expt-app-56"}]},"expt-app-48":{"name":"expt-app-48","priority":4,"resources":{"cpuCore":1,"memory":256,"storage":6},"dependencies":[{"appName":"expt-app-15"}]},"expt-app-49":{"name":"expt-app-49","priority":8,"resources":{"cpuCore":4,"memory":16384,"storage":100},"dependencies":[{"appName":"expt-app-52"}]},"expt-app-5":{"name":"expt-app-5","priority":10,"resources":{"cpuCore":2,"memory":1024,"storage":8},"dependencies":[{"appName":"expt-app-16"},{"appName":"expt-app-43"}]},"expt-app-50":{"name":"expt-app-50","priority":5,"resources":{"cpuCore":2,"memory":8192,"storage":128},"dependencies":[{"appName":"expt-app-54"},{"appName":"expt-app-45"}]},"expt-app-51":{"name":"expt-app-51","priority":5,"resources":{"cpuCore":2,"memory":8192,"storage":128},"dependencies":[{"appName":"expt-app-14"}]},"expt-app-52":{"name":"expt-app-52","priority":8,"resources":{"cpuCore":2,"memory":1024,"storage":2},"dependencies":null},"expt-app-53":{"name":"expt-app-53","priority":8,"resources":{"cpuCore":1,"memory":256,"storage":6},"dependencies":null},"expt-app-54":{"name":"expt-app-54","priority":5,"resources":{"cpuCore":2,"memory":1024,"storage":8},"dependencies":[{"appName":"expt-app-4"}]},"expt-app-55":{"name":"expt-app-55","priority":3,"resources":{"cpuCore":4,"memory":2048,"storage":20},"dependencies":[{"appName":"expt-app-57"},{"appName":"expt-app-54"}]},"expt-app-56":{"name":"expt-app-56","priority":6,"resources":{"cpuCore":1,"memory":500,"storage":0},"dependencies":[{"appName":"expt-app-45"}]},"expt-app-57":{"name":"expt-app-57","priority":5,"resources":{"cpuCore":1,"memory":500,"storage":0},"dependencies":[{"appName":"expt-app-4"}]},"expt-app-58":{"name":"expt-app-58","priority":2,"resources":{"cpuCore":2,"memory":2048,"storage":1},"dependencies":[{"appName":"expt-app-8"},{"appName":"expt-app-33"}]},"expt-app-59":{"name":"expt-app-59","priority":8,"resources":{"cpuCore":8,"memory":8192,"storage":255},"dependencies":[{"appName":"expt-app-5"}]},"expt-app-6":{"name":"expt-app-6","priority":6,"resources":{"cpuCore":2,"memory":1024,"storage":4},"dependencies":[{"appName":"expt-app-12"}]},"expt-app-7":{"name":"expt-app-7","priority":7,"resources":{"cpuCore":2,"memory":2048,"storage":1},"dependencies":null},"expt-app-8":{"name":"expt-app-8","priority":4,"resources":{"cpuCore":2,"memory":8192,"storage":128},"dependencies":[{"appName":"expt-app-38"},{"appName":"expt-app-29"}]},"expt-app-9":{"name":"expt-app-9","priority":8,"resources":{"cpuCore":1,"memory":500,"storage":0},"dependencies":null}}`

	appsOrderJson string = `["expt-app-35","expt-app-56","expt-app-38","expt-app-13","expt-app-18","expt-app-0","expt-app-16","expt-app-14","expt-app-15","expt-app-34","expt-app-48","expt-app-51","expt-app-55","expt-app-3","expt-app-12","expt-app-23","expt-app-40","expt-app-33","expt-app-44","expt-app-52","expt-app-54","expt-app-1","expt-app-6","expt-app-32","expt-app-5","expt-app-31","expt-app-2","expt-app-27","expt-app-58","expt-app-29","expt-app-41","expt-app-50","expt-app-53","expt-app-22","expt-app-28","expt-app-24","expt-app-26","expt-app-42","expt-app-45","expt-app-47","expt-app-49","expt-app-11","expt-app-17","expt-app-57","expt-app-30","expt-app-4","expt-app-8","expt-app-37","expt-app-19","expt-app-20","expt-app-9","expt-app-46","expt-app-36","expt-app-21","expt-app-25","expt-app-39","expt-app-43","expt-app-59","expt-app-7","expt-app-10"]`
)

var (
	clouds    map[string]asmodel.Cloud
	apps      map[string]asmodel.Application
	appsOrder []string
)

func init() {
	if err := json.Unmarshal([]byte(cloudsJson), &clouds); err != nil {
		panic("error json.Unmarshal clouds: " + err.Error())
	}
	if err := json.Unmarshal([]byte(appsJson), &apps); err != nil {
		panic("error json.Unmarshal clouds: " + err.Error())
	}
	if err := json.Unmarshal([]byte(appsOrderJson), &appsOrder); err != nil {
		panic("error json.Unmarshal clouds: " + err.Error())
	}
}

type cpmpRecord struct {
	cp      float64
	mp      float64
	fitness float64
}

type cpmpRecordSlice []cpmpRecord

func (cmrs cpmpRecordSlice) Len() int {
	return len(cmrs)
}

func (cmrs cpmpRecordSlice) Swap(i, j int) {
	cmrs[i], cmrs[j] = cmrs[j], cmrs[i]
}

func (cmrs cpmpRecordSlice) Less(i, j int) bool {
	return cmrs[i].fitness < cmrs[j].fitness
}

func main() {
	var testRecords []cpmpRecord

	var samplesPerScenario int = 10

	for cp := 0.0; cp <= 1.0; cp += 0.1 {
		for mp := 0.001; mp <= 0.02; mp += 0.001 {

			var fitnessSum float64 = 0
			for i := 0; i < samplesPerScenario; i++ {
				fmt.Printf("testing cp: %g, mp: %g, round %d\n", cp, mp, i)
				thisFitness := testParameters(cp, mp)
				fitnessSum += thisFitness
			}

			testRecords = append(testRecords, cpmpRecord{
				cp:      cp,
				mp:      mp,
				fitness: fitnessSum / float64(samplesPerScenario),
			})

		}
	}

	sort.Sort(cpmpRecordSlice(testRecords))
	for i := 0; i < len(testRecords); i++ {
		fmt.Printf("cp: %g, mp: %g, fitness: %g\n", testRecords[i].cp, testRecords[i].mp, testRecords[i].fitness)
	}
}

func testParameters(crossoverProbability float64, mutationProbability float64) float64 {
	mcssgaInstance := algorithms.NewMcssga(200, 5000, crossoverProbability, mutationProbability, 200, algorithms.DefaultExpAppCompuTimeOneCpu)
	solution, err := mcssgaInstance.Schedule(clouds, apps, appsOrder)
	if err != nil {
		panic(fmt.Sprintf("mcssgaInstance.Schedule, crossoverProbability: %g, mutationProbability: %g, error: %s", crossoverProbability, mutationProbability, err.Error()))
	}
	return mcssgaInstance.Fitness(clouds, apps, solution)
}
