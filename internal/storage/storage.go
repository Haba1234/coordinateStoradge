package storage

import (
	"sync"
)

type Storage struct {
	mu     sync.RWMutex
	points map[uint64]bool
}

func NewStorage() *Storage {
	return &Storage{
		mu:     sync.RWMutex{},
		points: make(map[uint64]bool),
	}
}

func (s *Storage) AddPoint(z uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.points[z] = true
}

func (s *Storage) ReadPoint(z uint64) (bool, bool) {
	val, ok := s.points[z]
	return val, ok
}

func (s *Storage) RLock() {
	s.mu.RLock()
}

func (s *Storage) RUnlock() {
	s.mu.RUnlock()
}

func (s *Storage) Len() int {
	return len(s.points)
}
