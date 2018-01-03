package rimu

import "github.com/srackham/rimu-go/api"

type CallbackMessage struct {
	kind string
	text string
}

// CallbackFunction TODO
type CallbackFunction func(message CallbackMessage)

// DoNothing TODO
var DoNothing CallbackFunction = func(message CallbackMessage) {} // Default render() callback.

// RenderOptions TODO
type RenderOptions struct {
	safeMode        int
	htmlReplacement string
	reset           bool
	callback        CallbackFunction
}

// Render is public API to translate Rimu Markup to HTML.
func Render(text string, options RenderOptions) string {
	// Force object instantiation before Options.update().
	// Otherwise the ensuing Api.render() will instanitate Api and the Api init{} block will reset Options.

	// Api // Ensure Api is instantiated.
	return api.Render(text)
}
