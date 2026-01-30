package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           uuid.UUID `bun:"id,pk,type:uuid"`
	Email        string    `bun:"email,unique,notnull"`
	PasswordHash string    `bun:"password_hash,notnull"`
	CreatedAt    time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}
