package lineblocks

import (
	"testing"

	"github.com/srackham/rimu-go/api"
	"github.com/srackham/rimu-go/iotext"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	api.Init()
	input := "# Test"
	reader := iotext.NewReader(input)
	writer := iotext.NewWriter()
	Render(reader, writer, nil)
	assert.Equal(t, "<h1>Test</h1>", writer.String())
}
