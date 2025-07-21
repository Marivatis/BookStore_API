package repository

import (
	"BookStore_API/internal/entity"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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
}

type Magazine interface {
	Create(ctx context.Context, mag entity.Magazine) (int, error)
	GetById(ctx context.Context, id int) (entity.Magazine, error)
	Update(ctx context.Context, mag entity.Magazine) error
	Delete(ctx context.Context, id int) error
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
		Book: NewBookRepository(db, logger),
	}
}
