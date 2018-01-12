package rimu

import (
	"testing"
)

func TestRimuBasic(t *testing.T) {
	source := "*Hello World!*"
	expected := "<p><em>Hello World!</em></p>"
	got := Render(source, RenderOptions{})
	if got != expected {
		t.Errorf("TestRimuBasic(%q) == %q, expected %q", source, got, expected)
	}
}
