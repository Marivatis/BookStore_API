package repository

import (
	"BookStore_API/internal/entity"
	"BookStore_API/internal/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type ProductRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewProductRepository(db *pgxpool.Pool, logger *zap.Logger) *ProductRepository {
	return &ProductRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ProductRepository) GetByIds(ctx context.Context, ids []int) ([]entity.BaseProduct, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// preparing zap fields to log
	logFields := make([]zap.Field, 0, len(ids)+1)
	logFields = append(logFields, zap.String("operation", "get_by_ids"))
	for _, id := range ids {
		logFields = append(logFields, zap.Int(fmt.Sprintf("id_%d", id), id))
	}

	r.logger.Debug("Starting repository products operation...", logFields...)

	// get products by ids
	rows, err := r.db.Query(ctx, postgres.GetByIdsProductsSQL, ids)
	if err != nil {
		return nil, handleDBError(r.logger, err, "get_by_ids_products", start, "failed to get products by ids")
	}
	defer rows.Close()

	products := make([]entity.BaseProduct, 0, len(ids))

	// rows parsing
	for rows.Next() {
		var product entity.BaseProduct
		var plug string // placeholder for unused 'type' column

		err = rows.Scan(
			&product.Id,
			&plug, // 'type' column ignored
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
		)

		if err != nil {
			return nil, handleDBError(r.logger, err, "scan_product", start, "failed to scan product")
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, handleDBError(r.logger, err, "rows_err", start, "failed during rows iteration")
	}

	r.logInfoProductOperation("get_by_ids", start)
	return products, nil
}

func (r *ProductRepository) logInfoProductOperation(operation string, start time.Time) {
	fields := append(
		[]zap.Field{
			zap.String("operation", operation),
			zap.Duration("elapsed", time.Since(start)),
		},
	)
	r.logger.Info("Finished repository product operation", fields...)
}
