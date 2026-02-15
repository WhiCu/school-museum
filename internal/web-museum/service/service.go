package service

import (
	"context"
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
)

type Storage interface {
	GetAllNews(ctx context.Context) ([]model.News, error)
	GetNewsByID(ctx context.Context, id uuid.UUID) (model.News, error)
	GetAllExhibitions(ctx context.Context) ([]model.Exhibition, error)
	GetExhibitionByID(ctx context.Context, id uuid.UUID) (model.Exhibition, error)
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

func (s *Service) GetAllNews(ctx context.Context) ([]model.News, error) {
	return s.storage.GetAllNews(ctx)
}

func (s *Service) GetNewsByID(ctx context.Context, id uuid.UUID) (model.News, error) {
	return s.storage.GetNewsByID(ctx, id)
}

func (s *Service) GetAllExhibitions(ctx context.Context) ([]model.Exhibition, error) {
	return s.storage.GetAllExhibitions(ctx)
}

func (s *Service) GetExhibitionByID(ctx context.Context, id uuid.UUID) (model.Exhibition, error) {
	return s.storage.GetExhibitionByID(ctx, id)
}
