package main

import (
	"testing"
)

func TestKVStore(t *testing.T) {
	store := NewStore()
	store.kvStore = make(map[string]string)
	store.kvStore["hello"] = "world"
	store.kvStore["to be"] = "deleted"

	t.Run("returns data from the key-value store", func(t *testing.T) {
		got, _ := store.Get("hello")
		want := "world"

		if got != want {
			t.Errorf("got %v, wanted %v\n", got, want)
		}
	})

	t.Run("sets data in the key-value store", func(t *testing.T) {
		store.Set("key", "value")
		got, _ := store.Get("key")
		want := "value"

		if got != want {
			t.Errorf("got %v, wanted %v\n", got, want)
		}
	})

	t.Run("deletes data in the key-value store", func(t *testing.T) {
		store.Delete("to be")
		value, exists := store.Get("to be")

		if exists {
			t.Errorf("%v should not exist\n", value)
		}
	})
}
