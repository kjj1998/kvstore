package tests

import (
	"time"

	"github.com/kjj1998/kvstore/models"
)

var expiryDate = time.Date(2025, time.June, 20, 15, 30, 0, 0, time.UTC)

type StubStore struct {
	kvStore map[string]models.Value
}

func (s *StubStore) Get(key string) (string, bool) {
	value, exists := s.kvStore[key]

	if value.Expiry.After(expiryDate) {
		s.Delete(key)
		return "", false
	}

	return value.Value, exists
}

func (s *StubStore) Set(key, value string) {
	s.kvStore[key] = models.Value{
		Value:  value,
		Expiry: expiryDate.Add(-2 * time.Minute),
	}
}

func (s *StubStore) Delete(key string) {
	delete(s.kvStore, key)
}
