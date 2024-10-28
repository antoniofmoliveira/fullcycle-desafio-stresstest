package pool

import (
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

func GetDb() *sql.DB {

	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")

	if err != nil {
		slog.Error("pool.GetDb", "msg", err.Error())
		panic(err)
	}
	return db
}
