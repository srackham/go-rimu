/*
  Command-lne app to convert Rimu source to HTML.
  Run 'node rimu.js --help' for details.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/srackham/go-rimu/rimu"
	"github.com/srackham/go-rimu/utils/stringlist"
)

var osExit = os.Exit // Mocked by tests.

const MANPAGE = `NAME
  rimuc - convert Rimu source to HTML

SYNOPSIS
  rimuc [OPTIONS...] [FILES...]

DESCRIPTION
  Reads Rimu source markup from stdin, converts them to HTML
  then writes the HTML to stdout. If FILES are specified
  the Rimu source is read from FILES. The contents of files
  with an .html extension are passed directly to the output.
  An input file named '-' is read from stdin.

  If a file named .rimurc exists in the user's home directory
  then its contents is processed (with --safe-mode 0).
  This behavior can be disabled with the --no-rimurc option.

  Inputs are processed in the following order: .rimurc file,
  --prepend-file options, --prepend options, FILES...

OPTIONS
  -h, --help
    Display help message.

  --html-replacement TEXT
    Embedded HTML is replaced by TEXT when --safe-mode is set to 2.
    Defaults to '<mark>replaced HTML</mark>'.

  --layout LAYOUT
    Generate a styled HTML document. rimuc includes the
    following built-in document layouts:

    'classic': Desktop-centric layout.
    'flex':    Flexbox mobile layout (experimental).
    'sequel':  Responsive cross-device layout.

    If only one source file is specified and the --output
    option is not specified then the output is written to a
    same-named file with an .html extension.
    This option enables --header-ids.

  -o, --output OUTFILE
    Write output to file OUTFILE instead of stdout.
    If OUTFILE is a hyphen '-' write to stdout.

  --pass
    Pass the stdin input verbatim to the output.

  -p, --prepend SOURCE
    Process the SOURCE text before all other inputs.
    Rendered with --safe-mode 0.

  --prepend-file PREPEND_FILE
    Process the PREPEND_FILE contents immediately after --prepend
    and .rimurc processing.
    Rendered with --safe-mode 0.

  --no-rimurc
    Do not process .rimurc from the user's home directory.

  --safe-mode NUMBER
    Non-zero safe modes ignore: Definition elements; API option elements;
    HTML attributes in Block Attributes elements.
    Also specifies how to process HTML elements:

    --safe-mode 0 renders HTML (default).
    --safe-mode 1 ignores HTML.
    --safe-mode 2 replaces HTML with --html-replacement option value.
    --safe-mode 3 renders HTML as text.

    Add 4 to --safe-mode to ignore Block Attribute elements.
    Add 8 to --safe-mode to allow Macro Definitions.

  --theme THEME, --lang LANG, --title TITLE, --highlightjs, --mathjax,
  --no-toc, --custom-toc, --section-numbers, --header-ids, --header-links
    Shortcuts for the following prepended macro definitions:

    --prepend "{--custom-toc}='true'"
    --prepend "{--header-ids}='true'"
    --prepend "{--header-links}='true'"
    --prepend "{--highlightjs}='true'"
    --prepend "{--lang}='LANG'"
    --prepend "{--mathjax}='true'"
    --prepend "{--no-toc}='true'"
    --prepend "{--section-numbers}='true'"
    --prepend "{--theme}='THEME'"
    --prepend "{--title}='TITLE'"

LAYOUT OPTIONS
  The following options are available when the --layout option
  specifies a built-in layout:

  Option             Description
  _______________________________________________________________
  --custom-toc       Set to a non-blank value if a custom table
                     of contents is used.
  --header-links     Set to a non-blank value to generate h2 and
                     h3 header header links.
  --highlightjs      Set to non-blank value to enable syntax
                     highlighting with Highlight.js.
  --lang             HTML document language attribute value.
  --mathjax          Set to a non-blank value to enable MathJax.
  --no-toc           Set to a non-blank value to suppress table of
                     contents generation.
  --section-numbers  Apply h2 and h3 section numbering.
  --theme            Styling theme. Theme names:
                     'legend', 'graystone', 'vintage'.
  --title            HTML document title.
  _______________________________________________________________
  These options are translated by rimuc to corresponding layout
  macro definitions using the --prepend option.

LAYOUT CLASSES
  The following CSS classes are available for use in Rimu Block
  Attributes elements when the --layout option specifies a
  built-in layout:

  CSS class        Description
  ______________________________________________________________
  align-center     Text alignment center.
  align-left       Text alignment left.
  align-right      Text alignment right.
  bordered         Adds table borders.
  cite             Quote and verse attribution.
  dl-horizontal    Format labeled lists horizontally.
  dl-numbered      Number labeled list items.
  dl-counter       Prepend dl item counter to element content.
  ol-counter       Prepend ol item counter to element content.
  ul-counter       Prepend ul item counter to element content.
  no-auto-toc      Exclude heading from table of contents.
  no-page-break    Avoid page break inside the element.
  no-print         Do not print.
  page-break       Force page break before the element.
  preserve-breaks  Honor line breaks in source text.
  sidebar          Sidebar format (paragraphs, division blocks).
  verse            Verse format (paragraphs, division blocks).
  ______________________________________________________________

PREDEFINED MACROS
  Macro name         Description
  _______________________________________________________________
  --                 Blank macro (empty string).
                     The Blank macro cannot be redefined.
  --header-ids       Set to a non-blank value to generate h1, h2
                     and h3 header id attributes.
  _______________________________________________________________
`
const STDIN = "-"

// rimurcPath returns path of $HOME/.rimurc file.
// Return "" if $HOME not found.
func rimurcPath() (result string) {
	if user, err := user.Current(); err == nil {
		result = filepath.Join(user.HomeDir, ".rimurc")
	}
	return
}

// MockExit type used by rimuc_test osExit() mock.
type MockExit struct{}

// Helpers.
func die(message string) {
	if message != "" {
		// TODO: this impacts rimuc tests.
		// fmt.Fprintln(os.Stderr, "error: "+message)
		fmt.Fprintln(os.Stderr, message)
	}
	osExit(1)
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func readResourceFile(name string) string {
	//TODO
	return ""
}

func importLayoutFile(name string) string {
	// Imports not supported in go-rimu.
	die("missing --layout: " + name)
	return ""
}

func main() {
	defer func() {
		r, ok := recover().(MockExit)
		if !ok {
			panic(r)
		}
	}()
	args := stringlist.StringList(os.Args)
	args.Shift() // Skip program name.
	nextArg := func(errMsg string) string {
		if len(args) == 0 {
			die(errMsg)
		}
		return args.Shift()
	}
	var safeMode interface{}
	var htmlReplacement interface{}
	layout := ""
	noRimurc := false
	var prependFiles stringlist.StringList
	pass := false
	// Parse command-line options.
	prepend := ""
	outfile := ""

outer:
	for len(args) > 0 {
		arg := args.Shift()
		switch arg {
		case "--help", "-h":
			fmt.Printf("\n" + MANPAGE)
			osExit(0)
		case "--lint", "-l": // Deprecated in Rimu 10.0.0
			break
		case "--output", "-o":
			outfile = nextArg("missing --output file name")
		case "--pass":
			pass = true
		case "--prepend", "-p":
			prepend += nextArg("missing --prepend value") + "\n"
		case "--prepend-file":
			prependFiles.Push(nextArg("missing --prepend-file file name"))
		case "--no-rimurc":
			noRimurc = true
		case "--safe-mode",
			"--safeMode": // Deprecated in Rimu 7.1.0.
			s := nextArg("missing --safe-mode value")
			n, err := strconv.ParseInt(s, 10, strconv.IntSize)
			if err != nil {
				die("illegal --safe-mode option value: " + s)
			}
			safeMode = int(n)
		case "--html-replacement",
			"--htmlReplacement": // Deprecated in Rimu 7.1.0.
			htmlReplacement = nextArg("missing --html-replacement value")
		case "--styled", "-s": // Deprecated in Rimu 10.0.0
			prepend += "{--header-ids}=\"true\"\n"
			if layout == "" {
				layout = "classic"
			}
			// Styling macro definitions shortcut options.
		case "--highlightjs",
			"--mathjax",
			"--section-numbers",
			"--theme",
			"--title",
			"--lang",
			"--toc", // Deprecated in Rimu 8.0.0
			"--no-toc",
			"--sidebar-toc",  // Deprecated in Rimu 10.0.0
			"--dropdown-toc", // Deprecated in Rimu 10.0.0
			"--custom-toc",
			"--header-ids",
			"--header-links":
			macroValue := ""
			if strings.Contains("--lang|--title|--theme", arg) {
				macroValue = nextArg("missing " + arg + " value")
			} else {
				macroValue = "true"
			}
			prepend += "{" + arg + "}=\"" + macroValue + "\"\n"
		case "--layout",
			"--styled-name": // Deprecated in Rimu 10.0.0
			layout = nextArg("missing --layout value")
			prepend += "{--header-ids}=\"true\"\n"
		default:
			args.Unshift(arg) // argv contains source file names.
			break outer
		}
	}
	// args contains the list of source files.
	files := args
	if len(files) == 0 {
		files.Push(STDIN)
	} else if len(files) == 1 && layout != "" && files[0] != "-" && outfile != "" {
		// Use the source file name with .html extension for the output file.
		ext := path.Ext(files[0])
		outfile = files[0][:len(files[0])-len(ext)] + ".html"
	}
	const RESOURCE_TAG = "resource:"    // Tag for resource files.
	const PREPEND = "--prepend options" // Tag for --prepend source.
	if layout != "" {
		// Envelope source files with header and footer.
		files.Unshift(RESOURCE_TAG + layout + "-header.rmu")
		files.Push(RESOURCE_TAG + layout + "-footer.rmu")
	}
	// Prepend $HOME/.rimurc file if it exists.
	if !noRimurc && fileExists(rimurcPath()) {
		prependFiles.Unshift(rimurcPath())
	}
	if prepend != "" {
		prependFiles.Push(PREPEND)
	}
	files = append(prependFiles, files...)
	// Convert Rimu source files to HTML.
	output := ""
	errors := 0
	var opts rimu.RenderOptions
	if htmlReplacement != nil {
		opts.HtmlReplacement = htmlReplacement
	}
	for _, infile := range files {
		var source string
		switch {
		case strings.HasPrefix(infile, RESOURCE_TAG):
			infile = infile[len(RESOURCE_TAG):]
			if (stringlist.StringList{"classic", "flex", "sequel", "v8"}).IndexOf(layout) >= 0 {
				source = readResourceFile(infile)
			} else {
				source = importLayoutFile(infile)
			}
			opts.SafeMode = 0 // Resources are trusted.
		case infile == STDIN:
			bytes, _ := ioutil.ReadAll(os.Stdin)
			source = string(bytes)
			opts.SafeMode = safeMode
		case infile == PREPEND:
			source = prepend
			opts.SafeMode = 0 // --prepend options are trusted.
		default:
			if !fileExists(infile) {
				die("source file does not exist: " + infile)
			}
			bytes, err := ioutil.ReadFile(infile)
			if err != nil {
				die(err.Error())
			}
			source = string(bytes)
			// Prepended and ~/.rimurc files are trusted.
			if prependFiles.IndexOf(infile) > -1 {
				opts.SafeMode = 0
			} else {
				opts.SafeMode = safeMode
			}
		}
		// Skip .html and pass-through inputs.
		if !(strings.HasSuffix(infile, ".html") || (pass && infile == STDIN)) {
			opts.Callback = func(message rimu.CallbackMessage) {
				f := infile
				if infile == STDIN {
					f = "/dev/stdin"
				}
				msg := message.Kind + ": " + f + ": " + message.Text
				if len(msg) > 120 {
					msg = msg[:117] + "..."
				}
				fmt.Fprintln(os.Stderr, msg)
				if message.Kind == "error" {
					errors++
				}
			}
			source = rimu.Render(source, opts)
		}
		source = strings.Trim(source, " \n")
		if source != "" {
			output += source + "\n"
		}
	}
	output = strings.Trim(output, " \n")
	if outfile == "" || outfile == "-" {
		fmt.Print(output)
	} else {
		err := ioutil.WriteFile(outfile, []byte(output), 0644)
		if err != nil {
			die(err.Error())
		}
	}
	if errors > 0 {
		osExit(1)
	}
	osExit(0)
}
