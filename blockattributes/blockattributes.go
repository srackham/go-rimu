package blockattributes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/srackham/rimu-go/utils/stringlist"

	"github.com/srackham/rimu-go/utils"
)

var Classes string    // Space separated HTML class names.
var Id string         // HTML element id.
var Css string        // HTML CSS styles.
var Attributes string // Other HTML element attributes.
var Options utils.ExpansionOptions

var ids stringlist.StringList // List of allocated HTML ids.

// Init resets options to default values.
func Init() {
	// TODO
	Classes = ""
	Id = ""
	Css = ""
	Attributes = ""
	Options = utils.ExpansionOptions{}
	ids = nil
}

// Inject HTML attributes from attrs into the opening tag.
// Consume HTML attributes unless the 'tag' argument is blank.
func Inject(tag string) string {
	return tag
}

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
