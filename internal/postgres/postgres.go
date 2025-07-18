package postgres

import (
	"BookStore_API/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	InsertProductSQL = `INSERT INTO products (type, name, price, stock, created_at)
						VALUES ($1, $2, $3, $4, $5)
						RETURNING id`
)

const (
	InsertBookSQL = `INSERT INTO books (product_id, author, isbn)
					 VALUES ($1, $2, $3)`
)

func NewPostgresDB(ctx context.Context, cfg *config.DBConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	cfgPool, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfgPool.MaxConns = 25
	cfgPool.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx, cfgPool)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to pong pgx pool: %w", err)
	}

	return pool, nil
}
