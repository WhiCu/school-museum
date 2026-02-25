package storage

import (
	"context"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ExhibitionStorage struct {
	db *bun.DB
}

var _ Storage[model.Exhibition] = (*ExhibitionStorage)(nil)

func NewExhibitionStorage(db *bun.DB) *ExhibitionStorage {
	return &ExhibitionStorage{
		db: db,
	}
}

func (s *ExhibitionStorage) Read(ctx context.Context, id uuid.UUID) (model.Exhibition, error) {
	var ex model.Exhibition
	err := s.db.NewSelect().
		Model(&ex).
		Relation("Exhibits").
		Where("ex.id = ?", id).
		Scan(ctx)
	if err != nil {
		return model.Exhibition{}, err
	}
	return ex, nil
}

func (s *ExhibitionStorage) Create(ctx context.Context, ex model.Exhibition) (uuid.UUID, error) {
	_, err := s.db.NewInsert().Model(&ex).Returning("id").Exec(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return ex.ID, nil
}

func (s *ExhibitionStorage) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*model.Exhibition)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *ExhibitionStorage) Update(ctx context.Context, ex model.Exhibition) (model.Exhibition, error) {
	_, err := s.db.NewUpdate().
		Model(&ex).
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return model.Exhibition{}, err
	}
	return ex, nil
}

func (s *ExhibitionStorage) List(ctx context.Context) (exhibitions []model.Exhibition, err error) {
	err = s.db.NewSelect().
		Model(&exhibitions).
		Relation("Exhibits").
		Scan(ctx)
	return exhibitions, err
}

func (s *ExhibitionStorage) First(ctx context.Context, f func(model.Exhibition) bool) (model.Exhibition, error) {
	var exhibitions []model.Exhibition
	err := s.db.NewSelect().Model(&exhibitions).Scan(ctx)
	if err != nil {
		return model.Exhibition{}, err
	}
	for _, ex := range exhibitions {
		if f(ex) {
			return ex, nil
		}
	}
	return model.Exhibition{}, ErrNotFound
}

// SetPreview sets the preview exhibit for an exhibition.
// Pass nil to clear the preview.
func (s *ExhibitionStorage) SetPreview(ctx context.Context, exhibitionID uuid.UUID, exhibitID *uuid.UUID) error {
	_, err := s.db.NewUpdate().
		Model((*model.Exhibition)(nil)).
		Set("preview_exhibit_id = ?", exhibitID).
		Where("id = ?", exhibitionID).
		Exec(ctx)
	return err
}
