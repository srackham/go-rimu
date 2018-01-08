package iotext

import (
	"regexp"
	"strings"
)

// Reader TODO
type Reader struct {
	lines []string
	pos   int // Line index of current line.
}

// NewReader TODO
func NewReader(text string) *Reader {
	r := new(Reader)
	r.lines = regexp.MustCompile(`\r\n|\r|\n`).Split(text, -1)
	return r
}

// Eof TODO
func (r *Reader) Eof() bool {
	return r.pos >= len(r.lines)
}

// SetCursor TODO
func (r *Reader) SetCursor(value string) {
	if r.Eof() {
		panic("unexpected eof")
	}
	r.lines[r.pos] = value
}

// Cursor TODO
func (r *Reader) Cursor() string {
	if r.Eof() {
		panic("unexpected eof")
	}
	return r.lines[r.pos]
}

// Next moves cursor to next input line.
func (r *Reader) Next() {
	if !r.Eof() {
		r.pos++
	}
}

// SkipBlankLines TODO
func (r *Reader) SkipBlankLines() {
	for !r.Eof() && r.Cursor() == "" {
		r.Next()
	}
}

// Writer TODO
type Writer struct {
	buffer []string // Appending an array is faster than string concatenation.
}

// NewWriter TODO
func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Write(s string) {
	w.buffer = append(w.buffer, s)
}

func (w *Writer) String() string {
	return strings.Join(w.buffer, "")
}
