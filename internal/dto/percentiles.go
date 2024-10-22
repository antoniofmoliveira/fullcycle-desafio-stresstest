package dto

import "time"

type Percentiles struct {
	P10 time.Duration
	P25 time.Duration
	P50 time.Duration
	P75 time.Duration
	P90 time.Duration
	P99 time.Duration
}
