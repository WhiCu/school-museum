package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type News struct {
	bun.BaseModel `bun:"table:news,alias:n"`

	ID      uuid.UUID `json:"id" bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Title   []byte    `json:"title" bun:"title,type:bytea"`
	Content []byte    `json:"content" bun:"content,type:bytea"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}
