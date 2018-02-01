package quotes

import (
	"regexp"
	"strings"
)

func init() {
	Init()
}

type Definition struct {
	Quote    string // Single quote character.
	OpenTag  string
	CloseTag string
	Spans    bool // Allow span elements inside quotes.
	re       *regexp.Regexp
}

var defs []Definition // Mutable definitions initialized by DEFAULT_DEFS.

var DEFAULT_DEFS = []Definition{
	{
		Quote:    "**",
		OpenTag:  "<strong>",
		CloseTag: "</strong>",
		Spans:    true,
	},
	{
		Quote:    "*",
		OpenTag:  "<em>",
		CloseTag: "</em>",
		Spans:    true,
	},
	{
		Quote:    "__",
		OpenTag:  "<strong>",
		CloseTag: "</strong>",
		Spans:    true,
	},
	{
		Quote:    "_",
		OpenTag:  "<em>",
		CloseTag: "</em>",
		Spans:    true,
	},
	{
		Quote:    "``",
		OpenTag:  "<code>",
		CloseTag: "</code>",
		Spans:    false,
	},
	{
		Quote:    "`",
		OpenTag:  "<code>",
		CloseTag: "</code>",
		Spans:    false,
	},
	{
		Quote:    "~~",
		OpenTag:  "<del>",
		CloseTag: "</del>",
		Spans:    true,
	},
}

// Reset definitions to defaults.
func Init() {
	defs = make([]Definition, len(DEFAULT_DEFS))
	for i, def := range DEFAULT_DEFS {
		defs[i] = def
	}
	initRegExps()
}

// Synthesise re's to find quotes.
func initRegExps() {
	// $1 is quote character(s), $2 is quoted text.
	// Quoted text cannot begin or end with whitespace.
	// Quoted can span multiple lines.
	// Quoted text cannot end with a backslash.
	for i, def := range defs {
		defs[i].re = regexp.MustCompile(`\\?(` + regexp.QuoteMeta(def.Quote) + `)([^\s\\]|\S[\s\S]*?[^\s\\])` + regexp.QuoteMeta(def.Quote))
	}
}

// Return the quote definition corresponding to 'quote', return nil if not found.
func GetDefinition(quote string) *Definition {
	for _, def := range defs {
		if def.Quote == quote {
			return &def
		}
	}
	return nil
}

// Update existing or add new quote definition.
func SetDefinition(def Definition) {
	for i := range defs {
		if defs[i].Quote == def.Quote {
			// Update existing definition.
			defs[i].OpenTag = def.OpenTag
			defs[i].CloseTag = def.CloseTag
			defs[i].Spans = def.Spans
			return
		}
	}
	// Double-quote definitions are prepended to the array so they are matched
	// before single-quote definitions (which are appended to the array).
	if len(def.Quote) == 2 {
		defs = append([]Definition{def}, defs...)
	} else {
		defs = append(defs, def)
	}
	initRegExps()
}

// Strip backslashes from quote characters.
func Unescape(s string) string {
	for _, def := range defs {
		s = strings.Replace(s, "\\"+def.Quote, def.Quote, -1)
	}
	return s
}

// Find looks for the first quote in `text` starting from `start`.
// Quotes prefixed with a backslash are ignored.
// Returns slice holding thre index pairs identifying:
// - The entire match: s[loc[0]:loc[1]]
// - The left quote    s[loc[2]:loc[3]]
// - The quoted text   s[loc[4]:loc[5]]
// Returns nil if not found.
func Find(text string) []int {
	// This function is necessary because Go regexp does not support backreferences.
	var match []int
	for _, def := range defs {
		allMatch := def.re.FindAllStringSubmatchIndex(text, -1)
		if allMatch == nil {
			continue
		}
		for _, nextMatch := range allMatch {
			if match == nil || nextMatch[0] < match[0] {
				match = nextMatch
			}
		}
	}
	return match
}
