package config

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

// Wraps sql.Open() and returns a sql.DB connection pool for a given DSN
func OpenDB() (*sql.DB, error) {

	dsn := os.Getenv("DSN")

	//dsn := mysql + "?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
