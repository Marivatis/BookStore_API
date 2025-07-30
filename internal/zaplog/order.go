package zaplog

import (
	"BookStore_API/internal/entity"
	"fmt"
	"go.uber.org/zap"
)

func OrderFields(order entity.Order) []zap.Field {
	fields := []zap.Field{
		zap.String("status", order.Status),
		zap.Int("item_count", len(order.Items)),
	}

	for i, item := range order.Items {
		prefix := fmt.Sprintf("item_%d_", i+1)
		for _, f := range OrderItemFields(item) {
			fields = append(fields, zap.Any(prefix+f.Key, f.Interface))
		}
	}

	return fields
}

func OrderItemFields(item entity.OrderItem) []zap.Field {
	return []zap.Field{
		zap.String("name", item.Product.Name),
		zap.Float64("price", item.Product.Price),
		zap.Int("quantity", item.Quantity),
	}
}
