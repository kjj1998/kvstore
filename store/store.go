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
	value, exists := s.kvStore[key]

	expired := exists && !value.Expiry.IsZero() && value.Expiry.Before(time.Now())
	s.mutex.RUnlock()

	if expired {
		s.Delete(key)
		return "", false
	}

	return value.Value, exists
}

func (s *Store) Set(key, value string, timeToLive int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var expiry = time.Time{}

	if timeToLive == 0 {
		expiry = time.Time{}
	} else {
		expiry = time.Now().Add(time.Duration(timeToLive) * time.Second)
	}

	s.kvStore[key] = models.Value{
		Value:  value,
		Expiry: expiry,
	}
}

func (s *Store) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.kvStore, key)
}
