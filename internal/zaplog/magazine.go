package zaplog

import (
	"BookStore_API/internal/entity"
	"go.uber.org/zap"
)

func MagazineFields(mag entity.Magazine) []zap.Field {
	return []zap.Field{
		zap.String("name", mag.Name),
		zap.Int("issueNumber", mag.IssueNumber),
		zap.Time("publicationDate", mag.PublicationDate),
		zap.Float64("price", mag.Price),
		zap.Int("stock", mag.Stock),
	}
}
