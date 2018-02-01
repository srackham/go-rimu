package utils

import (
	"strings"
)

func ReplaceSpecialChars(s string) (result string) {
	result = strings.Replace(s, "&", "&amp;", -1)
	result = strings.Replace(result, ">", "&gt;", -1)
	result = strings.Replace(result, "<", "&lt;", -1)
	return
}
