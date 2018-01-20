package rimu

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	_ "github.com/srackham/rimu-go/delimitedblocks"
	_ "github.com/srackham/rimu-go/lineblocks"
	_ "github.com/srackham/rimu-go/lists"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, c.Expected, got)
	}
}
