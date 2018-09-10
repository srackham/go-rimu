package expansion

import (
	"regexp"
	"strings"

	"github.com/srackham/go-rimu/v11/internal/options"
)

// Processing priority (highest to lowest): container, skip, spans and specials.
// If spans is true then both spans and specials are processed.
// They are assumed false if they are not explicitly defined.
// If a custom filter is specified their use depends on the filter.
type Options struct {
	Container bool
	Macros    bool
	Skip      bool
	Spans     bool // Span substitution also expands special characters.
	Specials  bool
	// xxxMerge specify if the Xxx field has been set.
	containerMerge bool
	macrosMerge    bool
	skipMerge      bool
	spansMerge     bool
	specialsMerge  bool
}

// Merge copies expansion options that are set from from to to.
func (to *Options) Merge(from Options) {
	if from.containerMerge {
		to.Container = from.Container
		to.containerMerge = true
	}
	if from.macrosMerge {
		to.Macros = from.Macros
		to.macrosMerge = true
	}
	if from.skipMerge {
		to.Skip = from.Skip
		to.skipMerge = true
	}
	if from.spansMerge {
		to.Spans = from.Spans
		to.spansMerge = true
	}
	if from.specialsMerge {
		to.Specials = from.Specials
		to.specialsMerge = true
	}
}

// Parse block-options string and return ExpansionOptions.
func Parse(optsString string) (result Options) {
	if optsString != "" {
		opts := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(optsString), -1)
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
					result.containerMerge = true
				case "macros":
					result.Macros = value
					result.macrosMerge = true
				case "skip":
					result.Skip = value
					result.skipMerge = true
				case "specials":
					result.Specials = value
					result.specialsMerge = true
				case "spans":
					result.Spans = value
					result.spansMerge = true
				}
			} else {
				options.ErrorCallback("illegal block option: " + opt)
			}
		}
	}
	return
}
