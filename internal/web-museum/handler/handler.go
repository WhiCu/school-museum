package handler

import (
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
)

type service interface {
	GetAllNews() []model.News
	GetNewsByID(id uuid.UUID) (model.News, bool)
	GetAllExhibitions() []model.Exhibition
	GetExhibitionByID(id uuid.UUID) (model.Exhibition, bool)
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
