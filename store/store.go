package store

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kjj1998/kvstore/errors"
	"github.com/kjj1998/kvstore/models"
)

type Store struct {
	kvStore map[string]models.Value
	mutex   sync.RWMutex
	walChan chan []string
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewStore(ctx context.Context, cancel context.CancelFunc) *Store {
	return &Store{
		kvStore: make(map[string]models.Value),
		walChan: make(chan []string, 1000),
		ctx:     ctx,
		cancel:  cancel,
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
		expiry = time.Now().UTC().Add(time.Duration(timeToLive) * time.Second)
	}

	s.walChan <- []string{"SET", key, value, expiry.Format(time.RFC3339)}

	s.kvStore[key] = models.Value{
		Value:  value,
		Expiry: expiry,
	}
}

func (s *Store) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.walChan <- []string{"DEL", key}
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

func (s *Store) BackgroundCleanUpService(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.CleanUpExpiredKeys()
			case <-s.ctx.Done():
				return
			}
		}
	}()
}

func (s *Store) RecoverFromLog() {
	f, err := os.Open("persistent_log.txt")
	errors.LogError(err, "Error opening the persistent log file: ")
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		s.writeToStore(line)
	}

	err = scanner.Err()
	errors.LogError(err, "Error occurred during scanning: ")
}

func (s *Store) StartWALWriterGoroutine() {
	f, err := os.OpenFile("persistent_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	errors.LogError(err, "Error opening the persistent log file: ")

	go func() {
		for {
			select {
			case commands := <-s.walChan:
				logLine := strings.Join(commands, " ")
				_, err := fmt.Fprintf(f, "%s\n", logLine)
				errors.LogError(err, "Error when writing to file: ")
			case <-s.ctx.Done():
				err = f.Sync()
				errors.LogError(err, "Error syncing file: ")
				f.Close()
				return
			}
		}
	}()
}

func (s *Store) writeToStore(line string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	commands := strings.Fields(line)
	switch commands[0] {
	case "SET":
		key := commands[1]
		value := commands[2]
		expiry := time.Time{}

		if len(commands) == 5 {
			parsedTime, err := time.Parse(time.RFC3339, commands[4])
			errors.LogError(err, "Error parsing time string: ")
			expiry = parsedTime
		}

		s.kvStore[key] = models.Value{
			Value:  value,
			Expiry: expiry,
		}
	case "DEL":
		key := commands[1]
		delete(s.kvStore, key)
	default:
		break
	}
}
