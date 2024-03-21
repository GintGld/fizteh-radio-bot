package storage

type Storage[T any] struct {
	vals map[int64]T
}

func New[T any]() Storage[T] {
	return Storage[T]{
		vals: make(map[int64]T),
	}
}

func (s *Storage[T]) Get(id int64) T {
	if t, ok := s.vals[id]; ok {
		return t
	} else {
		return *new(T)
	}
}

func (s *Storage[T]) Set(id int64, t T) {
	s.vals[id] = t
}

func (s *Storage[T]) Del(id int64) {
	delete(s.vals, id)
}
