package config

import (
	"database/sql"
	"fmt"
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
		return nil, fmt.Errorf("Error al abrir conexión en supabase: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
