package service

import (
	"BookStore_API/internal/entity"
	"BookStore_API/internal/repository"
	"context"
	"go.uber.org/zap"
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

type Service struct {
	Book
	Magazine
	Order
}

func NewService(r *repository.Repository, logger *zap.Logger) *Service {
	return &Service{
		Book:     NewBookService(r, logger),
		Magazine: NewMagazineService(r, logger),
		Order:    NewOrderService(r, logger),
	}
}
