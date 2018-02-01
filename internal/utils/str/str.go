package str

import "strings"

func ReplaceSpecialChars(s string) (result string) {
	result = strings.Replace(s, "&", "&amp;", -1)
	result = strings.Replace(result, ">", "&gt;", -1)
	result = strings.Replace(result, "<", "&lt;", -1)
	return
}

// TrimQuotes removes leading and trailing quote and returns result.
// If string is quoted then return it unchanged.
func TrimQuotes(s string, quote string) string {
	if len(s) >= 2*len(quote) && strings.HasPrefix(s, quote) && strings.HasSuffix(s, quote) {
		return strings.TrimPrefix(strings.TrimSuffix(s, quote), quote)
	}
	return s
}
