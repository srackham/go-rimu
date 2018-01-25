package rimu

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	_ "github.com/srackham/rimu-go/spans"
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
	Reset bool `json:"reset"`
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
		if tt.Unsupported != "" {
			continue
		}
		opts := RenderOptions{Reset: tt.Options.Reset}
		got := Render(tt.Input, opts)
		assert.Equal(t, tt.Expected, got)
	}
}
