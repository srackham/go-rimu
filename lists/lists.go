package lists

import (
	"github.com/srackham/rimu-go/api"
	"github.com/srackham/rimu-go/iotext"
)

func init() {
	api.RegisterInit(Init)
	api.RegisterRender(Render)
}

// TODO
// Stubs

func Init() {
	// TODO
}

// func Render(reader iotext.Reader, writer iotext.Writer) bool {
func Render(reader *iotext.Reader, writer *iotext.Writer) bool {
	return false
}
