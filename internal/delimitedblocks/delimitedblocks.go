package delimitedblocks

import (
	"regexp"
	"strings"

	"github.com/srackham/go-rimu/v11/internal/blockattributes"
	"github.com/srackham/go-rimu/v11/internal/expansion"
	"github.com/srackham/go-rimu/v11/internal/iotext"
	"github.com/srackham/go-rimu/v11/internal/macros"
	"github.com/srackham/go-rimu/v11/internal/options"
	"github.com/srackham/go-rimu/v11/internal/spans"
	"github.com/srackham/go-rimu/v11/internal/utils/stringlist"
)

// document package dependency injection.
var ApiRender func(source string) string

func init() {
	Init()
}

var MATCH_INLINE_TAG = regexp.MustCompile(`(?i)^(a|abbr|acronym|address|b|bdi|bdo|big|blockquote|br|cite|code|del|dfn|em|i|img|ins|kbd|mark|q|s|samp|small|span|strike|strong|sub|sup|time|tt|u|var|wbr)$`)

// Multi-line block element definition.
type Definition struct {
	name            string         // Unique identifier.
	openMatch       *regexp.Regexp // $1 (if defined) is prepended to block content.
	closeMatch      *regexp.Regexp
	openTag         string
	closeTag        string
	verify          func(match []string) bool                    // Additional match verification checks.
	delimiterFilter func(match []string, def *Definition) string // Process opening delimiter. Return any delimiter content.
	contentFilter   func(text string, match []string, opts expansion.Options) string
	options         expansion.Options
}

var defs []Definition // Mutable definitions initialized by DEFAULT_DEFS.

var DEFAULT_DEFS = []Definition{
	// Delimited blocks cannot be escaped with a backslash.

	// Multi-line macro literal value definition.
	{
		name:       "macro-definition",
		openMatch:  macros.LITERAL_DEF_OPEN, // $1 is first line of macro.
		closeMatch: macros.LITERAL_DEF_CLOSE,
		openTag:    "",
		closeTag:   "",
		options: expansion.Options{
			Macros: true,
		},
		delimiterFilter: delimiterTextFilter,
		contentFilter:   macroDefContentFilter,
	},
	// Multi-line macro expression value definition.
	// DEPRECATED as of 11.0.0.
	{
		name:       "deprecated-macro-expression",
		openMatch:  macros.EXPRESSION_DEF_OPEN, // $1 is first line of macro.
		closeMatch: macros.EXPRESSION_DEF_CLOSE,
		openTag:    "",
		closeTag:   "",
		options: expansion.Options{
			Macros: true,
		},
		delimiterFilter: delimiterTextFilter,
		contentFilter:   macroDefContentFilter,
	},
	// Comment block.
	{
		name:       "comment",
		openMatch:  regexp.MustCompile(`^\\?\/\*+$`),
		closeMatch: regexp.MustCompile(`^\*+\/$`),
		openTag:    "",
		closeTag:   "",
		options: expansion.Options{
			Skip:     true,
			Specials: true, // Fall-back if skip is disabled.
		},
	},
	// Division block.
	{
		name:      "division",
		openMatch: regexp.MustCompile(`^\\?(\.{2,})([\w\s-]*)$`), // $1 is delimiter text, $2 is optional class names.
		openTag:   "<div>",
		closeTag:  "</div>",
		options: expansion.Options{
			Container: true,
			Specials:  true, // Fall-back if container is disabled.
		},
		delimiterFilter: classInjectionFilter,
	},
	// Quote block.
	{
		name:      "quote",
		openMatch: regexp.MustCompile(`^\\?("{2,}|>{2,})([\w\s-]*)$`), // $1 is delimiter text, $2 is optional class names.
		openTag:   "<blockquote>",
		closeTag:  "</blockquote>",
		options: expansion.Options{
			Container: true,
			Specials:  true, // Fall-back if container is disabled.
		},
		delimiterFilter: classInjectionFilter,
	},
	// Code block.
	{
		name:      "code",
		openMatch: regexp.MustCompile(`^\\?(-{2,}|` + "`" + `{2,})([\w\s-]*)$`), // $1 is delimiter text, $2 is optional class names.
		openTag:   "<pre><code>",
		closeTag:  "</code></pre>",
		options: expansion.Options{
			Macros:   false,
			Specials: true,
		},
		verify: func(match []string) bool {
			// The deprecated '-' delimiter does not support appended class names.
			return !(match[1][0] == '-' && strings.TrimSpace(match[2]) != "")
		},
		delimiterFilter: classInjectionFilter,
	},
	// HTML block.
	{
		name: "html",
		// Block starts with HTML comment, DOCTYPE directive or block-level HTML start or end tag.
		// $1 is first line of block.
		// $2 is the alphanumeric tag name.
		openMatch:  regexp.MustCompile(`(?i)^(<!--.*|<!DOCTYPE(?:\s.*)?|<\/?([a-z][a-z0-9]*)(?:[\s>].*)?)$`),
		closeMatch: regexp.MustCompile(`^$`),
		openTag:    "",
		closeTag:   "",
		options: expansion.Options{
			Macros: true,
		},
		verify: func(match []string) bool {
			// Return false if the HTML tag is an inline (non-block) HTML tag.
			if match[2] != "" { // Matched alphanumeric tag name.
				return !MATCH_INLINE_TAG.MatchString(match[2])
			} else {
				return true // Matched HTML comment or doctype tag.
			}
		},
		delimiterFilter: delimiterTextFilter,
		// contentFilter:   options.HtmlSafeModeFilter,
		contentFilter: func(text string, _ []string, _ expansion.Options) string {
			return options.HtmlSafeModeFilter(text)
		},
	},
	// Indented paragraph.
	{
		name:       "indented",
		openMatch:  regexp.MustCompile(`^\\?(\s+\S.*)$`), // $1 is first line of block.
		closeMatch: regexp.MustCompile(`^$`),
		openTag:    "<pre><code>",
		closeTag:   "</code></pre>",
		options: expansion.Options{
			Specials: true,
		},
		delimiterFilter: delimiterTextFilter,
		contentFilter: func(text string, _ []string, _ expansion.Options) string {
			// Strip indent from start of each line.
			firstIndent := regexp.MustCompile(`\S`).FindStringIndex(text)[0]
			result := ""
			for _, line := range strings.Split(text, "\n") {
				// Strip first line indent width or up to first non-space character.
				indent := regexp.MustCompile(`\S|$`).FindStringIndex(line)[0]
				if indent > firstIndent {
					indent = firstIndent
				}
				result += line[indent:] + "\n"
			}
			return strings.TrimSuffix(result, "\n")
		},
	},
	// Quote paragraph.
	{
		name:       "quote-paragraph",
		openMatch:  regexp.MustCompile(`^\\?(>.*)$`), // $1 is first line of block.
		closeMatch: regexp.MustCompile(`^$`),
		openTag:    "<blockquote><p>",
		closeTag:   "</p></blockquote>",
		options: expansion.Options{
			Macros:   true,
			Spans:    true,
			Specials: true, // Fall-back if spans is disabled.
		},
		delimiterFilter: delimiterTextFilter,
		contentFilter: func(text string, _ []string, _ expansion.Options) string {
			// Strip leading > from start of each line and unescape escaped leading >.
			result := ""
			for _, line := range strings.Split(text, "\n") {
				line = regexp.MustCompile(`^>`).ReplaceAllString(line, "")
				line = regexp.MustCompile(`^\\>`).ReplaceAllString(line, "")
				result += line + "\n"
			}
			return strings.TrimSuffix(result, "\n")
		},
	},
	// Paragraph (lowest priority, cannot be escaped).
	{
		name:       "paragraph",
		openMatch:  regexp.MustCompile(`(.*)`), // $1 is first line of block.
		closeMatch: regexp.MustCompile(`^$`),
		openTag:    "<p>",
		closeTag:   "</p>",
		options: expansion.Options{
			Macros:   true,
			Spans:    true,
			Specials: true, // Fall-back if spans is disabled.
		},
		delimiterFilter: delimiterTextFilter,
	},
}

// Reset definitions to defaults.
func Init() {
	defs = make([]Definition, len(DEFAULT_DEFS))
	for i, def := range DEFAULT_DEFS {
		defs[i] = def
		defs[i].options = expansion.Options(def.options) // Clone expansion options.
		if def.closeMatch == nil {
			defs[i].closeMatch = def.openMatch
		}
	}
}

// If the next element in the reader is a valid delimited block render it
// and return true, else return false.
func Render(reader *iotext.Reader, writer *iotext.Writer, allowed []string) bool {
	if reader.Eof() {
		panic("premature eof")
	}
	for _, def := range defs {
		if len(allowed) > 0 && stringlist.StringList(allowed).IndexOf(def.name) == -1 {
			continue
		}
		matches := def.openMatch.FindAllStringSubmatch(reader.Cursor(), 1)
		if matches != nil {
			match := matches[0]
			// Escape non-paragraphs.
			if match[0][0] == '\\' && def.name != "paragraph" {
				// Drop backslash escape and continue.
				reader.SetCursor(reader.Cursor()[1:])
				continue
			}
			if def.verify != nil && !def.verify(match) {
				continue
			}
			// Process opening delimiter.
			delimiterText := ""
			if def.delimiterFilter != nil {
				delimiterText = def.delimiterFilter(match, &def)
			}
			// Read block content into lines.
			lines := []string{}
			if delimiterText != "" {
				lines = append(lines, delimiterText)
			}
			// Read content up to the closing delimiter.
			reader.Next()
			content := reader.ReadTo(def.closeMatch)
			if reader.Eof() && stringlist.StringList([]string{"code", "comment", "division", "quote"}).IndexOf(def.name) > -1 {
				options.ErrorCallback("unterminated " + def.name + " block: " + match[0])
			}
			reader.Next() // Skip closing delimiter.
			lines = append(lines, content...)
			// Calculate block expansion options.
			opts := def.options
			opts.Merge(blockattributes.Attrs.Options)
			// Translate block.
			if !opts.Skip {
				text := strings.Join(lines, "\n")
				if def.contentFilter != nil {
					text = def.contentFilter(text, match, opts)
				}
				opentag := def.openTag
				if def.name == "html" {
					text = blockattributes.Inject(text)
				} else {
					opentag = blockattributes.Inject(opentag)
				}
				if opts.Container {
					blockattributes.Attrs.Options.Container = false // Consume before recursing.
					text = ApiRender(text)
				} else {
					text = spans.ReplaceInline(text, opts)
				}
				closetag := def.closeTag
				if def.name == "division" && opentag == "<div>" {
					// Drop div tags if the opening div has no attributes.
					opentag = ""
					closetag = ""
				}
				text = opentag + text + closetag
				writer.Write(text)
				if text != "" && !reader.Eof() {
					// Add a trailing "\n" if we"ve written a non-blank line and there are more source lines left.
					writer.Write("\n")
				}
			}
			// Reset consumed Block Attributes expansion options.
			blockattributes.Attrs.Options = expansion.Options{}
			return true
		}
	}
	return false // No matching delimited block found.
}

// Return block definition or nil if not found.
func GetDefinition(name string) *Definition {
	for i, def := range defs {
		if def.name == name {
			return &defs[i]
		}
	}
	return nil
}

// Update existing named definition.
// Value syntax: <open-tag>|<close-tag> block-options
func SetDefinition(name string, value string) {
	def := GetDefinition(name)
	if def == nil {
		options.ErrorCallback("illegal delimited block name: " + name + ": |" + name + "|='" + value + "'")
		return
	}
	match := regexp.MustCompile(`^(?:(<[a-zA-Z].*>)\|(<[a-zA-Z/].*>))?(?:\s*)?([+-][ \w+-]+)?$`).FindStringSubmatch(strings.TrimSpace(value))
	if match == nil {
		options.ErrorCallback("illegal delimited block definition: |" + name + "|='" + value + "'")
		return
	}
	if strings.Contains(value, "|") {
		def.openTag = match[1]
		def.closeTag = match[2]
	}
	if match[3] != "" {
		def.options.Merge(expansion.Parse(match[3]))
	}
}

// delimiterFilter that returns opening delimiter line text from match group $1.
func delimiterTextFilter(match []string, _ *Definition) string {
	return match[1]
}

// delimiterFilter for code, division and quote blocks.
// Inject $2 into block class attribute, set close delimiter to $1.
func classInjectionFilter(match []string, def *Definition) string {
	if p1 := strings.TrimSpace(match[2]); p1 != "" {
		blockattributes.Attrs.Classes = p1
	}
	// closeMatch must be set at runtime so we correctly match closing delimiter
	def.closeMatch = regexp.MustCompile("^" + regexp.QuoteMeta(match[1]) + "$")
	return ""
}

// contentFilter for multi-line macro definitions.
func macroDefContentFilter(text string, match []string, opts expansion.Options) string {
	quote := string(match[0][len(match[0])-len(match[1])-1])                           // The leading macro value quote character.
	name := regexp.MustCompile(`^{([\w\-]+\??)}`).FindStringSubmatch(match[0])[1]      // Extract macro name from opening delimiter.
	text = regexp.MustCompile("("+quote+`) *\\\n`).ReplaceAllString(text, "$1\n")      // Unescape line-continuations.
	text = regexp.MustCompile("("+quote+` *[\\]+)\\\n`).ReplaceAllString(text, "$1\n") // Unescape escaped line-continuations.
	text = spans.ReplaceInline(text, opts)                                             // Expand macro invocations.
	macros.SetValue(name, text, quote)
	return ""
}
