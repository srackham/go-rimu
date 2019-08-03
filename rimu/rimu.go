package rimu

import (
	"github.com/srackham/go-rimu/v11/internal/api"
	"github.com/srackham/go-rimu/v11/internal/options"
)

// CallbackFunction is the API callback function type.
type CallbackFunction = options.CallbackFunction

// CallbackMessage contains the callback message passed to the callback function.
type CallbackMessage = options.CallbackMessage

// RenderOptions contains the API render options.
type RenderOptions = options.RenderOptions

// Render is public API to translate Rimu Markup to HTML.
func Render(text string, opts RenderOptions) string {
	options.UpdateOptions(opts)
	return api.Render(text)
}
