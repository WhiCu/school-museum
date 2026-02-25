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

func (s *Storage) GetAllNews(ctx context.Context) ([]model.News, error) {
	news, err := s.News.List(ctx)
	if err != nil {
		s.log.Error("failed to get all news", slog.String("error", err.Error()))
		return nil, err
	}
	return news, nil
}

func (s *Storage) GetNewsByID(ctx context.Context, id uuid.UUID) (model.News, error) {
	n, err := s.News.Read(ctx, id)
	if err != nil {
		s.log.Error("failed to get news by id", slog.String("id", id.String()), slog.String("error", err.Error()))
		return model.News{}, err
	}
	return n, nil
}

// --- Exhibitions ---

func (s *Storage) GetAllExhibitions(ctx context.Context) ([]model.Exhibition, error) {
	exhibitions, err := s.Exhibitions.List(ctx)
	if err != nil {
		s.log.Error("failed to get all exhibitions", slog.String("error", err.Error()))
		return nil, err
	}
	return exhibitions, nil
}

func (s *Storage) GetExhibitionByID(ctx context.Context, id uuid.UUID) (model.Exhibition, error) {
	ex, err := s.Exhibitions.Read(ctx, id)
	if err != nil {
		s.log.Error("failed to get exhibition by id", slog.String("id", id.String()), slog.String("error", err.Error()))
		return model.Exhibition{}, err
	}
	return ex, nil
}

// --- Visits ---

func (s *Storage) RecordVisit(ctx context.Context, v model.Visitor) error {
	if err := s.Visits.Record(ctx, v); err != nil {
		s.log.Error("failed to record visit", slog.String("error", err.Error()))
		return err
	}
	return nil
}
