package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

type rimucTest struct {
	Description string `json:"description"`
	Args        string `json:"args"`
	Input       string `json:"input"`
	Expected    string `json:"expectedOutput"`
	Predicate   string `json:"predicate"`
	ExitCode    int    `json:"exitCode,omitempty"`
	Unsupported string `json:"unsupported,omitempty"`
	Layouts     bool   `json:"layouts,omitempty"`
}

func TestRimuc(t *testing.T) {
	// Execute tests specified in JSON file.
	raw, err := os.ReadFile("./testdata/rimuc-tests.json")
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
			tt.Args = strings.Replace(tt.Args, "./examples/example-rimurc.rmu", "./testdata/example-rimurc.rmu", -1)
			command := "rimugo --no-rimurc"
			if layout != "" {
				command += " --layout " + layout
			}
			command += " " + tt.Args
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("PowerShell.exe", "-Command", command)
			} else {
				cmd = exec.Command("bash", "-c", command)
			}
			var outb, errb bytes.Buffer
			cmd.Stdout = &outb
			cmd.Stderr = &errb
			stdin, err := cmd.StdinPipe()
			if err != nil {
				panic(err.Error())
			}
			go func() {
				defer stdin.Close()
				io.WriteString(stdin, tt.Input)
			}()
			err = cmd.Run()
			exitCode := 0
			if err != nil {
				exitCode = 1
			}
			out := errb.String() + outb.String()
			out = strings.Replace(out, "\r", "", -1) // Strip Windows return characters.
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
