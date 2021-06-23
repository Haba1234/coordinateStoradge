package storage

import (
	"sync"

	"github.com/Haba1234/coordinateStoradge/internal/app"
)

type storage struct {
	mu     sync.RWMutex
	points map[uint64]bool
}

func NewStorage() app.Storage {
	return &storage{
		mu:     sync.RWMutex{},
		points: make(map[uint64]bool),
	}
}

func (s *storage) AddPoint(z uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.points[z] = true
}

func (s *storage) ReadPoint(z uint64) (bool, bool) {
	val, ok := s.points[z]
	return val, ok
}

func (s *storage) RLock() {
	s.mu.RLock()
}

func (s *storage) RUnlock() {
	s.mu.RUnlock()
}

func (s *storage) Len() int {
	return len(s.points)
}
