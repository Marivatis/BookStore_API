package repository

import (
	"BookStore_API/internal/entity"
	"BookStore_API/internal/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type BookRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewBookRepository(db *pgxpool.Pool, logger *zap.Logger) *BookRepository {
	return &BookRepository{
		db:     db,
		logger: logger,
	}
}

func (r *BookRepository) Create(ctx context.Context, book entity.Book) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer r.finalizeTx(ctx, tx, &err)

	r.logger.Debug("Starting book entity insert...",
		zap.String("operation", "insert"),
		zap.String("book_title", book.Name),
	)

	var id int

	// product insert, returning 'id'
	err = tx.QueryRow(ctx, postgres.InsertProductSQL,
		"book", book.Name, book.Price, book.Stock, start,
	).Scan(&id)
	if err != nil {
		return 0, r.handleDBError(err, "insert_product", start, "failed to insert product")
	}

	// book insert
	_, err = tx.Exec(ctx, postgres.InsertBookSQL,
		id, book.Author, book.Isbn,
	)
	if err != nil {
		return 0, r.handleDBError(err, "insert_book", start, "failed to insert book")
	}

	r.logger.Info("Book inserted successfully",
		zap.String("operation", "insert_book"),
		zap.Int("id", id),
		zap.Duration("duration", time.Since(start)),
	)
	return id, nil
}

func (r *BookRepository) GetById(ctx context.Context, id int) (entity.Book, error) {
	return entity.Book{}, nil
}
func (r *BookRepository) GetAll(ctx context.Context) ([]entity.Book, error) {
	return []entity.Book{}, nil
}
func (r *BookRepository) Update(ctx context.Context, book entity.Book) error {
	return nil
}
func (r *BookRepository) Delete(ctx context.Context, id int) error {
	return nil
}

func (r *BookRepository) finalizeTx(ctx context.Context, tx pgx.Tx, err *error) {
	if *err != nil {
		if rollErr := tx.Rollback(ctx); rollErr != nil {
			r.logger.Error("rollback failed", zap.Error(rollErr))
		}
		return
	}
	if commitErr := tx.Commit(ctx); commitErr != nil {
		r.logger.Error("commit failed", zap.Error(commitErr))
		*err = commitErr
	}
}

func (r *BookRepository) handleDBError(err error, operation string, start time.Time, msg string) error {
	if errors.Is(err, context.DeadlineExceeded) {
		r.logger.Error(msg,
			zap.String("operation", operation),
			zap.Duration("elapsed", time.Since(start)),
			zap.String("reason", "timeout"),
			zap.Error(err),
		)
		return fmt.Errorf("%s: timeout: %w", operation, err)
	}

	r.logger.Error(msg,
		zap.String("operation", operation),
		zap.Duration("elapsed", time.Since(start)),
		zap.Error(err),
	)
	return fmt.Errorf("%s: failed: %w", operation, err)
}
