package lineblocks

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/srackham/go-rimu/v11/internal/blockattributes"
	"github.com/srackham/go-rimu/v11/internal/delimitedblocks"
	"github.com/srackham/go-rimu/v11/internal/expansion"
	"github.com/srackham/go-rimu/v11/internal/iotext"
	"github.com/srackham/go-rimu/v11/internal/macros"
	"github.com/srackham/go-rimu/v11/internal/options"
	"github.com/srackham/go-rimu/v11/internal/quotes"
	"github.com/srackham/go-rimu/v11/internal/replacements"
	"github.com/srackham/go-rimu/v11/internal/spans"
	"github.com/srackham/go-rimu/v11/internal/utils/stringlist"
)

type Definition struct {
	match       *regexp.Regexp
	replacement string
	name        string // Optional unique identifier.
	filter      LineBlockFilter
	verify      LineBlockVerify // Additional match verification checks.
}

type LineBlockFilter = func(match []string, reader *iotext.Reader, def Definition) string
type LineBlockVerify = func(match []string, reader *iotext.Reader) bool // Additional match verification checks.

var defs = []Definition{
	// Prefix match with backslash to allow escaping.

	// Comment line.
	{
		match: regexp.MustCompile(`^\\?\/{2}(.*)$`),
	},
	// Expand lines prefixed with a macro invocation prior to all other processing.
	// macro name = $1, macro value = $2
	{
		match: macros.MATCH_LINE,
		verify: func(match []string, reader *iotext.Reader) bool {
			if macros.LITERAL_DEF_OPEN.MatchString(match[0]) || macros.EXPRESSION_DEF_OPEN.MatchString(match[0]) {
				// Do not process macro definitions.
				return false
			}
			// Silent because any macro expansion errors will be subsequently addressed downstream.
			value := macros.Render(match[0], true)
			if strings.HasPrefix(value, match[0]) || strings.Contains(value, "\n"+match[0]) {
				// The leading macro invocation expansion failed or contains itself.
				// This stops infinite recursion.
				return false
			}
			// Insert the macro value into the reader just ahead of the cursor.
			reader.Lines = stringlist.StringList(reader.Lines).InsertAt(reader.Pos+1, strings.Split(value, "\n")...)
			return true
		},
		filter: func(_ []string, _ *iotext.Reader, _ Definition) string {
			return "" // Already processed in the `verify` function.
		},
	},
	// Delimited Block definition.
	// name = $1, definition = $2
	{
		match: regexp.MustCompile(`^\\?\|([\w\-]+)\|\s*=\s*'(.*)'$`),
		filter: func(match []string, _ *iotext.Reader, _ Definition) string {
			if options.IsSafeModeNz() {
				return "" // Skip if a safe mode is set.
			}
			match[2] = spans.ReplaceInline(match[2], expansion.Options{Macros: true})
			delimitedblocks.SetDefinition(match[1], match[2])
			return ""
		},
	},
	// Quote definition.
	// quote = $1, openTag = $2, separator = $3, closeTag = $4
	{
		match: regexp.MustCompile(`^(\S{1,2})\s*=\s*'([^|]*)(\|{1,2})(.*)'$`),
		filter: func(match []string, _ *iotext.Reader, _ Definition) string {
			if options.IsSafeModeNz() {
				return "" // Skip if a safe mode is set.
			}
			quotes.SetDefinition(quotes.Definition{
				Quote:    match[1],
				OpenTag:  spans.ReplaceInline(match[2], expansion.Options{Macros: true}),
				CloseTag: spans.ReplaceInline(match[4], expansion.Options{Macros: true}),
				Spans:    match[3] == "|",
			})
			return ""
		},
	},
	// Replacement definition.
	// pattern = $1, flags = $2, replacement = $3
	{
		match: regexp.MustCompile(`^\\?\/(.+)\/([igm]*)\s*=\s*'(.*)'$`),
		filter: func(match []string, _ *iotext.Reader, _ Definition) string {
			if options.IsSafeModeNz() {
				return "" // Skip if a safe mode is set.
			}
			pattern := match[1]
			flags := match[2]
			replacement := match[3]
			replacement = spans.ReplaceInline(replacement, expansion.Options{Macros: true})
			replacements.SetDefinition(pattern, flags, replacement)
			return ""
		},
	},
	// Macro definition.
	// name = $1, value = $2
	{
		match: macros.LINE_DEF,
		verify: func(match []string, reader *iotext.Reader) bool {
			// Necessary because Go regexps do not support regexp backreferences,
			return match[2] == match[4] // Leading and trailing quote must match.
		},
		filter: func(match []string, _ *iotext.Reader, _ Definition) string {
			name := match[1]
			quote := match[2]
			value := match[3]
			value = spans.ReplaceInline(value, expansion.Options{Macros: true})
			macros.SetValue(name, value, quote)
			return ""
		},
	},
	// Headers.
	// $1 is ID, $2 is header text, $3 is the optional trailing ID.
	{
		match:       regexp.MustCompile(`^\\?([#=]{1,6})\s+(.+?)(?:\s+([#=]{1,6}))?$`),
		replacement: "<h$1>$$2</h$1>",
		verify: func(match []string, reader *iotext.Reader) bool {
			// Necessary because Go regexps do not support regexp backreferences,
			return match[3] == "" || match[3] == match[1] // Leading and trailing IDs must match.
		},
		filter: func(match []string, _ *iotext.Reader, def Definition) string {
			match[1] = fmt.Sprint(len(match[1])) // Replace $1 with header number.
			if macros.IsNotBlank("--header-ids") && blockattributes.Id == "" {
				blockattributes.Id = blockattributes.Slugify(match[2])
			}
			return spans.ReplaceMatch(match, def.replacement, expansion.Options{Macros: true})
		},
	},
	// Block image: <image:src|alt>
	// src = $1, alt = $2
	{
		match:       regexp.MustCompile(`^\\?<image:([^\s|]+)\|(.+?)>$`),
		replacement: "<img src=\"$1\" alt=\"$2\">",
	},
	// Block image: <image:src>
	// src = $1, alt = $1
	{
		match:       regexp.MustCompile(`^\\?<image:([^\s|]+?)>$`),
		replacement: "<img src=\"$1\" alt=\"$1\">",
	},
	// DEPRECATED as of 3.4.0.
	// Block anchor: <<#id>>
	// id = $1
	{
		match:       regexp.MustCompile(`^\\?<<#([a-zA-Z][\w\-]*)>>$`),
		replacement: "<div id=\"$1\"></div>",
		filter: func(match []string, _ *iotext.Reader, def Definition) string {
			if options.SkipBlockAttributes() {
				return ""
			} else {
				// Default (non-filter) replacement processing.
				return spans.ReplaceMatch(match, def.replacement, expansion.Options{Macros: true})
			}
		},
	},
	// Block Attributes.
	// Syntax: .class-names #id [html-attributes] block-options
	{
		name:  "attributes",
		match: regexp.MustCompile(`^\\?\.[a-zA-Z#"\[+-].*$`), // A loose match because Block Attributes can contain macro references.
		verify: func(match []string, _ *iotext.Reader) bool {
			return blockattributes.Parse(match[0])
		},
	},
	// API Option.
	// name = $1, value = $2
	{
		match: regexp.MustCompile(`^\\?\.(\w+)\s*=\s*'(.*)'$`),
		filter: func(match []string, _ *iotext.Reader, _ Definition) string {
			if !options.IsSafeModeNz() {
				value := spans.ReplaceInline(match[2], expansion.Options{Macros: true})
				options.SetOption(match[1], value)
			}
			return ""
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
		if len(allowed) > 0 && !allowed.Contains(def.name) {
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
				text = spans.ReplaceMatch(match, def.replacement, expansion.Options{Macros: true})
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
