package stringlist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackMutators(t *testing.T) {
	list := StringList{}
	assert.EqualValues(t, []string{}, list)
	list.Push("foo")
	assert.EqualValues(t, []string{"foo"}, list)
	list.Push("bar")
	assert.EqualValues(t, []string{"foo", "bar"}, list)
	s := list.Pop()
	assert.Equal(t, "bar", s)
	assert.EqualValues(t, []string{"foo"}, list)
	list = list.Concat([]string{"foo", "bar"}...)
	assert.EqualValues(t, []string{"foo", "foo", "bar"}, list)
	list.Unshift("pre")
	assert.EqualValues(t, []string{"pre", "foo", "foo", "bar"}, list)
	s = list.Shift()
	assert.Equal(t, "pre", s)
	assert.EqualValues(t, []string{"foo", "foo", "bar"}, list)
}

func TestCollectionFunctions(t *testing.T) {
	list := StringList{"x", "y", "z", "x", "x"}
	got := list.Filter(func(s string) bool { return s != "x" })
	assert.EqualValues(t, []string{"y", "z"}, got)
	got = got.Map(func(s string) string { return s + s })
	assert.EqualValues(t, []string{"yy", "zz"}, got)
	assert.EqualValues(t, 2, list.IndexOf("z"))
	assert.EqualValues(t, -1, list.IndexOf("XXX"))
	assert.True(t, !list.Contains("XXX"))
	assert.True(t, !list.Any(func(s string) bool { return s == "XXX" }))
	assert.True(t, list.Any(func(s string) bool { return s == "z" }))
	assert.True(t, !list.All(func(s string) bool { return s == "XXX" }))
	assert.True(t, list.All(func(s string) bool { return len(s) == 1 }))
	assert.EqualValues(t, []string{"x", "y", "z", "x", "x", "x", "y", "z", "x", "x"}, list.Concat(list...))
	list = []string{"x", "y", "z"}
	assert.EqualValues(t, []string{"x", "A", "B", "y", "z"}, list.InsertAt(1, "A", "B"))
	assert.EqualValues(t, []string{"A", "x", "y", "z"}, list.InsertAt(0, "A"))
	assert.EqualValues(t, []string{"x", "y", "z", "A", "B"}, list.InsertAt(3, []string{"A", "B"}...))
}
