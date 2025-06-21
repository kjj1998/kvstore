package store

import (
	"sync"
	"time"

	"github.com/kjj1998/kvstore/models"
)

type Store struct {
	kvStore map[string]models.Value
	mutex   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		kvStore: make(map[string]models.Value),
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, exists := s.kvStore[key]

	if value.Expiry.After(time.Now()) {
		s.Delete(key)
		return "", false
	}

	return value.Value, exists
}

func (s *Store) Set(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.kvStore[key] = models.Value{Value: value, Expiry: time.Now().Add(-2 * time.Minute)}
}

func (s *Store) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.kvStore, key)
}
