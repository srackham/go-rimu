package options

import (
	"fmt"
	"strconv"

	"github.com/srackham/go-rimu/internal/utils/str"
)

// api package dependency injection.
var ApiInit func()

func init() {
	Init()
}

// RenderOptions sole use is for passing options into the public API.
// All fields can be nil so that options can be selectively updated (if nil then don't update).
type RenderOptions struct {
	SafeMode        interface{} // nil or int
	HtmlReplacement interface{} // nil or string
	Reset           interface{} // nil or bool
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

// Return true if safeMode is non-zero.
func IsSafeModeNz() bool {
	return safeMode != 0
}

// Return true if Macro Definitions are ignored.
func SkipMacroDefs() bool {
	return safeMode != 0 && safeMode&0x8 == 0
}

// Return true if Block Attribute elements are ignored.
func SkipBlockAttributes() bool {
	return safeMode != 0 && safeMode&0x4 != 0
}

// UpdateOptions processes non-nil opts fields.
// Error callback option values are illegal.
func UpdateOptions(opts RenderOptions) {
	// Install callback first to ensure option errors are logged.
	if opts.Callback != nil {
		callback = opts.Callback
	}
	// Reset takes priority.
	if opts.Reset != nil {
		SetOption("reset", fmt.Sprintf("%v", opts.Reset))
	}
	// Install callback again in case it has been reset.
	if opts.Callback != nil {
		callback = opts.Callback
	}
	if opts.SafeMode != nil {
		SetOption("safeMode", fmt.Sprintf("%v", opts.SafeMode))
	}
	if opts.HtmlReplacement != nil {
		SetOption("htmlReplacement", fmt.Sprintf("%v", opts.HtmlReplacement))
	}
}

// SetOption parses a named API option value.
// Error callback if option values are illegal.
func SetOption(name string, value string) {
	switch name {
	case "safeMode":
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil || n < 0 || n > 15 {
			ErrorCallback("illegal safeMode API option value: " + value)
		} else {
			safeMode = int(n)
		}
	case "htmlReplacement":
		htmlReplacement = value
	case "reset":
		b, err := strconv.ParseBool(value)
		if err != nil {
			ErrorCallback("illegal reset API option value: " + value)
		} else {
			if b {
				ApiInit()
			}
		}
	default:
		ErrorCallback("illegal API option name: " + name)
	}
}

// Filter HTML based on current safeMode.
func HtmlSafeModeFilter(html string) string {
	switch safeMode & 0x3 {
	case 0: // Raw HTML (default behavior).
		return html
	case 1: // Drop HTML.
		return ""
	case 2: // Replace HTML with 'htmlReplacement' option string.
		return htmlReplacement
	case 3: // Render HTML as text.
		return str.ReplaceSpecialChars(html)
		return html
	default:
		return ""
	}
}

func ErrorCallback(message string) {
	if callback != nil {
		callback(CallbackMessage{Kind: "error", Text: message})
	}
}
