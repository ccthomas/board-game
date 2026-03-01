package helper

import "time"

func AreTimesEqual(expected *time.Time, actual *time.Time) bool {
	if expected == nil && actual == nil {
		return true
	}
	if expected == nil || actual == nil {
		return false
	}

	return expected.UnixNano() == actual.UnixNano()
}
