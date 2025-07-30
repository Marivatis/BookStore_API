package entity

import "time"

const (
	OrderStatusCreated   = "created"
	OrderStatusAccepted  = "accepted"
	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusShipped   = "shipped"
	OrderStatusDelivered = "delivered"
	OrderStatusCanceled  = "canceled"
)

type Order struct {
	Id        int
	Items     []OrderItem
	Status    string
	CreatedAt time.Time
}

type OrderItem struct {
	Product  BaseProduct
	Quantity int
}
