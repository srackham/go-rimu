package delimitedblocks

import (
	"testing"

	"github.com/srackham/go-rimu-mod/internal/iotext"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.Equal(t, len(DEFAULT_DEFS), len(defs))
	assert.NotEqual(t, DEFAULT_DEFS, defs)
}

func TestRender(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"foo", "<p>foo</p>"},
		{"  foo", "<pre><code>foo</code></pre>"},
	}
	var reader *iotext.Reader
	var writer *iotext.Writer
	for _, tt := range tests {
		reader = iotext.NewReader(tt.in)
		writer = iotext.NewWriter()
		assert.True(t, Render(reader, writer, nil))
		assert.Equal(t, tt.want, writer.String())
	}
}

func TestGetDefinition(t *testing.T) {
	Init()
	def := GetDefinition("paragraph")
	assert.Equal(t, "<p>", def.openTag)
	def = GetDefinition("foo")
	assert.Nil(t, def)
}

func TestSetDefinition(t *testing.T) {
	Init()
	SetDefinition("indented", "<foo>|</foo>")
	def := GetDefinition("indented")
	assert.Equal(t, "<foo>", def.openTag)
	assert.Equal(t, "</foo>", def.closeTag)
}
