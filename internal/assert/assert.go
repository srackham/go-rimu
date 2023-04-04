/*
	Simple assertions package.
*/

package assert

import (
	"regexp"
	"strings"
	"testing"
)

// PassIf fails and prints formatted message if not ok.
func PassIf(t *testing.T, ok bool, format string, args ...any) {
	t.Helper()
	if strings.Contains(format, "%s") {
		t.Logf("use '%%#v' instead of '%%s' to report possibly nil values: '%s'", format)
		t.FailNow()
	}
	if !ok {
		t.Errorf(format, args...)
	}
}

func Equal[T comparable](t *testing.T, wanted, got T) {
	t.Helper()
	PassIf(t, got == wanted, "wanted %#v, got %#v", wanted, got)
}

func NotEqual[T comparable](t *testing.T, wanted, got T) {
	t.Helper()
	PassIf(t, wanted != got, "should not be %#v", got)
}

func True(t *testing.T, got bool) {
	t.Helper()
	PassIf(t, got, "should be true")
}

func False(t *testing.T, got bool) {
	t.Helper()
	PassIf(t, !got, "should be false")
}

func EqualValues[T comparable](t *testing.T, wanted, got []T) {
	t.Helper()
	PassIf(t, len(got) == len(wanted), "wanted %#v, got %#v", wanted, got)
	for k := range got {
		PassIf(t, got[k] == wanted[k], "wanted %#v, got %#v", wanted, got)
	}
}

func Panics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		recover()
	}()
	f()
	t.Error("should have panicked")
}

func Contains(t *testing.T, s, substr string) {
	t.Helper()
	PassIf(t, strings.Contains(s, substr), "%q does not contain %q", s, substr)
}

func ContainsPattern(t *testing.T, s, pattern string) {
	t.Helper()
	matched, _ := regexp.MatchString(pattern, s)
	PassIf(t, matched, "%q does not contain pattern %q", s, pattern)
}
