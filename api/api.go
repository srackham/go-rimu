package api

import (
	"github.com/srackham/rimu-go/options"
)

// Init TODO
func Init() {
	options.Init()
}

// Render TODO
func Render(text string) string {
	return "<p>Hello <em>rimu-go!</em></p>"
}
