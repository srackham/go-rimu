package utils

import (
	"testing"
)

func TestStringList(t *testing.T) {
	list := StringList{}
	list.Push("foo")
	if len(list.values) != 1 {
		t.Errorf("len(list.values) == %q, expected %q", len(list.values), 1)
	}
	list.Push("bar")
	if len(list.values) != 2 {
		t.Errorf("len(list.values) == %q, expected %q", len(list.values), 2)
	}
	s := list.Pop()
	if s != "bar" {
		t.Errorf("Pop() == %q, expected %q", s, "bar")
	}
	if len(list.values) != 1 {
		t.Errorf("len(list.values) == %q, expected %q", len(list.values), 1)
	}
	s = list.Pop()
	if s != "foo" {
		t.Errorf("Pop() == %q, expected %q", s, "foo")
	}
	if len(list.values) != 0 {
		t.Errorf("len(list.values) == %q, expected %q", len(list.values), 0)
	}
	list.AppendSlice([]string{"foo", "bar"})
	if len(list.values) != 2 {
		t.Errorf("len(list.values) == %q, expected %q", len(list.values), 2)
	}
	list.Unshift("pre")
	if len(list.values) != 3 {
		t.Errorf("len(list.values) == %q, expected %q", len(list.values), 3)
	}
	s = list.Shift()
	if s != "pre" {
		t.Errorf("Shift() == %q, expected %q", s, "pre")
	}
	if len(list.values) != 2 {
		t.Errorf("len(list.values) == %q, expected %q", len(list.values), 2)
	}
}
