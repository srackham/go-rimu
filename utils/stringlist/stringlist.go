package stringlist

// StringList type to host stack and collection methods.
type StringList []string

/*
  Stack mutator functions.
  See https://github.com/golang/go/wiki/SliceTricks
*/

// Push onto end of list.
func (list *StringList) Push(s string) {
	*list = append(*list, s)
}

// Pop from end of list.
func (list *StringList) Pop() (result string) {
	result, *list = (*list)[len(*list)-1], (*list)[:len(*list)-1]
	return
}

// Shift from start of list.
func (list *StringList) Shift() (result string) {
	result, *list = (*list)[0], (*list)[1:]
	return
}

// Unshift onto start of list.
func (list *StringList) Unshift(s string) {
	*list = append([]string{s}, *list...)
}

/*
  Collection functions.
  See https://gobyexample.com/collection-functions
*/

// Returns the first index of the target string `t`, or
// -1 if no match is found.
func (list StringList) IndexOf(t string) int {
	for i, v := range list {
		if v == t {
			return i
		}
	}
	return -1
}

// Returns `true` if the target string t is in the
// slice.
func (list StringList) Contains(t string) bool {
	return list.IndexOf(t) >= 0
}

// Returns `true` if one of the strings in the slice
// satisfies the predicate `f`.
func (list StringList) Any(f func(string) bool) bool {
	for _, v := range list {
		if f(v) {
			return true
		}
	}
	return false
}

// Returns `true` if all of the strings in the slice
// satisfy the predicate `f`.
func (list StringList) All(f func(string) bool) bool {
	for _, v := range list {
		if !f(v) {
			return false
		}
	}
	return true
}

// Returns a new slice containing all strings in the
// slice that satisfy the predicate `f`.
func (list StringList) Filter(f func(string) bool) StringList {
	result := make([]string, 0)
	for _, v := range list {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// Returns a new slice containing the results of applying
// the function `f` to each string in the original slice.
func (list StringList) Map(f func(string) string) StringList {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = f(v)
	}
	return result
}

// Returns a new slice containing the concatenation of the receiver and values.
func (list StringList) Concat(values ...string) StringList {
	return append(list, values...)
}

// Returns a new slice containing the receiver with values inserted at index.
func (list StringList) InsertAt(index int, values ...string) StringList {
	return append(list[:index], append(values, list[index:]...)...)
}
