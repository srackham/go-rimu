package expansion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		opts string
		want ExpansionOptions
	}{
		{"", ExpansionOptions{}},
		{"+skip +macros +container +specials +spans", ExpansionOptions{true, true, true, true, true, true, true, true, true, true}},
	}
	for _, tt := range tests {
		got := Parse(tt.opts)
		assert.Equal(t, tt.want, got)
	}
}

func TestReplaceMatch(t *testing.T) {
	tests := []struct {
		match            []string
		replacement      string
		expansionOptions ExpansionOptions
		want             string
	}{
		{nil, "", ExpansionOptions{}, ""},
		{[]string{"foo bar", "foo", "bar"}, "$2 $1", ExpansionOptions{}, "bar foo"},
	}
	for _, tt := range tests {
		got := ReplaceMatch(tt.match, tt.replacement, tt.expansionOptions)
		assert.Equal(t, tt.want, got)
	}
}

/*
func TestReplaceInline(t *testing.T) {
	tests := []struct {
		text string
		opts ExpansionOptions
		want string
	}{
		{"", ExpansionOptions{}, ""},
		{"<>& _foo_", ExpansionOptions{}, "<>& _foo_"},
		{"<>& _foo_", ExpansionOptions{Specials: true}, "&lt;&gt;&amp; _foo_"},
		{"<>& _foo_", ExpansionOptions{Spans: true}, "&lt;&gt;&amp; <em>foo</em>"},
	}
	for _, tt := range tests {
		got := ReplaceInline(tt.text, tt.opts)
		assert.Equal(t, tt.want, got)
	}
}
*/
