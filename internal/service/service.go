package service

import (
	"BookStore_API/internal/entity"
	"BookStore_API/internal/repository"
	"context"
)

type BookService interface {
	Create(ctx context.Context, book entity.Book) (int, error)
	GetByID(ctx context.Context, id int) (entity.Book, error)
	GetAll(ctx context.Context) ([]entity.Book, error)
	Update(ctx context.Context, book entity.Book) error
	Delete(ctx context.Context, id int) error
}

type MagazineService interface {
	Create(ctx context.Context, mag entity.Magazine) (int, error)
	GetByID(ctx context.Context, id int) (entity.Magazine, error)
	GetAll(ctx context.Context) ([]entity.Magazine, error)
	Update(ctx context.Context, mag entity.Magazine) error
	Delete(ctx context.Context, id int) error
}

type OrderService interface {
	Create(ctx context.Context, order entity.Order) (int, error)
	GetByID(ctx context.Context, id int) (entity.Order, error)
	GetAll(ctx context.Context) ([]entity.Order, error)
	Update(ctx context.Context, order entity.Order) error
	Delete(ctx context.Context, id int) error
}

type Service struct {
	BookService
	MagazineService
	OrderService
}

func NewService(r *repository.Repository) *Service {
	return &Service{}
}
