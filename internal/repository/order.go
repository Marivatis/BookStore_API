package repository

import (
	"BookStore_API/internal/entity"
	"BookStore_API/internal/postgres"
	"BookStore_API/internal/zaplog"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type OrderRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewOrderRepository(db *pgxpool.Pool, logger *zap.Logger) *OrderRepository {
	return &OrderRepository{
		db:     db,
		logger: logger,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order entity.Order) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logDebugOrderOperation("insert", order)

	var orderId int

	// order insert, returning 'orderId'
	err = tx.QueryRow(ctx, postgres.InsertOrdersSQL,
		order.Status, start,
	).Scan(&orderId)
	if err != nil {
		return 0, handleDBError(r.logger, err, "insert_order", start, "failed to insert order")
	}

	// order item insert
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, postgres.InsertOrderItemsSQL,
			orderId, item.Product.Id, item.Quantity, item.Product.Price)
		if err != nil {
			return 0, handleDBError(r.logger, err, "insert_order_item", start, "failed to insert order item")
		}
	}

	r.logger.Info("Order inserted successfully",
		zap.String("operation", "insert"),
		zap.Int("orderId", orderId),
		zap.Duration("duration", time.Since(start)),
	)
	return orderId, nil
}
func (r *OrderRepository) GetById(ctx context.Context, id int) (entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entity.Order{}, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logger.Debug("Starting repository order operation...",
		zap.String("operation", "get_by_id"),
		zap.Int("id", id),
	)

	var order entity.Order

	// order get by id
	err = tx.QueryRow(ctx, postgres.GetByIdOrdersSQL, id).
		Scan(&order.Status, &order.CreatedAt)
	if err != nil {
		return entity.Order{}, handleDBError(r.logger, err, "get_by_id_order", start, "failed to get order by id")
	}

	order.Id = id

	rows, err := tx.Query(ctx, postgres.GetByOrderIdOrderItemsSQL, id)
	if err != nil {
		return entity.Order{}, handleDBError(r.logger, err, "get_by_order_id_order_items", start, "failed to get order items by order id")
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.OrderItem

		err = rows.Scan(
			&item.Product.Id,
			&item.Quantity,
			&item.Product.Price,
		)
		if err != nil {
			return entity.Order{}, handleDBError(r.logger, err, "scan_order_item", start, "failed to scan order item")
		}

		order.Items = append(order.Items, item)
	}

	if err = rows.Err(); err != nil {
		return entity.Order{}, handleDBError(r.logger, err, "rows_err", start, "failed during rows iteration")
	}

	r.logInfoOrderOperation("get_by_id", start, order)
	return order, nil
}
func (r *OrderRepository) Update(ctx context.Context, order entity.Order) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logDebugOrderOperation("update", order)

	// order update by id
	tag, err := tx.Exec(ctx, postgres.UpdateOrdersSQL, order.Id, order.Status)
	if err != nil {
		return handleDBError(r.logger, err, "update_order", start, "failed to update order by id")
	}

	// order update result check
	if tag.RowsAffected() == 0 {
		r.logger.Warn("no order affected - possibly it does not exist",
			zap.Int("id", order.Id),
		)
	}

	// order items update
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, postgres.UpsertOrderItemsSQL,
			order.Id, item.Product.Id, item.Quantity, item.Product.Price)
		if err != nil {
			return handleDBError(r.logger, err, "update_order_item", start, "failed to update order item")
		}
	}

	// preparing items id array
	ids := make([]int, len(order.Items))
	for i, item := range order.Items {
		ids[i] = item.Product.Id
	}

	// delete unnecessary order items
	_, err = tx.Exec(ctx, postgres.DeleteByOrderIdOrderItemsSQL, order.Id, ids)
	if err != nil {
		return handleDBError(r.logger, err, "delete_order_items", start, "failed to delete unnecessary order items")
	}

	r.logInfoOrderOperation("update", start, order)
	return nil
}
func (r *OrderRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logger.Debug("Starting repository order operation...",
		zap.String("operation", "delete_by_id"),
		zap.Int("id", id),
	)

	// delete order by id
	tag, err := tx.Exec(ctx, postgres.DeleteByIdOrdersSQL, id)
	if err != nil {
		return handleDBError(r.logger, err, "delete_by_id_order", start, "failed to delete order by id")
	}

	// order delete result check
	if tag.RowsAffected() == 0 {
		r.logger.Warn("no order affected - possibly it does not exist",
			zap.Int("id", id),
		)
	}

	var count int

	// order items delete result check
	err = tx.QueryRow(ctx, postgres.ExistsOrderItemsWithOrderId, id).Scan(&count)
	if err != nil {
		return handleDBError(r.logger, err, "exists_oder_items", start, "failed to check order items existence")
	}
	if count > 0 {
		r.logger.Warn("order_items were not deleted as expected",
			zap.Int("remaining_items", count))
	}

	r.logger.Info("Finished repository order operation",
		zap.String("operation", "delete"),
		zap.Int("id", id),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}

func (r *OrderRepository) logDebugOrderOperation(operation string, order entity.Order) {
	fields := append(
		[]zap.Field{zap.String("operation", operation)},
		zaplog.OrderFields(order)...,
	)
	r.logger.Debug("Starting repository order operation...", fields...)
}
func (r *OrderRepository) logInfoOrderOperation(operation string, start time.Time, order entity.Order) {
	fields := append(
		[]zap.Field{
			zap.String("operation", operation),
			zap.Duration("elapsed", time.Since(start)),
		},
		zaplog.OrderFields(order)...,
	)
	r.logger.Info("Finished repository order operation", fields...)
}
