package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/wahyusahajaa/mulo-api-go/app/config"
)

type DB struct {
	*sql.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("postgres", cfg.DbURL)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Pooling setting
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)

	return &DB{db}, nil
}
