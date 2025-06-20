package main

import (
	"testing"
	"time"
)

type StubStore struct {
	kvStore map[string]Value
}

var expiry = time.Date(2025, time.June, 20, 15, 30, 0, 0, time.UTC)

func (s *StubStore) Get(key string) (string, bool) {
	value, exists := s.kvStore[key]

	if value.expiry.After(expiry) {
		s.Delete(key)
		return "", false
	}

	return value.val, exists
}

func (s *StubStore) Set(key, value string) {
	s.kvStore[key] = Value{
		val:    value,
		expiry: expiry.Add(-2 * time.Minute),
	}
}

func (s *StubStore) Delete(key string) {
	delete(s.kvStore, key)
}

func assertEquality(t testing.TB, got, want string) {
	if got != want {
		t.Errorf("got %v, but wanted %v\n", got, want)
	}
}

func TestKvStoreGet(t *testing.T) {
	store := StubStore{
		kvStore: map[string]Value{
			"hello":   {val: "world", expiry: expiry.Add(-time.Minute)},
			"goodbye": {val: "world", expiry: expiry.Add(time.Minute)},
		},
	}

	t.Run("Get fresh value from the key-value store", func(t *testing.T) {
		got, _ := store.Get("hello")
		want := "world"

		assertEquality(t, got, want)
	})

	t.Run("Get expired value from the key-value store", func(t *testing.T) {
		got, exists := store.Get("goodbye")
		want := ""

		if exists {
			t.Errorf("value exists, expected value to not exist\n")
		}

		assertEquality(t, got, want)
	})
}

func TestKvStoreSet(t *testing.T) {
	store := StubStore{
		kvStore: map[string]Value{
			"hello": {val: "world", expiry: expiry.Add(-time.Minute)},
		},
	}

	t.Run("Set value in key-value store", func(t *testing.T) {
		store.Set("hello", "world")

		got, _ := store.Get("hello")
		want := "world"

		assertEquality(t, got, want)
	})
}

func TestKvStoreDelete(t *testing.T) {
	store := StubStore{
		kvStore: map[string]Value{
			"hello": {val: "world", expiry: expiry.Add(-time.Minute)},
		},
	}

	t.Run("Delete value stored in key-value store", func(t *testing.T) {
		store.Delete("hello")

		got, exists := store.Get("hello")
		want := ""

		if exists {
			t.Errorf("value exists, expected value to not exist\n")
		}
		assertEquality(t, got, want)
	})
}

// func TestKVStore(t *testing.T) {
// 	store := NewStore()
// 	store.kvStore = make(map[string]Value)
// 	store.kvStore["hello"] = Value{val: "world", expiry: time.Date(2025, time.June, 20, 15, 30, 30, 0, time.UTC)}
// 	store.kvStore["to be"] = "deleted"

// 	t.Run("returns data from the key-value store", func(t *testing.T) {
// 		got, _ := store.Get("hello")
// 		want := "world"

// 		if got != want {
// 			t.Errorf("got %v, wanted %v\n", got, want)
// 		}
// 	})

// 	t.Run("sets data in the key-value store", func(t *testing.T) {
// 		store.Set("key", "value")
// 		got, _ := store.Get("key")
// 		want := "value"

// 		if got != want {
// 			t.Errorf("got %v, wanted %v\n", got, want)
// 		}
// 	})

// 	t.Run("deletes data in the key-value store", func(t *testing.T) {
// 		store.Delete("to be")
// 		value, exists := store.Get("to be")

// 		if exists {
// 			t.Errorf("%v should not exist\n", value)
// 		}
// 	})
// }
