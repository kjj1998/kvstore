package tests

import (
	"testing"
	"time"

	"github.com/kjj1998/kvstore/models"
)

func TestKvStoreGet(t *testing.T) {
	store := StubStore{
		kvStore: map[string]models.Value{
			"hello":   {Value: "world", Expiry: CurrentTime.Add(1 * time.Minute)},
			"goodbye": {Value: "world", Expiry: CurrentTime.Add(-1 * time.Minute)},
		},
	}

	t.Run("Get fresh value from the key-value store", func(t *testing.T) {
		got := store.Get("hello")
		want := "world"

		AssertEquality(t, got, want)
	})

	t.Run("Get expired value from the key-value store", func(t *testing.T) {
		got := store.Get("goodbye")
		want := "EXPIRED"

		AssertEquality(t, got, want)
	})
}

func TestKvStoreSet(t *testing.T) {
	store := StubStore{
		kvStore: map[string]models.Value{
			"hello": {Value: "world", Expiry: CurrentTime.Add(-time.Minute)},
		},
	}

	t.Run("Set value in key-value store", func(t *testing.T) {
		store.Set("hello", "world", 0)

		got := store.Get("hello")
		want := "world"

		AssertEquality(t, got, want)
	})

	t.Run("Override existing value in key-value store", func(t *testing.T) {
		store.Set("go", "is cool", 0)
		store.Set("go", "is super cool", 0)

		got := store.Get("go")
		want := "is super cool"

		AssertEquality(t, got, want)
	})

	t.Run("Set value in the key-value store with a time-to-live in seconds", func(t *testing.T) {
		store.Set("set", "expiry", 30)

		got := store.Get("set")
		want := "expiry"

		AssertEquality(t, got, want)
	})
}

func TestKvStoreDelete(t *testing.T) {
	store := StubStore{
		kvStore: map[string]models.Value{
			"hello": {Value: "world", Expiry: CurrentTime.Add(1 * time.Minute)},
		},
	}

	t.Run("Delete value stored in key-value store", func(t *testing.T) {
		store.Delete("hello")

		got := store.Get("hello")
		want := "NULL"

		AssertEquality(t, got, want)
	})
}
