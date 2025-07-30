package repository

import (
	"BookStore_API/internal/entity"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

var (
	ErrInvalidProductType = errors.New("invalid product type")
	ErrDBOperation        = errors.New("database operation failed")
)

type Book interface {
	Create(ctx context.Context, book entity.Book) (int, error)
	GetById(ctx context.Context, id int) (entity.Book, error)
	Update(ctx context.Context, book entity.Book) error
	Delete(ctx context.Context, id int) error
	IsbnExists(ctx context.Context, isbn string) (bool, error)
}

type Magazine interface {
	Create(ctx context.Context, mag entity.Magazine) (int, error)
	GetById(ctx context.Context, id int) (entity.Magazine, error)
	Update(ctx context.Context, mag entity.Magazine) error
	Delete(ctx context.Context, id int) error
	ExistsIssueNumber(ctx context.Context, issueNumber int) (bool, error)
}

type Order interface {
	Create(ctx context.Context, order entity.Order) (int, error)
	GetById(ctx context.Context, id int) (entity.Order, error)
	Update(ctx context.Context, order entity.Order) error
	Delete(ctx context.Context, id int) error
}

type Repository struct {
	Book
	Magazine
	Order
}

func NewRepository(db *pgxpool.Pool, logger *zap.Logger) *Repository {
	return &Repository{
		Book:     NewBookRepository(db, logger),
		Magazine: NewMagazineRepository(db, logger),
		Order:    NewOrderRepository(db, logger),
	}
}

func finalizeTx(logger *zap.Logger, ctx context.Context, tx pgx.Tx, err *error) {
	if *err != nil {
		if rollErr := tx.Rollback(ctx); rollErr != nil {
			logger.Error("rollback failed", zap.Error(rollErr))
		}
		return
	}
	if commitErr := tx.Commit(ctx); commitErr != nil {
		logger.Error("commit failed", zap.Error(commitErr))
		*err = commitErr
	}
}

func handleDBError(logger *zap.Logger, err error, operation string, start time.Time, msg string) error {
	if errors.Is(err, context.DeadlineExceeded) {
		logger.Error(msg,
			zap.String("operation", operation),
			zap.Duration("elapsed", time.Since(start)),
			zap.String("reason", "timeout"),
			zap.Error(err),
		)
		return fmt.Errorf("%w(%s): timeout: %w", ErrDBOperation, operation, err)
	}

	logger.Error(msg,
		zap.String("operation", operation),
		zap.Duration("elapsed", time.Since(start)),
		zap.Error(err),
	)
	return fmt.Errorf("%w(%s): failed: %w", ErrDBOperation, operation, err)
}
