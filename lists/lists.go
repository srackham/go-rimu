package lists

import (
	"github.com/srackham/rimu-go/api"
	"github.com/srackham/rimu-go/iotext"
)

func init() {
	api.ListsRender = Render
}

// TODO
// Stubs

func Render(reader *iotext.Reader, writer *iotext.Writer) bool {
	return false
}
