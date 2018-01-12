package main

import (
	"io/ioutil"
	"os"
	"regexp"
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

func TestRimuc(t *testing.T) {
	cases := []struct {
		description string
		args        string
		input       string
		expected    string
		exitCode    int
	}{
		{"rimuc basic test", `-p "foo bar"`, "*Hello World!", "<p><em>Hello World!</em></p>", 0},
	}
	for _, c := range cases {
		exitCode := 0
		savedExit := osExit
		// defer func() { osExit = savedExit }()
		osExit = func(code int) {
			exitCode = code
		}
		// Save and set command-line arguments.
		savedArgs := os.Args
		os.Args = parseArgs(c.args)
		// Capture sdtout from main() (see https://stackoverflow.com/a/29339052).
		savedStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		main()
		w.Close()
		bytes, _ := ioutil.ReadAll(r)
		os.Stdout = savedStdout
		os.Args = savedArgs
		osExit = savedExit
		out := string(bytes)
		if out != c.expected {
			t.Errorf("%s: got %q, expected %q", c.description, out, c.expected)
		}
		if exitCode != c.exitCode {
			t.Errorf("%s: exit code %d, expected %d", c.description, exitCode, c.exitCode)
		}
	}
}
