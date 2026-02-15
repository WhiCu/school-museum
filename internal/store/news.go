package store

import (
	"sync"
	"time"

	"github.com/WhiCu/school-museum/internal/model"
	"github.com/google/uuid"
)

type NewsStore struct {
	mu   sync.RWMutex
	data []model.News
}

func NewNewsStore() *NewsStore {
	return &NewsStore{}
}

func (s *NewsStore) GetAll() []model.News {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]model.News, len(s.data))
	copy(result, s.data)
	return result
}

func (s *NewsStore) GetByID(id uuid.UUID) (model.News, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, n := range s.data {
		if n.ID == id {
			return n, true
		}
	}
	return model.News{}, false
}

func (s *NewsStore) Create(title, content string) model.News {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := model.News{
		ID:        uuid.New(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}
	s.data = append(s.data, n)
	return n
}

func (s *NewsStore) Delete(id uuid.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, n := range s.data {
		if n.ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return true
		}
	}
	return false
}
