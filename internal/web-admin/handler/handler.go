package handler

import (
	"context"
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
)

type service interface {
	CreateNews(ctx context.Context, n model.News) (model.News, error)
	UpdateNews(ctx context.Context, n model.News) (model.News, error)
	DeleteNews(ctx context.Context, id uuid.UUID) error

	CreateExhibition(ctx context.Context, ex model.Exhibition) (model.Exhibition, error)
	UpdateExhibition(ctx context.Context, ex model.Exhibition) (model.Exhibition, error)
	DeleteExhibition(ctx context.Context, id uuid.UUID) error

	CreateExhibit(ctx context.Context, e model.Exhibit) (model.Exhibit, error)
	UpdateExhibit(ctx context.Context, e model.Exhibit) (model.Exhibit, error)
	DeleteExhibit(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	service service
	log     *slog.Logger
}

func NewHandler(service service, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}
