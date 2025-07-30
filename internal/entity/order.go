package entity

import "time"

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
