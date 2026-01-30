package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type Ledger struct {
	bun.BaseModel `bun:"table:ledger,alias:l"`

	ID            int64           `bun:"id,pk,autoincrement"`
	AccountID     uuid.UUID       `bun:"account_id,type:uuid,notnull"`
	TransactionID uuid.UUID       `bun:"transaction_id,type:uuid,notnull"`
	Amount        decimal.Decimal  `bun:"amount,type:numeric(19,2),notnull"`
	Currency      Currency        `bun:"currency,type:currency,notnull"`
	CreatedAt     time.Time       `bun:"created_at,nullzero,notnull,default:current_timestamp"`

	// Relations
	Account     *Account     `bun:"rel:belongs-to,join:account_id=id"`
	Transaction *Transaction `bun:"rel:belongs-to,join:transaction_id=id"`
}
