package replacements

import (
	"regexp"
	"strings"

	"github.com/srackham/rimu-go/options"
)

func init() {
	Init()
}

type Definition struct {
	Match       *regexp.Regexp
	Replacement string
	Filter      func(match []string) string
}

var Defs []Definition // Mutable definitions initialized by DEFAULT_DEFS.

var DEFAULT_DEFS = []Definition{
	// Begin match with \\? to allow the replacement to be escaped.
	// Global flag must be set on match re's so that the RegExp lastIndex property is set.
	// Replacements and special characters are expanded in replacement groups ($1..).
	// Replacement order is important.

	// DEPRECATED as of 3.4.0.
	// Anchor: <<#id>>
	{
		Match:       regexp.MustCompile(`\\?<<#([a-zA-Z][\w\-]*)>>`),
		Replacement: `<span id="$1"></span>`,
	},

	// Image: <image:src|alt>
	// src = $1, alt = $2
	{
		Match:       regexp.MustCompile(`(?m)\\?<image:([^\s|]+)\|(.*?)>`),
		Replacement: `<img src="$1" alt="$2">`,
	},

	// Image: <image:src>
	// src = $1, alt = $1
	{
		Match:       regexp.MustCompile(`\\?<image:([^\s|]+?)>`),
		Replacement: `<img src="$1" alt="$1">`,
	},

	// Image: ![alt](url)
	// alt = $1, url = $2
	{
		Match:       regexp.MustCompile(`\\?!\[([^[]*?)]\((\S+?)\)`),
		Replacement: `<img src="$2" alt="$1">`,
	},

	// Email: <address|caption>
	// address = $1, caption = $2
	{
		Match:       regexp.MustCompile(`(?m)\\?<(\S+@[\w.\-]+)\|(.+?)>`),
		Replacement: `<a href="mailto:$1">$$2</a>`,
	},

	// Email: <address>
	// address = $1, caption = $1
	{
		Match:       regexp.MustCompile(`\\?<(\S+@[\w.\-]+)>`),
		Replacement: `<a href="mailto:$1">$1</a>`,
	},

	// Link: [caption](url)
	// caption = $1, url = $2
	{
		Match:       regexp.MustCompile(`\\?\[([^[]*?)]\((\S+?)\)`),
		Replacement: `<a href="$2">$$1</a>`,
	},

	// Link: <url|caption>
	// url = $1, caption = $2
	{
		Match:       regexp.MustCompile(`(?m)\\?<(\S+?)\|(.*?)>`),
		Replacement: `<a href="$1">$$2</a>`,
	},

	// HTML inline tags.
	// Match HTML comment or HTML tag.
	// $1 = tag, $2 = tag name
	{
		Match:       regexp.MustCompile(`(?i)\\?(<!--(?:[^<>&]*)?-->|<\/?([a-z][a-z0-9]*)(?:\s+[^<>&]+)?>)`),
		Replacement: "",
		Filter: func(match []string) string {
			return options.HtmlSafeModeFilter(match[1]) // Matched HTML comment or inline tag.
		},
	},

	// Link: <url>
	// url = $1
	{
		Match:       regexp.MustCompile(`\\?<([^|\s]+?)>`),
		Replacement: `<a href="$1">$1</a>`,
	},

	// Auto-encode (most) raw HTTP URLs as links.
	{
		Match:       regexp.MustCompile(`\\?((?:http|https):\/\/[^\s"']*[A-Za-z0-9/#])`),
		Replacement: `<a href="$1">$1</a>`,
	},

	// Character entity.
	{
		Match:       regexp.MustCompile(`\\?(&[\w#][\w]+;)`),
		Replacement: "",
		Filter: func(match []string) string {
			return match[1] // Pass the entity through verbatim.
		},
	},

	// Line-break (space followed by \ at end of line).
	{
		Match:       regexp.MustCompile(`[\\ ]\\(\n|$)`),
		Replacement: `<br>$1`,
	},

	// This hack ensures backslashes immediately preceding closing code quotes are rendered
	// verbatim (Markdown behaviour).
	// Works by finding escaped closing code quotes and replacing the backslash and the character
	// preceding the closing quote with itself.
	// NOTE: match differs from rimu-js and rimu-kt because regxp does not support (?=X) look-ahead.
	{
		Match:       regexp.MustCompile(`(\S\\` + "`" + `)`),
		Replacement: `$1`,
	},

	// This hack ensures underscores within words rendered verbatim and are not treated as
	// underscore emphasis quotes (GFM behaviour).
	// NOTE: match differs from rimu-js and rimu-kt because regxp does not support (?=X) look-ahead.
	{
		Match:       regexp.MustCompile(`([a-zA-Z0-9]_[a-zA-Z0-9])`),
		Replacement: `$1`,
	},
}

// Reset definitions to defaults.
func Init() {
	Defs = make([]Definition, len(DEFAULT_DEFS))
	for i, def := range DEFAULT_DEFS {
		Defs[i] = def
	}
}

// Update existing or add new replacement definition.
func SetDefinition(pattern string, flags string, replacement string) {

	if strings.Contains(flags, "i") {
		pattern = `(?i)` + pattern
	}
	if strings.Contains(flags, "m") {
		pattern = `(?m)` + pattern
	}
	for i, def := range Defs {
		if def.Match.String() == pattern {
			// Update existing definition.
			Defs[i].Replacement = replacement
			return
		}

	}
	// Append new definition to end of defs list (custom definitons have lower precedence).
	Defs = append(Defs, Definition{Match: regexp.MustCompile(pattern), Replacement: replacement})
}
