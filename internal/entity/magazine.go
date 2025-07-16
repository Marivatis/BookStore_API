package entity

import "time"

type Magazine struct {
	BaseProduct
	IssueNumber     int
	PublicationDate time.Time
}
