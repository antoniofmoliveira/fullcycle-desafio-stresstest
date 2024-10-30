package stats

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/dto"
)

func TestCalculateRed(t *testing.T) {
	type args struct {
		recs []*dto.Red
	}
	tests := []struct {
		name string
		args args
		want map[string]*dto.ResultRed
	}{
		{
			name: "Success",
			args: args{
				recs: mockReds,
			},
			want: map[string]*dto.ResultRed{
				fmt.Sprintf("%d", now.Minute()*100+now.Second()): {
					NumRequestPerSecond:          10,
					NumRequestWithErrorPerSecond: 0,
					NumNetworkErrorPerSecond:     0,
					AverageDuration:              time.Duration(10 * time.Second),
					MaxDuration:                  time.Duration(time.Second),
					MinDuration:                  time.Duration(time.Second),
				},
				fmt.Sprintf("%d", now2.Minute()*100+now2.Second()): {
					NumRequestPerSecond:          17,
					NumRequestWithErrorPerSecond: 10,
					NumNetworkErrorPerSecond:     0,
					AverageDuration:              time.Duration(17 * time.Second),
					MaxDuration:                  time.Duration(time.Second),
					MinDuration:                  time.Duration(time.Second),
				},
				fmt.Sprintf("%d", now3.Minute()*100+now3.Second()): {
					NumRequestPerSecond:          3,
					NumRequestWithErrorPerSecond: 3,
					NumNetworkErrorPerSecond:     0,
					AverageDuration:              time.Duration(3 * time.Second),
					MaxDuration:                  time.Duration(time.Second),
					MinDuration:                  time.Duration(time.Second),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateRed(tt.args.recs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateRed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateErrors(t *testing.T) {
	type args struct {
		recs []*dto.Red
	}
	tests := []struct {
		name string
		args args
		want map[int]*dto.ResultError
	}{
		{
			name: "Success",
			args: args{
				recs: mockReds,
			},
			want: map[int]*dto.ResultError{
				200: {
					ErrorType:                    200,
					NumRequestWithErrorPerSecond: 17,
				},
				500: {
					ErrorType:                    500,
					NumRequestWithErrorPerSecond: 13,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateErrors(tt.args.recs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculatePercentile(t *testing.T) {
	type args struct {
		recs []*dto.Red
	}
	tests := []struct {
		name string
		args args
		want dto.Percentiles
	}{
		{
			name: "Success",
			args: args{
				recs: mockReds,
			},
			want: dto.Percentiles{
				P10: time.Duration(100 * time.Nanosecond),
				P25: time.Duration(100 * time.Nanosecond),
				P50: time.Duration(100 * time.Nanosecond),
				P75: time.Duration(500 * time.Nanosecond),
				P90: time.Duration(500 * time.Nanosecond),
				P99: time.Duration(1 * time.Microsecond),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculatePercentile(tt.args.recs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculatePercentile() = %v, want %v", got, tt.want)
			}
		})
	}
}

var now = time.Now()
var now2 = now.Add(1 * time.Second)
var now3 = now.Add(2 * time.Second)
var now4 = now.Add(3 * time.Second)

var mockReds = []*dto.Red{
	{
		Target:     "test",
		SentAt:     now,
		ReceivedAt: now2,
		StatusCode: 200,
		Duration:   1000,
	},
	{
		Target:     "test",
		SentAt:     now,
		ReceivedAt: now2,
		StatusCode: 200,
		Duration:   1000,
	},
	{
		Target:     "test",
		SentAt:     now,
		ReceivedAt: now2,
		StatusCode: 200,
		Duration:   time.Duration(500),
	},
	{
		Target:     "test",
		SentAt:     now,
		ReceivedAt: now2,
		StatusCode: 200,
		Duration:   time.Duration(500),
	},
	{
		Target:     "test",
		SentAt:     now,
		ReceivedAt: now2,
		StatusCode: 200,
		Duration:   time.Duration(500),
	},
	{
		Target:     "test",
		SentAt:     now,
		ReceivedAt: now2,
		StatusCode: 200,
		Duration:   time.Duration(500),
	},
	{
		Target:     "test",
		SentAt:     now,
		ReceivedAt: now2,
		StatusCode: 200,
		Duration:   time.Duration(500),
	},
	{
		Target:     "test",
		SentAt:     now,
		ReceivedAt: now2,
		StatusCode: 200,
		Duration:   time.Duration(500),
	},
	{
		Target:     "test",
		SentAt:     now,
		ReceivedAt: now2,
		StatusCode: 200,
		Duration:   time.Duration(500),
	},
	{
		Target:     "test",
		SentAt:     now,
		ReceivedAt: now2,
		StatusCode: 200,
		Duration:   time.Duration(500),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 200,
		Duration:   time.Duration(500),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 200,
		Duration:   time.Duration(500),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 200,
		Duration:   time.Duration(500),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 200,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 200,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 200,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 200,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 500,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 500,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 500,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 500,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 500,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 500,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 500,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 500,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 500,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now2,
		ReceivedAt: now3,
		StatusCode: 500,
		Duration:   time.Duration(100),
	},
	{
		Target:     "test",
		SentAt:     now3,
		ReceivedAt: now4,
		StatusCode: 500,
		Duration:   time.Duration(50),
	},
	{
		Target:     "test",
		SentAt:     now3,
		ReceivedAt: now4,
		StatusCode: 500,
		Duration:   time.Duration(50),
	},
	{
		Target:     "test",
		SentAt:     now3,
		ReceivedAt: now4,
		StatusCode: 500,
		Duration:   time.Duration(50),
	},
}
