package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type News struct {
	bun.BaseModel `bun:"table:news,alias:n"`

	ID        uuid.UUID `json:"id" bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Title     string    `json:"title" bun:"title,type:text"`
	Content   string    `json:"content" bun:"content,type:text"`
	ImageURLs []string  `json:"image_urls" bun:"image_urls,type:text[],default:'{}'"`

	CreatedAt time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `json:"-" bun:",soft_delete,nullzero"`
}

// ImageURL returns the first image URL for backward compatibility.
func (n News) ImageURL() string {
	if len(n.ImageURLs) > 0 {
		return n.ImageURLs[0]
	}
	return ""
}
