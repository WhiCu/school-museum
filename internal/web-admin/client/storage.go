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

// --- News ---

func (s *Storage) CreateNews(title, content string) model.News {
	// return s.store.News.Create(title, content)
	return model.News{}
}

func (s *Storage) DeleteNews(id uuid.UUID) bool {
	return s.store.News.Delete(id)
}

// --- Exhibitions ---

func (s *Storage) CreateExhibition(title, description string) model.Exhibition {
	return s.store.Exhibitions.Create(title, description)
}

func (s *Storage) UpdateExhibition(id uuid.UUID, title, description string) (model.Exhibition, bool) {
	return s.store.Exhibitions.Update(id, title, description)
}

func (s *Storage) DeleteExhibition(id uuid.UUID) bool {
	ok := s.store.Exhibitions.Delete(id)
	if ok {
		s.store.Exhibits.DeleteByExhibitionID(id)
	}
	return ok
}

func (s *Storage) ExhibitionExists(id uuid.UUID) bool {
	_, ok := s.store.Exhibitions.GetByID(id)
	return ok
}

// --- Exhibits ---

func (s *Storage) CreateExhibit(exhibitionID uuid.UUID, title, description, imageURL string) model.Exhibit {
	return s.store.Exhibits.Create(exhibitionID, title, description, imageURL)
}

func (s *Storage) UpdateExhibit(id uuid.UUID, title, description, imageURL string) (model.Exhibit, bool) {
	return s.store.Exhibits.Update(id, title, description, imageURL)
}

func (s *Storage) DeleteExhibit(id uuid.UUID) bool {
	return s.store.Exhibits.Delete(id)
}
