package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// Wraps sql.Open() and returns a sql.DB connection pool for a given DSN
func OpenDB() (*sql.DB, error) {

	dsn := os.Getenv("DSN")

	//dsn := mysql + "?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, fmt.Errorf("Error al abrir conexión en sqlite: %w", err)
	}

	// SQLite: una conexión para evitar conflictos de escritura
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
