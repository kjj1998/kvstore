package tests

import (
	"testing"
	"time"

	"github.com/kjj1998/kvstore/models"
)

func TestKvStoreGet(t *testing.T) {
	store := StubStore{
		kvStore: map[string]models.Value{
			"hello":   {Value: "world", Expiry: expiryDate.Add(-time.Minute)},
			"goodbye": {Value: "world", Expiry: expiryDate.Add(time.Minute)},
		},
	}

	t.Run("Get fresh value from the key-value store", func(t *testing.T) {
		got, _ := store.Get("hello")
		want := "world"

		AssertEquality(t, got, want)
	})

	t.Run("Get expired value from the key-value store", func(t *testing.T) {
		got, exists := store.Get("goodbye")
		want := ""

		if exists {
			t.Errorf("value exists, expected value to not exist\n")
		}

		AssertEquality(t, got, want)
	})
}

func TestKvStoreSet(t *testing.T) {
	store := StubStore{
		kvStore: map[string]models.Value{
			"hello": {Value: "world", Expiry: expiryDate.Add(-time.Minute)},
		},
	}

	t.Run("Set value in key-value store", func(t *testing.T) {
		store.Set("hello", "world")

		got, _ := store.Get("hello")
		want := "world"

		AssertEquality(t, got, want)
	})

	t.Run("Override existing value in key-value store", func(t *testing.T) {
		store.Set("hello", "world")
		store.Set("hello", "earth")

		got, _ := store.Get("hello")
		want := "earth"

		AssertEquality(t, got, want)
	})
}

func TestKvStoreDelete(t *testing.T) {
	store := StubStore{
		kvStore: map[string]models.Value{
			"hello": {Value: "world", Expiry: expiryDate.Add(-time.Minute)},
		},
	}

	t.Run("Delete value stored in key-value store", func(t *testing.T) {
		store.Delete("hello")

		got, exists := store.Get("hello")
		want := ""

		if exists {
			t.Errorf("value exists, expected value to not exist\n")
		}
		AssertEquality(t, got, want)
	})
}
