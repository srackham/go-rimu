package options

type RenderOptions struct {
	SafeMode        int
	HtmlReplacement string
	Reset           bool
	Callback        CallbackFunction
}

type CallbackMessage struct {
	Kind string
	Text string
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
