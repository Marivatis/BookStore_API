package dto

import (
	"BookStore_API/internal/entity"
	"github.com/go-playground/validator/v10"
	"time"
)

var validate = validator.New()

type BookCreateRequest struct {
	Name   string  `json:"name" validate:"required"`
	Price  float64 `json:"price"`
	Stock  int     `json:"stock"`
	Author string  `json:"author" validate:"required"`
	Isbn   string  `json:"isbn" validate:"required,len=14"`
}

type BookUpdateRequest struct {
	Name   *string  `json:"name"`
	Price  *float64 `json:"price"`
	Stock  *int     `json:"stock"`
	Author *string  `json:"author"`
	Isbn   *string  `json:"isbn" validate:"len=14"`
}

type BookResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	Author    string    `json:"author"`
	Isbn      string    `json:"isbn"`
	CreatedAt time.Time `json:"createdAt"`
}

func (r *BookCreateRequest) Validate() error {
	return validate.Struct(r)
}

func (r *BookUpdateRequest) Validate() error {
	return validate.Struct(r)
}

func FromEntityBook(b entity.Book) BookResponse {
	return BookResponse{
		Id:        b.Id,
		Name:      b.Name,
		Price:     b.Price,
		Stock:     b.Stock,
		Author:    b.Author,
		Isbn:      b.Isbn,
		CreatedAt: b.CreatedAt,
	}
}

// ToEntity method has pointer receiver in case future logic mutates the receiver.
func (r *BookCreateRequest) ToEntity() entity.Book {
	return entity.Book{
		BaseProduct: entity.BaseProduct{
			Name:  r.Name,
			Price: r.Price,
			Stock: r.Stock,
		},
		Author: r.Author,
		Isbn:   r.Isbn,
	}
}

// ApplyToEntity method has pointer receiver in case future logic mutates the receiver.
func (r *BookUpdateRequest) ApplyToEntity(b *entity.Book) {
	if r.Name != nil {
		b.Name = *r.Name
	}
	if r.Price != nil {
		b.Price = *r.Price
	}
	if r.Stock != nil {
		b.Stock = *r.Stock
	}
	if r.Author != nil {
		b.Author = *r.Author
	}
	if r.Isbn != nil {
		b.Isbn = *r.Isbn
	}
}
