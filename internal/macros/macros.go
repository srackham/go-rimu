package macros

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/srackham/go-rimu-mod/internal/options"
	"github.com/srackham/go-rimu-mod/internal/spans"
	"github.com/srackham/go-rimu-mod/internal/utils/re"
)

func init() {
	Init()
	spans.MacrosRender = Render
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

// Return true if macro is defined.
func IsDefined(name string) bool {
	for _, def := range defs {
		if def.name == name {
			return true
		}
	}
	return false
}

// Return named macro value. If it is not defined found is false.
func Value(name string) (value string, found bool) {
	for _, def := range defs {
		if def.name == name {
			return def.value, true
		}
	}
	return "", false
}

// Return true if macro value is non-blank.
func IsNotBlank(name string) bool {
	value, found := Value(name)
	return found && value != ""
}

// Set named macro value or add it if it doesn't exist.
// If the name ends with '?' then don't set the macro if it already exists.
// `quote` is a single character: ' if a literal value, ` if an expression value.
func SetValue(name string, value string, quote string) {
	if options.SkipMacroDefs() {
		return // Skip if a safe mode is set.
	}
	existential := false
	if strings.HasSuffix(name, "?") {
		name = strings.TrimSuffix(name, "?")
		existential = true
	}
	if name == "--" && value != "" {
		options.ErrorCallback("the predefined blank '--' macro cannot be redefined")
		return
	}
	if quote == "`" {
		options.ErrorCallback("unsupported: expression macro values: `" + value + "`")
	}
	for i, def := range defs {
		if def.name == name {
			if !existential {
				defs[i].value = value
			}
			return
		}
	}
	defs = append(defs, Macro{name: name, value: value})
}

// Render all macro invocations in text string.
// Render Simple invocations first, followed by Parametized, Inclusion and Exclusion invocations.
func Render(text string, silent bool) (result string) {
	MATCH_COMPLEX := regexp.MustCompile(`(?s)\\?\{([\w\-]+)([!=|?](?:|.*?[^\\]))}`) // Parametrized, Inclusion and Exclusion invocations.
	MATCH_SIMPLE := regexp.MustCompile(`\\?\{([\w\-]+)()}`)                         // Simple macro invocation.
	var savedSimple []string
	result = text
	for _, find := range []*regexp.Regexp{MATCH_SIMPLE, MATCH_COMPLEX} {
		result = re.ReplaceAllStringSubmatchFunc(find, result, func(match []string) string {
			if match[0][0] == '\\' {
				return match[0][1:]
			}
			params := match[2]
			if params != "" && params[0] == '?' { // DEPRECATED: Existential macro invocation.
				if !silent {
					options.ErrorCallback("existential macro invocations are deprecated: " + match[0])
				}
				return match[0]
			}
			name := match[1]
			value, found := Value(name)
			if !found {
				if !silent {
					options.ErrorCallback("undefined macro: " + match[0] + ": " + text)
				}
				return match[0]
			}
			if find == MATCH_SIMPLE {
				savedSimple = append(savedSimple, value)
				return "\u0002"
			}
			// Process non-simple macro.
			params = strings.Replace(params, "\\}", "}", -1) // Unescape escaped } characters.
			switch params[0] {
			case '|': // Parametrized macro.
				paramsList := strings.Split(params[1:], "|")
				// Substitute macro parameters.
				// Matches macro definition formal parameters [$]$<param-number>[[\]:<default-param-value>$]
				// 1st group: [$]$
				// 2nd group: <param-number> (1, 2..)
				// 3rd group: :[\]<default-param-value>$
				// 4th group: <default-param-value>
				PARAM_RE := regexp.MustCompile(`(?s)\\?(\$\$?)(\d+)(\\?:(|.*?[^\\])\$)?`)
				value = re.ReplaceAllStringSubmatchFunc(PARAM_RE, value, func(mr []string) string {
					if mr[0][0] == '\\' { // Unescape escaped macro parameters.
						return mr[0][1:]
					}
					p1 := mr[1]
					p2, _ := strconv.ParseInt(mr[2], 10, strconv.IntSize)
					if p2 == 0 {
						return mr[0] // $0 is not a valid parameter name.
					}
					p3 := mr[3]
					p4 := mr[4]
					var param string
					if len(paramsList) < int(p2) {
						// Unassigned parameters are replaced with a blank string.
						param = ""
					} else {
						param = paramsList[p2-1]
					}
					if p3 != "" {
						if p3[0] == '\\' { // Unescape escaped default parameter.
							param += p3[1:]
						} else {
							if param == "" {
								param = p4                                     // Assign default parameter value.
								param = strings.Replace(param, "\\$", "$", -1) // Unescape escaped $ characters in the default value.
							}
						}
					}
					if p1 == "$$" {
						param = spans.Render(param)
					}
					return param
				}, -1)
				return value
			case '!', '=': // Exclusion and Inclusion macro.
				pattern := params[1:]
				pre, err := regexp.Compile("^" + pattern + "$")
				if err != nil {
					if !silent {
						options.ErrorCallback("illegal macro regular expression: " + pattern + ": " + text)
					}
					return match[0]
				}
				skip := !pre.MatchString(value)
				if params[0] == '!' {
					skip = !skip
				}
				if skip {
					return "\u0003" // '\0' flags line for deletion.
				} else {
					return ""
				}
			default:
				options.ErrorCallback("illegal macro syntax: " + match[0])
				return ""
			}

		}, -1)
	}
	// Restore expanded Simple values.
	result = regexp.MustCompile(`\x{0002}`).ReplaceAllStringFunc(result, func(string) string {
		if len(savedSimple) == 0 {
			// This should not happen but there is a limitation: repeated macro substitution parameters
			// ($1, $2...) cannot contain simple macro invocations.
			options.ErrorCallback("repeated macro parameters: " + text)
			return ""
		}
		// Pop from start of list.
		first := savedSimple[0]
		savedSimple = append([]string{}, savedSimple[1:]...)
		return first
	})
	// Delete lines flagged by Inclusion/Exclusion macros.
	if strings.Index(result, "\u0003") >= 0 {
		s := ""
		for _, line := range strings.Split(result, "\n") {
			if !strings.Contains(line, "\u0003") {
				s += line + "\n"
			}
		}
		result = strings.TrimSuffix(s, "\n")
	}
	return
}
