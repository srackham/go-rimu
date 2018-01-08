package delimitedblocks

import "github.com/srackham/rimu-go/iotext"

// Init resets optrwns to default values.
func Init() {
	// TODO
}

// Render TODO
func Render(reader *iotext.Reader, writer *iotext.Writer) bool {
	reader.Next()
	writer.Write("<p>Hello <em>rimu-go!</em></p>")
	return true
}
