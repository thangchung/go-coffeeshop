package utils

import "os"

// ref: https://www.thorsten-hans.com/check-if-application-is-running-in-docker-container/
func IsRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err != nil {
		return false
	}

	return true
}
