package dto

import (
	"BookStore_API/internal/entity"
	"time"
)

type MagazineCreateRequest struct {
	Name            string    `json:"name" validate:"required"`
	Price           float64   `json:"price"`
	Stock           int       `json:"stock"`
	IssueNumber     int       `json:"issueNumber" validate:"required"`
	PublicationDate time.Time `json:"publicationDate" validate:"required"`
}

type MagazineUpdateRequest struct {
	Name            *string    `json:"name"`
	Price           *float64   `json:"price"`
	Stock           *int       `json:"stock"`
	IssueNumber     *int       `json:"issueNumber"`
	PublicationDate *time.Time `json:"publicationDate"`
}

type MagazineResponse struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Price           float64   `json:"price"`
	Stock           int       `json:"stock"`
	IssueNumber     int       `json:"issueNumber"`
	PublicationDate time.Time `json:"publicationDate"`
	CreatedAt       time.Time `json:"createdAt"`
}

func (r *MagazineCreateRequest) Validate() error {
	return validate.Struct(r)
}

func (r *MagazineUpdateRequest) Validate() error {
	return validate.Struct(r)
}

func FromEntityMagazine(m entity.Magazine) MagazineResponse {
	return MagazineResponse{
		Id:              m.Id,
		Name:            m.Name,
		Price:           m.Price,
		Stock:           m.Stock,
		IssueNumber:     m.IssueNumber,
		PublicationDate: m.PublicationDate,
		CreatedAt:       m.CreatedAt,
	}
}

// ToEntity method has pointer receiver in case future logic mutates the receiver.
func (r *MagazineCreateRequest) ToEntity() entity.Magazine {
	return entity.Magazine{
		BaseProduct: entity.BaseProduct{
			Name:  r.Name,
			Price: r.Price,
			Stock: r.Stock,
		},
		IssueNumber:     r.IssueNumber,
		PublicationDate: r.PublicationDate,
	}
}

// ApplyToEntity method has pointer receiver in case future logic mutates the receiver.
func (r *MagazineUpdateRequest) ApplyToEntity(m *entity.Magazine) {
	if r.Name != nil {
		m.Name = *r.Name
	}
	if r.Price != nil {
		m.Price = *r.Price
	}
	if r.Stock != nil {
		m.Stock = *r.Stock
	}
	if r.IssueNumber != nil {
		m.IssueNumber = *r.IssueNumber
	}
	if r.PublicationDate != nil {
		m.PublicationDate = *r.PublicationDate
	}
}
