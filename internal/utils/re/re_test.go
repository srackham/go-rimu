package re

import (
	"regexp"
	"testing"
)

func TestReplaceAllStringSubmatchFunc(t *testing.T) {
	tests := []struct {
		re   string
		src  string
		repl func(match []string) string
		want string
	}{
		{``, "",
			func(match []string) string {
				return "foo"
			},
			"foo",
		},
		{`Z`, "x",
			func(match []string) string {
				return "foo"
			},
			"x",
		},
		{`([a-z]+)(\d+)`, "xyz123",
			func(match []string) string {
				return match[0] + " " + match[2] + " " + match[1]
			},
			"xyz123 123 xyz",
		},
		{`([a-z]+)(\d+)`, "xyz123 ab98",
			func(match []string) string {
				return match[0] + " " + match[2] + " " + match[1]
			},
			"xyz123 123 xyz ab98 98 ab",
		},
	}
	for _, tt := range tests {
		re := regexp.MustCompile(tt.re)
		got := ReplaceAllStringSubmatchFunc(re, tt.src, tt.repl, -1)
		if got != tt.want {
			t.Errorf("TestReplaceAllStringSubmatchFunc(%q, %q) == %q, want %q", tt.re, tt.src, got, tt.want)
		}
	}
}
