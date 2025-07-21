package service

import (
	"BookStore_API/internal/entity"
	"BookStore_API/internal/repository"
	"context"
	"go.uber.org/zap"
)

type BookService struct {
	repo   *repository.Repository
	logger *zap.Logger
}

func NewBookService(repo *repository.Repository, logger *zap.Logger) *BookService {
	return &BookService{
		repo:   repo,
		logger: logger,
	}
}

func (s *BookService) Create(ctx context.Context, book entity.Book) (int, error) {
	return s.repo.Book.Create(ctx, book)
}
func (s *BookService) GetById(ctx context.Context, id int) (entity.Book, error) {
	return s.repo.Book.GetById(ctx, id)
}
func (s *BookService) Update(ctx context.Context, book entity.Book) error {
	return s.repo.Book.Update(ctx, book)
}
func (s *BookService) Delete(ctx context.Context, id int) error {
	return s.repo.Book.Delete(ctx, id)
}
