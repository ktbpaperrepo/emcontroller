package models

import (
	"fmt"
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

// The function to measure network performance between every two clouds
func MeasNetPerf() {
	beego.Info("Start to measure network performance between every two clouds.")
	for name, _ := range Clouds {
		beego.Info(fmt.Sprintf("We have cloud %s", name))
	}
	beego.Info("Finish measuring network performance between every two clouds.")
}
