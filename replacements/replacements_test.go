package replacements

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init()
	assert.Equal(t, len(DEFAULT_DEFS), len(defs))
	assert.NotEqual(t, DEFAULT_DEFS, defs)
}

func TestSetDefinition(t *testing.T) {
	Init()
	SetDefinition(`\\?<image:([^\s|]+?)>`, "", "foo")
	assert.Equal(t, len(DEFAULT_DEFS), len(defs))
	SetDefinition(`bar`, "mi", "foo")
	assert.Equal(t, len(DEFAULT_DEFS)+1, len(defs))
	assert.Equal(t, defs[len(defs)-1].match.String(), "(?m)(?i)bar")
}
