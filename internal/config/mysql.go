package config

import (
	"database/sql"
	"os"
)

var DB *sql.DB

// Wraps sql.Open() and returns a sql.DB connection pool for a given DSN
func OpenDB() (*sql.DB, error) {

	mysql := os.Getenv("DSN")
	dsn := mysql + "?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
