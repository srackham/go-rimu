package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/srackham/go-rimu/internal/api"
)

type rimucTest struct {
	Description string `json:"description"`
	Args        string `json:"args"`
	Input       string `json:"input"`
	Expected    string `json:"expectedOutput"`
	Predicate   string `json:"predicate"`
	ExitCode    int    `json:"exitCode"`
	Unsupported string `json:"unsupported"`
	Layouts     bool   `json:"layouts"`
}

// Convert command-line arguments string to array of arguments.
// Arguments can be double-quoted.
func parseArgs(args string) (result []string) {
	result = regexp.MustCompile(`".+?"|\S+`).FindAllString(args, -1)
	// Strip double quotes.
	for i, v := range result {
		if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
			result[i] = v[1 : len(v)-1]
		}
	}
	return
}

// TestMain is a test harness for main(). It mocks stdin, stdout, stderr and os.Exit().
// This allows debuggable tests to be run in the executable environment.
func TestMain(t *testing.T) {
	// Read JSON test cases.
	raw, err := ioutil.ReadFile("./testdata/rimuc-tests.json")
	if err != nil {
		t.Error(err.Error())
		return
	}
	var tests []rimucTest
	json.Unmarshal(raw, &tests)
	// Run test cases.
	for _, tt := range tests {
		if strings.Contains(tt.Unsupported, "go") {
			continue
		}
		for _, layout := range []string{"", "classic", "flex", "sequel"} {
			// Skip if not a layouts test and we have a layout, or if it is a layouts test but no layout is specified.
			if !tt.Layouts && layout != "" || tt.Layouts && layout == "" {
				continue
			}
			tt.Expected = strings.Replace(tt.Expected, "./test/fixtures/", "./testdata/", -1)
			tt.Args = strings.Replace(tt.Args, "./test/fixtures/", "./testdata/", -1)
			tt.Args = strings.Replace(tt.Args, "./src/examples/example-rimurc.rmu", "./testdata/example-rimurc.rmu", -1)
			// Save and set os.Exit mock to capture exit code
			// (see https://stackoverflow.com/a/40801733 and https://npf.io/2015/06/testing-exec-command/).
			exitCode := 0
			savedExit := osExit
			defer func() { osExit = savedExit }()
			osExit = func(code int) {
				exitCode = code
				panic(MockExit{})
			}
			// Save and set command-line arguments.
			savedArgs := os.Args
			defer func() { os.Args = savedArgs }()
			os.Args = []string{"rimugo", "--no-rimurc"}
			if layout != "" {
				os.Args = append(os.Args, "--layout", layout)
			}
			os.Args = append(os.Args, parseArgs(tt.Args)...)
			// Capture sdtout (see https://stackoverflow.com/a/29339052).
			savedStdout := os.Stdout
			defer func() { os.Stdout = savedStdout }()
			rout, wout, _ := os.Pipe()
			os.Stdout = wout
			// Capture sdterr.
			savedStderr := os.Stderr
			defer func() { os.Stderr = savedStderr }()
			rerr, werr, _ := os.Pipe()
			os.Stderr = werr
			// Mock stdin.
			savedStdin := os.Stdin
			defer func() { os.Stdin = savedStdin }()
			rin, win, _ := os.Pipe()
			os.Stdin = rin
			win.WriteString(tt.Input)
			win.Close()
			// Execute rimugo.
			api.Init()
			main()
			// Get stdout and stderr.
			wout.Close()
			werr.Close()
			bytes, _ := ioutil.ReadAll(rout)
			out := string(bytes)
			bytes, _ = ioutil.ReadAll(rerr)
			out += string(bytes)
			// Restore mocks.
			os.Stdin = savedStdin
			os.Stdout = savedStdout
			os.Stderr = savedStderr
			os.Args = savedArgs
			osExit = savedExit
			// Test outputs and exit code.
			passed := false
			switch tt.Predicate {
			case "contains":
				passed = strings.Contains(out, tt.Expected)
			case "!contains":
				passed = !strings.Contains(out, tt.Expected)
			case "equals":
				passed = out == tt.Expected
			case "!equals":
				passed = out != tt.Expected
			case "startsWith":
				passed = strings.HasPrefix(out, tt.Expected)
			default:
				panic(tt.Description + ": illegal predicate: " + tt.Predicate)
			}
			if !passed {
				t.Errorf("\n%-15s: %s\n%-15s: %s\n%-15s: %s\n%-15s: %s\n%-15s: %s\n%-15s: %s\n\n",
					"description", tt.Description, "args", tt.Args, "input", tt.Input,
					"predicate:", tt.Predicate, "expected", tt.Expected, "got", out)
			}
			if exitCode != tt.ExitCode {
				t.Errorf("\n%-15s: %s\n%-15s: %d (expected %d)\n\n",
					"description", tt.Description, "exitcode", exitCode, tt.ExitCode)
			}
		}
	}
}
