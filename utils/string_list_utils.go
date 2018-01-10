package utils

/*
  String slice functions.
*/
// TODO

/*
  StringList class.
  See https://github.com/golang/go/wiki/SliceTricks
*/
// StringList TODO
type StringList struct {
	values []string
}

// AppendSlice TODO
func (list *StringList) AppendSlice(slice []string) {
	list.values = append(list.values, slice...)
}

// Push TODO
func (list *StringList) Push(s string) {
	list.values = append(list.values, s)
}

// Pop TODO
func (list *StringList) Pop() (result string) {
	result, list.values = list.values[len(list.values)-1], list.values[:len(list.values)-1]
	return
}

// Unshift TODO
func (list *StringList) Unshift(s string) {
	list.values = append([]string{s}, list.values...)
}

// Shift TODO
func (list *StringList) Shift() (result string) {
	result, list.values = list.values[0], list.values[1:]
	return
}
