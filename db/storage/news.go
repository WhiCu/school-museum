package storage

import (
	"context"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type NewsStorage struct {
	db *bun.DB
}

var _ Storage[model.News] = (*NewStorage)(nil)

func NewNewsStorage(db *bun.DB) *NewsStorage {
	return &NewsStorage{
		db: db,
	}
}

func (s *NewsStorage) Read(ctx context.Context, id uuid.UUID) (model.News, error) {
	var n model.News
	err := s.db.NewSelect().Model(&n).Where("id = ?", id).Scan(ctx, &n)
	if err != nil {
		return model.News{}, err
	}
	return n, nil
}

func (s *NewsStorage) Create(ctx context.Context, n model.News) (uuid.UUID, error) {
	err := s.db.NewInsert().Model(&n).Returning("id").Scan(ctx, &n.ID)
	if err != nil {
		return uuid.Nil, err
	}
	return n.ID, nil
}

func (s *NewsStorage) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*model.News)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *NewsStorage) Update(ctx context.Context, n model.News) (model.News, error) {
	err := s.db.NewUpdate().Model(&n).Where("id = ?", n.ID).Scan(ctx, &n)
	if err != nil {
		return model.News{}, err
	}
	return n, nil
}

func (s *NewsStorage) List(ctx context.Context) (news []model.News, err error) {
	err = s.db.NewSelect().Model(&news).Scan(ctx)
	return news, err
}

func (s *NewsStorage) First(ctx context.Context, f func(model.News) bool) (model.News, error) {
	var news []model.News
	err := s.db.NewSelect().Model(&news).Scan(ctx)
	if err != nil {
		return model.News{}, err
	}
	for _, n := range news {
		if f(n) {
			return n, nil
		}
	}
	return model.News{}, ErrNotFound
}
