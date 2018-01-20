package blockattributes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/srackham/rimu-go/api"
	"github.com/srackham/rimu-go/expansion"
	"github.com/srackham/rimu-go/utils/stringlist"
)

func init() {
	api.RegisterInit(Init)
}

var Classes string    // Space separated HTML class names.
var Id string         // HTML element id.
var Css string        // HTML CSS styles.
var Attributes string // Other HTML element attributes.
var Options expansion.ExpansionOptions

var ids stringlist.StringList // List of allocated HTML ids.

// Init resets options to default values.
func Init() {
	// TODO
	Classes = ""
	Id = ""
	Css = ""
	Attributes = ""
	Options = expansion.ExpansionOptions{}
	ids = nil
}

func Parse(match []string) bool {
	return true
}

/*
  export function parse(match: RegExpExecArray): boolean {
    // Parse Block Attributes.
    // class names = $1, id = $2, css-properties = $3, html-attributes = $4, block-options = $5
    let text = match[0]
    text = replaceInline(text, {macros: true})
    let m = /^\\?\.((?:\s*[a-zA-Z][\w\-]*)+)*(?:\s*)?(#[a-zA-Z][\w\-]*\s*)?(?:\s*)?(?:"(.+?)")?(?:\s*)?(\[.+])?(?:\s*)?([+-][ \w+-]+)?$/.exec(text)
    if (!m) {
      return false
    }
    if (!Options.skipBlockAttributes()) {
      if (m[1]) { // HTML element class names.
        classes += ' ' + m[1].trim()
        classes = classes.trim()
      }
      if (m[2]) { // HTML element id.
        id = m[2].trim().slice(1)
      }
      if (m[3]) { // CSS properties.
        if (css && css.substr(-1) !== ';') css += ';'
        css += ' ' + m[3].trim()
        css = css.trim()
      }
      if (m[4] && !Options.isSafeModeNz()) { // HTML attributes.
        attributes += ' ' + m[4].slice(1, m[4].length - 1).trim()
        attributes = attributes.trim()
      }
      DelimitedBlocks.setBlockOptions(options, m[5])
    }
    return true
  }
*/

// Inject HTML attributes from attrs into the opening tag.
// Consume HTML attributes unless the 'tag' argument is blank.
func Inject(tag string) string {
	// TODO
	return tag
}

/*
  // Inject HTML attributes from attrs into the opening tag.
  // Consume HTML attributes unless the 'tag' argument is blank.
  export function inject(tag: string): string {
    if (!tag) {
      return tag
    }
    let attrs = ''
    if (classes) {
      if (/class=".*?"/i.test(tag)) {
        // Inject class names into existing class attribute.
        tag = tag.replace(/class="(.*?)"/i, `class="${classes} $1"`)
      }
      else {
        attrs = `class="${classes}"`
      }
    }
    if (id) {
      id = id.toLowerCase()
      let has_id = /id=".*?"/i.test(tag)
      if (has_id || ids.indexOf(id) > -1) {
        Options.errorCallback(`duplicate 'id' attribute: ${id}`)
      }
      else {
        ids.push(id)
      }
      if (!has_id) {
        attrs += ` id="${id}"`
      }
    }
    if (css) {
      if (/style=".*?"/i.test(tag)) {
        // Inject CSS styles into existing style attribute.
        tag = tag.replace(/style="(.*?)"/i, function (match: string, p1: string): string {
          p1 = p1.trim()
          if (p1 && p1.substr(-1) !== ';') p1 += ';'
          return `style="${p1} ${css}"`
        })
      }
      else {
        attrs += ` style="${css}"`
      }
    }
    if (attributes) {
      attrs += ' ' + attributes
    }
    attrs = attrs.trim()
    if (attrs) {
      let match = tag.match(/^<([a-zA-Z]+|h[1-6])(?=[ >])/)
      if (match) {
        let before = tag.slice(0, match[0].length)
        let after = tag.slice(match[0].length)
        tag = before + ' ' + attrs + after
      }
    }
    // Consume the attributes.
    classes = ''
    id = ''
    css = ''
    attributes = ''
    return tag
  }
*/

func Slugify(text string) string {
	slug := regexp.MustCompile(`\W+`).ReplaceAllString(text, "-") // Replace non-alphanumeric characters with dashes.
	slug = regexp.MustCompile(`-+`).ReplaceAllString(text, "-")   // Replace multiple dashes with single dash.
	slug = strings.Trim(slug, "-")                                // Trim leading and trailing dashes.
	slug = strings.ToLower(slug)
	if slug == "" {
		slug = "x"
	}
	if ids.IndexOf(slug) > -1 { // Another element already has that id.
		i := 2
		for ids.IndexOf(slug+"-"+fmt.Sprint(i)) > -1 {
			i++
		}
		slug += "-" + fmt.Sprint(i)
	}
	return slug
}
