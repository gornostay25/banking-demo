package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"banking/internal/database"
	"banking/internal/models"
	"banking/internal/repositories"
	"banking/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type AccountService struct {
	db          database.Service
	accountRepo *repositories.AccountRepository
}

// AccountResponse represents an account in API responses
type AccountResponse struct {
	ID        string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Currency  models.Currency `json:"currency" enums:"USD,EUR" example:"USD"`
	Balance   string          `json:"balance" example:"1000.00"`
	CreatedAt time.Time       `json:"created_at" example:"2026-01-30T12:00:00Z"`
	UpdatedAt time.Time       `json:"updated_at" example:"2026-01-30T12:00:00Z"`
}

// AccountBalanceResponse represents account balance response
type AccountBalanceResponse struct {
	ID       string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Currency models.Currency `json:"currency" enums:"USD,EUR" example:"USD"`
	Balance  string          `json:"balance" example:"1000.00"`
}

type AccountBalanceURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type ErrorResponse = utils.ErrorResponse

func NewAccountService(db database.Service) *AccountService {
	accountRepo := repositories.NewAccountRepository(db.DB())
	return &AccountService{
		db:          db,
		accountRepo: accountRepo,
	}
}

// ListAccountsHandler godoc
// @Summary List user's accounts
// @Description Returns all accounts for the authenticated user with balances
// @Tags accounts
// @Accept json
// @Produce json
// @Success 200 {array} services.AccountResponse "List of user accounts"
// @Failure 401 {object} services.ErrorResponse
// @Failure 500 {object} services.ErrorResponse
// @Router /api/accounts [get]
// @Security BearerAuth
func (s *AccountService) ListAccountsHandler(c *gin.Context) {
	// User is guaranteed to be authenticated by middleware
	user := utils.MustGetUserFromContext(c)

	ctx := c.Request.Context()
	accounts, err := s.accountRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		utils.RespondWithInternalError(c, "Failed to retrieve accounts", err)
		return
	}

	// Convert to response format
	response := make([]AccountResponse, 0, len(accounts))
	for _, account := range accounts {
		response = append(response, AccountResponse{
			ID:        account.ID.String(),
			Currency:  account.Currency,
			Balance:   account.Balance.String(),
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetAccountBalanceHandler godoc
// @Summary Get account balance
// @Description Returns balance for a specific account belonging to the authenticated user
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID" example:"550e8400-e29b-41d4-a716-446655440000"
// @Success 200 {object} services.AccountBalanceResponse "Account balance information"
// @Failure 400 {object} services.ErrorResponse
// @Failure 401 {object} services.ErrorResponse
// @Failure 404 {object} services.ErrorResponse
// @Failure 500 {object} services.ErrorResponse
// @Router /api/accounts/{id}/balance [get]
// @Security BearerAuth
func (s *AccountService) GetAccountBalanceHandler(c *gin.Context) {
	// User is guaranteed to be authenticated by middleware
	user := utils.MustGetUserFromContext(c)

	var uri AccountBalanceURI
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.RespondWithBadRequest(c, "Invalid request parameters", utils.FormatValidationErrors(err))
		return
	}

	accountID, err := uuid.Parse(uri.ID)
	if err != nil {
		utils.RespondWithBadRequest(c, "Invalid account ID format", err)
		return
	}

	ctx := c.Request.Context()
	account, err := s.accountRepo.GetByIDAndUserID(ctx, accountID, user.ID)
	if err != nil {
		if err == repositories.ErrAccountNotFound {
			utils.RespondWithNotFound(c, "Account not found or does not belong to user")
			return
		}
		utils.RespondWithInternalError(c, "Failed to retrieve account", err)
		return
	}

	response := AccountBalanceResponse{
		ID:       account.ID.String(),
		Currency: account.Currency,
		Balance:  account.Balance.String(),
	}

	c.JSON(http.StatusOK, response)
}

func lockAccountsForUpdate(ctx context.Context, tx bun.Tx, accountRepo *repositories.AccountRepository, account1ID, account2ID uuid.UUID) (*models.Account, *models.Account, error) {
	var first, second *models.Account
	var err error

	// Lock accounts in consistent order (by ID) to prevent deadlocks
	if account1ID.String() < account2ID.String() {
		first, err = accountRepo.GetByIDForUpdate(ctx, tx, account1ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to lock first account: %w", err)
		}
		second, err = accountRepo.GetByIDForUpdate(ctx, tx, account2ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to lock second account: %w", err)
		}
	} else {
		second, err = accountRepo.GetByIDForUpdate(ctx, tx, account2ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to lock second account: %w", err)
		}
		first, err = accountRepo.GetByIDForUpdate(ctx, tx, account1ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to lock first account: %w", err)
		}
	}

	return first, second, nil
}
