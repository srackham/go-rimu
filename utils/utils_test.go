package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceSpecialChars(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"<>&", "&lt;&gt;&amp;"},
	}
	for _, tt := range tests {
		got := ReplaceSpecialChars(tt.in)
		assert.Equal(t, tt.want, got)
	}
}
