package document

import (
	"github.com/srackham/go-rimu/v11/internal/blockattributes"
	"github.com/srackham/go-rimu/v11/internal/delimitedblocks"
	"github.com/srackham/go-rimu/v11/internal/iotext"
	"github.com/srackham/go-rimu/v11/internal/lineblocks"
	"github.com/srackham/go-rimu/v11/internal/lists"
	"github.com/srackham/go-rimu/v11/internal/macros"
	"github.com/srackham/go-rimu/v11/internal/options"
	"github.com/srackham/go-rimu/v11/internal/quotes"
	"github.com/srackham/go-rimu/v11/internal/replacements"
)

func init() {
	// Dependency injectiion so we can use document functions in imported packages without incuring import cycle errors.
	options.ApiInit = Init
	delimitedblocks.ApiRender = Render
}

// Init initialises Rimu state.
func Init() {
	blockattributes.Init()
	options.Init()
	delimitedblocks.Init()
	macros.Init()
	quotes.Init()
	replacements.Init()
}

// Render source text to HTML string.
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
