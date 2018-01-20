package api

import (
	"github.com/srackham/rimu-go/iotext"
)

func init() {
	// 	// So we can use these functions in imported packages without incuring import cycle errors.
	// 	proxies.ApiInit = Init
	// 	proxies.ApiRender = Render
}

var initializers []func()
var renderers []func(reader *iotext.Reader, writer *iotext.Writer) bool

func RegisterInit(init func()) {
	initializers = append(initializers, init)
}

func RegisterRender(render func(reader *iotext.Reader, writer *iotext.Writer) bool) {
	renderers = append(renderers, render)
}

// Init TODO
func Init() {
	for _, init := range initializers {
		init()
	}
	// blockattributes.Init()
	// options.Init()
	// delimitedblocks.Init()
	// macros.Init()
	// quotes.Init()
	// replacements.Init()
}

// Render TODO
func Render(source string) string {
	reader := iotext.NewReader(source)
	writer := iotext.NewWriter()
outer:
	for !reader.Eof() {
		if reader.Eof() {
			break
		}
		reader.SkipBlankLines()
		for _, render := range renderers {
			if render(reader, writer) {
				continue outer
			}
		}
		panic("no matching delimited block found")

		// if reader.Eof() {
		// 	break
		// }
		// if lineblocks.Render(reader, writer, []string{}) {
		// 	continue
		// }
		// if lists.Render(reader, writer) {
		// 	continue
		// }
		// if delimitedblocks.Render(reader, writer, nil) {
		// 	continue
		// }
		// This code should never be executed (normal paragraphs should match anything).
		// panic("no matching delimited block found")
	}
	return writer.String()
}
