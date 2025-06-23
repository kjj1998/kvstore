package tests

import (
	"time"

	"github.com/kjj1998/kvstore/models"
)

var CurrentTime = time.Date(2025, time.June, 20, 15, 30, 0, 0, time.UTC)

type StubStore struct {
	kvStore map[string]models.Value
}

func (s *StubStore) Get(key string) string {
	value, exists := s.kvStore[key]

	expired := exists && !value.Expiry.IsZero() && value.Expiry.Before(CurrentTime)

	if expired {
		s.Delete(key)
		return "EXPIRED"
	}

	if !exists {
		return "NULL"
	}

	return value.Value
}

func (s *StubStore) Set(key, value string, timeToLive int) {
	var expiry = time.Time{}

	if timeToLive == 0 {
		expiry = CurrentTime.Add(2 * time.Minute)
	} else {
		expiry = CurrentTime.Add(time.Duration(timeToLive) * time.Second)
	}

	s.kvStore[key] = models.Value{
		Value:  value,
		Expiry: expiry,
	}
}

func (s *StubStore) Delete(key string) {
	delete(s.kvStore, key)
}
