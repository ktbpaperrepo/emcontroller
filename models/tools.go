package models

import "strings"

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
