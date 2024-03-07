package session

import "sync"

// MapCache is a simple
// implementation for
// controller.Session
// interface.
type MapCache struct {
	// Users routing
	paths map[int64]string

	// Session cache
	users map[int64]map[string]string

	// Mutex per user
	mutexes map[int64]sync.Mutex
}

func New() *MapCache {
	return &MapCache{
		paths:   make(map[int64]string),
		users:   make(map[int64]map[string]string),
		mutexes: make(map[int64]sync.Mutex),
	}
}

func (s *MapCache) Status(id int64) string {
	if res, ok := s.paths[id]; ok {
		return res
	} else {
		return ""
	}
}

func (s *MapCache) Redirect(id int64, path string) {
	s.paths[id] = path
}

func (s *MapCache) Get(id int64, key string) string {
	if user, ok := s.users[id]; ok {
		if res, exists := user[key]; exists {
			return res
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func (s *MapCache) Set(id int64, key string, value string) {
	if _, ok := s.users[id]; ok {
		s.users[id][key] = value
	} else {
		s.users[id] = map[string]string{key: value}
	}
}

func (s *MapCache) Del(id int64, key string) {
	if _, ok := s.users[id]; ok {
		delete(s.users[id], key)
	}
}
