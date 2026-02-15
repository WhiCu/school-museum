package handler

import (
	"context"
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
)

type service interface {
	GetAllNews(ctx context.Context) ([]model.News, error)
	GetNewsByID(ctx context.Context, id uuid.UUID) (model.News, error)
	GetAllExhibitions(ctx context.Context) ([]model.Exhibition, error)
	GetExhibitionByID(ctx context.Context, id uuid.UUID) (model.Exhibition, error)
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
