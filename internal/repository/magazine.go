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

type MagazineRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewMagazineRepository(db *pgxpool.Pool, logger *zap.Logger) *MagazineRepository {
	return &MagazineRepository{
		db:     db,
		logger: logger,
	}
}

func (r *MagazineRepository) Create(ctx context.Context, mag entity.Magazine) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logDebugMagazineOperation("insert", mag)

	var id int

	// product insert, returning 'id'
	err = tx.QueryRow(ctx, postgres.InsertProductsSQL,
		"magazine", mag.Name, mag.Price, mag.Stock, start,
	).Scan(&id)
	if err != nil {
		return 0, handleDBError(r.logger, err, "insert_product", start, "failed to insert product")
	}

	// magazine insert
	_, err = tx.Exec(ctx, postgres.InsertMagazinesSQL,
		id, mag.IssueNumber, mag.PublicationDate,
	)
	if err != nil {
		return 0, handleDBError(r.logger, err, "insert_magazine", start, "failed to insert magazine")
	}

	r.logger.Info("Magazine inserted successfully",
		zap.String("operation", "insert"),
		zap.Int("id", id),
		zap.Duration("duration", time.Since(start)),
	)
	return id, nil
}
func (r *MagazineRepository) GetById(ctx context.Context, id int) (entity.Magazine, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return entity.Magazine{}, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logger.Debug("Starting repository magazine operation...",
		zap.String("operation", "get_by_id"),
		zap.Int("id", id),
	)

	var mag entity.Magazine
	var productType string

	// product get by id
	err = tx.QueryRow(ctx, postgres.GetByIdProductsSQL, id).
		Scan(&mag.Id, &productType, &mag.Name, &mag.Price, &mag.Stock, &mag.CreatedAt)
	if err != nil {
		return entity.Magazine{}, handleDBError(r.logger, err, "get_by_id_product", start, "failed to get product by id")
	}

	// product type check
	if productType != "magazine" {
		return entity.Magazine{}, fmt.Errorf("%w: expected 'magazine', got: '%s'", ErrInvalidProductType, productType)
	}

	// mag get by id
	err = tx.QueryRow(ctx, postgres.GetByIdMagazinesSQL, id).
		Scan(&mag.IssueNumber, &mag.PublicationDate)
	if err != nil {
		return entity.Magazine{}, handleDBError(r.logger, err, "get_by_id_magazine", start, "failed to get magazine by id")
	}

	r.logInfoMagazineOperation("get_by_id", start, mag)
	return mag, nil
}
func (r *MagazineRepository) Update(ctx context.Context, mag entity.Magazine) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logDebugMagazineOperation("update", mag)

	// product update by id
	tag, err := tx.Exec(ctx, postgres.UpdateProductsSQL, mag.Id, mag.Name, mag.Price, mag.Stock)
	if err != nil {
		return handleDBError(r.logger, err, "update_product", start, "failed to update product by id")
	}

	// product update result check
	if tag.RowsAffected() == 0 {
		r.logger.Warn("no product affected - possibly it does not exist",
			zap.Int("id", mag.Id),
		)
	}

	// magazine update by id
	tag, err = tx.Exec(ctx, postgres.UpdateMagazinesSQL,
		mag.Id, mag.IssueNumber, mag.PublicationDate)
	if err != nil {
		return handleDBError(r.logger, err, "update_magazine", start, "failed to update magazine by id")
	}

	// magazine update result check
	if tag.RowsAffected() == 0 {
		r.logger.Warn("no magazine affected - possibly it does not exist",
			zap.Int("id", mag.Id),
		)
	}

	r.logInfoMagazineOperation("update", start, mag)
	return nil
}
func (r *MagazineRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	// transaction initialization
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer finalizeTx(r.logger, ctx, tx, &err)

	r.logger.Debug("Starting repository magazine operation...",
		zap.String("operation", "delete_by_id"),
		zap.Int("id", id),
	)

	// delete magazine by id
	tag, err := tx.Exec(ctx, postgres.DeleteByIdMagazinesSQL, id)
	if err != nil {
		return handleDBError(r.logger, err, "delete_magazine", start, "failed to delete magazine by id")
	}

	// magazine delete result check
	if tag.RowsAffected() == 0 {
		r.logger.Warn("no magazine affected - possibly it does not exist",
			zap.Int("id", id),
		)
	}

	// delete product by id
	tag, err = tx.Exec(ctx, postgres.DeleteByIdProductsSQL, id)
	if err != nil {
		return handleDBError(r.logger, err, "delete_product", start, "failed to delete product by id")
	}

	// product delete result check
	if tag.RowsAffected() == 0 {
		r.logger.Warn("no product affected - possibly it does not exist",
			zap.Int("id", id),
		)
	}

	r.logger.Info("Finished repository magazine operation",
		zap.String("operation", "delete"),
		zap.Int("id", id),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}
func (r *MagazineRepository) ExistsIssueNumber(ctx context.Context, issueNumber int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	start := time.Now()

	r.logger.Debug("Starting repository book operation...",
		zap.String("operation", "exists_isbn"),
		zap.Int("issue_number", issueNumber),
	)

	var exists bool

	// check issue number existence
	err := r.db.QueryRow(ctx, postgres.ExistsIssueNumberMagazinesSQL, issueNumber).
		Scan(&exists)
	if err != nil {
		return false, handleDBError(r.logger, err, "exists_issue_number", start, "failed to check isbn existence")
	}

	r.logger.Info("Finished repository magazine operation",
		zap.String("operation", "exists_issue_number"),
		zap.Int("issue_number", issueNumber),
		zap.Duration("duration", time.Since(start)),
	)
	return exists, nil
}

func (r *MagazineRepository) logDebugMagazineOperation(operation string, mag entity.Magazine) {
	fields := append(
		[]zap.Field{zap.String("operation", operation)},
		zaplog.MagazineFields(mag)...,
	)
	r.logger.Debug("Starting repository magazine operation...", fields...)
}
func (r *MagazineRepository) logInfoMagazineOperation(operation string, start time.Time, mag entity.Magazine) {
	fields := append(
		[]zap.Field{
			zap.String("operation", operation),
			zap.Duration("elapsed", time.Since(start)),
		},
		zaplog.MagazineFields(mag)...,
	)
	r.logger.Info("Finished repository magazine operation", fields...)
}
