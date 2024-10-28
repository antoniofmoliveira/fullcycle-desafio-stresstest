package dto

import "time"

type Red struct {
	Target     string
	SentAt     time.Time
	ReceivedAt time.Time
	StatusCode int
	Duration   time.Duration
}
