package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"well_track/internal/config"
)

func NewPostgresDB(cfg *config.PgConfig) (*sql.DB, error) {
	
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
