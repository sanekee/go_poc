package main

import (
	"testing"
	"time"

	_ "unsafe"
)

var theTime time.Time

//go:linkname timeNow time.Now
func timeNow() time.Time {
	return theTime
}

func TestIsValid(t *testing.T) {
	t.Run("is valid before 2021-01-01", func(t *testing.T) {
		theTime = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		theTime = theTime.Add((time.Duration(-1) * time.Nanosecond))

		if !isValid() {
			t.Log("Current Time", time.Now())
			t.Fail()
		}
	})

	t.Run("is not valid after 2021-01-01", func(t *testing.T) {
		theTime = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

		if isValid() {
			t.Log("Current Time", time.Now())
			t.Fail()
		}
	})
}
