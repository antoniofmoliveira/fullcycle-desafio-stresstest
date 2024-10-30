package db

import (
	"context"
	"database/sql"
	"log/slog"
	"reflect"
	"testing"
	"time"

	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/dto"
	_ "github.com/mattn/go-sqlite3"
)

func TestNewDB(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Success",
			args: args{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := sql.Open("sqlite3", "file::memory:")
			if err != nil {
				slog.Error("test sql.Open", "msg", err.Error())
			}
			db := NewDB(d, make(chan *dto.Red))
			now := time.Now()
			s := &dto.Red{
				Target:     "test",
				SentAt:     now,
				ReceivedAt: now,
				StatusCode: 200,
				Duration:   1000,
			}

			go func() {
				db.input <- s
			}()

			query := "INSERT INTO red (target, sent_at, received_at, status_code, duration) VALUES (?, ?, ?, ?, ?)"
			db.db.Exec(query, s.Target, s.SentAt, s.ReceivedAt, s.StatusCode, s.Duration)

			query = "SELECT target, sent_at, received_at, status_code, duration FROM red"

			rows, err := db.db.Query(query)
			if err != nil {
				slog.Error("db.getReds", "msg", err.Error())
			}
			defer rows.Close()
			var reds []*dto.Red
			for rows.Next() {
				r := &dto.Red{}
				err := rows.Scan(&r.Target, &r.SentAt, &r.ReceivedAt, &r.StatusCode, &r.Duration)
				if err != nil {
					slog.Error("db.getReds scan", "msg", err.Error())
				}
				reds = append(reds, r)
			}

			if len(reds) != 1 {
				t.Errorf("Expected 1 red, got %d", len(reds))
			}

			db.db.Exec("DELETE FROM red")
			db.db.Close()

			for {
				select {
				case r := <-db.input:
					if !reflect.DeepEqual(r, s) {
						t.Errorf("Expected %v, got %v", s, r)
					}
					return
				default:
					time.Sleep(1 * time.Second)
					t.Errorf("Timeout waiting for red")
				}
			}
		})
	}
}

func TestDB_Store(t *testing.T) {

	type args struct {
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := sql.Open("sqlite3", "file::memory:")
			if err != nil {
				slog.Error("test sql.Open", "msg", err.Error())
			}
			db := NewDB(d, make(chan *dto.Red))

			go db.Store(context.Background())

			now := time.Now()
			s := &dto.Red{
				Target:     "test",
				SentAt:     now,
				ReceivedAt: now,
				StatusCode: 200,
				Duration:   1000,
			}

			db.input <- s

			time.Sleep(1 * time.Second)

			query := "SELECT target, sent_at, received_at, status_code, duration FROM red"
			reds := db.getReds(query)
			if len(reds) != 1 {
				t.Errorf("Expected 1 red, got %d", len(reds))
			}

			db.db.Exec("DELETE FROM red")
			db.db.Close()

		})
	}
}

func TestDB_GetAllReds(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "Success",
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := sql.Open("sqlite3", "file::memory:")
			if err != nil {
				slog.Error("test sql.Open", "msg", err.Error())
			}
			db := NewDB(d, make(chan *dto.Red))

			go db.Store(context.Background())
			now := time.Now()
			s := &dto.Red{
				Target:     "test",
				SentAt:     now,
				ReceivedAt: now,
				StatusCode: 200,
				Duration:   1000,
			}
			s2 := &dto.Red{
				Target:     "test2",
				SentAt:     now,
				ReceivedAt: now,
				StatusCode: 500,
				Duration:   1000,
			}

			go func() {
				db.input <- s
				db.input <- s2
			}()

			time.Sleep(1 * time.Second)

			reds := db.GetAllReds()
			if len(reds) != tt.want {
				t.Errorf("Expected %d reds, got %d", tt.want, len(reds))
			}

			db.db.Exec("DELETE FROM red")
			db.db.Close()

		})
	}
}

func TestDB_GetRedsWithoutErrors(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "Success",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := sql.Open("sqlite3", "file::memory:")
			if err != nil {
				slog.Error("test sql.Open", "msg", err.Error())
			}
			db := NewDB(d, make(chan *dto.Red))

			go db.Store(context.Background())
			now := time.Now()
			s := &dto.Red{
				Target:     "test",
				SentAt:     now,
				ReceivedAt: now,
				StatusCode: 200,
				Duration:   1000,
			}
			s2 := &dto.Red{
				Target:     "test2",
				SentAt:     now,
				ReceivedAt: now,
				StatusCode: 500,
				Duration:   1000,
			}
			db.input <- s
			db.input <- s2

			time.Sleep(1 * time.Second)

			reds := db.GetRedsWithoutErrors()
			if len(reds) != tt.want {
				t.Errorf("Expected %d reds, got %d", tt.want, len(reds))
			}
			db.db.Exec("DELETE FROM red")
			db.db.Close()

		})
	}
}

func TestDB_GetRedWithErrors(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "Success",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := sql.Open("sqlite3", "file::memory:")
			if err != nil {
				slog.Error("test sql.Open", "msg", err.Error())
			}
			db := NewDB(d, make(chan *dto.Red))

			go db.Store(context.Background())
			now := time.Now()
			s := &dto.Red{
				Target:     "test",
				SentAt:     now,
				ReceivedAt: now,
				StatusCode: 200,
				Duration:   1000,
			}
			s2 := &dto.Red{
				Target:     "test2",
				SentAt:     now,
				ReceivedAt: now,
				StatusCode: 500,
				Duration:   1000,
			}
			db.input <- s
			db.input <- s2

			time.Sleep(1 * time.Second)

			reds := db.GetRedWithErrors()
			if len(reds) != tt.want {
				t.Errorf("Expected %d reds, got %d", tt.want, len(reds))
			}

			db.db.Exec("DELETE FROM red")
			db.db.Close()
		})
	}
}

func TestDB_Close(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := sql.Open("sqlite3", "file::memory:")
			if err != nil {
				slog.Error("test sql.Open", "msg", err.Error())
			}
			db := NewDB(d, make(chan *dto.Red))

			err = db.Close()

			if err != nil {
				t.Errorf("Expected nil, got %v", err)
			}

		})
	}
}
