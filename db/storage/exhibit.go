package storage

import (
	"context"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ExhibitStorage struct {
	db *bun.DB
}

var _ Storage[model.Exhibit] = (*ExhibitStorage)(nil)

func NewExhibitStorage(db *bun.DB) *ExhibitStorage {
	return &ExhibitStorage{
		db: db,
	}
}

func (s *ExhibitStorage) Read(ctx context.Context, id uuid.UUID) (model.Exhibit, error) {
	var e model.Exhibit
	err := s.db.NewSelect().Model(&e).Where("id = ?", id).Scan(ctx, &e)
	if err != nil {
		return model.Exhibit{}, err
	}
	return e, nil
}

func (s *ExhibitStorage) Create(ctx context.Context, e model.Exhibit) (uuid.UUID, error) {
	err := s.db.NewInsert().Model(&e).Returning("id").Scan(ctx, &e.ID)
	if err != nil {
		return uuid.Nil, err
	}
	return e.ID, nil
}

func (s *ExhibitStorage) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*model.Exhibit)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *ExhibitStorage) Update(ctx context.Context, e model.Exhibit) (model.Exhibit, error) {
	_, err := s.db.NewUpdate().
		Model(&e).
		WherePK().
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return model.Exhibit{}, err
	}
	return e, nil
}

func (s *ExhibitStorage) List(ctx context.Context) (exhibits []model.Exhibit, err error) {
	err = s.db.NewSelect().Model(&exhibits).Scan(ctx)
	return exhibits, err
}

func (s *ExhibitStorage) First(ctx context.Context, f func(model.Exhibit) bool) (model.Exhibit, error) {
	var exhibits []model.Exhibit
	err := s.db.NewSelect().Model(&exhibits).Scan(ctx)
	if err != nil {
		return model.Exhibit{}, err
	}
	for _, e := range exhibits {
		if f(e) {
			return e, nil
		}
	}
	return model.Exhibit{}, ErrNotFound
}
