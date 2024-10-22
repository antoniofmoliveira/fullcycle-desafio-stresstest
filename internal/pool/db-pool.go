package pool

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func GetDb() *sql.DB {

	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")

	if err != nil {
		panic(err)
	}
	return db
}
