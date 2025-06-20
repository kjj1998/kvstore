package main

import (
	"sync"
	"time"
)

type Store struct {
	kvStore map[string]Value
	mutex   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		kvStore: make(map[string]Value),
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, exists := s.kvStore[key]

	if value.expiry.After(time.Now()) {
		s.Delete(key)
		return "", false
	}

	return value.val, exists
}

func (s *Store) Set(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.kvStore[key] = Value{val: value, expiry: time.Now().Add(2 * time.Minute)}
}

func (s *Store) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.kvStore, key)
}
