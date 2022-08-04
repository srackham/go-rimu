package lineblocks

import (
	"testing"

	"github.com/srackham/go-rimu/v11/internal/assert"
	"github.com/srackham/go-rimu/v11/internal/iotext"
)

func TestRender(t *testing.T) {
	tests := []struct {
		source string
		want   string
	}{
		{`# foo`, `<h1>foo</h1>`},
		{`// foo`, ``},
		{`<image:foo|bar>`, `<img src="foo" alt="bar">`},
		{`<<#foo>>`, `<div id="foo"></div>`},
		{`.class #id "css"`, ``},
		{`.safeMode='0'`, ``},
		{`|code|='<code>|</code>'`, ``},
		{`^='<sup>|</sup>'`, ``},
		{`/\.{3}/i = '&hellip;'`, ``},
		{`{foo}='bar'`, ``},
	}
	for _, tt := range tests {
		reader := iotext.NewReader(tt.source)
		writer := iotext.NewWriter()
		Render(reader, writer, nil)
		got := writer.String()
		assert.Equal(t, tt.want, got)
	}
}
