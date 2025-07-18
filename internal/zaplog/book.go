package zaplog

import (
	"BookStore_API/internal/entity"
	"go.uber.org/zap"
)

func BookFields(book entity.Book) []zap.Field {
	return []zap.Field{
		zap.String("name", book.Name),
		zap.String("author", book.Author),
		zap.String("isbn", book.Isbn),
		zap.Float64("price", book.Price),
		zap.Int("stock", book.Stock),
	}
}
