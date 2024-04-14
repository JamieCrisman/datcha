package pipe

import "time"

type TimeRanged interface {
	// TimeRange returns the start and end time range
	//
	// NOTE: could return a zero time?
	TimeRange() (time.Time, time.Time)
}
