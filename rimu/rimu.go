package rimu

import (
	"github.com/srackham/go-rimu/api"
	"github.com/srackham/go-rimu/options"
)

// CallbackMessage TODO
type CallbackMessage = options.CallbackMessage

// CallbackFunction TODO
type CallbackFunction = options.CallbackFunction

// RenderOptions TODO
type RenderOptions = options.RenderOptions

// Render is public API to translate Rimu Markup to HTML.
func Render(text string, opts RenderOptions) string {
	options.UpdateOptions(opts)
	return api.Render(text)
}
