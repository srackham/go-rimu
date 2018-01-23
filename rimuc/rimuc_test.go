package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	_ "github.com/srackham/rimu-go/spans"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Description string `json:"description"`
	Args        string `json:"args"`
	Input       string `json:"input"`
	Expected    string `json:"expectedOutput"`
	Predicate   string `json:"predicate"`
	ExitCode    int    `json:"exitCode"`
}

// Convert command-line arguments string to array of arguments.
// Arguments can be double-quoted.
func parseArgs(args string) []string {
	result := regexp.MustCompile(`".+?"|\S+`).FindAllString(args, -1)
	// Strip double quotes.
	for i, v := range result {
		if len(v) >= 2 && v[0] == '"' && v[len(v)-1] == '"' {
			result[i] = v[1 : len(v)-1]
		}
	}
	return result
}

// ??? TODO put in local package utils/str
func TrimQuotes(s string, quote string) string {
	if len(s) >= 2*len(quote) && strings.HasPrefix(s, quote) && strings.HasSuffix(s, quote) {
		return strings.TrimPrefix(strings.TrimSuffix(s, quote), quote)
	}
	return s
}

func TestRimuc(t *testing.T) {
	// Read JSON test cases.
	raw, err := ioutil.ReadFile("./testdata/rimuc-tests.json")
	if err != nil {
		t.Error(err.Error())
		return
	}
	var cases []TestCase
	json.Unmarshal(raw, &cases)
	// Run test cases.
	for _, c := range cases {
		// Save and set os.Exit mock to capture exit code
		// (see https://stackoverflow.com/a/40801733 and https://npf.io/2015/06/testing-exec-command/).
		exitCode := 0
		savedExit := osExit
		osExit = func(code int) {
			exitCode = code
		}
		// Save and set command-line arguments.
		savedArgs := os.Args
		os.Args = append([]string{"rimuc"}, parseArgs(c.Args)...)
		// Capture sdtout (see https://stackoverflow.com/a/29339052).
		savedStdout := os.Stdout
		defer func() { os.Stdout = savedStdout }() // Ensures restore after a panic.
		rout, wout, _ := os.Pipe()
		os.Stdout = wout
		// Capture sdterr.
		savedStderr := os.Stderr
		defer func() { os.Stderr = savedStderr }() // Ensures restore after a panic.
		rerr, werr, _ := os.Pipe()
		os.Stderr = werr
		// Execute rimuc.
		main()
		// Get stdout and stderr.
		wout.Close()
		werr.Close()
		bytes, _ := ioutil.ReadAll(rout)
		out := string(bytes)
		bytes, _ = ioutil.ReadAll(rerr)
		out += string(bytes)
		// Restore mocks.
		os.Stdout = savedStdout
		os.Stderr = savedStderr
		os.Args = savedArgs
		osExit = savedExit
		// Test outputs and exit code.
		assert.Equal(t, c.Expected, out)
		assert.Equal(t, c.ExitCode, exitCode)
	}
}
