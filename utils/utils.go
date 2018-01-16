package utils

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/srackham/rimu-go/macros"
	"github.com/srackham/rimu-go/spans"
	"github.com/srackham/rimu-go/utils/re"
)

// Processing priority (highest to lowest): container, skip, spans and specials.
// If spans is true then both spans and specials are processed.
// They are assumed false if they are not explicitly defined.
// If a custom filter is specified their use depends on the filter.
type ExpansionOptions struct {
	Macros    bool
	Container bool
	Skip      bool
	Spans     bool // Span substitution also expands special characters.
	Specials  bool
}

// http://stackoverflow.com/questions/3561493/is-there-a-regexp-escape-function-in-javascript
func EscapeRegExp(s string) string {
	return regexp.QuoteMeta(s)
}

func ReplaceSpecialChars(s string) string {
	result := strings.Replace(s, "&", "&amp;", -1)
	result = strings.Replace(result, ">", "&gt;", -1)
	result = strings.Replace(result, "<", "&lt;", -1)
	return result
}

// Replace pattern "$1" or "$$1", "$2" or "$$2"... in `replacement` with corresponding match groups
// from `match`. If pattern starts with one "$" character add specials to `expansionOptions`,
// if it starts with two "$" characters add spans to `expansionOptions`.
func ReplaceMatch(match []string, replacement string, expansionOptions ExpansionOptions) string {
	return re.ReplaceAllStringSubmatchFunc(regexp.MustCompile(`(\${1,2})(\d)`), replacement, func(arguments []string) string {
		// Replace $1, $2 ... with corresponding match groups.
		switch {
		case arguments[1] == "$$":
			expansionOptions.Spans = true
		default:
			expansionOptions.Specials = true
		}
		i, _ := strconv.ParseInt(arguments[2], 10, 64) // match group number.
		text := match[i]                               // match group text.
		return ReplaceInline(text, expansionOptions)
	})
}

// Replace the inline elements specified in options in text and return the result.
func ReplaceInline(text string, expansionOptions ExpansionOptions) string {
	if expansionOptions.Macros {
		text = macros.Render(text, false)
	}
	// Spans also expand special characters.
	switch {
	case expansionOptions.Spans:
		text = spans.Render(text)
	case expansionOptions.Specials:
		text = ReplaceSpecialChars(text)
	}
	return text
}
