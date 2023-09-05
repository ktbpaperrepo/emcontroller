package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"emcontroller/auto-schedule/experiments/server/cpu"
	"emcontroller/auto-schedule/experiments/server/dependency"
	"emcontroller/auto-schedule/experiments/server/memory"
	"emcontroller/auto-schedule/experiments/server/storage"
)

var cpuNumToUse int

// parameter 1: workload, the input value of cumulative sum
// parameter 2: the number of CPU cores to use in this program
// parameter 3: the size of memory that this program can use, unit MiB
// parameter 4: the size of storage that this program can use, unit GiB
// other parameters: dependent URLs for this app to call
// test commands:
// go run main.go 5000000 4 2048 15 http://192.168.100.96:3333/experiment http://192.168.100.97:3333/experiment http://192.168.100.98:3333/experiment
// time curl -i http://localhost:3333/experiment
func main() {

	// read input parameters
	if workload, err := strconv.Atoi(os.Args[1]); err != nil {
		panic(fmt.Sprintf("parse %s to int error: %s", os.Args[1], err.Error()))
	} else {
		cpu.Workload = workload
	}
	if cpus, err := strconv.Atoi(os.Args[2]); err != nil {
		panic(fmt.Sprintf("parse %s to int error: %s", os.Args[2], err.Error()))
	} else {
		cpuNumToUse = cpus
	}
	if mem, err := strconv.ParseFloat(os.Args[3], 64); err != nil {
		panic(fmt.Sprintf("parse %s to float64 error: %s", os.Args[3], err.Error()))
	} else {
		memory.MemCanUse = mem
	}
	if disk, err := strconv.ParseFloat(os.Args[4], 64); err != nil {
		panic(fmt.Sprintf("parse %s to float64 error: %s", os.Args[4], err.Error()))
	} else {
		storage.DiskCanUse = disk
	}

	// set default
	if cpu.Workload == 0 {
		cpu.Workload = 10000000
	}
	if cpuNumToUse == 0 {
		cpuNumToUse = 2
	}

	// all parameters after the storage should be the dependent URLs for this app to call
	dependency.DepUrls = os.Args[5:]

	fmt.Println("finish handling the input parameters.")

	// occupy storage
	storage.Exec()
	fmt.Println("finish occupying storage.")

	// open a goroutine to occupy memory
	fmt.Println("start occupying memory.")
	go memory.Exec()

	fmt.Println("start open http server.")
	// set the maximum number of CPUs to use
	runtime.GOMAXPROCS(cpuNumToUse)

	//// run some routine work in a goroutine, or otherwise this application's workload is too small. The effect is rubbish, so I commented it.
	//for i := 0; i < cpuNumToUse; i++ {
	//	go routineWork()
	//}

	http.HandleFunc("/experiment", reqFunc)

	if err := http.ListenAndServe(":3333", nil); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server closed")
			fmt.Printf("error starting server: %s\n", err)
		} else {
			fmt.Printf("error starting server: %s\n", err)
		}
	}

}

func reqFunc(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	fmt.Println("A request is received.")
	// call all dependent urls first
	if err := dependency.Exec(); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// then execute this app's work
	cpu.Exec()
	end := time.Now()
	exeTimeMs := fmt.Sprintf("%g", float64(end.Sub(start).Microseconds())/1000)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(exeTimeMs))
	return
}

// The effect is rubbish, so I commented it.
func routineWork() {
	counter := 1
	for {
		counter++
		if counter >= 10000 {
			counter = 1
		}
	}
}
