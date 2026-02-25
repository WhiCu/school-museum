package client

import (
	"context"
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/WhiCu/school-museum/db/storage"
	"github.com/google/uuid"
)

type Storage struct {
	News        storage.Storage[model.News]
	Exhibitions storage.Storage[model.Exhibition]
	Exhibits    storage.Storage[model.Exhibit]
	Visits      *storage.VisitStorage
	log         *slog.Logger
}

func NewStorage(news storage.Storage[model.News], exhibitions storage.Storage[model.Exhibition], exhibits storage.Storage[model.Exhibit], visits *storage.VisitStorage, log *slog.Logger) *Storage {
	return &Storage{
		News:        news,
		Exhibitions: exhibitions,
		Exhibits:    exhibits,
		Visits:      visits,
		log:         log,
	}
}

// --- News ---

func (s *Storage) CreateNews(ctx context.Context, n model.News) (model.News, error) {
	id, err := s.News.Create(ctx, n)
	if err != nil {
		s.log.Error("failed to create news", slog.String("error", err.Error()))
		return model.News{}, err
	}
	n.ID = id
	return n, nil
}

func (s *Storage) UpdateNews(ctx context.Context, n model.News) (model.News, error) {
	updated, err := s.News.Update(ctx, n)
	if err != nil {
		s.log.Error("failed to update news", slog.String("id", n.ID.String()), slog.String("error", err.Error()))
		return model.News{}, err
	}
	return updated, nil
}

func (s *Storage) DeleteNews(ctx context.Context, id uuid.UUID) error {
	if err := s.News.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete news", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}
	return nil
}

// --- Exhibitions ---

func (s *Storage) CreateExhibition(ctx context.Context, ex model.Exhibition) (model.Exhibition, error) {
	id, err := s.Exhibitions.Create(ctx, ex)
	if err != nil {
		s.log.Error("failed to create exhibition", slog.String("error", err.Error()))
		return model.Exhibition{}, err
	}
	ex.ID = id
	return ex, nil
}

func (s *Storage) UpdateExhibition(ctx context.Context, ex model.Exhibition) (model.Exhibition, error) {
	updated, err := s.Exhibitions.Update(ctx, ex)
	if err != nil {
		s.log.Error("failed to update exhibition", slog.String("id", ex.ID.String()), slog.String("error", err.Error()))
		return model.Exhibition{}, err
	}
	return updated, nil
}

func (s *Storage) DeleteExhibition(ctx context.Context, id uuid.UUID) error {
	if err := s.Exhibitions.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete exhibition", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *Storage) ExhibitionExists(ctx context.Context, id uuid.UUID) bool {
	_, err := s.Exhibitions.Read(ctx, id)
	return err == nil
}

// --- Exhibits ---

func (s *Storage) CreateExhibit(ctx context.Context, e model.Exhibit) (model.Exhibit, error) {
	id, err := s.Exhibits.Create(ctx, e)
	if err != nil {
		s.log.Error("failed to create exhibit", slog.String("error", err.Error()))
		return model.Exhibit{}, err
	}
	e.ID = id
	return e, nil
}

func (s *Storage) UpdateExhibit(ctx context.Context, e model.Exhibit) (model.Exhibit, error) {
	updated, err := s.Exhibits.Update(ctx, e)
	if err != nil {
		s.log.Error("failed to update exhibit", slog.String("id", e.ID.String()), slog.String("error", err.Error()))
		return model.Exhibit{}, err
	}
	return updated, nil
}

func (s *Storage) DeleteExhibit(ctx context.Context, id uuid.UUID) error {
	if err := s.Exhibits.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete exhibit", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}
	return nil
}

// --- Stats ---

func (s *Storage) GetStats(ctx context.Context) (model.VisitStats, error) {
	stats, err := s.Visits.Stats(ctx)
	if err != nil {
		s.log.Error("failed to get stats", slog.String("error", err.Error()))
		return model.VisitStats{}, err
	}
	return stats, nil
}
