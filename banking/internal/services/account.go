package services

import (
	"net/http"
	"time"

	"banking/internal/database"
	"banking/internal/models"
	"banking/internal/repositories"
	"banking/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad request"`
	Error   string `json:"error,omitempty" example:"Invalid input"`
}

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
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve accounts",
			Error:   err.Error(),
		})
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

	accountIDStr := c.Param("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid account ID format",
		})
		return
	}

	ctx := c.Request.Context()
	account, err := s.accountRepo.GetByIDAndUserID(ctx, accountID, user.ID)
	if err != nil {
		if err == repositories.ErrAccountNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Account not found or does not belong to user",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve account",
			Error:   err.Error(),
		})
		return
	}

	response := AccountBalanceResponse{
		ID:       account.ID.String(),
		Currency: account.Currency,
		Balance:  account.Balance.String(),
	}

	c.JSON(http.StatusOK, response)
}
