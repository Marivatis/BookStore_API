package service

import (
	"BookStore_API/internal/entity"
	"BookStore_API/internal/repository"
	"context"
	"fmt"
	"go.uber.org/zap"
)

type OrderService struct {
	repo   *repository.Repository
	logger *zap.Logger
}

func NewOrderService(repo *repository.Repository, logger *zap.Logger) *OrderService {
	return &OrderService{
		repo:   repo,
		logger: logger,
	}
}

func (s *OrderService) Create(ctx context.Context, order entity.Order) (int, error) {
	ids := make([]int, len(order.Items))
	for i, item := range order.Items {
		ids[i] = item.Product.Id
	}

	products, err := s.repo.Product.GetByIds(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to get products by ids: %w", err)
	}

	productMap := make(map[int]entity.BaseProduct)
	for _, p := range products {
		productMap[p.Id] = p
	}

	for i, item := range order.Items {
		if prod, ok := productMap[item.Product.Id]; ok {
			order.Items[i].Product.Price = prod.Price
		} else {
			return 0, fmt.Errorf("order creation failed: product with id %d not found in database", item.Product.Id)
		}
	}

	return s.repo.Order.Create(ctx, order)
}
func (s *OrderService) GetById(ctx context.Context, id int) (entity.Order, error) {
	order, err := s.repo.Order.GetById(ctx, id)
	if err != nil {
		return entity.Order{}, fmt.Errorf("order get failed: %w", err)
	}

	ids := make([]int, len(order.Items))
	for i, item := range order.Items {
		ids[i] = item.Product.Id
	}

	products, err := s.repo.Product.GetByIds(ctx, ids)
	if err != nil {
		return entity.Order{}, fmt.Errorf("failed to get products by ids: %w", err)
	}

	productMap := make(map[int]entity.BaseProduct)
	for _, p := range products {
		productMap[p.Id] = p
	}

	for i, item := range order.Items {
		if prod, ok := productMap[item.Product.Id]; ok {
			order.Items[i].Product = prod
		} else {
			return entity.Order{}, fmt.Errorf("order creation failed: product with id %d not found in database", item.Product.Id)
		}
	}

	return order, nil
}
func (s *OrderService) Update(ctx context.Context, order entity.Order) error {
	return s.repo.Order.Update(ctx, order)
}
func (s *OrderService) Delete(ctx context.Context, id int) error {
	return s.repo.Order.Delete(ctx, id)
}
