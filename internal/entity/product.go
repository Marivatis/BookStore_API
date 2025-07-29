package entity

import "time"

type BaseProduct struct {
	Id        int
	Name      string
	Price     float64
	Stock     int
	CreatedAt time.Time
}
