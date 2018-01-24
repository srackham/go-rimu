package macros

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValues(t *testing.T) {
	Init()
	assert.Equal(t, 2, len(defs))
	got, err := Value("--")
	assert.Nil(t, err)
	assert.Equal(t, "", got)

	SetValue("foo", "bar", "'")
	assert.Equal(t, 3, len(defs))
	got, err = Value("foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", got)

	SetValue("foo?", "baz", "'")
	assert.Equal(t, 3, len(defs))
	got, err = Value("foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", got)

	SetValue("foo", "baz", "'")
	assert.Equal(t, 3, len(defs))
	got, err = Value("foo")
	assert.Nil(t, err)
	assert.Equal(t, "baz", got)
}

func TestRender(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{"", ""},
		{"{--}{--header-ids}", ""},
	}
	for _, tt := range tests {
		got := Render(tt.text, false)
		assert.Equal(t, tt.want, got)
	}
}
