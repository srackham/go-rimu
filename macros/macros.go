package macros

import (
	"regexp"
	"strings"

	"github.com/srackham/rimu-go/options"
)

func init() {
	Init()
}

// Matches a line starting with a macro invocation. $1 = macro invocation.
var MATCH_LINE = regexp.MustCompile(`^({(?:[\w\-]+)(?:[!=|?](?:|.*?[^\\]))?}).*$`)

// Match single-line macro definition. $1 = name, $2 = delimiter, $3 = value, $4 trailing delimiter.
var LINE_DEF = regexp.MustCompile(`^\\?{([\w\-]+\??)}\s*=\s*` + "(['`])" + `(.*)` + "(['`])" + `$`)

// Match multi-line macro definition literal value open delimiter. $1 is first line of macro.
var LITERAL_DEF_OPEN = regexp.MustCompile(`^\\?{[\w\-]+\??}\s*=\s*'(.*)$`)
var LITERAL_DEF_CLOSE = regexp.MustCompile(`^(.*)'$`)

// Match multi-line macro definition expression value open delimiter. $1 is first line of macro.
var EXPRESSION_DEF_OPEN = regexp.MustCompile(`^\\?{[\w\-]+\??}\s*=\s*` + "`" + `(.*)$`)
var EXPRESSION_DEF_CLOSE = regexp.MustCompile("^(.*)`$")

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

// Set named macro value or add it if it doesn't exist.
// If the name ends with '?' then don't set the macro if it already exists.
// `quote` is a single character: ' if a literal value, ` if an expression value.
func SetValue(name string, value string, quote string) {
	// TODO: Implement this as Options.skipMacroDefs() c.f. rimu-kt
	if options.SkipMacroDefs() {
		return // Skip if a safe mode is set.
	}
	existential := false
	if strings.HasSuffix(name, "?") {
		name = strings.TrimSuffix(name, "?")
		existential = true
	}
	if name == "--" && value != "" {
		options.ErrorCallback("the predefined blank \"--\" macro cannot be redefined")
		return
	}
	if quote == "`" {
		options.ErrorCallback("unsupported: expression macro values: `" + value + "`")
	}
	for _, def := range defs {
		if def.name == name {
			if !existential {
				def.value = value
			}
			return
		}
	}
	defs = append(defs, Macro{name: name, value: value})
}

// Render all macro invocations in text string.
// Render Simple invocations first, followed by Parametized, Inclusion and Exclusion invocations.
func Render(text string, silent bool) string {
	// TODO
	return text
}
