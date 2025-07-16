package entity

import "time"

type Product interface {
	GetId() int
	GeyName() string
	GetPrice() float64
	GetStock() int
}

type BaseProduct struct {
	Id        int
	Name      string
	Price     float64
	Stock     int
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

func (b *BaseProduct) GetStock() int {
	return b.Stock
}
