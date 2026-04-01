package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Exhibit struct {
	bun.BaseModel `bun:"table:exhibits,alias:e"`
	ID            uuid.UUID `json:"id" bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	ExhibitionID  uuid.UUID `json:"exhibition_id" bun:"exhibition_id,type:uuid"`
	Title         string    `json:"title" bun:"title"`
	Description   string    `json:"description" bun:"description"`
	ImageURLs     []string  `json:"image_urls" bun:"image_urls,type:text[],default:'{}'"`

	CreatedAt time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `json:"-" bun:",soft_delete,nullzero"`
}

// ImageURL returns the first image URL for backward compatibility.
func (e Exhibit) ImageURL() string {
	if len(e.ImageURLs) > 0 {
		return e.ImageURLs[0]
	}
	return ""
}

type Exhibition struct {
	bun.BaseModel `bun:"table:exhibitions,alias:ex"`

	ID               uuid.UUID  `json:"id" bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Title            string     `json:"title" bun:"title"`
	Description      string     `json:"description" bun:"description"`
	PreviewExhibitID *uuid.UUID `json:"preview_exhibit_id,omitempty" bun:"preview_exhibit_id,type:uuid,nullzero"`
	Exhibits         []Exhibit  `json:"exhibits,omitempty" bun:"rel:has-many,join:id=exhibition_id"`

	CreatedAt time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `json:"-" bun:",soft_delete,nullzero"`
}
