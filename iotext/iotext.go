package iotext

import (
	"regexp"
	"strings"
)

/*
  Reader class.
*/
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

// ReadTo reads to the first line matching the re.
// Return the array of lines preceding the match plus a line containing
// the $1 match group (if it exists).
// Return nil if an EOF is encountered.
// Exit with the reader pointing to the line following the match.
func (r *Reader) ReadTo(re *regexp.Regexp) []string {
	result := []string{}
	var match []string
	for !r.Eof() {
		match = re.FindStringSubmatch(r.Cursor())
		if match != nil {
			if len(match) > 1 {
				result = append(result, match[1]) // $1
			}
			r.Next()
			break
		}
		result = append(result, r.Cursor())
		r.Next()
	}
	// Blank line matches EOF.
	if match != nil || re.String() == "^$" && r.Eof() {
		return result
	}
	return nil
}

// SkipBlankLines TODO
func (r *Reader) SkipBlankLines() {
	for !r.Eof() && r.Cursor() == "" {
		r.Next()
	}
}

/*
  Writer class.
*/
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
