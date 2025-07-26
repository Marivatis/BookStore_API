package service

import (
	"BookStore_API/internal/entity"
	"BookStore_API/internal/repository"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type MagazineService struct {
	repo   *repository.Repository
	logger *zap.Logger
}

func NewMagazineService(repo *repository.Repository, logger *zap.Logger) *MagazineService {
	return &MagazineService{
		repo:   repo,
		logger: logger,
	}
}

func (s *MagazineService) Create(ctx context.Context, mag entity.Magazine) (int, error) {
	exists, err := s.repo.Magazine.ExistsIssueNumber(ctx, mag.IssueNumber)
	if err != nil {
		return 0, fmt.Errorf("check issue number exists: %w", err)
	}
	if exists {
		return 0, fmt.Errorf("magazine with the same issue number already exists")
	}

	id, err := s.repo.Magazine.Create(ctx, mag)
	if err != nil {
		return 0, fmt.Errorf("create magazine: %w", err)
	}

	return id, nil
}
func (s *MagazineService) GetById(ctx context.Context, id int) (entity.Magazine, error) {
	return s.repo.Magazine.GetById(ctx, id)
}
func (s *MagazineService) Update(ctx context.Context, mag entity.Magazine) error {
	exists, err := s.repo.Magazine.ExistsIssueNumber(ctx, mag.IssueNumber)
	if err != nil {
		return fmt.Errorf("check issue number exists: %w", err)
	}
	if exists {
		return fmt.Errorf("magazine with the same issue number already exists")
	}

	err = s.repo.Magazine.Update(ctx, mag)
	if err != nil {
		return fmt.Errorf("update magazine: %w", err)
	}

	return nil
}
func (s *MagazineService) Delete(ctx context.Context, id int) error {
	return s.repo.Magazine.Delete(ctx, id)
}
