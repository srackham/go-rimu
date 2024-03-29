package lists

import (
	"testing"

	"github.com/srackham/go-rimu/v11/internal/assert"
	"github.com/srackham/go-rimu/v11/internal/iotext"
)

func TestRender(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{`- foo`, `<ul><li>foo</li></ul>`},
	}
	for _, tt := range tests {
		reader := iotext.NewReader(tt.in)
		writer := iotext.NewWriter()
		Render(reader, writer)
		got := writer.String()
		assert.Equal(t, tt.want, got)
	}
}
