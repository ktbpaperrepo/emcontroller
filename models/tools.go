package models

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

// When getting the IPs of a VM, we do not need to IPs about containers. We can also require that the interface name must be something
func IsIfNeeded(ifName string) bool {
	var allowedPrefixes []string = []string{
		"en",
		"et",
	}
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(ifName, prefix) {
			return true
		}
	}
	return false
}

// a polling function to wait for conditions
func MyWaitFor(timeoutSec int, intervalSec int, predicate func() (bool, error)) error {
	var success bool
	var err error

	start := time.Now().Unix()

	for i := 0; ; i++ {
		// our expected state does not appear before timeout
		if timeoutSec > 0 && time.Now().Unix()-start > int64(timeoutSec) {
			return fmt.Errorf("MyWaitFor timeout")
		}
		beego.Info(fmt.Sprintf("try %d time", i+1))
		// check once every 5 seconds
		beego.Info(fmt.Sprintf("Sleep %v", time.Duration(intervalSec)*time.Second))
		time.Sleep(time.Duration(intervalSec) * time.Second)

		ch := make(chan struct{}, 0)
		go func() {
			defer close(ch)
			success, err = predicate()
		}()

		select {
		case <-ch:
			if err != nil {
				return fmt.Errorf("error in MyWaitFor [%w]", err)
			}
			if success {
				return nil
			}
		case <-time.After(time.Duration(timeoutSec) * time.Second):
			return fmt.Errorf("MyWaitFor a predicate does not return before timeout")
		}
	}
}

// combine multiple errors into one error
func HandleErrSlice(errs []error) error {
	var sumErr string = "\n\r"

	for _, err := range errs {
		sumErr += err.Error() + "\n\r"
	}

	return fmt.Errorf(sumErr)
}

// Execute a task periodically using time.Timer
func CronTaskTimer(f func(), period time.Duration) {
	var t *time.Timer
	for {
		t = time.NewTimer(period)
		<-t.C

		// execute f after the time of "period"
		f()
	}
}

// --------------

// Calculate the available Vcpu of a VM
func CalcVmAvailVcpu(totalVcpu float64) float64 {
	provisionalResult := math.Floor(totalVcpu - ReservedCPULogicCore)
	if provisionalResult < 0 {
		provisionalResult = 0
	}
	return provisionalResult
}

// Calculate the available Memory (MiB) of a VM
func CalcVmAvailRamMiB(totalRamMiB float64) float64 {
	provisionalResult := math.Floor(totalRamMiB - (totalRamMiB*ReservedRamMiBPercentage + ReservedRamMiB))
	if provisionalResult < 0 {
		provisionalResult = 0
	}
	return provisionalResult
}

// Calculate the available Storage (GiB) of a VM
func CalcVmAvailStorGiB(totalStorGiB float64) float64 {
	provisionalResult := math.Floor(totalStorGiB - (totalStorGiB*ReservedStoragePercentage + ReservedStorageGiB))
	if provisionalResult < 0 {
		provisionalResult = 0
	}
	return provisionalResult
}

/*
In this part, the following functions are the inverse functions of the above ones. The only difference is that the above ones use math.Floor, but the following ones use math.Ceil, which makes the inverse functions not strict, but we can see them as the inverse functions.
*/

// Calculate the needed total Vcpu of a VM from its needed available Vcpu
func CalcVmTotalVcpu(availVcpu float64) float64 {
	return math.Ceil(availVcpu + ReservedCPULogicCore)
}

// Calculate the needed total Memory (MiB) of a VM from its needed available Memory (MiB)
func CalcVmTotalRamMiB(availRamMiB float64) float64 {
	return math.Ceil((availRamMiB + ReservedRamMiB) / (1.0 - ReservedRamMiBPercentage))
}

// Calculate the needed total Storage (GiB) of a VM from its needed available Storage (GiB)
func CalcVmTotalStorGiB(availStorGiB float64) float64 {
	return math.Ceil((availStorGiB + ReservedStorageGiB) / (1.0 - ReservedStoragePercentage))
}

// --------------

// for debug, to show the line number of the code, we do not print the log inside this function, and do it outside.
func JsonString(obj interface{}) string {
	if solnBytes, err := json.Marshal(obj); err != nil {
		return fmt.Sprintf("Error [%s] when json.Marshal [%+v]", err.Error(), obj)
	} else {
		return string(solnBytes)
	}
}
