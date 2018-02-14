package rimu

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type renderTest struct {
	Description string     `json:"description"`
	Input       string     `json:"input"`
	Expected    string     `json:"expectedOutput"`
	Callback    string     `json:"expectedCallback"`
	Options     apiOptions `json:"options"`
	Unsupported string     `json:"unsupported,omitempty"`
}

type apiOptions struct {
	SafeMode        int    `json:"safeMode,omitempty"`
	HtmlReplacement string `json:"htmlReplacement,omitempty"`
	Reset           bool   `json:"reset,omitempty"`
}

func TestRender(t *testing.T) {
	// Read JSON test cases.
	raw, err := ioutil.ReadFile("./testdata/rimu-tests.json")
	if err != nil {
		t.Error(err.Error())
		return
	}
	var tests []renderTest
	json.Unmarshal(raw, &tests)
	// Append test with invalid UTF-8 input because JSON does not support binary data (all strings are valid UTF-8)).
	tests = append(tests, renderTest{
		Description: "Invalid UTF-8 input",
		Input:       string([]byte("\xbb")),
		Expected:    "",
		Callback:    "error: invalid UTF-8 input",
		Options:     apiOptions{Reset: true},
	})
	// Run test cases.
	for _, tt := range tests {
		// ioutil.WriteFile(fmt.Sprintf("./testdata/fuzz-samples/sample-%03d.txt", i), []byte(tt.Input), 0644)
		// if tt.Description != "illegal render() safeMode value" {
		// 	continue
		// }
		if strings.Contains(tt.Unsupported, "go") {
			continue
		}
		msg := ""
		opts := RenderOptions{
			Reset:           tt.Options.Reset,
			SafeMode:        tt.Options.SafeMode,
			HtmlReplacement: tt.Options.HtmlReplacement,
			Callback:        func(message CallbackMessage) { msg += message.Kind + ": " + message.Text + "\n" },
		}
		// fmt.Println("Description: ", tt.Description)
		got := Render(tt.Input, opts)
		if got != tt.Expected {
			t.Errorf("\n%-15s: %s\n%-15s: %s\n%-15s: %s\n%-15s: %s\n\n",
				"description", tt.Description, "input", tt.Input, "expected", tt.Expected, "got", got)
		}
		switch {
		case tt.Callback != "":
			assert.Equal(t, tt.Callback, strings.TrimSpace(msg))
		case msg != "":
			t.Errorf("unexpected callback: %s", msg)
		}
	}
}

func BenchmarkSmall(b *testing.B) {
	text, err := ioutil.ReadFile("./testdata/benchmark-small.rmu")
	if err != nil {
		panic(err.Error())
	}
	for n := 0; n < b.N; n++ {
		Render(string(text), RenderOptions{})
	}
}
