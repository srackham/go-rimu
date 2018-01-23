package spans

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/srackham/rimu-go/expansion"
	"github.com/srackham/rimu-go/quotes"
	"github.com/srackham/rimu-go/replacements"
	"github.com/srackham/rimu-go/utils"
	"github.com/srackham/rimu-go/utils/re"
)

// macros and spans package dependency injections.
var MacrosRender func(text string, silent bool) string

type fragment struct {
	text     string
	done     bool
	verbatim string // Replacements text rendered verbatim.
}

func Render(source string) string {
	result := preReplacements(source)
	frags := []fragment{{text: result, done: false}}
	frags = fragQuotes(frags)
	frags = fragSpecials(frags)
	result = defrag(frags)
	return postReplacements(result)
}

// Converts fragments to a string.
func defrag(frags []fragment) string {
	result := ""
	for _, frag := range frags {
		result += frag.text
	}
	return result
}

// Fragment quotes in all fragments and return resulting fragments array.
func fragQuotes(frags []fragment) []fragment {
	result := []fragment{}
	for _, frag := range frags {
		result = append(result, fragQuote(frag)...)
	}
	// Strip backlash from escaped quotes in non-done fragments.
	for _, frag := range frags {
		if !frag.done {
			frag.text = quotes.Unescape(frag.text)
		}
	}
	return result
}

// Fragment quotes in a single fragment and return resulting fragments array.
func fragQuote(frag fragment) (result []fragment) {
	if frag.done {
		return []fragment{frag}
	}
	match := quotes.Find(frag.text)
	if match == nil {
		return []fragment{frag}
	}
	quote := frag.text[match[2]:match[3]]
	quoted := frag.text[match[4]:match[5]]
	startIndex := match[0]
	endIndex := match[1]
	// Check for same closing quote one character further to the right.
	for endIndex < len(frag.text) && frag.text[endIndex] == quote[0] {
		// Move to closing quote one character to right.
		quoted += string(quote[0])
		endIndex += 1
	}
	// Arrive here if we have a matched quote.
	// The quote splits the input fragment into 5 or more output fragments:
	// Text before the quote, left quote tag, quoted text, right quote tag and text after the quote.
	def := quotes.GetDefinition(quote)
	before := frag.text[:startIndex]
	after := frag.text[endIndex:]
	result = append(result, fragment{text: before, done: false})
	result = append(result, fragment{text: def.OpenTag, done: true})
	if !def.Spans {
		// Spans are disabled so render the quoted text verbatim.
		quoted = utils.ReplaceSpecialChars(quoted)
		quoted = strings.Replace(quoted, string('\u0000'), string('\u0001'), -1) // Substitute verbatim replacement placeholder.
		result = append(result, fragment{text: quoted, done: true})
	} else {
		// Recursively process the quoted text.
		result = append(result, fragQuote(fragment{text: quoted, done: false})...)
	}
	result = append(result, fragment{text: def.CloseTag, done: true})
	// Recursively process the following text.
	result = append(result, fragQuote(fragment{text: after, done: false})...)
	return result
}

// Stores placeholder replacement fragments saved by `preReplacements()` and restored by `postReplacements()`.
var savedReplacements []fragment

// Return text with replacements replaced with placeholders (see `postReplacements()`).
func preReplacements(text string) (result string) {
	savedReplacements = nil
	frags := fragReplacements([]fragment{{text: text, done: false}})
	// Reassemble text with replacement placeholders.
	for _, frag := range frags {
		if frag.done {
			savedReplacements = append(savedReplacements, frag) // Save replaced text.
			result += string('\u0000')                          // Placeholder for replaced text.
		} else {
			result += frag.text
		}
	}
	return result
}

// Replace replacements placeholders with replacements text from savedReplacements[].
func postReplacements(text string) string {
	return regexp.MustCompile(`[\x{0000}\x{0001}]`).ReplaceAllStringFunc(text, func(match string) string {
		var frag fragment
		frag, savedReplacements = savedReplacements[0], savedReplacements[1:] // Remove frag from start of list.
		if match == string('\u0000') {
			return frag.text
		} else {
			return utils.ReplaceSpecialChars(frag.verbatim)
		}

	})
}

// Fragment replacements in all fragments and return resulting fragments array.
func fragReplacements(frags []fragment) (result []fragment) {
	result = frags
	for _, def := range replacements.Defs {
		var tmp []fragment
		for _, frag := range result {
			tmp = append(tmp, fragReplacement(frag, def)...)
		}
		result = tmp
	}
	return result
}

// Fragment replacements in a single fragment for a single replacement definition.
// Return resulting fragments list.
func fragReplacement(frag fragment, def replacements.Definition) (result []fragment) {
	if frag.done {
		return []fragment{frag}
	}
	match := def.Match.FindStringIndex(frag.text)
	if match == nil {
		return []fragment{frag}
	}
	// Arrive here if we have a matched replacement.
	// The replacement splits the input fragment into 3 output fragments:
	// Text before the replacement, replaced text and text after the replacement.
	before := frag.text[:match[0]]
	matched := frag.text[match[0]:match[1]]
	after := frag.text[match[1]:]
	result = append(result, fragment{text: before, done: false})
	var replacement string
	if strings.HasPrefix(matched, "\\") {
		// Remove leading backslash.
		replacement = utils.ReplaceSpecialChars(matched[1:])
	} else {
		submatches := def.Match.FindStringSubmatch(matched)
		if def.Filter == nil {
			replacement = ReplaceMatch(submatches, def.Replacement, expansion.Options{})
		} else {
			replacement = def.Filter(submatches)
		}
	}
	result = append(result, fragment{text: replacement, done: true, verbatim: matched})
	// Recursively process the remaining text.
	result = append(result, fragReplacement(fragment{text: after, done: false}, def)...)
	return result
}

func fragSpecials(frags []fragment) (result []fragment) {
	// Replace special characters in all non-done fragments.
	result = make([]fragment, len(frags))
	for i, frag := range frags {
		if !frag.done {
			frag.text = utils.ReplaceSpecialChars(frag.text)
		}
		result[i] = frag
	}
	return result
}

// Replace pattern "$1" or "$$1", "$2" or "$$2"... in `replacement` with corresponding match groups
// from `match`. If pattern starts with one "$" character add specials to `opts`,
// if it starts with two "$" characters add spans to `opts`.
func ReplaceMatch(match []string, replacement string, opts expansion.Options) string {
	return re.ReplaceAllStringSubmatchFunc(regexp.MustCompile(`(\${1,2})(\d)`), replacement, func(arguments []string) string {
		// Replace $1, $2 ... with corresponding match groups.
		switch {
		case arguments[1] == "$$":
			opts.Spans = true
		default:
			opts.Specials = true
		}
		i, _ := strconv.ParseInt(arguments[2], 10, 64) // match group number.
		text := match[i]                               // match group text.
		return ReplaceInline(text, opts)
	})
}

// Replace the inline elements specified in options in text and return the result.
func ReplaceInline(text string, opts expansion.Options) string {
	if opts.Macros {
		text = MacrosRender(text, false)
	}
	// Spans also expand special characters.
	switch {
	case opts.Spans:
		text = Render(text)
	case opts.Specials:
		text = utils.ReplaceSpecialChars(text)
	}
	return text
}
