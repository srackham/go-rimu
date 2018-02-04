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
	Unsupported string     `json:"unsupported"`
}

type apiOptions struct {
	SafeMode        int    `json:"safeMode"`
	HtmlReplacement string `json:"htmlReplacement"`
	Reset           bool   `json:"reset"`
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
	// Run test cases.
	for _, tt := range tests {
		// ioutil.WriteFile(fmt.Sprintf("./testdata/fuzz-samples/sample-%03d.txt", i), []byte(tt.Input), 0644)
		// if tt.Description != "Block Attributes on Fenced Block attached to list item" {
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
		if tt.Callback != "" {
			assert.Equal(t, tt.Callback, strings.TrimSpace(msg))
		}
	}
}
