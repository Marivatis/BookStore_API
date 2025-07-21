package postgres

import (
	"BookStore_API/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// products table sql queries
const (
	InsertProductsSQL = `INSERT INTO products (type, name, price, stock, created_at)
						 VALUES ($1, $2, $3, $4, $5)
						 RETURNING id`
	GetByIdProductsSQL = `SELECT id, type, name, price, stock, created_at
						  FROM products
						  WHERE id = $1`
	UpdateProductsSQL = `UPDATE products
						 SET name = $2,
						 	 price = $3,
						 	 stock = $4
						 WHERE id = $1`
	DeleteByIdProductsSQL = `DELETE FROM products
							 WHERE id = $1`
)

// books table sql queries
const (
	InsertBooksSQL = `INSERT INTO books (product_id, author, isbn)
					  VALUES ($1, $2, $3)`
	GetByIdBooksSQL = `SELECT product_id, author, isbn
					   FROM books
					   WHERE product_id = $1`
	UpdateBooksSQL = `UPDATE books
					  SET author = $2,
					      isbn = $3
					  WHERE product_id = $1`
	DeleteByIdBooksSQL = `DELETE FROM books
				  		  WHERE product_id = $1`
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
