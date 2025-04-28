package database

import (
	"database/sql"
	"fmt"

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

	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)
	// db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("database connected")

	return &DB{db}, nil
}
