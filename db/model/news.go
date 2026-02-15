package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type News struct {
	bun.BaseModel `bun:"table:news,alias:n"`

	ID       uuid.UUID `json:"id" bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Title    string    `json:"title" bun:"title,type:text"`
	Content  string    `json:"content" bun:"content,type:text"`
	ImageURL string    `json:"image_url" bun:"image_url"`

	CreatedAt time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `json:"-" bun:",soft_delete,nullzero"`
}
