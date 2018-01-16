package lineblocks

import (
	"fmt"
	"regexp"

	"github.com/srackham/rimu-go/blockattributes"
	"github.com/srackham/rimu-go/macros"

	"github.com/srackham/rimu-go/utils/stringlist"

	"github.com/srackham/rimu-go/iotext"
	"github.com/srackham/rimu-go/utils"
)

type LineBlockFilter = func(match []string, reader *iotext.Reader, def Definition) string
type LineBlockVerify = func(match []string, reader *iotext.Reader) bool // Additional match verification checks.

type Definition struct {
	match       *regexp.Regexp
	replacement string
	name        string // Optional unique identifier.
	filter      LineBlockFilter
	verify      LineBlockVerify // Additional match verification checks.
}

var defs = []Definition{
	// Headers.
	// $1 is ID, $2 is header text.
	{
		// The Go regexp package does not support regexp backreferences,
		// see https://stackoverflow.com/questions/23968992/how-to-match-a-regex-with-backreference-in-go
		// match:       regexp.MustCompile(`^\\?([#=]{1,6})\s+(.+?)(?:\s+\1)?$`),
		match:       regexp.MustCompile(`^\\?([#=]{1,6})\s+(.+?)(?:\s+[#=]{1,6})?$`),
		replacement: "<h$1>$$2</h$1>",
		filter: func(match []string, _ *iotext.Reader, def Definition) string {
			match[1] = fmt.Sprint(len(match[1])) // Replace $1 with header number.
			if macros.IsDefined("--header-ids") && blockattributes.Id == "" {
				blockattributes.Id = blockattributes.Slugify(match[2])
			}
			return utils.ReplaceMatch(match, def.replacement, utils.ExpansionOptions{Macros: true})
		},
	},
}

// If the next element in the reader is a valid line block render it
// and return true, else return false.
func Render(reader *iotext.Reader, writer *iotext.Writer, allowed stringlist.StringList) bool {
	if reader.Eof() {
		panic("premature eof")
	}
	for _, def := range defs {
		if len(allowed) > 0 && allowed.Contains(def.name) {
			continue
		}
		match := def.match.FindStringSubmatch(reader.Cursor())
		if match != nil {
			if match[0][0] == '\\' {
				// Drop backslash escape and continue.
				reader.SetCursor(reader.Cursor()[1:])
				continue
			}
			if def.verify != nil && !def.verify(match, reader) {
				continue
			}
			var text string
			if def.filter == nil {
				text = utils.ReplaceMatch(match, def.replacement, utils.ExpansionOptions{Macros: true})
			} else {
				text = def.filter(match, reader, def)
			}
			if text != "" {
				text = blockattributes.Inject(text)
				writer.Write(text)
				reader.Next()
				if !reader.Eof() {
					writer.Write("\n") // Add a trailing '\n' if there are more lines.
				}
			} else {
				reader.Next()
			}
			return true
		}
	}
	return false
}
