package iotext

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/srackham/go-rimu/v11/internal/options"
)

/*
  Reader class.
*/
// Reader state.
type Reader struct {
	Lines []string
	Pos   int // Line index of current line.
}

// NewReader returns a new reader for text string.
func NewReader(text string) *Reader {
	if !utf8.ValidString(text) {
		options.ErrorCallback("invalid UTF-8 input")
		text = ""
	}
	r := new(Reader)
	text = strings.Replace(text, "\u0000", " ", -1) // Used internally by spans package.
	text = strings.Replace(text, "\u0001", " ", -1) // Used internally by spans package.
	text = strings.Replace(text, "\u0002", " ", -1) // Used internally by macros package.
	r.Lines = regexp.MustCompile(`\r\n|\r|\n`).Split(text, -1)
	return r
}

// Eof returns true is reader is at end of text.
func (r *Reader) Eof() bool {
	return r.Pos >= len(r.Lines)
}

// SetCursor sets the reader cursor position.
func (r *Reader) SetCursor(value string) {
	if r.Eof() {
		panic("unexpected eof")
	}
	r.Lines[r.Pos] = value
}

// Cursor returns the cursor position.
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
// If an EOF is encountered return all lines.
// Exit with the reader pointing to the line containing the matched line.
func (r *Reader) ReadTo(re *regexp.Regexp) (result []string) {
	result = []string{}
	var match []string
	for !r.Eof() {
		match = re.FindStringSubmatch(r.Cursor())
		if match != nil {
			if len(match) > 1 {
				result = append(result, match[1]) // $1
			}
			break
		}
		result = append(result, r.Cursor())
		r.Next()
	}
	return
}

// SkipBlankLines advances cursor to next non-blank line.
func (r *Reader) SkipBlankLines() {
	for !r.Eof() && strings.TrimSpace(r.Cursor()) == "" {
		r.Next()
	}
}

/*
  Writer class.
*/
// Writer is a container for lines of text.
type Writer struct {
	Buffer []string // Appending an array is faster than string concatenation.
}

// NewWriter return a new empty Writer.
func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Write(s string) {
	w.Buffer = append(w.Buffer, s)
}

func (w *Writer) String() string {
	return strings.Join(w.Buffer, "")
}
