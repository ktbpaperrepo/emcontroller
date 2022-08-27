package models

import "time"

func Hello(in string) (out string) {
	out = in + "world"
	return
}

func UnixToDate(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}
