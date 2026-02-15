package store

import (
	"sync"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
)

type ExhibitStore struct {
	mu   sync.RWMutex
	data []model.Exhibit
}

func NewExhibitStore() *ExhibitStore {
	return &ExhibitStore{}
}

func (s *ExhibitStore) GetByExhibitionID(exhibitionID uuid.UUID) []model.Exhibit {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []model.Exhibit
	for _, e := range s.data {
		if e.ExhibitionID == exhibitionID {
			result = append(result, e)
		}
	}
	return result
}

func (s *ExhibitStore) GetByID(id uuid.UUID) (model.Exhibit, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.data {
		if e.ID == id {
			return e, true
		}
	}
	return model.Exhibit{}, false
}

func (s *ExhibitStore) Create(exhibitionID uuid.UUID, title, description, imageURL string) model.Exhibit {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := model.Exhibit{
		ID:           uuid.New(),
		ExhibitionID: exhibitionID,
		Title:        title,
		Description:  description,
		ImageURL:     imageURL,
	}
	s.data = append(s.data, e)
	return e
}

func (s *ExhibitStore) Update(id uuid.UUID, title, description, imageURL string) (model.Exhibit, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.data {
		if e.ID == id {
			s.data[i].Title = title
			s.data[i].Description = description
			s.data[i].ImageURL = imageURL
			return s.data[i], true
		}
	}
	return model.Exhibit{}, false
}

func (s *ExhibitStore) Delete(id uuid.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.data {
		if e.ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return true
		}
	}
	return false
}

func (s *ExhibitStore) DeleteByExhibitionID(exhibitionID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	filtered := s.data[:0]
	for _, e := range s.data {
		if e.ExhibitionID != exhibitionID {
			filtered = append(filtered, e)
		}
	}
	s.data = filtered
}
