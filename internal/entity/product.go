package entity

import "time"

type Product interface {
	GetId() int
	SetId(id int)
	GetName() string
	SetName(name string)
	GetPrice() float64
	SetPrice(price float64)
	GetStock() int
	SetStock(stock int)
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
func (b *BaseProduct) SetId(id int) {
	b.Id = id
}

func (b *BaseProduct) GetName() string {
	return b.Name
}
func (b *BaseProduct) SetName(name string) {
	b.Name = name
}

func (b *BaseProduct) GetPrice() float64 {
	return b.Price
}
func (b *BaseProduct) SetPrice(price float64) {
	b.Price = price
}

func (b *BaseProduct) GetStock() int {
	return b.Stock
}
func (b *BaseProduct) SetStock(stock int) {
	b.Stock = stock
}
