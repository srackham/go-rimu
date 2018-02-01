package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimQuotes(t *testing.T) {
	tests := []struct {
		s     string
		quote string
		want  string
	}{
		{`"foo"`, `"`, `foo`},
		{`""foo"`, `"`, `"foo`},
		{`""foo""`, `""`, `foo`},
		{`"foo`, `"`, `"foo`},
		{`foo"`, `"`, `foo"`},
	}
	for _, tt := range tests {
		got := TrimQuotes(tt.s, tt.quote)
		assert.Equal(t, tt.want, got)
	}
}
