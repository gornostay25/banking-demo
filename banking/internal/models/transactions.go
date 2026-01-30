package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`

	ID        uuid.UUID     `bun:"id,pk,type:uuid"`
	UserID    uuid.UUID     `bun:"user_id,type:uuid,notnull"`
	Type      TransactionType `bun:"type,type:transaction_type,notnull"`
	CreatedAt time.Time     `bun:"created_at,nullzero,notnull,default:current_timestamp"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}
