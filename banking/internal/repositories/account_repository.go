package repositories

import (
	"context"
	"database/sql"
	"errors"

	"banking/internal/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var ErrAccountNotFound = errors.New("account not found")

type AccountRepository struct {
	db *bun.DB
}

func NewAccountRepository(db *bun.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	account := new(models.Account)
	err := r.db.NewSelect().
		Model(account).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return account, nil
}

func (r *AccountRepository) GetByIDForUpdate(ctx context.Context, tx bun.Tx, id uuid.UUID) (*models.Account, error) {
	account := new(models.Account)
	err := tx.NewSelect().
		Model(account).
		Where("id = ?", id).
		For("UPDATE").
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return account, nil
}

func (r *AccountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Account, error) {
	var accounts []*models.Account
	err := r.db.NewSelect().
		Model(&accounts).
		Where("user_id = ?", userID).
		Order("currency ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *AccountRepository) GetByIDAndUserID(ctx context.Context, accountID, userID uuid.UUID) (*models.Account, error) {
	account := new(models.Account)
	err := r.db.NewSelect().
		Model(account).
		Where("id = ? AND user_id = ?", accountID, userID).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return account, nil
}
