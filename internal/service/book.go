package service

import (
	"BookStore_API/internal/entity"
	"BookStore_API/internal/repository"
	"context"
	"fmt"
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
	exists, err := s.repo.Book.IsbnExists(ctx, book.Isbn)
	if err != nil {
		return 0, fmt.Errorf("check isbn exists: %w", err)
	}
	if exists {
		return 0, fmt.Errorf("book with the same ISBN already exists")
	}

	id, err := s.repo.Book.Create(ctx, book)
	if err != nil {
		return 0, fmt.Errorf("create book: %w", err)
	}

	return id, nil
}
func (s *BookService) GetById(ctx context.Context, id int) (entity.Book, error) {
	return s.repo.Book.GetById(ctx, id)
}
func (s *BookService) Update(ctx context.Context, book entity.Book) error {
	exists, err := s.repo.Book.IsbnExists(ctx, book.Isbn)
	if err != nil {
		return fmt.Errorf("check isbn exists: %w", err)
	}
	if exists {
		return fmt.Errorf("book with the same ISBN already exists")
	}

	err = s.repo.Book.Update(ctx, book)
	if err != nil {
		return fmt.Errorf("update book: %w", err)
	}

	return nil
}
func (s *BookService) Delete(ctx context.Context, id int) error {
	return s.repo.Book.Delete(ctx, id)
}
