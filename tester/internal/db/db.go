package db

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/antoniofmoliveira/fullcycle-desafio-stresstest/internal/dto"
)

type DB struct {
	db    *sql.DB
	input chan *dto.Red
}

func NewDB(db *sql.DB, input chan *dto.Red) *DB {
	db.Exec("CREATE TABLE IF NOT EXISTS red (target text, sent_at timestamp, received_at timestamp, status_code int, duration int)")
	return &DB{
		db:    db,
		input: input,
	}
}

func (d *DB) Store(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			r := <-d.input
			_, err := d.db.Exec("INSERT INTO red (target, sent_at, received_at, status_code, duration) VALUES ( ?, ?, ?, ?, ?)", r.Target, r.SentAt, r.ReceivedAt, r.StatusCode, r.Duration)
			if err != nil {
				slog.Error("db.Store", "msg", err.Error())
			}
		}
	}
}

func (d *DB) getReds(query string) []*dto.Red {
	rows, err := d.db.Query(query)
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
	return reds
}

func (d *DB) GetAllReds() []*dto.Red {
	return d.getReds("SELECT target, sent_at, received_at, status_code, duration FROM red")
}

func (d *DB) GetRedsWithoutErrors() []*dto.Red {
	return d.getReds("SELECT target, sent_at, received_at, status_code, duration FROM red where status_code = 200")
}

func (d *DB) GetRedWithErrors() []*dto.Red {
	return d.getReds("SELECT target, sent_at, received_at, status_code, duration FROM red WHERE status_code != 200")
}

func (d *DB) Close() error {
	return d.db.Close()
}
