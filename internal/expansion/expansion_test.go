package expansion

import (
	"testing"

	"github.com/srackham/go-rimu/v11/internal/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		opts string
		want Options
	}{
		{"", Options{}},
		{"+skip +macros +container +specials +spans", Options{true, true, true, true, true, true, true, true, true, true}},
	}
	for _, tt := range tests {
		got := Parse(tt.opts)
		assert.Equal(t, tt.want, got)
	}
}
