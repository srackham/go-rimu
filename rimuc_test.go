package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

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

// TODO put in local package utils/str
func TrimQuotes(s string, quote string) string {
	if len(s) >= 2*len(quote) && strings.HasPrefix(s, quote) && strings.HasSuffix(s, quote) {
		return strings.TrimPrefix(strings.TrimSuffix(s, quote), quote)
	}
	return s
}

func TestRimuc(t *testing.T) {
	cases := []struct {
		description string
		args        string
		input       string
		expected    string
		exitCode    int
	}{
		{"rimuc basic test", `-p "foo bar" "baz qux"`, "*Hello World!", "<p><em>Hello World!</em></p>", 0},
		{"rimuc illegal option test", `--illegal`, "*Hello World!", "illegal option: --illegal", 1},
	}
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
		os.Args = append([]string{"rimuc"}, parseArgs(c.args)...)
		// Capture sdtout from main() (see https://stackoverflow.com/a/29339052).
		savedStdout := os.Stdout
		defer func() { os.Stdout = savedStdout }() // Ensures restore after a panic.
		rout, wout, _ := os.Pipe()
		os.Stdout = wout

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
		if out != c.expected {
			t.Errorf("%s: got %q, expected %q", c.description, out, c.expected)
		}
		if exitCode != c.exitCode {
			t.Errorf("%s: exit code %d, expected %d", c.description, exitCode, c.exitCode)
		}
	}
}
