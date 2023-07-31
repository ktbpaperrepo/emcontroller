package model

import (
	"encoding/json"
	"fmt"
)

// for debug, to show the line number of the code, we do not print the log inside this function, and do it outside.
func JsonString(obj interface{}) string {
	if solnBytes, err := json.Marshal(obj); err != nil {
		return fmt.Sprintf("Error [%s] when json.Marshal [%+v]", err.Error(), obj)
	} else {
		return string(solnBytes)
	}
}
