package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"time"
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
