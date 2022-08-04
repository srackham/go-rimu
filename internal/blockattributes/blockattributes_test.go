package blockattributes

import (
	"testing"

	"github.com/srackham/go-rimu/v11/internal/assert"

	_ "github.com/srackham/go-rimu/v11/internal/macros"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in   string
		want attrs
	}{
		{".class #id", attrs{Classes: "class", ID: "id"}},
		{".\"css\"", attrs{css: "css"}},
	}
	for _, tt := range tests {
		Init()
		Parse(tt.in)
		got := Attrs
		assert.Equal(t, tt.want, got)
	}
}

func TestInject(t *testing.T) {
	tests := []struct {
		tag                          string
		want                         string
		classes, id, css, attributes string
	}{
		{tag: `<p>`, id: `id`, want: `<p id="id">`},
		{tag: `<p>`, classes: `class`, want: `<p class="class">`},
		{tag: `<p class="class">`, classes: `class2`, want: `<p class="class2 class">`},
	}
	for _, tt := range tests {
		Init()
		Attrs.ID = tt.id
		Attrs.Classes = tt.classes
		Attrs.css = tt.css
		Attrs.attributes = tt.attributes
		got := Inject(tt.tag)
		assert.Equal(t, tt.want, got)
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Foo Bar", "foo-bar"},
		{"Foo Bar", "foo-bar-2"},
		{"--", "x"},
	}
	Init()
	for _, tt := range tests {
		got := Slugify(tt.in)
		assert.Equal(t, tt.want, got)
		ids.Push("foo-bar")
	}
}
