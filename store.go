package main

import "sync"

type Store struct {
	kvStore map[string]string
	mutex   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		kvStore: make(map[string]string),
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	value, exists := s.kvStore[key]
	return value, exists
}

func (s *Store) Set(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.kvStore[key] = value
}

func (s *Store) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.kvStore, key)
}
