package repositories

import (
	"context"
	"time"

	"banking/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type TransactionRepository struct {
	db *bun.DB
}

func NewTransactionRepository(db *bun.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

type UserTransactionDetail struct {
	Transaction *models.Transaction
	Amount      string // decimal.Decimal as string
	Currency    models.Currency
	Direction   string // "send" or "receive"
}

func (r *TransactionRepository) GetUserTransactionDetails(ctx context.Context, userID uuid.UUID, transactionType *models.TransactionType, page, limit int) ([]UserTransactionDetail, int, error) {
	type result struct {
		TransactionID uuid.UUID              `bun:"transaction_id"`
		UserID        uuid.UUID              `bun:"user_id"`
		Type          models.TransactionType `bun:"type"`
		CreatedAt     time.Time              `bun:"created_at"`
		Amount        decimal.Decimal        `bun:"amount"`
		Currency      models.Currency        `bun:"currency"`
	}

	var results []result

	query := r.db.NewSelect().
		TableExpr("transactions t").
		ColumnExpr("t.id AS transaction_id").
		ColumnExpr("t.user_id").
		ColumnExpr("t.type").
		ColumnExpr("t.created_at").
		ColumnExpr("l.amount").
		ColumnExpr("l.currency").
		Join("JOIN ledger l ON l.transaction_id = t.id").
		Join("JOIN accounts a ON a.id = l.account_id").
		Where("a.user_id = ?", userID)

	if transactionType != nil {
		query = query.Where("t.type = ?", *transactionType)
	}

	// Get total count
	countQuery := r.db.NewSelect().
		TableExpr("ledger l").
		Join("JOIN transactions t ON l.transaction_id = t.id").
		Join("JOIN accounts a ON a.id = l.account_id").
		Where("a.user_id = ?", userID)

	if transactionType != nil {
		countQuery = countQuery.Where("t.type = ?", *transactionType)
	}

	count, err := countQuery.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// pagination
	if limit > 0 {
		query = query.Limit(limit)
	}

	if page > 0 && limit > 0 {
		query = query.Offset((page - 1) * limit)
	}

	// Order by created_at descending
	err = query.Order("t.created_at DESC").Scan(ctx, &results)
	if err != nil {
		return nil, 0, err
	}

	details := make([]UserTransactionDetail, 0, len(results))
	for _, res := range results {
		direction := "receive"
		amount := res.Amount
		if amount.IsNegative() {
			direction = "send"
			amount = amount.Neg() // Make positive for display
		}

		details = append(details, UserTransactionDetail{
			Transaction: &models.Transaction{
				ID:        res.TransactionID,
				UserID:    res.UserID,
				Type:      res.Type,
				CreatedAt: res.CreatedAt,
			},
			Amount:    amount.String(),
			Currency:  res.Currency,
			Direction: direction,
		})
	}

	return details, count, nil
}
