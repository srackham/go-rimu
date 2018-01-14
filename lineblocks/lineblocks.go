package lineblocks

import "github.com/srackham/rimu-go/iotext"

// Render TODO
func Render(reader *iotext.Reader, writer *iotext.Writer) bool {
	writer.Write("<h1>Test</h1>")
	return true
}
