package options

import (
	"strconv"

	"github.com/srackham/rimu-go/proxies"
	"github.com/srackham/rimu-go/utils"
	// "github.com/srackham/rimu-go/utils"
)

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

// UpdateOptions processes non-nil ops fields. Panics if non-nil types are incorrect.
func UpdateOptions(opts RenderOptions) {
	// Reset takes priority.
	if opts.Reset != nil {
		if opts.Reset.(bool) {
			proxies.ApiInit()
		}
	}
	if opts.SafeMode != nil {
		safeMode = opts.SafeMode.(int)
	}
	if opts.HtmlReplacement != nil {
		htmlReplacement = opts.HtmlReplacement.(string)
	}
	if opts.Callback != nil {
		callback = opts.Callback
	}
}

// SetOption parses a named API option value. Panics if there is an error.
func SetOption(name string, value string) {
	switch name {
	case "safeMode":
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(err)
		}
		safeMode = int(n)
	case "htmlReplacement":
		htmlReplacement = value
	case "reset":
		b, err := strconv.ParseBool(value)
		if err != nil {
			panic(err)
		}
		if b {
			proxies.ApiInit()
		}
	default:
		panic("illegal API option name: " + name)
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
		return utils.ReplaceSpecialChars(html)
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
