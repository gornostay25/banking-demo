package models

// Currency represents the currency enum type
type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

// TransactionType represents the transaction type enum
type TransactionType string

const (
	TransactionTypeTransfer TransactionType = "transfer"
	TransactionTypeExchange TransactionType = "exchange"
)
