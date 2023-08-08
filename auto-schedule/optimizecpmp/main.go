package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"emcontroller/auto-schedule/algorithms"
	asmodel "emcontroller/auto-schedule/model"
)

const (
	cloudsJson string = `{"CLAAUDIAweifan":{"name":"CLAAUDIAweifan","type":"openstack","resources":{"limit":{"vcpu":20,"ram":512000,"vm":5,"volume":10,"storage":10000,"port":500},"inUse":{"vcpu":16,"ram":86016,"vm":5,"volume":5,"storage":640,"port":11}},"netState":{"CLAAUDIAweifan":{"rtt":0.592},"HPE1":{"rtt":3.179},"NOKIA10":{"rtt":2.372},"NOKIA4":{"rtt":2.317},"NOKIA7":{"rtt":2.328},"NOKIA8":{"rtt":2.478}},"k8sNodes":[{"name":"nctest","residualResources":{"cpuCore":2.9,"memory":10764,"storage":150}}]},"HPE1":{"name":"HPE1","type":"proxmox","resources":{"limit":{"vcpu":128,"ram":515499.56640625,"vm":-1,"volume":-1,"storage":893.7522621154785,"port":-1},"inUse":{"vcpu":54,"ram":59392,"vm":-1,"volume":-1,"storage":485,"port":-1}},"netState":{"CLAAUDIAweifan":{"rtt":4.063},"HPE1":{"rtt":3.027},"NOKIA10":{"rtt":1.735},"NOKIA4":{"rtt":2.202},"NOKIA7":{"rtt":1.928},"NOKIA8":{"rtt":1.883}},"k8sNodes":null},"NOKIA10":{"name":"NOKIA10","type":"proxmox","resources":{"limit":{"vcpu":40,"ram":64288.671875,"vm":-1,"volume":-1,"storage":893.7522621154785,"port":-1},"inUse":{"vcpu":38,"ram":59392,"vm":-1,"volume":-1,"storage":849,"port":-1}},"netState":{"CLAAUDIAweifan":{"rtt":2.575},"HPE1":{"rtt":1.766},"NOKIA10":{"rtt":0.735},"NOKIA4":{"rtt":0.741},"NOKIA7":{"rtt":0.942},"NOKIA8":{"rtt":0.919}},"k8sNodes":null},"NOKIA4":{"name":"NOKIA4","type":"proxmox","resources":{"limit":{"vcpu":56,"ram":128796.75390625,"vm":-1,"volume":-1,"storage":1396.5185890197754,"port":-1},"inUse":{"vcpu":22,"ram":55296,"vm":-1,"volume":-1,"storage":529,"port":-1}},"netState":{"CLAAUDIAweifan":{"rtt":3.068},"HPE1":{"rtt":2.777},"NOKIA10":{"rtt":0.9},"NOKIA4":{"rtt":0.709},"NOKIA7":{"rtt":0.882},"NOKIA8":{"rtt":1.137}},"k8sNodes":[{"name":"n4test","residualResources":{"cpuCore":2.9,"memory":6668,"storage":150}}]},"NOKIA7":{"name":"NOKIA7","type":"proxmox","resources":{"limit":{"vcpu":56,"ram":128796.75390625,"vm":-1,"volume":-1,"storage":931.012393951416,"port":-1},"inUse":{"vcpu":38,"ram":67180,"vm":-1,"volume":-1,"storage":177,"port":-1}},"netState":{"CLAAUDIAweifan":{"rtt":2.868},"HPE1":{"rtt":1.624},"NOKIA10":{"rtt":0.842},"NOKIA4":{"rtt":0.814},"NOKIA7":{"rtt":0.706},"NOKIA8":{"rtt":0.759}},"k8sNodes":null},"NOKIA8":{"name":"NOKIA8","type":"proxmox","resources":{"limit":{"vcpu":56,"ram":128796.75390625,"vm":-1,"volume":-1,"storage":931.012393951416,"port":-1},"inUse":{"vcpu":18,"ram":30720,"vm":-1,"volume":-1,"storage":449,"port":-1}},"netState":{"CLAAUDIAweifan":{"rtt":2.562},"HPE1":{"rtt":1.603},"NOKIA10":{"rtt":0.825},"NOKIA4":{"rtt":0.843},"NOKIA7":{"rtt":0.838},"NOKIA8":{"rtt":0.67}},"k8sNodes":[{"name":"n8test","residualResources":{"cpuCore":2.9,"memory":6668,"storage":150}}]}}`

	appsJson string = `{"test-app-0":{"name":"test-app-0","priority":5,"resources":{"cpuCore":8,"memory":548,"storage":86},"dependencies":[{"appName":"test-app-39"},{"appName":"test-app-31"},{"appName":"test-app-27"},{"appName":"test-app-4"}]},"test-app-1":{"name":"test-app-1","priority":7,"resources":{"cpuCore":12.8,"memory":8091,"storage":1},"dependencies":[{"appName":"test-app-36"},{"appName":"test-app-34"}]},"test-app-10":{"name":"test-app-10","priority":8,"resources":{"cpuCore":11.9,"memory":9996,"storage":61},"dependencies":[{"appName":"test-app-14"},{"appName":"test-app-4"}]},"test-app-11":{"name":"test-app-11","priority":6,"resources":{"cpuCore":4.3,"memory":8937,"storage":115},"dependencies":[{"appName":"test-app-1"},{"appName":"test-app-10"},{"appName":"test-app-14"}]},"test-app-12":{"name":"test-app-12","priority":5,"resources":{"cpuCore":8.8,"memory":7259,"storage":20},"dependencies":[{"appName":"test-app-0"},{"appName":"test-app-39"},{"appName":"test-app-31"},{"appName":"test-app-28"},{"appName":"test-app-20"},{"appName":"test-app-23"},{"appName":"test-app-37"},{"appName":"test-app-10"},{"appName":"test-app-22"},{"appName":"test-app-6"}]},"test-app-13":{"name":"test-app-13","priority":4,"resources":{"cpuCore":10.3,"memory":2215,"storage":111},"dependencies":[{"appName":"test-app-19"},{"appName":"test-app-2"},{"appName":"test-app-27"},{"appName":"test-app-36"},{"appName":"test-app-34"},{"appName":"test-app-37"},{"appName":"test-app-14"},{"appName":"test-app-26"},{"appName":"test-app-22"},{"appName":"test-app-6"},{"appName":"test-app-9"}]},"test-app-14":{"name":"test-app-14","priority":9,"resources":{"cpuCore":17.3,"memory":3909,"storage":22},"dependencies":null},"test-app-15":{"name":"test-app-15","priority":5,"resources":{"cpuCore":15.6,"memory":4332,"storage":47},"dependencies":[{"appName":"test-app-12"},{"appName":"test-app-19"},{"appName":"test-app-31"},{"appName":"test-app-10"},{"appName":"test-app-9"}]},"test-app-16":{"name":"test-app-16","priority":2,"resources":{"cpuCore":14.7,"memory":4148,"storage":135},"dependencies":[{"appName":"test-app-13"},{"appName":"test-app-19"},{"appName":"test-app-27"},{"appName":"test-app-17"},{"appName":"test-app-37"},{"appName":"test-app-21"},{"appName":"test-app-4"},{"appName":"test-app-9"}]},"test-app-17":{"name":"test-app-17","priority":7,"resources":{"cpuCore":3.1,"memory":2791,"storage":62},"dependencies":[{"appName":"test-app-34"},{"appName":"test-app-21"},{"appName":"test-app-4"}]},"test-app-18":{"name":"test-app-18","priority":2,"resources":{"cpuCore":8.2,"memory":1259,"storage":110},"dependencies":[{"appName":"test-app-32"},{"appName":"test-app-24"},{"appName":"test-app-13"},{"appName":"test-app-15"},{"appName":"test-app-20"},{"appName":"test-app-27"},{"appName":"test-app-36"},{"appName":"test-app-14"},{"appName":"test-app-6"}]},"test-app-19":{"name":"test-app-19","priority":5,"resources":{"cpuCore":11.6,"memory":9516,"storage":152},"dependencies":[{"appName":"test-app-0"},{"appName":"test-app-2"},{"appName":"test-app-1"},{"appName":"test-app-17"},{"appName":"test-app-34"},{"appName":"test-app-37"},{"appName":"test-app-10"}]},"test-app-2":{"name":"test-app-2","priority":6,"resources":{"cpuCore":18.1,"memory":9731,"storage":116},"dependencies":[{"appName":"test-app-28"},{"appName":"test-app-27"},{"appName":"test-app-37"},{"appName":"test-app-21"},{"appName":"test-app-4"}]},"test-app-20":{"name":"test-app-20","priority":7,"resources":{"cpuCore":8.2,"memory":10937,"storage":159},"dependencies":[{"appName":"test-app-27"},{"appName":"test-app-6"}]},"test-app-21":{"name":"test-app-21","priority":9,"resources":{"cpuCore":10.9,"memory":14340,"storage":41},"dependencies":[{"appName":"test-app-4"},{"appName":"test-app-9"}]},"test-app-22":{"name":"test-app-22","priority":9,"resources":{"cpuCore":5.4,"memory":6021,"storage":110},"dependencies":[{"appName":"test-app-9"}]},"test-app-23":{"name":"test-app-23","priority":7,"resources":{"cpuCore":4.1,"memory":13122,"storage":149},"dependencies":[{"appName":"test-app-34"},{"appName":"test-app-37"},{"appName":"test-app-9"}]},"test-app-24":{"name":"test-app-24","priority":4,"resources":{"cpuCore":7.5,"memory":12247,"storage":107},"dependencies":[{"appName":"test-app-12"},{"appName":"test-app-20"},{"appName":"test-app-17"},{"appName":"test-app-34"},{"appName":"test-app-9"}]},"test-app-25":{"name":"test-app-25","priority":4,"resources":{"cpuCore":13.5,"memory":5635,"storage":29},"dependencies":[{"appName":"test-app-15"},{"appName":"test-app-12"},{"appName":"test-app-0"},{"appName":"test-app-28"},{"appName":"test-app-27"},{"appName":"test-app-38"},{"appName":"test-app-37"},{"appName":"test-app-10"}]},"test-app-26":{"name":"test-app-26","priority":9,"resources":{"cpuCore":5.4,"memory":9680,"storage":164},"dependencies":null},"test-app-27":{"name":"test-app-27","priority":7,"resources":{"cpuCore":9.8,"memory":9196,"storage":62},"dependencies":[{"appName":"test-app-17"},{"appName":"test-app-9"}]},"test-app-28":{"name":"test-app-28","priority":6,"resources":{"cpuCore":16.7,"memory":6606,"storage":156},"dependencies":[{"appName":"test-app-23"},{"appName":"test-app-34"},{"appName":"test-app-21"},{"appName":"test-app-4"}]},"test-app-29":{"name":"test-app-29","priority":1,"resources":{"cpuCore":6,"memory":4432,"storage":132},"dependencies":[{"appName":"test-app-33"},{"appName":"test-app-35"},{"appName":"test-app-18"},{"appName":"test-app-32"},{"appName":"test-app-39"},{"appName":"test-app-23"},{"appName":"test-app-21"}]},"test-app-3":{"name":"test-app-3","priority":1,"resources":{"cpuCore":6.7,"memory":810,"storage":163},"dependencies":[{"appName":"test-app-33"},{"appName":"test-app-35"},{"appName":"test-app-32"},{"appName":"test-app-12"},{"appName":"test-app-0"},{"appName":"test-app-20"},{"appName":"test-app-27"},{"appName":"test-app-23"},{"appName":"test-app-17"},{"appName":"test-app-38"},{"appName":"test-app-37"},{"appName":"test-app-21"},{"appName":"test-app-4"}]},"test-app-30":{"name":"test-app-30","priority":2,"resources":{"cpuCore":10.8,"memory":7361,"storage":105},"dependencies":[{"appName":"test-app-39"},{"appName":"test-app-27"}]},"test-app-31":{"name":"test-app-31","priority":6,"resources":{"cpuCore":7.3,"memory":7275,"storage":123},"dependencies":[{"appName":"test-app-1"},{"appName":"test-app-23"},{"appName":"test-app-17"},{"appName":"test-app-36"},{"appName":"test-app-38"},{"appName":"test-app-14"},{"appName":"test-app-26"}]},"test-app-32":{"name":"test-app-32","priority":4,"resources":{"cpuCore":15.9,"memory":6472,"storage":145},"dependencies":[{"appName":"test-app-20"},{"appName":"test-app-27"},{"appName":"test-app-37"},{"appName":"test-app-10"},{"appName":"test-app-22"},{"appName":"test-app-4"}]},"test-app-33":{"name":"test-app-33","priority":2,"resources":{"cpuCore":14.9,"memory":3004,"storage":139},"dependencies":[{"appName":"test-app-8"},{"appName":"test-app-7"},{"appName":"test-app-32"},{"appName":"test-app-24"},{"appName":"test-app-15"},{"appName":"test-app-19"},{"appName":"test-app-0"},{"appName":"test-app-27"},{"appName":"test-app-37"},{"appName":"test-app-14"},{"appName":"test-app-22"}]},"test-app-34":{"name":"test-app-34","priority":8,"resources":{"cpuCore":7.5,"memory":15721,"storage":106},"dependencies":[{"appName":"test-app-10"},{"appName":"test-app-14"}]},"test-app-35":{"name":"test-app-35","priority":2,"resources":{"cpuCore":16.4,"memory":5556,"storage":129},"dependencies":[{"appName":"test-app-18"},{"appName":"test-app-7"},{"appName":"test-app-19"}]},"test-app-36":{"name":"test-app-36","priority":7,"resources":{"cpuCore":14.1,"memory":659,"storage":77},"dependencies":[{"appName":"test-app-38"},{"appName":"test-app-21"},{"appName":"test-app-6"}]},"test-app-37":{"name":"test-app-37","priority":8,"resources":{"cpuCore":22.2,"memory":3545,"storage":100},"dependencies":[{"appName":"test-app-10"},{"appName":"test-app-21"},{"appName":"test-app-26"},{"appName":"test-app-22"}]},"test-app-38":{"name":"test-app-38","priority":8,"resources":{"cpuCore":10.2,"memory":14480,"storage":175},"dependencies":[{"appName":"test-app-26"},{"appName":"test-app-22"}]},"test-app-39":{"name":"test-app-39","priority":5,"resources":{"cpuCore":3.5,"memory":589,"storage":119},"dependencies":[{"appName":"test-app-38"},{"appName":"test-app-6"}]},"test-app-4":{"name":"test-app-4","priority":10,"resources":{"cpuCore":7.2,"memory":215,"storage":154},"dependencies":[{"appName":"test-app-9"}]},"test-app-5":{"name":"test-app-5","priority":1,"resources":{"cpuCore":10.4,"memory":9590,"storage":87},"dependencies":[{"appName":"test-app-32"},{"appName":"test-app-19"},{"appName":"test-app-31"},{"appName":"test-app-2"},{"appName":"test-app-17"},{"appName":"test-app-21"},{"appName":"test-app-4"},{"appName":"test-app-9"}]},"test-app-6":{"name":"test-app-6","priority":10,"resources":{"cpuCore":4.6,"memory":8174,"storage":18},"dependencies":[{"appName":"test-app-4"}]},"test-app-7":{"name":"test-app-7","priority":4,"resources":{"cpuCore":4,"memory":2630,"storage":125},"dependencies":[{"appName":"test-app-15"},{"appName":"test-app-2"},{"appName":"test-app-1"},{"appName":"test-app-27"},{"appName":"test-app-23"},{"appName":"test-app-17"},{"appName":"test-app-36"},{"appName":"test-app-37"}]},"test-app-8":{"name":"test-app-8","priority":3,"resources":{"cpuCore":6.6,"memory":5422,"storage":36},"dependencies":[{"appName":"test-app-12"},{"appName":"test-app-28"},{"appName":"test-app-1"},{"appName":"test-app-23"},{"appName":"test-app-17"},{"appName":"test-app-36"},{"appName":"test-app-37"}]},"test-app-9":{"name":"test-app-9","priority":10,"resources":{"cpuCore":7.4,"memory":9679,"storage":11},"dependencies":null}}`

	appsOrderJson string = `["test-app-0","test-app-1","test-app-10","test-app-11","test-app-12","test-app-13","test-app-14","test-app-15","test-app-16","test-app-17","test-app-18","test-app-19","test-app-2","test-app-20","test-app-21","test-app-22","test-app-23","test-app-24","test-app-25","test-app-26","test-app-27","test-app-28","test-app-29","test-app-3","test-app-30","test-app-31","test-app-32","test-app-33","test-app-34","test-app-35","test-app-36","test-app-37","test-app-38","test-app-39","test-app-4","test-app-5","test-app-6","test-app-7","test-app-8","test-app-9"]`
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
	mcssgaInstance := algorithms.NewMcssga(200, 5000, crossoverProbability, mutationProbability, 200)
	solution, err := mcssgaInstance.Schedule(clouds, apps, appsOrder)
	if err != nil {
		panic(fmt.Sprintf("mcssgaInstance.Schedule, crossoverProbability: %g, mutationProbability: %g, error: %s", crossoverProbability, mutationProbability, err.Error()))
	}
	return mcssgaInstance.Fitness(clouds, apps, solution)
}
