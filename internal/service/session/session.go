package session

// MapCache is a simple
// implementation for
// controller.Session
// interface.
type MapCache[T comparable] struct {
	// Users routing
	paths map[int64]T
}

func New[T comparable]() *MapCache[T] {
	return &MapCache[T]{
		paths: make(map[int64]T),
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
