/*
  Command-lne app to convert Rimu source to HTML.
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

	"github.com/srackham/go-rimu/v11/internal/utils/stringlist"
	"github.com/srackham/go-rimu/v11/rimu"
)

var osExit = os.Exit // Mocked by tests.

const VERSION = "11.1.5"
const STDIN = "-"

// rimurcPath returns path of $HOME/.rimurc file.
// Return "" if $HOME not found.
func rimurcPath() (result string) {
	if user, err := user.Current(); err == nil {
		result = filepath.Join(user.HomeDir, ".rimurc")
	}
	return
}

// MockExit type used by rimugo_test osExit() mock.
type MockExit struct{}

// Helpers.
func die(message string) {
	if message != "" {
		fmt.Fprintln(os.Stderr, message)
	}
	osExit(1)
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func readResourceFile(name string) (result string) {
	data, err := Asset("resources/" + name)
	if err != nil {
		panic("missing resource file: " + name)
	}
	result = string(data)
	return
}

func importLayoutFile(name string) string {
	// External layouts not supported in go-rimu.
	die("missing --layout: " + name)
	return ""
}

func main() {
	defer func() {
		// Ignore the panic if it's a MockExit (i.e. we're running a test).
		r, ok := recover().(MockExit)
		if !ok {
			panic(r)
		}
	}()
	args := stringlist.StringList(os.Args)
	args.Shift() // Skip program name.
	nextArg := func(err string) string {
		if len(args) == 0 {
			die(err)
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
			fmt.Printf("\n" + readResourceFile("manpage.txt") + "\n")
			osExit(0)
		case "--version":
			fmt.Printf(VERSION + "\n")
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
			prepend += "{" + arg + "}='" + macroValue + "'\n"
		case "--layout",
			"--styled-name": // Deprecated in Rimu 10.0.0
			layout = nextArg("missing --layout value")
			prepend += "{--header-ids}='true'\n"
		case "--styled", "-s":
			prepend += "{--header-ids}='true'\n"
			prepend += "{--no-toc}='true'\n"
			layout = "sequel"
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
			if (stringlist.StringList{"classic", "flex", "plain", "sequel", "v8"}).IndexOf(layout) >= 0 {
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
		source = strings.TrimSpace(source)
		if source != "" {
			output += source + "\n"
		}
	}
	output = strings.TrimSpace(output)
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
