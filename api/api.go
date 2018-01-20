package api

import (
	"github.com/srackham/rimu-go/blockattributes"
	"github.com/srackham/rimu-go/delimitedblocks"
	"github.com/srackham/rimu-go/iotext"
	"github.com/srackham/rimu-go/lineblocks"
	"github.com/srackham/rimu-go/lists"
	"github.com/srackham/rimu-go/macros"
	"github.com/srackham/rimu-go/options"
	"github.com/srackham/rimu-go/quotes"
	"github.com/srackham/rimu-go/replacements"
)

func init() {
	// Dependency injectiion so we can use api functions in imported packages without incuring import cycle errors.
	options.ApiInit = Init
	delimitedblocks.ApiRender = Render
}

// Init TODO
func Init() {
	blockattributes.Init()
	options.Init()
	delimitedblocks.Init()
	macros.Init()
	quotes.Init()
	replacements.Init()
}

// Render TODO
func Render(source string) string {
	reader := iotext.NewReader(source)
	writer := iotext.NewWriter()
	for !reader.Eof() {
		reader.SkipBlankLines()
		if reader.Eof() {
			break
		}
		if lineblocks.Render(reader, writer, nil) {
			continue
		}
		if lists.Render(reader, writer) {
			continue
		}
		if delimitedblocks.Render(reader, writer, nil) {
			continue
		}
		// This code should never be executed (normal paragraphs should match anything).
		panic("no matching delimited block found")
	}
	return writer.String()
}
