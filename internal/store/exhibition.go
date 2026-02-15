package store

import (
	"sync"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/google/uuid"
)

type ExhibitionStore struct {
	mu   sync.RWMutex
	data []model.Exhibition
}

func NewExhibitionStore() *ExhibitionStore {
	return &ExhibitionStore{}
}

func (s *ExhibitionStore) GetAll() []model.Exhibition {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.Exhibition, len(s.data))
	copy(result, s.data)
	return result
}

func (s *ExhibitionStore) GetByID(id uuid.UUID) (model.Exhibition, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, e := range s.data {
		if e.ID == id {
			return e, true
		}
	}
	return model.Exhibition{}, false
}

func (s *ExhibitionStore) Create(title, description string) model.Exhibition {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := model.Exhibition{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
	}
	s.data = append(s.data, e)
	return e
}

func (s *ExhibitionStore) Update(id uuid.UUID, title, description string) (model.Exhibition, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, e := range s.data {
		if e.ID == id {
			s.data[i].Title = title
			s.data[i].Description = description
			return s.data[i], true
		}
	}
	return model.Exhibition{}, false
}

func (s *ExhibitionStore) Delete(id uuid.UUID) bool {
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
