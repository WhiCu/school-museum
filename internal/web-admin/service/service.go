package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
)

var ErrExhibitionNotFound = errors.New("exhibition not found")

type Storage interface {
	CreateNews(ctx context.Context, n model.News) (model.News, error)
	UpdateNews(ctx context.Context, n model.News) (model.News, error)
	DeleteNews(ctx context.Context, id uuid.UUID) error

	CreateExhibition(ctx context.Context, ex model.Exhibition) (model.Exhibition, error)
	UpdateExhibition(ctx context.Context, ex model.Exhibition) (model.Exhibition, error)
	DeleteExhibition(ctx context.Context, id uuid.UUID) error
	ExhibitionExists(ctx context.Context, id uuid.UUID) bool
	SetExhibitionPreview(ctx context.Context, exhibitionID uuid.UUID, exhibitID *uuid.UUID) (model.Exhibition, error)

	CreateExhibit(ctx context.Context, e model.Exhibit) (model.Exhibit, error)
	UpdateExhibit(ctx context.Context, e model.Exhibit) (model.Exhibit, error)
	DeleteExhibit(ctx context.Context, id uuid.UUID) error

	GetStats(ctx context.Context) (model.VisitStats, error)
}

type Service struct {
	storage Storage
	log     *slog.Logger
}

func NewService(storage Storage, log *slog.Logger) *Service {
	return &Service{
		storage: storage,
		log:     log,
	}
}

// --- News ---

func (s *Service) CreateNews(ctx context.Context, n model.News) (model.News, error) {
	return s.storage.CreateNews(ctx, n)
}

func (s *Service) UpdateNews(ctx context.Context, n model.News) (model.News, error) {
	return s.storage.UpdateNews(ctx, n)
}

func (s *Service) DeleteNews(ctx context.Context, id uuid.UUID) error {
	return s.storage.DeleteNews(ctx, id)
}

// --- Exhibitions ---

func (s *Service) CreateExhibition(ctx context.Context, ex model.Exhibition) (model.Exhibition, error) {
	return s.storage.CreateExhibition(ctx, ex)
}

func (s *Service) UpdateExhibition(ctx context.Context, ex model.Exhibition) (model.Exhibition, error) {
	return s.storage.UpdateExhibition(ctx, ex)
}

func (s *Service) DeleteExhibition(ctx context.Context, id uuid.UUID) error {
	return s.storage.DeleteExhibition(ctx, id)
}

func (s *Service) SetExhibitionPreview(ctx context.Context, exhibitionID uuid.UUID, exhibitID *uuid.UUID) (model.Exhibition, error) {
	return s.storage.SetExhibitionPreview(ctx, exhibitionID, exhibitID)
}

// --- Exhibits ---

func (s *Service) CreateExhibit(ctx context.Context, e model.Exhibit) (model.Exhibit, error) {
	if !s.storage.ExhibitionExists(ctx, e.ExhibitionID) {
		return model.Exhibit{}, fmt.Errorf("exhibition %s: %w", e.ExhibitionID, ErrExhibitionNotFound)
	}
	return s.storage.CreateExhibit(ctx, e)
}

func (s *Service) UpdateExhibit(ctx context.Context, e model.Exhibit) (model.Exhibit, error) {
	return s.storage.UpdateExhibit(ctx, e)
}

func (s *Service) DeleteExhibit(ctx context.Context, id uuid.UUID) error {
	return s.storage.DeleteExhibit(ctx, id)
}

// --- Stats ---

func (s *Service) GetStats(ctx context.Context) (model.VisitStats, error) {
	return s.storage.GetStats(ctx)
}
