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
	// return s.store.Exhibitions.GetByID(id)
	return model.Exhibition{}, false
}

func (s *Storage) GetExhibitsByExhibitionID(exhibitionID uuid.UUID) []model.Exhibit {
	// return s.store.Exhibits.GetByExhibitionID(exhibitionID)
	return []model.Exhibit{}
}
