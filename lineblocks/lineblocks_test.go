package lineblocks

import (
	"testing"

	"github.com/srackham/rimu-go/iotext"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	input := "# Test"
	reader := iotext.NewReader(input)
	writer := iotext.NewWriter()
	Render(reader, writer)
	assert.Equal(t, "<h1>Test</h1>", writer.String())
}
