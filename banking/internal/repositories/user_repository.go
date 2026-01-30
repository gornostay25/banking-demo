package repositories

import (
	"context"
	"errors"
	"database/sql"

	"banking/internal/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByEmail retrieves a user by email address
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := new(models.User)
	err := r.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return err
}

// ExistsByEmail checks if a user with the given email exists
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.db.NewSelect().
		Model((*models.User)(nil)).
		Where("email = ?", email).
		Exists(ctx)
	return exists, err
}
