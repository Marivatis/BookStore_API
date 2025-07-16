package entity

import "time"

type Product interface {
	GetId() int
	GeyName() string
	GetPrice() float64
	GetQuantity() int
}

type BaseProduct struct {
	Id        int
	Name      string
	Price     float64
	Quantity  int
	CreatedAt time.Time
}

func (b *BaseProduct) GetId() int {
	return b.Id
}

func (b *BaseProduct) GetName() string {
	return b.Name
}

func (b *BaseProduct) GetPrice() float64 {
	return b.Price
}

func (b *BaseProduct) GetQuantity() int {
	return b.Quantity
}
