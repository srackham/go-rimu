package api

import (
	"github.com/srackham/rimu-go/iotext"
)

type initFunc = func()

// Initialisation functions registered by low-level packages.
var (
	BlockAttributesInit initFunc
	OptionsInit         initFunc
	DelimitedBlocksInit initFunc
	MacrosInit          initFunc
	QuotesInit          initFunc
	ReplacementsInit    initFunc
)

// Init performs global API initialisation.
func Init() {
	BlockAttributesInit()
	OptionsInit()
	DelimitedBlocksInit()
	MacrosInit()
	QuotesInit()
	ReplacementsInit()
}

type renderFunc = func(reader *iotext.Reader, writer *iotext.Writer) bool

// Render functions registered by low-level packages.
var (
	DelimitedBlocksRender renderFunc
	ListsRender           renderFunc
	LineBlocksRender      renderFunc
)

// Render converts Rimu source markup and returns HTML.
func Render(source string) string {
	reader := iotext.NewReader(source)
	writer := iotext.NewWriter()
	// outer:
	for !reader.Eof() {
		if reader.Eof() {
			break
		}
		if LineBlocksRender(reader, writer) {
			continue
		}
		if ListsRender(reader, writer) {
			continue
		}
		if DelimitedBlocksRender(reader, writer) {
			continue
		}
		panic("no matching delimited block found")
	}
	return writer.String()
}
