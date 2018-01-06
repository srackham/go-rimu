package options

type RenderOptions struct {
	safeMode        int
	htmlReplacement string
	reset           bool
	callback        CallbackFunction
}

type CallbackMessage struct {
	kind string
	text string
}

// CallbackFunction TODO
type CallbackFunction func(message CallbackMessage)

// Global option values.
var safeMode int
var htmlReplacement string
var callback CallbackFunction

// Init resets options to default values.
func Init() {
	safeMode = 0
	htmlReplacement = "<mark>replaced HTML</mark>"
	callback = nil
}

// Update TODO
func Update(options RenderOptions) {

}
