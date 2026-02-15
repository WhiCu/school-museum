package service

import (
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
)

type Storage interface {
	GetAllNews() []model.News
	GetNewsByID(id uuid.UUID) (model.News, bool)
	GetAllExhibitions() []model.Exhibition
	GetExhibitionByID(id uuid.UUID) (model.Exhibition, bool)
	GetExhibitsByExhibitionID(exhibitionID uuid.UUID) []model.Exhibit
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

func (s *Service) GetAllNews() []model.News {
	return s.storage.GetAllNews()
}

func (s *Service) GetNewsByID(id uuid.UUID) (model.News, bool) {
	return s.storage.GetNewsByID(id)
}

func (s *Service) GetAllExhibitions() []model.Exhibition {
	return s.storage.GetAllExhibitions()
}

func (s *Service) GetExhibitionByID(id uuid.UUID) (model.Exhibition, bool) {
	ex, ok := s.storage.GetExhibitionByID(id)
	if !ok {
		return model.Exhibition{}, false
	}
	exhibits := s.storage.GetExhibitsByExhibitionID(id)
	if exhibits == nil {
		exhibits = []model.Exhibit{}
	}
	return ex, true
}
