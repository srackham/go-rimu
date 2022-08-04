package spans

import (
	"testing"

	"github.com/srackham/go-rimu/v11/internal/assert"
)

func TestRender(t *testing.T) {
	tests := []struct {
		source string
		want   string
	}{
		{"", ""},
		{"foo", "foo"},
		// Quotes.
		{"*foo*", "<em>foo</em>"},
		{"**foo**", "<strong>foo</strong>"},
		{"*foo* **bar**", "<em>foo</em> <strong>bar</strong>"},
		{"*foo __bar__*", "<em>foo <strong>bar</strong></em>"},
		{"`**foo**`", "<code>**foo**</code>"},
		// Replacements.
		{"<image:foo|bar>", `<img src="foo" alt="bar">`},
		{"<image:foo|bar\nboo>", "<img src=\"foo\" alt=\"bar\nboo\">"},
	}
	for _, tt := range tests {
		got := Render(tt.source)
		assert.Equal(t, tt.want, got)
	}
}

func Test_defrag(t *testing.T) {
	tests := []struct {
		frags []fragment
		want  string
	}{
		{
			frags: []fragment{
				{text: ""},
			},
			want: "",
		},
		{
			frags: []fragment{
				{text: "foo"},
				{text: "bar"},
			},
			want: "foobar",
		},
	}
	for _, tt := range tests {
		got := defrag(tt.frags)
		assert.Equal(t, tt.want, got)
	}
}
