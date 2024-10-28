package dto

import "time"

type ResultRed struct {
	NumRequestPerSecond          int
	NumRequestWithErrorPerSecond int
	NumNetworkErrorPerSecond     int
	AverageDuration              time.Duration
	MaxDuration                  time.Duration
	MinDuration                  time.Duration
}
