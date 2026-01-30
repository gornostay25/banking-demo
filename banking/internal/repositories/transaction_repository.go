package repositories

import (
	"context"
	"database/sql"
	"errors"

	"banking/internal/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var ErrTransactionNotFound = errors.New("transaction not found")

type TransactionRepository struct {
	db *bun.DB
}

func NewTransactionRepository(db *bun.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	transaction := new(models.Transaction)
	err := r.db.NewSelect().
		Model(transaction).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	return transaction, nil
}

// GetByUserID retrieves transactions for a user with pagination and filters
func (r *TransactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, transactionType *models.TransactionType, page, limit int) ([]*models.Transaction, int, error) {
	var transactions []*models.Transaction

	query := r.db.NewSelect().
		Model(&transactions).
		Where("user_id = ?", userID)

	if transactionType != nil {
		query = query.Where("type = ?", *transactionType)
	}

	// Get total count
	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}

	if page > 0 && limit > 0 {
		query = query.Offset((page - 1) * limit)
	}

	// Order by created_at descending
	err = query.Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}
