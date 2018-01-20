package quotes

import "github.com/srackham/rimu-go/api"

func init() {
	api.RegisterInit(Init)
}

type Definition struct {
	Quote    string // Single quote character.
	OpenTag  string
	CloseTag string
	Spans    bool // Allow span elements inside quotes.
}

var defs []Definition // Mutable definitions initialized by DEFAULT_DEFS.

var DEFAULT_DEFS = []Definition{
	{
		Quote:    "**",
		OpenTag:  "<strong>",
		CloseTag: "</strong>",
		Spans:    true,
	},
	{
		Quote:    "*",
		OpenTag:  "<em>",
		CloseTag: "</em>",
		Spans:    true,
	},
	{
		Quote:    "__",
		OpenTag:  "<strong>",
		CloseTag: "</strong>",
		Spans:    true,
	},
	{
		Quote:    "_",
		OpenTag:  "<em>",
		CloseTag: "</em>",
		Spans:    true,
	},
	{
		Quote:    "``",
		OpenTag:  "<code>",
		CloseTag: "</code>",
		Spans:    false,
	},
	{
		Quote:    "`",
		OpenTag:  "<code>",
		CloseTag: "</code>",
		Spans:    false,
	},
	{
		Quote:    "~~",
		OpenTag:  "<del>",
		CloseTag: "</del>",
		Spans:    true,
	},
}

// TODO
// Stubs
func Init() {
	// TODO
}

func SetDefinition(def Definition) {
}

/*

export let quotesRe: RegExp // Searches for quoted text.
let unescapeRe: RegExp      // Searches for escaped quotes.

// Reset definitions to defaults.
export function init(): void {
  defs = DEFAULT_DEFS.map(def => Utils.copy(def))
  initializeRegExps()
}

// Synthesise re's to find and unescape quotes.
export function initializeRegExps(): void {
  let quotes = defs.map(def => Utils.escapeRegExp(def.Quote))
  // $1 is quote character(s), $2 is quoted text.
  // Quoted text cannot begin or end with whitespace.
  // Quoted can span multiple lines.
  // Quoted text cannot end with a backslash.
  quotesRe = RegExp('\\\\?(' + quotes.join('|') + ')([^\\s\\\\]|\\S[\\s\\S]*?[^\\s\\\\])\\1', 'g')
  // $1 is quote character(s).
  unescapeRe = RegExp('\\\\(' + quotes.join('|') + ')', 'g')
}

// Return the quote definition corresponding to 'quote' character, return undefined if not found.
export function getDefinition(quote: string): Definition {
  return defs.filter(def => def.Quote === quote)[0]
}

// Strip backslashes from quote characters.
export function unescape(s: string): string {
  return s.replace(unescapeRe, '$1')
}

// Update existing or add new quote definition.
export function setDefinition(def: Definition): void {
  for (let d of defs) {
    if (d.quote === def.quote) {
      // Update existing definition.
      d.openTag = def.openTag
      d.closeTag = def.closeTag
      d.spans = def.spans
      return
    }
  }
  // Double-quote definitions are prepended to the array so they are matched
  // before single-quote definitions (which are appended to the array).
  if (def.quote.length === 2) {
    defs.unshift(def)
  }
  else {
    defs.push(def)
  }
  initializeRegExps()
}

*/
