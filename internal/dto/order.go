package dto

import (
	"BookStore_API/internal/entity"
	"fmt"
	"time"
)

var validOrderStatuses = map[string]struct{}{
	entity.OrderStatusCreated:   {},
	entity.OrderStatusAccepted:  {},
	entity.OrderStatusPending:   {},
	entity.OrderStatusPaid:      {},
	entity.OrderStatusShipped:   {},
	entity.OrderStatusDelivered: {},
	entity.OrderStatusCanceled:  {},
}

type OrderItemRequest struct {
	ProductId int `json:"productId" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

type OrderItemResponse struct {
	ProductId int     `json:"productId" validate:"required"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
}

type OrderCreateRequest struct {
	Items  []OrderItemRequest `json:"items" validate:"required,dive"`
	Status string             `json:"status" validate:"required"`
}

type OrderUpdateRequest struct {
	Items  *[]OrderItemRequest `json:"items"`
	Status *string             `json:"status"`
}

type OrderResponse struct {
	Id        int                 `json:"id"`
	Status    string              `json:"status"`
	Items     []OrderItemResponse `json:"items"`
	CreatedAt time.Time           `json:"createdAt"`
}

func (r *OrderCreateRequest) Validate() error {
	if _, ok := validOrderStatuses[r.Status]; !ok {
		return fmt.Errorf("invalid order status: %s", r.Status)
	}
	return validate.Struct(r)
}

func (r *OrderUpdateRequest) Validate() error {
	if r.Status != nil {
		if _, ok := validOrderStatuses[*r.Status]; !ok {
			return fmt.Errorf("invalid order status: %s", *r.Status)
		}
	}
	return validate.Struct(r)
}

func FromEntityOrderItem(p entity.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		ProductId: p.Product.Id,
		Name:      p.Product.Name,
		Price:     p.Product.Price,
		Quantity:  p.Quantity,
	}
}

func FromEntityOrder(o entity.Order) OrderResponse {
	items := make([]OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = FromEntityOrderItem(item)
	}

	return OrderResponse{
		Id:        o.Id,
		Items:     items,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
	}
}

func (r *OrderItemRequest) ToEntity() entity.OrderItem {
	return entity.OrderItem{
		Product: entity.BaseProduct{
			Id: r.ProductId,
		},
		Quantity: r.Quantity,
	}
}

func (r *OrderCreateRequest) ToEntity() entity.Order {
	items := make([]entity.OrderItem, len(r.Items))
	for i, item := range r.Items {
		items[i] = item.ToEntity()
	}
	return entity.Order{
		Items:  items,
		Status: r.Status,
	}
}

func (r *OrderUpdateRequest) ApplyToEntity(o *entity.Order) {
	if r.Items != nil {
		items := make([]entity.OrderItem, len(*r.Items))
		for i, item := range *r.Items {
			items[i] = item.ToEntity()
		}
		o.Items = items
	}
	if r.Status != nil {
		o.Status = *r.Status
	}
}
