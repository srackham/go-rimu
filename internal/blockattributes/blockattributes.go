package blockattributes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/srackham/go-rimu/v11/internal/expansion"
	"github.com/srackham/go-rimu/v11/internal/options"
	"github.com/srackham/go-rimu/v11/internal/spans"
	"github.com/srackham/go-rimu/v11/internal/utils/stringlist"
)

type attrs struct {
	Classes    string // Space separated HTML class names.
	ID         string // HTML element id.
	css        string // HTML CSS styles.
	attributes string // Other HTML element attributes.
	Options    expansion.Options
}

// Attrs are those of the last parsed Block Attributes element.
var Attrs attrs

var ids stringlist.StringList // List of allocated HTML ids.

func init() {
	Init()
}

// Init resets options to default values.
func Init() {
	Attrs.Classes = ""
	Attrs.ID = ""
	Attrs.css = ""
	Attrs.attributes = ""
	Attrs.Options = expansion.Options{}
	ids = nil
}

// Parse text to Attrs block attributes.
func Parse(text string) bool {
	// class names = $1, id = $2, css-properties = $3, html-attributes = $4, block-options = $5
	text = spans.ReplaceInline(text, expansion.Options{Macros: true})
	m := regexp.MustCompile(`^\\?\.((?:\s*[a-zA-Z][\w\-]*)+)*(?:\s*)?(#[a-zA-Z][\w\-]*\s*)?(?:\s*)?(?:"(.+?)")?(?:\s*)?(\[.+])?(?:\s*)?([+-][ \w+-]+)?$`).FindStringSubmatch(text)
	if m == nil {
		return false
	}
	for i, v := range m {
		m[i] = strings.TrimSpace(v)
	}
	if !options.SkipBlockAttributes() {
		if m[1] != "" { // HTML element class names.
			if Attrs.Classes != "" {
				Attrs.Classes += " "
			}
			Attrs.Classes += m[1]
		}
		if m[2] != "" { // HTML element id.
			Attrs.ID = m[2][1:]
		}
		if m[3] != "" { // CSS properties.
			if Attrs.css != "" && !strings.HasSuffix(Attrs.css, ";") {
				Attrs.css += ";"
			}
			if Attrs.css != "" {
				Attrs.css += " "
			}
			Attrs.css += m[3]
		}
		if m[4] != "" && !options.IsSafeModeNz() { // HTML attributes.
			if Attrs.attributes != "" {
				Attrs.attributes += " "
			}
			Attrs.attributes += strings.TrimSpace(m[4][1 : len(m[4])-1])
		}
		if m[5] != "" {
			Attrs.Options.Merge(expansion.Parse(m[5]))
		}
	}
	return true
}

// Inject HTML attributes into the HTML `tag` and return result.
// Consume HTML attributes unless the `tag` argument is blank.
func Inject(tag string) string {
	if tag == "" {
		return tag
	}
	attrs := ""
	if Attrs.Classes != "" {
		m := regexp.MustCompile(`(?i)^<[^>]*class="`).FindStringIndex(tag)
		if m != nil {
			// Inject class names into first existing class attribute in first tag.
			before := tag[:m[1]]
			after := tag[m[1]:]
			tag = before + Attrs.Classes + " " + after
		} else {
			attrs = "class=\"" + Attrs.Classes + "\""
		}
	}
	if Attrs.ID != "" {
		Attrs.ID = strings.ToLower(Attrs.ID)
		hasID := regexp.MustCompile(`(?i)^<[^<]*id=".*?"`).MatchString(tag)
		if hasID || ids.IndexOf(Attrs.ID) >= 0 {
			options.ErrorCallback("duplicate 'id' attribute: " + Attrs.ID)
		} else {
			ids.Push(Attrs.ID)
		}
		if !hasID {
			attrs += " id=\"" + Attrs.ID + "\""
		}
	}
	if Attrs.css != "" {
		m := regexp.MustCompile(`(?i)^<[^<]*style="(.*?)"`).FindStringSubmatchIndex(tag)
		if m != nil {
			// Inject CSS styles into first existing style attribute in first tag.
			before := tag[:m[2]]
			after := tag[m[3]:]
			css := tag[m[2]:m[3]]
			css = strings.TrimSpace(css)
			if !strings.HasSuffix(css, ";") {
				css += ";"
			}
			tag = before + css + " " + Attrs.css + after
		} else {
			attrs += " style=\"" + Attrs.css + "\""
		}
	}
	if Attrs.attributes != "" {
		attrs += " " + Attrs.attributes
	}
	attrs = strings.TrimLeft(attrs, " \n")
	if attrs != "" {
		m := regexp.MustCompile(`(?i)^(<[a-z]+|<h[1-6])(?:[ >])`).FindStringSubmatch(tag) // Match start tag.
		if m != nil {
			before := m[1]
			after := tag[len(m[1]):]
			tag = before + " " + attrs + after
		}
	}
	// Consume the attributes.
	Attrs.Classes = ""
	Attrs.ID = ""
	Attrs.css = ""
	Attrs.attributes = ""
	return tag
}

// Slugify converts text to a slug.
func Slugify(text string) string {
	slug := text
	slug = regexp.MustCompile(`\W+`).ReplaceAllString(slug, "-") // Replace non-alphanumeric characters with dashes.
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")  // Replace multiple dashes with single dash.
	slug = strings.Trim(slug, "-")                               // Trim leading and trailing dashes.
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
