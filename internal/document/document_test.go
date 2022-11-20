package document

import (
	"testing"

	"github.com/srackham/go-rimu/v11/internal/assert"
)

func TestInit(t *testing.T) {
	Init()

}

func TestRender(t *testing.T) {
	in := "# Title\nParagraph **bold** `code` _emphasised text_\n\n.test-class [title=\"Code\"]\n  Indented `paragraph`\n\n- Item 1\n\"\"\nQuoted\n\"\"\n- Item 2\n . Nested 1\n\n{x} = '1$$1$$2'\n{x?} = '2'\n\\{x}={x|}\n{x|2|3}"
	want := "<h1>Title</h1>\n<p>Paragraph <strong>bold</strong> <code>code</code> <em>emphasised text</em></p>\n<pre class=\"test-class\" title=\"Code\"><code>Indented `paragraph`</code></pre>\n<ul><li>Item 1<blockquote><p>Quoted</p></blockquote>\n</li><li>Item 2<ol><li>Nested 1</li></ol></li></ul><p>{x}=1\n123</p>"
	got := Render(in)
	assert.Equal(t, want, got)
}
