package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts,alias:a"`

	ID        uuid.UUID       `bun:"id,pk,type:uuid"`
	UserID    uuid.UUID       `bun:"user_id,type:uuid,notnull"`
	Currency  Currency        `bun:"currency,type:currency,notnull"`
	Balance   decimal.Decimal `bun:"balance,type:numeric(19,2),notnull,default:0.00"`
	CreatedAt time.Time       `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time       `bun:"updated_at,nullzero,notnull,default:current_timestamp"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}
