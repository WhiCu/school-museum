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
	ImageURL      string    `json:"image_url" bun:"image_url"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}

type Exhibition struct {
	ID          uuid.UUID `json:"id" bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Title       string    `json:"title" bun:"title"`
	Description string    `json:"description" bun:"description"`
	Exhibit     []Exhibit `bun:"rel:has-many,join:id=exhibition_id"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete,nullzero"`
}
