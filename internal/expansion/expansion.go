package expansion

import (
	"regexp"
	"strings"

	"github.com/srackham/go-rimu/internal/options"
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
	// xxxSet specify if the Xxx field is included blockattributes.Options and only in the Merge method.
	containerSet bool
	macrosSet    bool
	skipSet      bool
	spansSet     bool
	specialsSet  bool
}

// Merge copies expansion options that are set from from to to.
func (to *Options) Merge(from Options) {
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
	return
}
