package utils

import (
	"regexp"
	"testing"
)

func TestReplaceAllStringSubmatchFunc(t *testing.T) {
	cases := []struct {
		re       string
		src      string
		expected string
		repl     func(match []string) string
	}{
		{``, "", "foo",
			func(match []string) string {
				return "foo"
			},
		},
		{`Z`, "x", "x",
			func(match []string) string {
				return "foo"
			},
		},
		{`([a-z]+)(\d+)`, "xyz123", "xyz123 123 xyz",
			func(match []string) string {
				return match[0] + " " + match[2] + " " + match[1]
			},
		},
	}
	for _, c := range cases {
		re := regexp.MustCompile(c.re)
		got := ReplaceAllStringSubmatchFunc(re, c.src, c.repl)
		if got != c.expected {
			t.Errorf("TestReplaceAllStringSubmatchFunc(%q, %q) == %q, expected %q", c.re, c.src, got, c.expected)
		}
	}
}
