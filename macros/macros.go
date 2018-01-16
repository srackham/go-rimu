package macros

/*
// Matches a line starting with a macro invocation. $1 = macro invocation.
export const MATCH_LINE = /^({(?:[\w\-]+)(?:[!=|?](?:|.*?[^\\]))?}).*$/
// Match single-line macro definition. $1 = name, $2 = delimiter, $3 = value.
export const LINE_DEF = /^\\?{([\w\-]+\??)}\s*=\s*(['`])(.*)\2$/
// Match multi-line macro definition literal value open delimiter. $1 is first line of macro.
export const LITERAL_DEF_OPEN = /^\\?{[\w\-]+\??}\s*=\s*'(.*)$/
export const LITERAL_DEF_CLOSE = /^(.*)'$/
// Match multi-line macro definition expression value open delimiter. $1 is first line of macro.
export const EXPRESSION_DEF_OPEN = /^\\?{[\w\-]+\??}\s*=\s*`(.*)$/
export const EXPRESSION_DEF_CLOSE = /^(.*)`$/
*/

type Macro struct {
	name  string
	value string
}

var defs []Macro

// Reset definitions to defaults.
func Init() {
	// Initialize predefined macros.
	defs = []Macro{
		{name: "--", value: ""},
		{name: "--header-ids", value: ""},
	}
}

// Return named macro value or nil if it doesn't exist.
func IsDefined(name string) bool {
	for _, def := range defs {
		if def.name == name {
			return true
		}
	}
	return false
}

// Return named macro value or nil if it doesn't exist.
func GetValue(name string) (value string, found bool) {
	for _, def := range defs {
		if def.name == name {
			return def.value, true
		}
	}
	return "", false
}

// Render all macro invocations in text string.
// Render Simple invocations first, followed by Parametized, Inclusion and Exclusion invocations.
func Render(text string, silent bool) string {
	// TODO
	return text
}
