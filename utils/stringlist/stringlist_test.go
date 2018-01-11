package stringlist

import (
	"testing"
)

/* Helpers */
// Compares list to expected list and reports error if they are not equal.
func eqTest(t *testing.T, list, expected []string) {
	if len(list) != len(expected) {
		t.Errorf("len(list) == %q, expected %q", len(list), len(expected))
	} else {
		for i, v := range list {
			if v != expected[i] {
				t.Errorf("list[%d] == %q, expected %q", i, v, expected[i])
			}
		}
	}
}

func TestStackMutators(t *testing.T) {
	list := StringList{}
	eqTest(t, list, StringList{})
	list.Push("foo")
	eqTest(t, list, StringList{"foo"})
	list.Push("bar")
	eqTest(t, list, StringList{"foo", "bar"})
	s := list.Pop()
	if s != "bar" {
		t.Errorf("Pop() == %q, expected %q", s, "bar")
	}
	eqTest(t, list, StringList{"foo"})
	list = list.Concat([]string{"foo", "bar"})
	eqTest(t, list, StringList{"foo", "foo", "bar"})
	list.Unshift("pre")
	eqTest(t, list, StringList{"pre", "foo", "foo", "bar"})
	s = list.Shift()
	if s != "pre" {
		t.Errorf("Shift() == %q, expected %q", s, "pre")
	}
	eqTest(t, list, StringList{"foo", "foo", "bar"})
	list = list.Concat(list)
	eqTest(t, list, StringList{"foo", "foo", "bar", "foo", "foo", "bar"})
}

func TestCollectionFunctions(t *testing.T) {
	list := StringList{"x", "y", "z", "x", "x"}
	got := list.Filter(func(s string) bool { return s != "x" })
	eqTest(t, got, []string{"y", "z"})
	got = got.Map(func(s string) string { return s + s })
	eqTest(t, got, []string{"yy", "zz"})
	if list.IndexOf("z") != 2 {
		t.Errorf("IndexOf(\"z\") == %q, expected %q", list.IndexOf("z"), 2)
	}
	if list.IndexOf("XXX") != -1 {
		t.Errorf("IndexOf(\"XXX\") == %q, expected %q", list.IndexOf("XXX"), -1)
	}
	if list.Contains("XXX") {
		t.Errorf("Contains(\"XXX\") == %t, expected %t", list.Contains("XXX"), false)
	}
	b := list.Any(func(s string) bool { return s == "XXX" })
	if b {
		t.Errorf("Any(\"XXX\") == %t, expected %t", b, false)
	}
	b = !list.Any(func(s string) bool { return s == "z" })
	if b {
		t.Errorf("Any(\"z\") == %t, expected %t", b, true)
	}
	b = list.All(func(s string) bool { return s == "XXX" })
	if b {
		t.Errorf("All(s == \"XXX\") == %t, expected %t", b, false)
	}
	b = !list.All(func(s string) bool { return len(s) == 1 })
	if b {
		t.Errorf("All(len(s) ==1) == %t, expected %t", b, true)
	}
}
