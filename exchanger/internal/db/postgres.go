package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

func New(dbURL string, log *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Error("Failed to open database connection", "error", err)
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if err := db.Ping(); err != nil {
		log.Error("Failed to ping database", "error", err)
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	log.Info("Successfully connected to the database")
	return db, nil
}
