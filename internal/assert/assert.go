package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, want, got T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", want, got)
	}
}

func NotEqual[T comparable](t *testing.T, want, got T) {
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

func EqualValues[T comparable](t *testing.T, want, got []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("got %v, want %v", want, got)
	} else {
		for k, _ := range got {
			if got[k] != want[k] {
				t.Errorf("got %v, want %v", want, got)
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

func Contains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("%q does not contain %q", s, substr)
	}
}
