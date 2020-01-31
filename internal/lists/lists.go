package lists

import (
	"regexp"
	"strings"

	"github.com/srackham/go-rimu/v11/internal/blockattributes"
	"github.com/srackham/go-rimu/v11/internal/delimitedblocks"
	"github.com/srackham/go-rimu/v11/internal/expansion"
	"github.com/srackham/go-rimu/v11/internal/iotext"
	"github.com/srackham/go-rimu/v11/internal/lineblocks"
	"github.com/srackham/go-rimu/v11/internal/spans"
	"github.com/srackham/go-rimu/v11/internal/utils/stringlist"
)

// Definition of List element.
type Definition struct {
	match        *regexp.Regexp
	listOpenTag  string
	listCloseTag string
	itemOpenTag  string
	itemCloseTag string
	termOpenTag  string // Definition lists only.
	termCloseTag string // Definition lists only.
}

// ItemInfo contains information about a matched list item element.
type ItemInfo struct {
	match []string
	def   Definition
	id    string // List ID.
}

var defs = []Definition{
	// Prefix match with backslash to allow escaping.

	// Unordered lists.
	// $1 is list ID $2 is item text.
	{
		match:        regexp.MustCompile(`^\\?\s*(-|\+|\*{1,4})\s+(.*)$`),
		listOpenTag:  "<ul>",
		listCloseTag: "</ul>",
		itemOpenTag:  "<li>",
		itemCloseTag: "</li>",
	},
	// Ordered lists.
	// $1 is list ID $2 is item text.
	{
		match:        regexp.MustCompile(`^\\?\s*(?:\d*)(\.{1,4})\s+(.*)$`),
		listOpenTag:  "<ol>",
		listCloseTag: "</ol>",
		itemOpenTag:  "<li>",
		itemCloseTag: "</li>",
	},
	// Definition lists.
	// $1 is term, $2 is list ID, $3 is definition.
	{
		match:        regexp.MustCompile(`^\\?\s*(.*[^:])(:{2,4})(|\s+.*)$`),
		listOpenTag:  "<dl>",
		listCloseTag: "</dl>",
		itemOpenTag:  "<dd>",
		itemCloseTag: "</dd>",
		termOpenTag:  "<dt>",
		termCloseTag: "</dt>",
	},
}

const noMatch = "NO_MATCH" // "no matching list item found" list ID constant.

// noMatchItem returns "no matching list item found" constant.
func noMatchItem() ItemInfo { return ItemInfo{id: noMatch} }

var ids []string // Stack of open list IDs.

// Render list item in reader to writer.
func Render(reader *iotext.Reader, writer *iotext.Writer) bool {
	if reader.Eof() {
		panic("premature eof")
	}
	startItem := matchItem(reader)
	if startItem.id == noMatch {
		return false
	}
	ids = nil
	renderList(startItem, reader, writer)
	// ids should now be empty.
	if len(ids) != 0 {
		panic("list stack failure")
	}
	return true
}

func renderList(item ItemInfo, reader *iotext.Reader, writer *iotext.Writer) ItemInfo {
	ids = append(ids, item.id)
	writer.Write(blockattributes.Inject(item.def.listOpenTag))
	for {
		nextItem := renderListItem(item, reader, writer)
		if nextItem.id == noMatch || nextItem.id != item.id {
			// End of list or next item belongs to ancestor.
			writer.Write(item.def.listCloseTag)
			ids = ids[:len(ids)-1] // pop
			return nextItem
		}
		item = nextItem
	}
}

// Render the current list item, return the next list item or null if there are no more items.
func renderListItem(item ItemInfo, reader *iotext.Reader, writer *iotext.Writer) ItemInfo {
	def := item.def
	match := item.match
	var text string
	if len(match) == 4 { // 3 match groups => definition list.
		attrs := blockattributes.Attrs
		writer.Write(blockattributes.Inject(def.termOpenTag))
		attrs.ID = ""
		blockattributes.Attrs = attrs // Restore consumed block attributes.
		text = spans.ReplaceInline(match[1], expansion.Options{Macros: true, Spans: true})
		writer.Write(text)
		writer.Write(def.termCloseTag)
	}
	writer.Write(blockattributes.Inject(def.itemOpenTag))
	// Process item text from first line.
	itemLines := iotext.NewWriter()
	text = match[len(match)-1]
	itemLines.Write(text + "\n")
	// Process remainder of list item i.e. item text, optional attached block, optional child list.
	reader.Next()
	attachedLines := iotext.NewWriter()
	blankLines := 0
	attachedDone := false
	var nextItem ItemInfo
	for {
		blankLines = consumeBlockAttributes(reader, attachedLines)
		if blankLines >= 2 || blankLines == -1 {
			// EOF or two or more blank lines terminates list.
			nextItem = noMatchItem()
			break
		}
		nextItem = matchItem(reader)
		if nextItem.id != noMatch {
			if stringlist.StringList(ids).IndexOf(nextItem.id) != -1 {
				// Next item belongs to current list or a parent list.
			} else {
				// Render child list.
				nextItem = renderList(nextItem, reader, attachedLines)
			}
			break
		}
		if attachedDone {
			break // Multiple attached blocks are not permitted.
		}
		if blankLines == 0 {
			savedIds := ids
			ids = nil
			if delimitedblocks.Render(reader, attachedLines, []string{"comment", "code", "division", "html", "quote"}) {
				attachedDone = true
			} else {
				// Item body line.
				itemLines.Write(reader.Cursor() + "\n")
				reader.Next()
			}
			ids = savedIds
		} else if blankLines == 1 {
			if delimitedblocks.Render(reader, attachedLines, []string{"indented", "quote-paragraph"}) {
				attachedDone = true
			} else {
				break
			}
		}
	}
	// Write item text.
	text = strings.TrimSpace(itemLines.String())
	text = spans.ReplaceInline(text, expansion.Options{Macros: true, Spans: true})
	writer.Write(text)
	// Write attachment and child list.
	writer.Buffer = append(writer.Buffer, attachedLines.Buffer...)
	// Close list item.
	writer.Write(def.itemCloseTag)
	return nextItem
}

// Consume blank lines and Block Attributes.
// Return number of blank lines read or -1 if EOF.
func consumeBlockAttributes(reader *iotext.Reader, writer *iotext.Writer) int {
	blanks := 0
	for {
		if reader.Eof() {
			return -1
		}
		if lineblocks.Render(reader, writer, []string{"attributes"}) {
			continue
		}
		if reader.Cursor() != "" {
			return blanks
		}
		blanks++
		reader.Next()
	}
}

// Check if the line at the reader cursor matches a list related element.
// Unescape escaped list items in reader.
// If it does not match a list related element return null.
func matchItem(reader *iotext.Reader) ItemInfo {
	if !reader.Eof() {
		var item ItemInfo // ItemInfo factory.
		for _, def := range defs {
			match := def.match.FindStringSubmatch(reader.Cursor())
			if match != nil {
				if match[0][0] == '\\' {
					reader.SetCursor(reader.Cursor()[1:]) // Drop backslash.
					break
				}
				item.match = match
				item.def = def
				item.id = match[len(match)-2] // The second to last match group is the list ID.
				return item
			}
		}
	}
	return noMatchItem()
}
