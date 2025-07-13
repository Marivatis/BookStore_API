package postgres

import (
	"BookStore_API/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func NewPostgresDB(cfg *config.DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
