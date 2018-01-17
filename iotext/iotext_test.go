package iotext

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	reader := NewReader("")
	assert.Equal(t, false, reader.Eof())
	assert.Equal(t, 1, len(reader.Lines))
	assert.Equal(t, "", reader.Cursor())
	reader.Next()
	assert.Equal(t, true, reader.Eof())

	reader = NewReader("Hello\nWorld!")
	assert.Equal(t, 2, len(reader.Lines))
	assert.Equal(t, "Hello", reader.Cursor())
	reader.Next()
	assert.Equal(t, "World!", reader.Cursor())
	assert.Equal(t, false, reader.Eof())
	reader.Next()
	assert.Equal(t, true, reader.Eof())

	reader = NewReader("\n\nHello")
	assert.Equal(t, 3, len(reader.Lines))
	reader.SkipBlankLines()
	assert.Equal(t, "Hello", reader.Cursor())
	assert.Equal(t, false, reader.Eof())
	reader.Next()
	assert.Equal(t, true, reader.Eof())

	reader = NewReader("Hello\n*\nWorld!\nHello\n< Goodbye >")
	assert.Equal(t, 5, len(reader.Lines))
	lines := reader.ReadTo(regexp.MustCompile(`\*`))
	assert.Equal(t, 1, len(lines))
	assert.Equal(t, "Hello", lines[0])
	assert.Equal(t, false, reader.Eof())
	lines = reader.ReadTo(regexp.MustCompile(`^<(.*)>$`))
	assert.Equal(t, 3, len(lines))
	assert.Equal(t, " Goodbye ", lines[2])
	assert.Equal(t, true, reader.Eof())

	reader = NewReader("\n\nHello\nWorld!")
	assert.Equal(t, 4, len(reader.Lines))
	reader.SkipBlankLines()
	lines = reader.ReadTo(regexp.MustCompile(`^$`))
	assert.Equal(t, 2, len(lines))
	assert.Equal(t, "World!", lines[1])
	assert.Equal(t, true, reader.Eof())
	assert.Panics(t, func() { reader.Cursor() })
	assert.Panics(t, func() { reader.SetCursor("foo") })
}

func TestWriter(t *testing.T) {
	writer := NewWriter()
	writer.Write("Hello")
	assert.Equal(t, "Hello", writer.buffer[0])
	writer.Write("World!")
	assert.Equal(t, "World!", writer.buffer[1])
	assert.Equal(t, "HelloWorld!", writer.String())
}
