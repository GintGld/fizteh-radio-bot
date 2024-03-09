package session

import "sync"

// MapCache is a simple
// implementation for
// controller.Session
// interface.
type MapCache[T comparable] struct {
	// Users routing
	paths map[int64]T

	// Session cache
	users map[int64]map[T]T

	// Mutex per user
	mutexes map[int64]sync.Mutex
}

func New[T comparable]() *MapCache[T] {
	return &MapCache[T]{
		paths:   make(map[int64]T),
		users:   make(map[int64]map[T]T),
		mutexes: make(map[int64]sync.Mutex),
	}
}

func (s *MapCache[T]) Status(id int64) T {
	if res, ok := s.paths[id]; ok {
		return res
	} else {
		return *new(T)
	}
}

func (s *MapCache[T]) Redirect(id int64, path T) {
	s.paths[id] = path
}
