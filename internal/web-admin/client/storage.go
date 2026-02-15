package client

import (
	"log/slog"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/WhiCu/school-museum/db/storage"
	"github.com/google/uuid"
)

type Storage struct {
	News        storage.Storage[model.News]
	Exhibitions storage.Storage[model.Exhibition]
	Exhibits    storage.Storage[model.Exhibit]
	log         *slog.Logger
}

func NewStorage(news storage.Storage[model.News], exhibitions storage.Storage[model.Exhibition], exhibits storage.Storage[model.Exhibit], log *slog.Logger) *Storage {
	return &Storage{
		News:        news,
		Exhibitions: exhibitions,
		Exhibits:    exhibits,
		log:         log}
}

// --- News ---

func (s *Storage) CreateNews(title, content string) model.News {
	// return s.store.News.Create(title, content)
	return model.News{}
}

func (s *Storage) DeleteNews(id uuid.UUID) bool {
	// return s.store.News.Delete(id)
	return false
}

// --- Exhibitions ---

func (s *Storage) CreateExhibition(title, description string) model.Exhibition {
	// return s.store.Exhibitions.Create(title, description)
	return model.Exhibition{}
}

func (s *Storage) UpdateExhibition(id uuid.UUID, title, description string) (model.Exhibition, bool) {
	// return s.store.Exhibitions.Update(id, title, description)
	return model.Exhibition{}, false
}

func (s *Storage) DeleteExhibition(id uuid.UUID) bool {
	return false
}

func (s *Storage) ExhibitionExists(id uuid.UUID) bool {
	return false
}

// --- Exhibits ---

func (s *Storage) CreateExhibit(exhibitionID uuid.UUID, title, description, imageURL string) model.Exhibit {
	// return s.store.Exhibits.Create(exhibitionID, title, description, imageURL)
	return model.Exhibit{}
}

func (s *Storage) UpdateExhibit(id uuid.UUID, title, description, imageURL string) (model.Exhibit, bool) {
	// return s.store.Exhibits.Update(id, title, description, imageURL)
	return model.Exhibit{}, false
}

func (s *Storage) DeleteExhibit(id uuid.UUID) bool {
	// return s.store.Exhibits.Delete(id)
	return false

}
