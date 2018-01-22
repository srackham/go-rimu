package expansion

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/srackham/rimu-go/options"
	"github.com/srackham/rimu-go/utils"
	"github.com/srackham/rimu-go/utils/re"
)

// macros and spans package dependency injections.
var MacrosRender func(text string, silent bool) string
var SpansRender func(text string) string

// Processing priority (highest to lowest): container, skip, spans and specials.
// If spans is true then both spans and specials are processed.
// They are assumed false if they are not explicitly defined.
// If a custom filter is specified their use depends on the filter.
type ExpansionOptions struct {
	Container bool
	Macros    bool
	Skip      bool
	Spans     bool // Span substitution also expands special characters.
	Specials  bool
	// xxxSet specify if the Xxx field is included blockattributes.Options and only in the Merge method.
	containerSet bool
	macrosSet    bool
	skipSet      bool
	spansSet     bool
	specialsSet  bool
}

// Merge copies expansion options that are set from from to to.
func (to *ExpansionOptions) Merge(from ExpansionOptions) {
	if from.containerSet {
		to.Container = from.Container
	}
	if from.macrosSet {
		to.Macros = from.Macros
	}
	if from.skipSet {
		to.Skip = from.Skip
	}
	if from.spansSet {
		to.Spans = from.Spans
	}
	if from.specialsSet {
		to.Specials = from.Specials
	}
}

// Parse block-options string and return ExpansionOptions.
func Parse(optionsString string) ExpansionOptions {
	result := ExpansionOptions{}
	if optionsString != "" {
		opts := regexp.MustCompile(`\s+`).Split(strings.Trim(optionsString, " "), -1)
		for _, opt := range opts {
			if options.IsSafeModeNz() && opt == "-specials" {
				options.ErrorCallback("-specials block option not valid in safeMode")
				continue
			}
			if regexp.MustCompile(`^[+-](macros|spans|specials|container|skip)$`).MatchString(opt) {
				value := opt[0] == '+'
				switch opt[1:] {
				case "container":
					result.Container = value
					result.containerSet = true
				case "macros":
					result.Macros = value
					result.macrosSet = true
				case "skip":
					result.Skip = value
					result.skipSet = true
				case "specials":
					result.Specials = value
					result.specialsSet = true
				case "spans":
					result.Spans = value
					result.spansSet = true
				}
			} else {
				options.ErrorCallback("illegal block option: " + opt)
			}
		}
	}
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
		text = MacrosRender(text, false)
	}
	// Spans also expand special characters.
	switch {
	case expansionOptions.Spans:
		text = SpansRender(text)
	case expansionOptions.Specials:
		text = utils.ReplaceSpecialChars(text)
	}
	return text
}
