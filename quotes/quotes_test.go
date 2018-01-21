package quotes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.Equal(t, len(DEFAULT_DEFS), len(defs))
	assert.NotEqual(t, DEFAULT_DEFS, defs)
}

func TestGetDefinition(t *testing.T) {
	Init()
	assert.NotNil(t, GetDefinition("*"))
	assert.Nil(t, GetDefinition("X"))
}

func TestSetDefinition(t *testing.T) {
	Init()

	SetDefinition(Definition{
		Quote:    "*",
		OpenTag:  "<strong>",
		CloseTag: "</strong>",
		Spans:    true,
	})
	assert.Equal(t, len(DEFAULT_DEFS), len(defs))
	def := GetDefinition("*")
	assert.Equal(t, "<strong>", def.OpenTag)

	SetDefinition(Definition{
		Quote:    "x",
		OpenTag:  "<del>",
		CloseTag: "</del>",
		Spans:    true,
	})
	assert.Equal(t, len(DEFAULT_DEFS)+1, len(defs))
	def = GetDefinition("x")
	assert.Equal(t, "<del>", def.OpenTag)
	assert.Equal(t, "<del>", defs[len(defs)-1].OpenTag)

	SetDefinition(Definition{
		Quote:    "xx",
		OpenTag:  "<u>",
		CloseTag: "</u>",
		Spans:    true,
	})
	assert.Equal(t, len(DEFAULT_DEFS)+2, len(defs))
	def = GetDefinition("xx")
	assert.Equal(t, "<u>", def.OpenTag)
	assert.Equal(t, "<u>", defs[0].OpenTag)
}

func TestUnescape(t *testing.T) {
	Init()
	assert.Equal(t, `* ~~ \x`, Unescape(`\* \~~ \x`))
}

func TestFind(t *testing.T) {
	tests := []struct {
		text string
		want []int
	}{
		{``, nil},
		{`*foo*`, []int{0, 5, 0, 1, 1, 4}},
		{`\*foo*`, nil},
		{`_bar_ *foo*`, []int{0, 5, 0, 1, 1, 4}},
		{`\_bar_ *foo*`, []int{7, 12, 7, 8, 8, 11}},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, Find(tt.text))
	}
}
