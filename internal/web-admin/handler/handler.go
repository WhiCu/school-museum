package handler

import (
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
)

type service interface {
	CreateNews(title, content string) model.News
	DeleteNews(id uuid.UUID) error

	CreateExhibition(title, description string) model.Exhibition
	UpdateExhibition(id uuid.UUID, title, description string) (model.Exhibition, error)
	DeleteExhibition(id uuid.UUID) error

	CreateExhibit(exhibitionID uuid.UUID, title, description, imageURL string) (model.Exhibit, error)
	UpdateExhibit(id uuid.UUID, title, description, imageURL string) (model.Exhibit, error)
	DeleteExhibit(id uuid.UUID) error
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
