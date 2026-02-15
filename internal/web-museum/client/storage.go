package client

import (
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/WhiCu/school-museum/internal/store"
	"github.com/google/uuid"
)

type Storage struct {
	store *store.Store
	log   *slog.Logger
}

func NewStorage(s *store.Store, log *slog.Logger) *Storage {
	return &Storage{store: s, log: log}
}

func (s *Storage) GetAllNews() []model.News {
	// return s.store.News.GetAll()
	return []model.News{}
}

func (s *Storage) GetNewsByID(id uuid.UUID) (model.News, bool) {
	// return s.store.News.GetByID(id)
	return model.News{}, false
}

func (s *Storage) GetAllExhibitions() []model.Exhibition {
	// return s.store.Exhibitions.GetAll()
	return []model.Exhibition{}
}

func (s *Storage) GetExhibitionByID(id uuid.UUID) (model.Exhibition, bool) {
	return s.store.Exhibitions.GetByID(id)
}

func (s *Storage) GetExhibitsByExhibitionID(exhibitionID uuid.UUID) []model.Exhibit {
	// return s.store.Exhibits.GetByExhibitionID(exhibitionID)
	return []model.Exhibit{}
}
