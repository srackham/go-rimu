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
	Lines []string
	Pos   int // Line index of current line.
}

// NewReader TODO
func NewReader(text string) *Reader {
	r := new(Reader)
	text = strings.Replace(text, "\u0000", " ", -1) // Used internally by macros and spans packages.
	text = strings.Replace(text, "\u0001", " ", -1) // Used internally by macros and spans packages.
	r.Lines = regexp.MustCompile(`\r\n|\r|\n`).Split(text, -1)
	return r
}

// Eof TODO
func (r *Reader) Eof() bool {
	return r.Pos >= len(r.Lines)
}

// SetCursor TODO
func (r *Reader) SetCursor(value string) {
	if r.Eof() {
		panic("unexpected eof")
	}
	r.Lines[r.Pos] = value
}

// Cursor TODO
func (r *Reader) Cursor() string {
	if r.Eof() {
		panic("unexpected eof")
	}
	return r.Lines[r.Pos]
}

// Next moves cursor to next input line.
func (r *Reader) Next() {
	if !r.Eof() {
		r.Pos++
	}
}

// ReadTo reads to the first line matching the re.
// Return the array of lines preceding the match plus a line containing
// the $1 match group (if it exists).
// Return nil if an EOF is encountered.
// Exit with the reader pointing to the line following the match.
func (r *Reader) ReadTo(re *regexp.Regexp) (result []string) {
	result = []string{}
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
		return
	}
	return nil
}

// SkipBlankLines TODO
func (r *Reader) SkipBlankLines() {
	for !r.Eof() && strings.TrimSpace(r.Cursor()) == "" {
		r.Next()
	}
}

/*
  Writer class.
*/
// Writer TODO
type Writer struct {
	Buffer []string // Appending an array is faster than string concatenation.
}

// NewWriter TODO
func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Write(s string) {
	w.Buffer = append(w.Buffer, s)
}

func (w *Writer) String() string {
	return strings.Join(w.Buffer, "")
}
