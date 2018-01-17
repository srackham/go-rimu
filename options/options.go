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
var SafeMode int
var htmlReplacement string
var callback CallbackFunction

// Init resets options to default values.
func Init() {
	SafeMode = 0
	htmlReplacement = "<mark>replaced HTML</mark>"
	callback = nil
}

// Return true if safeMode is non-zero.
func IsSafeModeNz() bool {
	return SafeMode != 0
}

/*
export function getSafeMode(): number {
  return safeMode
}
*/

// Return true if Block Attribute elements are ignored.
func SkipBlockAttributes() bool {
	return SafeMode != 0 && SafeMode&0x4 != 0
}

/*
function setSafeMode(value: number | string | undefined): void {
  let n = Number(value)
  if (!isNaN(n)) {
    safeMode = n
  }
}

function setHtmlReplacement(value: string | undefined): void {
  if (value === undefined) return
  htmlReplacement = value
}

function setReset(value: boolean | string | undefined): void {
  if (value === true || value === 'true') {
    Api.init()
  }
}
*/

// TODO
func UpdateOptions(options RenderOptions) {
}

/*
export function updateOptions(options: RenderOptions): void {
  if ('reset' in options) setReset(options.reset) // Reset takes priority.
  if ('safeMode' in options) setSafeMode(options.safeMode)
  if ('htmlReplacement' in options) setHtmlReplacement(options.htmlReplacement)
  if ('callback' in options) callback = options.callback
}
*/

// Set named option value.
func SetOption(name string, value string) {
	/*
			switch name {
			case "safeMode":
				n, err := strconv.ParseInt(value, 10, 64)
				SafeMode = int(n)
				// case "htmlReplacement": htmlReplacement = value
				// case "reset": if (value.toBoolean()) Api.init()
				// default: throw IllegalArgumentException()
			}
			if err != nil {
				painic("TODO")
		  }
	*/
}

/*
// Filter HTML based on current safeMode.
export function htmlSafeModeFilter(html: string): string {
  switch (safeMode & 0x3) {
    case 0:   // Raw HTML (default behavior).
      return html
    case 1:   // Drop HTML.
      return ''
    case 2:   // Replace HTML with 'htmlReplacement' option string.
      return htmlReplacement
    case 3:   // Render HTML as text.
      return Utils.replaceSpecialChars(html)
    default:
      return ''
  }
}
*/

func ErrorCallback(message string) {
	callback(CallbackMessage{Kind: "error", Text: message})
}

/*
export function errorCallback(message: string): void {
  if (callback) {
    callback({type: 'error', text: message})
  }
}

// Called when an unexpected program error occurs.
export function panic(message: string): void {
  let msg = 'panic: ' + message
  console.error(msg)
  errorCallback(msg)
}
*/
