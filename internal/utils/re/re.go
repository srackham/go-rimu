package re

import (
	"regexp"
)

// ReplaceAllStringSubmatchFunc returns a string with all re matches replaced by the repl
// callback function. repl is passed a slice containing the matched text (match[0]) and
// any submatches (match[1]...) (c.f. Regexp.ReplaceAllStringFunc)
// Unmatched groups return a blank string.
// if n >= 0, the function returns at most n matches/submatches.
// Code from: http://elliot.land/post/go-replace-string-with-regular-expression-callback
// See also:https://github.com/golang/go/issues/5690
func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, src string, repl func(match []string) string, n int) (result string) {
	lastIndex := 0
	for _, v := range re.FindAllStringSubmatchIndex(src, n) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			if v[i] == -1 {
				// Blank string for unmatch groups.
				groups = append(groups, "")
			} else {
				groups = append(groups, src[v[i]:v[i+1]])
			}
		}
		result += src[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}
	return result + src[lastIndex:]
}
