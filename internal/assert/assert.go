package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func NotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("did not want %v", got)
	}
}

func True(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Error("should be true")
	}
}

func False(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Error("should be false")
	}
}

func EqualValues[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("got %v, want %v", got, want)
	} else {
		for k, _ := range got {
			if got[k] != want[k] {
				t.Errorf("got %v, want %v", got, want)
			}
		}
	}
}

func Panics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		recover()
	}()
	f()
	t.Errorf("should have paniced")
}
