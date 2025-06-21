package tests

import "testing"

func AssertEquality(t testing.TB, got, want string) {
	if got != want {
		t.Errorf("got %v, but wanted %v\n", got, want)
	}
}
