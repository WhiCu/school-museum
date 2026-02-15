package service

import (
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/WhiCu/school-museum/db/storage"
	"github.com/google/uuid"
)

type Storage interface {
	CreateNews(title, content string) model.News
	DeleteNews(id uuid.UUID) bool

	CreateExhibition(title, description string) model.Exhibition
	UpdateExhibition(id uuid.UUID, title, description string) (model.Exhibition, bool)
	DeleteExhibition(id uuid.UUID) bool
	ExhibitionExists(id uuid.UUID) bool

	CreateExhibit(exhibitionID uuid.UUID, title, description, imageURL string) model.Exhibit
	UpdateExhibit(id uuid.UUID, title, description, imageURL string) (model.Exhibit, bool)
	DeleteExhibit(id uuid.UUID) bool
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

func (s *Service) CreateNews(title, content string) model.News {
	return s.storage.CreateNews(title, content)
}

func (s *Service) DeleteNews(id uuid.UUID) error {
	if !s.storage.DeleteNews(id) {
		return storage.ErrNotFound
	}
	return nil
}

// --- Exhibitions ---

func (s *Service) CreateExhibition(title, description string) model.Exhibition {
	return s.storage.CreateExhibition(title, description)
}

func (s *Service) UpdateExhibition(id uuid.UUID, title, description string) (model.Exhibition, error) {
	ex, ok := s.storage.UpdateExhibition(id, title, description)
	if !ok {
		return model.Exhibition{}, storage.ErrNotFound
	}
	return ex, nil
}

func (s *Service) DeleteExhibition(id uuid.UUID) error {
	if !s.storage.DeleteExhibition(id) {
		return storage.ErrNotFound
	}
	return nil
}

// --- Exhibits ---

func (s *Service) CreateExhibit(exhibitionID uuid.UUID, title, description, imageURL string) (model.Exhibit, error) {
	if !s.storage.ExhibitionExists(exhibitionID) {
		return model.Exhibit{}, storage.ErrNotFound
	}
	return s.storage.CreateExhibit(exhibitionID, title, description, imageURL), nil
}

func (s *Service) UpdateExhibit(id uuid.UUID, title, description, imageURL string) (model.Exhibit, error) {
	ex, ok := s.storage.UpdateExhibit(id, title, description, imageURL)
	if !ok {
		return model.Exhibit{}, storage.ErrNotFound
	}
	return ex, nil
}

func (s *Service) DeleteExhibit(id uuid.UUID) error {
	if !s.storage.DeleteExhibit(id) {
		return storage.ErrNotFound
	}
	return nil
}
