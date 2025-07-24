package repository

import (
	"BookStore_API/internal/entity"
	"BookStore_API/internal/postgres"
	"BookStore_API/internal/zaplog"
	"context"
	"fmt"
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

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logDebugBookOperation("insert", book)

	var id int

	// product insert, returning 'id'
	err = tx.QueryRow(ctx, postgres.InsertProductsSQL,
		"book", book.Name, book.Price, book.Stock, start,
	).Scan(&id)
	if err != nil {
		return 0, handleDBError(r.logger, err, "insert_product", start, "failed to insert product")
	}

	// book insert
	_, err = tx.Exec(ctx, postgres.InsertBooksSQL,
		id, book.Author, book.Isbn,
	)
	if err != nil {
		return 0, handleDBError(r.logger, err, "insert_book", start, "failed to insert book")
	}

	r.logger.Info("Book inserted successfully",
		zap.String("operation", "insert"),
		zap.Int("id", id),
		zap.Duration("duration", time.Since(start)),
	)
	return id, nil
}
func (r *BookRepository) GetById(ctx context.Context, id int) (entity.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entity.Book{}, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logger.Debug("Starting repository book operation...",
		zap.String("operation", "get_by_id"),
		zap.Int("id", id),
	)

	var book entity.Book
	var productType string

	// product get by id
	err = tx.QueryRow(ctx, postgres.GetByIdProductsSQL, id).
		Scan(&book.Id, &productType, &book.Name, &book.Price, &book.Stock, &book.CreatedAt)
	if err != nil {
		return entity.Book{}, handleDBError(r.logger, err, "get_by_id_product", start, "failed to get product by id")
	}

	// product type check
	if productType != "book" {
		return entity.Book{}, fmt.Errorf("%w: expected 'book', got: '%s'", ErrInvalidProductType, productType)
	}

	// book get by id
	err = tx.QueryRow(ctx, postgres.GetByIdBooksSQL, id).
		Scan(&book.Author, &book.Isbn)
	if err != nil {
		return entity.Book{}, handleDBError(r.logger, err, "get_by_id_book", start, "failed to get book by id")
	}

	r.logInfoBookOperation("get_by_id", start, book)
	return book, nil
}
func (r *BookRepository) Update(ctx context.Context, book entity.Book) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logDebugBookOperation("update", book)

	// product update by id
	tag, err := tx.Exec(ctx, postgres.UpdateProductsSQL, book.Id, book.Name, book.Price, book.Stock)
	if err != nil {
		return handleDBError(r.logger, err, "update_product", start, "failed to update product by id")
	}

	// product update result check
	if tag.RowsAffected() == 0 {
		r.logger.Warn("no product affected - possibly it does not exist",
			zap.Int("id", book.Id),
		)
	}

	// book update by id
	tag, err = tx.Exec(ctx, postgres.UpdateBooksSQL, book.Id, book.Author, book.Isbn)
	if err != nil {
		return handleDBError(r.logger, err, "update_book", start, "failed to update book by id")
	}

	// book update result check
	if tag.RowsAffected() == 0 {
		r.logger.Warn("no book affected - possibly it does not exist",
			zap.Int("id", book.Id),
		)
	}

	r.logInfoBookOperation("update", start, book)
	return nil
}
func (r *BookRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logger.Debug("Starting repository book operation...",
		zap.String("operation", "delete_by_id"),
		zap.Int("id", id),
	)

	// delete book by id
	tag, err := tx.Exec(ctx, postgres.DeleteByIdBooksSQL, id)
	if err != nil {
		return handleDBError(r.logger, err, "delete_book", start, "failed to delete book by id")
	}

	// book delete result check
	if tag.RowsAffected() == 0 {
		r.logger.Warn("no book affected - possibly it does not exist",
			zap.Int("id", id),
		)
	}

	// delete product by id
	tag, err = tx.Exec(ctx, postgres.DeleteByIdProductsSQL, id)
	if err != nil {
		return handleDBError(r.logger, err, "delete_product", start, "failed to delete product by id")
	}

	// product delete result check
	if tag.RowsAffected() == 0 {
		r.logger.Warn("no product affected - possibly it does not exist",
			zap.Int("id", id),
		)
	}

	r.logger.Info("Finished repository book operation",
		zap.String("operation", "delete"),
		zap.Int("id", id),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}

func (r *BookRepository) IsbnExists(ctx context.Context, isbn string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	r.logger.Debug("Starting repository book operation...",
		zap.String("operation", "exists_isbn"),
		zap.String("isbn", isbn),
	)

	var exists bool

	// check isbn existence
	err := r.db.QueryRow(ctx, postgres.ExistsIsbnBooksSQL, isbn).
		Scan(&exists)
	if err != nil {
		return false, handleDBError(r.logger, err, "exists_isbn", start, "failed to check isbn existence")
	}

	r.logger.Info("Finished repository book operation",
		zap.String("operation", "exists_isbn"),
		zap.String("isbn", isbn),
		zap.Duration("duration", time.Since(start)),
	)
	return exists, nil
}

func (r *BookRepository) logDebugBookOperation(operation string, book entity.Book) {
	fields := append(
		[]zap.Field{zap.String("operation", operation)},
		zaplog.BookFields(book)...,
	)
	r.logger.Debug("Starting repository book operation...", fields...)
}
func (r *BookRepository) logInfoBookOperation(operation string, start time.Time, book entity.Book) {
	fields := append(
		[]zap.Field{
			zap.String("operation", operation),
			zap.Duration("elapsed", time.Since(start)),
		},
		zaplog.BookFields(book)...,
	)
	r.logger.Info("Finished repository book operation", fields...)
}
