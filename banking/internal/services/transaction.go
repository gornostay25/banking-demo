package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"banking/internal/database"
	"banking/internal/models"
	"banking/internal/repositories"
	"banking/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

const (
	// ExchangeRate: 1 USD = 0.92 EUR
	ExchangeRate = 0.92
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type TransactionService struct {
	db              database.Service
	accountRepo     *repositories.AccountRepository
	transactionRepo *repositories.TransactionRepository
}

func NewTransactionService(db database.Service) *TransactionService {
	bunDB := db.DB()
	return &TransactionService{
		db:              db,
		accountRepo:     repositories.NewAccountRepository(bunDB),
		transactionRepo: repositories.NewTransactionRepository(bunDB),
	}
}

// TransferRequest represents a transfer transaction request
// @Description Transfer money between user accounts
type TransferRequest struct {
	ToAccountID uuid.UUID       `json:"to_account_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Amount      string          `json:"amount" binding:"required,amount" example:"100.50"`
	Currency    models.Currency `json:"currency" binding:"required,oneof=USD EUR" enums:"USD,EUR" example:"USD"`
}

// ExchangeRequest represents an exchange transaction request
// @Description Exchange currency within user's accounts
type ExchangeRequest struct {
	FromCurrency models.Currency `json:"from_currency" binding:"required,oneof=USD EUR" enums:"USD,EUR" example:"USD"`
	ToCurrency   models.Currency `json:"to_currency" binding:"required,oneof=USD EUR" enums:"USD,EUR" example:"EUR"`
	Amount       string          `json:"amount" binding:"required,amount" example:"100.00"`
}

// TransactionResponse represents a transaction in API responses
type TransactionResponse struct {
	ID        string                 `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type      models.TransactionType `json:"type" enums:"transfer,exchange" example:"transfer"`
	Amount    string                 `json:"amount" example:"100.50"`
	Currency  models.Currency        `json:"currency" enums:"USD,EUR" example:"USD"`
	Direction string                 `json:"direction" enums:"send,receive" example:"send"` // "send" for outgoing, "receive" for incoming
	CreatedAt time.Time              `json:"created_at" example:"2026-01-30T12:00:00Z"`
}

// TransactionListResponse represents paginated transaction list response
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int                   `json:"total" example:"25"`
	Page         int                   `json:"page" example:"1"`
	Limit        int                   `json:"limit" example:"10"`
}

// TransactionListQuery represents query parameters for transaction history.
type TransactionListQuery struct {
	Type  *string `form:"type" binding:"omitempty,oneof=transfer exchange"`
	Page  int     `form:"page" binding:"omitempty,min=1"`
	Limit int     `form:"limit" binding:"omitempty,min=1,max=100"`
}

// TransferHandler godoc
// @Summary Transfer money between users
// @Description Transfer money from authenticated user's account to another user's account (same currency). Uses double-entry ledger system.
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body services.TransferRequest true "Transfer request"
// @Success 200 {object} services.TransactionResponse "Transaction created successfully"
// @Failure 400 {object} services.ErrorResponse "Bad request - invalid input or insufficient funds"
// @Failure 401 {object} services.ErrorResponse "Unauthorized"
// @Failure 404 {object} services.ErrorResponse "Account not found"
// @Failure 500 {object} services.ErrorResponse "Internal server error"
// @Router /api/transactions/transfer [post]
// @Security BearerAuth
func (s *TransactionService) TransferHandler(c *gin.Context) {
	user := utils.MustGetUserFromContext(c)

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithBadRequest(c, "Invalid request data", utils.FormatValidationErrors(err))
		return
	}

	ctx := c.Request.Context()

	// Get user's account for the currency
	fromAccounts, err := s.accountRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		utils.RespondWithInternalError(c, "Failed to retrieve accounts", err)
		return
	}

	var fromAccount *models.Account
	for _, acc := range fromAccounts {
		if acc.Currency == req.Currency {
			fromAccount = acc
			break
		}
	}

	if fromAccount == nil {
		utils.RespondWithNotFound(c, fmt.Sprintf("Account with currency %s not found", req.Currency))
		return
	}

	// Get recipient account
	toAccount, err := s.accountRepo.GetByID(ctx, req.ToAccountID)
	if err != nil {
		if err == repositories.ErrAccountNotFound {
			utils.RespondWithNotFound(c, "Recipient account not found")
			return
		}
		utils.RespondWithInternalError(c, "Failed to retrieve recipient account", err)
		return
	}

	// Prevent self-transfer
	if toAccount.UserID == user.ID {
		utils.RespondWithBadRequest(c, "Cannot transfer to your own account", nil)
		return
	}

	// Validate same currency
	if toAccount.Currency != req.Currency {
		utils.RespondWithBadRequest(c, "Currency mismatch: transfer must be in the same currency", nil)
		return
	}

	// Validate amount
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		utils.RespondWithBadRequest(c, "Invalid amount format", err)
		return
	}

	// Check sufficient funds
	if fromAccount.Balance.LessThan(amount) {
		utils.RespondWithBadRequest(c, "Insufficient funds", nil)
		return
	}

	// Execute transfer in database transaction
	transaction, err := s.executeTransfer(ctx, user.ID, fromAccount, toAccount, amount, req.Currency)
	if err != nil {
		if errors.Is(err, ErrInsufficientFunds) {
			utils.RespondWithBadRequest(c, "Insufficient funds", nil)
			return
		}
		utils.RespondWithInternalError(c, "Failed to execute transfer", err)
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		ID:        transaction.ID.String(),
		Type:      transaction.Type,
		CreatedAt: transaction.CreatedAt,
	})
}

// ExchangeHandler godoc
// @Summary Exchange currency
// @Description Exchange currency within user's accounts. Fixed exchange rate: 1 USD = 0.92 EUR. Uses double-entry ledger system.
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body services.ExchangeRequest true "Exchange request"
// @Success 200 {object} services.TransactionResponse "Exchange transaction created successfully"
// @Failure 400 {object} services.ErrorResponse "Bad request - invalid input or insufficient funds"
// @Failure 401 {object} services.ErrorResponse "Unauthorized"
// @Failure 404 {object} services.ErrorResponse "Account not found"
// @Failure 500 {object} services.ErrorResponse "Internal server error"
// @Router /api/transactions/exchange [post]
// @Security BearerAuth
func (s *TransactionService) ExchangeHandler(c *gin.Context) {
	user := utils.MustGetUserFromContext(c)

	var req ExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithBadRequest(c, "Invalid request data", utils.FormatValidationErrors(err))
		return
	}

	// Validate currencies
	if req.FromCurrency == req.ToCurrency {
		utils.RespondWithBadRequest(c, "From and to currencies must be different", nil)
		return
	}

	ctx := c.Request.Context()

	// Get user's accounts
	accounts, err := s.accountRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		utils.RespondWithInternalError(c, "Failed to retrieve accounts", err)
		return
	}

	var fromAccount, toAccount *models.Account
	for _, acc := range accounts {
		if acc.Currency == req.FromCurrency {
			fromAccount = acc
		}
		if acc.Currency == req.ToCurrency {
			toAccount = acc
		}
	}

	if fromAccount == nil {
		utils.RespondWithNotFound(c, fmt.Sprintf("Account with currency %s not found", req.FromCurrency))
		return
	}

	if toAccount == nil {
		utils.RespondWithNotFound(c, fmt.Sprintf("Account with currency %s not found", req.ToCurrency))
		return
	}

	// Validate amount
	fromAmount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		utils.RespondWithBadRequest(c, "Invalid amount format", err)
		return
	}

	// Check sufficient funds
	if fromAccount.Balance.LessThan(fromAmount) {
		utils.RespondWithBadRequest(c, "Insufficient funds", nil)
		return
	}

	// Calculate exchange amount
	toAmount := s.calculateExchangeAmount(fromAmount, req.FromCurrency, req.ToCurrency)

	// Execute exchange in database transaction
	transaction, err := s.executeExchange(ctx, user.ID, fromAccount, toAccount, fromAmount, toAmount, req.FromCurrency, req.ToCurrency)
	if err != nil {
		if errors.Is(err, ErrInsufficientFunds) {
			utils.RespondWithBadRequest(c, "Insufficient funds", nil)
			return
		}
		utils.RespondWithInternalError(c, "Failed to execute exchange", err)
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		ID:        transaction.ID.String(),
		Type:      transaction.Type,
		CreatedAt: transaction.CreatedAt,
	})
}

// GetTransactionsHandler godoc
// @Summary Get transaction history
// @Description Returns paginated list of transactions for authenticated user with optional filters
// @Tags transactions
// @Accept json
// @Produce json
// @Param type query string false "Transaction type filter" Enums(transfer, exchange) example:"transfer"
// @Param page query int false "Page number" default(1) minimum(1) example:1
// @Param limit query int false "Items per page" default(10) minimum(1) maximum(100) example:10
// @Success 200 {object} services.TransactionListResponse "Paginated list of transactions"
// @Failure 401 {object} services.ErrorResponse "Unauthorized"
// @Failure 500 {object} services.ErrorResponse "Internal server error"
// @Router /api/transactions [get]
// @Security BearerAuth
func (s *TransactionService) GetTransactionsHandler(c *gin.Context) {
	user := utils.MustGetUserFromContext(c)

	var query TransactionListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid query parameters",
			Error:   utils.FormatValidationErrors(err),
		})
		return
	}

	page := query.Page
	limit := query.Limit
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 5
	}

	var transactionType *models.TransactionType
	if query.Type != nil {
		t := models.TransactionType(*query.Type)
		transactionType = &t
	}

	ctx := c.Request.Context()
	transactionDetails, total, err := s.transactionRepo.GetUserTransactionDetails(ctx, user.ID, transactionType, page, limit)
	if err != nil {
		utils.RespondWithInternalError(c, "Failed to retrieve transactions", err)
		return
	}

	response := TransactionListResponse{
		Transactions: make([]TransactionResponse, 0, len(transactionDetails)),
		Total:        total,
		Page:         page,
		Limit:        limit,
	}

	for _, detail := range transactionDetails {
		response.Transactions = append(response.Transactions, TransactionResponse{
			ID:        detail.Transaction.ID.String(),
			Type:      detail.Transaction.Type,
			Amount:    detail.Amount,
			Currency:  detail.Currency,
			Direction: detail.Direction,
			CreatedAt: detail.Transaction.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

// executeTransfer executes a transfer transaction with double-entry ledger
func (s *TransactionService) executeTransfer(ctx context.Context, userID uuid.UUID, fromAccount, toAccount *models.Account, amount decimal.Decimal, currency models.Currency) (*models.Transaction, error) {
	var transaction *models.Transaction

	err := s.db.DB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		lockedFrom, lockedTo, err := lockAccountsForUpdate(ctx, tx, s.accountRepo, fromAccount.ID, toAccount.ID)
		if err != nil {
			return err
		}

		// Ensure accounts are different
		if lockedFrom.ID == lockedTo.ID {
			return fmt.Errorf("cannot transfer to the same account")
		}

		if lockedFrom.Balance.LessThan(amount) {
			return ErrInsufficientFunds
		}

		// Create transaction record
		transaction = &models.Transaction{
			ID:     uuid.New(),
			UserID: userID,
			Type:   models.TransactionTypeTransfer,
		}

		_, err = tx.NewInsert().Model(transaction).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		// Create double-entry ledger entries
		// Debit from sender (negative amount)
		ledgerFrom := &models.Ledger{
			AccountID:     lockedFrom.ID,
			TransactionID: transaction.ID,
			Amount:        amount.Neg(),
			Currency:      currency,
		}

		// Credit to recipient (positive amount)
		ledgerTo := &models.Ledger{
			AccountID:     lockedTo.ID,
			TransactionID: transaction.ID,
			Amount:        amount,
			Currency:      currency,
		}

		_, err = tx.NewInsert().Model(ledgerFrom).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create ledger entry: %w", err)
		}

		_, err = tx.NewInsert().Model(ledgerTo).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create ledger entry: %w", err)
		}

		// Update account balances
		lockedFrom.Balance = lockedFrom.Balance.Sub(amount)
		lockedTo.Balance = lockedTo.Balance.Add(amount)
		lockedFrom.UpdatedAt = time.Now()
		lockedTo.UpdatedAt = time.Now()

		_, err = tx.NewUpdate().Model(lockedFrom).WherePK().Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to update sender balance: %w", err)
		}

		_, err = tx.NewUpdate().Model(lockedTo).WherePK().Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to update recipient balance: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// executeExchange executes an exchange transaction with double-entry ledger
func (s *TransactionService) executeExchange(ctx context.Context, userID uuid.UUID, fromAccount, toAccount *models.Account, fromAmount, toAmount decimal.Decimal, fromCurrency, toCurrency models.Currency) (*models.Transaction, error) {
	var transaction *models.Transaction

	err := s.db.DB().RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		lockedFrom, lockedTo, err := lockAccountsForUpdate(ctx, tx, s.accountRepo, fromAccount.ID, toAccount.ID)
		if err != nil {
			return err
		}

		if lockedFrom.Balance.LessThan(fromAmount) {
			return ErrInsufficientFunds
		}

		// Create transaction record
		transaction = &models.Transaction{
			ID:     uuid.New(),
			UserID: userID,
			Type:   models.TransactionTypeExchange,
		}

		_, err = tx.NewInsert().Model(transaction).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		// Create double-entry ledger entries
		// Debit from source account (negative amount)
		ledgerFrom := &models.Ledger{
			AccountID:     lockedFrom.ID,
			TransactionID: transaction.ID,
			Amount:        fromAmount.Neg(),
			Currency:      fromCurrency,
		}

		// Credit to destination account (positive amount)
		ledgerTo := &models.Ledger{
			AccountID:     lockedTo.ID,
			TransactionID: transaction.ID,
			Amount:        toAmount,
			Currency:      toCurrency,
		}

		_, err = tx.NewInsert().Model(ledgerFrom).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create ledger entry: %w", err)
		}

		_, err = tx.NewInsert().Model(ledgerTo).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create ledger entry: %w", err)
		}

		// Update account balances
		lockedFrom.Balance = lockedFrom.Balance.Sub(fromAmount)
		lockedTo.Balance = lockedTo.Balance.Add(toAmount)
		lockedFrom.UpdatedAt = time.Now()
		lockedTo.UpdatedAt = time.Now()

		_, err = tx.NewUpdate().Model(lockedFrom).WherePK().Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to update source balance: %w", err)
		}

		_, err = tx.NewUpdate().Model(lockedTo).WherePK().Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to update destination balance: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// calculateExchangeAmount calculates the exchange amount based on fixed rate
func (s *TransactionService) calculateExchangeAmount(amount decimal.Decimal, fromCurrency, toCurrency models.Currency) decimal.Decimal {
	exchangeRate := decimal.NewFromFloat(ExchangeRate)
	if fromCurrency == models.CurrencyUSD && toCurrency == models.CurrencyEUR {
		// 1 USD = 0.92 EUR
		return amount.Mul(exchangeRate)
	} else if fromCurrency == models.CurrencyEUR && toCurrency == models.CurrencyUSD {
		// 1 EUR = 1/0.92 USD
		return amount.Div(exchangeRate)
	}
	return amount
}
