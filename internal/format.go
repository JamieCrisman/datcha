package internal

import "time"

func Mbpx(count uint64, total time.Duration, scale time.Duration) float64 {
	if total < scale {
		total = scale
	}
	return (float64(count*8) / float64(total/scale)) / 1_000_000
}
