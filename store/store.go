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

func (s *Store) Get(key string) string {
	s.mutex.RLock()
	value, exists := s.kvStore[key]

	expired := exists && !value.Expiry.IsZero() && value.Expiry.Before(time.Now())
	s.mutex.RUnlock()

	if expired {
		s.Delete(key)
		return "EXPIRED"
	}

	if !exists {
		return "NULL"
	}

	return value.Value
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

func (s *Store) CleanUpExpiredKeys() {
	now := time.Now()
	var expiredKeys []string

	s.mutex.RLock()
	for key, value := range s.kvStore {
		if !value.Expiry.IsZero() && value.Expiry.Before(now) {
			expiredKeys = append(expiredKeys, key)
		}
	}
	s.mutex.RUnlock()

	s.mutex.Lock()
	for _, key := range expiredKeys {
		delete(s.kvStore, key)
	}
	s.mutex.Unlock()
}

func (s *Store) BackgroundCleanUpService(interval time.Duration, stop <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.CleanUpExpiredKeys()
			case <-stop:
				return
			}
		}
	}()
}
