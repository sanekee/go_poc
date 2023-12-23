package main

import (
	"time"
)

// isValid return true if current time is on or after begining of year 2021 UTC
func isValid() bool {
	return time.Now().UTC().Before(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
}
