package spans

import (
	"regexp"
	"strings"

	"github.com/srackham/rimu-go/expansion"
	"github.com/srackham/rimu-go/quotes"
	"github.com/srackham/rimu-go/replacements"
	"github.com/srackham/rimu-go/utils"
)

func init() {
	expansion.SpansRender = Render
}

type fragment struct {
	text     string
	done     bool
	verbatim string // Replacements text rendered verbatim.
}

func Render(source string) string {
	result := preReplacements(source)
	fragments := []fragment{{text: result, done: false}}
	fragments = fragQuotes(fragments)
	fragSpecials(fragments)
	result = defrag(fragments)
	return postReplacements(result)
}

// Converts fragments to a string.
func defrag(fragments []fragment) string {
	result := ""
	for _, f := range fragments {
		result += f.text
	}
	return result
}

// Fragment quotes in all fragments and return resulting fragments array.
func fragQuotes(fragments []fragment) []fragment {
	result := []fragment{}
	for _, f := range fragments {
		result = append(result, fragQuote(f)...)
	}
	// Strip backlash from escaped quotes in non-done fragments.
	for _, f := range fragments {
		if !f.done {
			f.text = quotes.Unescape(f.text)
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
	for frag.text[endIndex] == quote[0] {
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
	fragments := fragReplacements([]fragment{{text: text, done: false}})
	// Reassemble text with replacement placeholders.
	for _, frag := range fragments {
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
	return regexp.MustCompile(`[\u0000\u0001]`).ReplaceAllStringFunc(text, func(match string) string {
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
func fragReplacements(fragments []fragment) (result []fragment) {
	result = fragments
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
			replacement = expansion.ReplaceMatch(submatches, def.Replacement, expansion.ExpansionOptions{})
		} else {
			replacement = def.Filter(submatches)
		}
	}
	result = append(result, fragment{text: replacement, done: true, verbatim: matched})
	// Recursively process the remaining text.
	result = append(result, fragReplacement(fragment{text: after, done: false}, def)...)
	return result
}

func fragSpecials(fragments []fragment) {
	// Replace special characters in all non-done fragments.
	for _, frag := range fragments {
		if !frag.done {
			frag.text = utils.ReplaceSpecialChars(frag.text)
		}
	}
}
