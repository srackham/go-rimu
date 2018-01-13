package rimu

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

type TestCase struct {
	Description string      `json:"description"`
	Input       string      `json:"input"`
	Expected    string      `json:"expectedOutput"`
	Callback    string      `json:"expectedCallback"`
	Options     TestOptions `json:"options"`
}

type TestOptions struct {
	Reset bool `json:"reset"`
}

func TestRimuTestCases(t *testing.T) {
	// Read JSON test cases.
	raw, err := ioutil.ReadFile("./fixtures/rimu-tests.json")
	if err != nil {
		t.Error(err.Error())
		return
	}
	var cases []TestCase
	json.Unmarshal(raw, &cases)
	// Run test cases.
	for _, c := range cases {
		opts := RenderOptions{Reset: c.Options.Reset}
		got := Render(c.Input, opts)
		if got != c.Expected {
			t.Errorf("TestRimuTestCases(%q) == %q, expected %q", c.Input, got, c.Expected)
		}
	}
}
