package lineblocks

import (
	"testing"

	"github.com/srackham/rimu-go/iotext"
	_ "github.com/srackham/rimu-go/spans"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	tests := []struct {
		source string
		want   string
	}{
		{"# foo", "<h1>foo</h1>"},
		{"// foo", ""},
		{"<image:foo|bar>", `<img src="foo" alt="bar">`},
		{"<image:foo|bar>", `<img src="foo" alt="bar">`},
	}
	for _, tt := range tests {
		reader := iotext.NewReader(tt.source)
		writer := iotext.NewWriter()
		Render(reader, writer, nil)
		got := writer.String()
		assert.Equal(t, tt.want, got)
	}
}
