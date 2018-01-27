package rimu

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	_ "github.com/srackham/rimu-go/spans"
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
		// if tt.Description != "Block Attributes on Fenced Block attached to list item" {
		// 	continue
		// }
		if strings.Contains(tt.Unsupported, "go") {
			continue
		}
		opts := RenderOptions{Reset: tt.Options.Reset, SafeMode: tt.Options.SafeMode, HtmlReplacement: tt.Options.HtmlReplacement}
		// fmt.Println("Description: ", tt.Description)
		got := Render(tt.Input, opts)
		if got != tt.Expected {
			t.Errorf("\n%-15s: %s\n%-15s: %s\n%-15s: %s\n%-15s: %s\n\n",
				"description", tt.Description, "input", tt.Input, "expected", tt.Expected, "got", got)
		}
	}
}
