package spans

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	tests := []struct {
		source string
		want   string
	}{
		{"", ""},
		{"foo", "foo"},
		// Quotes.
		{"*foo*", "<em>foo</em>"},
		{"**foo**", "<strong>foo</strong>"},
		{"*foo* **bar**", "<em>foo</em> <strong>bar</strong>"},
		{"*foo __bar__*", "<em>foo <strong>bar</strong></em>"},
		{"`**foo**`", "<code>**foo**</code>"},
		// Replacements.
		{"<image:foo|bar>", `<img src="foo" alt="bar">`},
		{"<image:foo|bar\nboo>", "<img src=\"foo\" alt=\"bar\nboo\">"},
	}
	for _, tt := range tests {
		got := Render(tt.source)
		assert.Equal(t, tt.want, got)
	}
}

func Test_defrag(t *testing.T) {
	tests := []struct {
		frags []fragment
		want  string
	}{
		{
			frags: []fragment{
				{text: ""},
			},
			want: "",
		},
		{
			frags: []fragment{
				{text: "foo"},
				{text: "bar"},
			},
			want: "foobar",
		},
	}
	for _, tt := range tests {
		got := defrag(tt.frags)
		assert.Equal(t, tt.want, got)
	}
}

/*
func Test_fragQuotes(t *testing.T) {
	type args struct {
		fragments []fragment
	}
	tests := []struct {
		name string
		args args
		want []fragment
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fragQuotes(tt.args.fragments); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fragQuotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fragQuote(t *testing.T) {
	type args struct {
		frag fragment
	}
	tests := []struct {
		name       string
		args       args
		wantResult []fragment
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := fragQuote(tt.args.frag); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("fragQuote() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_preReplacements(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := preReplacements(tt.args.text); gotResult != tt.wantResult {
				t.Errorf("preReplacements() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_postReplacements(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := postReplacements(tt.args.text); got != tt.want {
				t.Errorf("postReplacements() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fragReplacements(t *testing.T) {
	type args struct {
		fragments []fragment
	}
	tests := []struct {
		name       string
		args       args
		wantResult []fragment
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := fragReplacements(tt.args.fragments); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("fragReplacements() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_fragReplacement(t *testing.T) {
	type args struct {
		frag fragment
		def  replacements.Definition
	}
	tests := []struct {
		name       string
		args       args
		wantResult []fragment
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := fragReplacement(tt.args.frag, tt.args.def); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("fragReplacement() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
*/

func Test_fragSpecials(t *testing.T) {
	tests := []struct {
		frags []fragment
		want  string
	}{
		{
			frags: []fragment{
				{text: "&<>"},
			},
			want: "&amp;&lt;&gt;",
		},
		{
			frags: []fragment{
				{text: "<foo>"},
				{text: "<bar>"},
			},
			want: "&lt;foo&gt;&lt;bar&gt;",
		},
	}
	for _, tt := range tests {
		frags := fragSpecials(tt.frags)
		got := defrag(frags)
		assert.Equal(t, tt.want, got)
	}
}
