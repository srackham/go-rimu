package replacements

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.Equal(t, len(DEFAULT_DEFS), len(Defs))
	assert.NotEqual(t, DEFAULT_DEFS, Defs)
}

func TestSetDefinition(t *testing.T) {
	Init()
	SetDefinition(`\\?<image:([^\s|]+?)>`, "", "foo")
	assert.Equal(t, len(DEFAULT_DEFS), len(Defs))
	SetDefinition(`bar`, "mi", "foo")
	assert.Equal(t, len(DEFAULT_DEFS)+1, len(Defs))
	assert.Equal(t, Defs[len(Defs)-1].Match.String(), "(?m)(?i)bar")
}
