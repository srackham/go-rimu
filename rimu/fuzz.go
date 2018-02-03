package rimu

// Fuzz is used by https://github.com/dvyukov/go-fuzz
func Fuzz(data []byte) int {
	err := false
	Render(string(data), RenderOptions{Callback: func(msg CallbackMessage) { err = true }})
	if !err {
		return 1 // Valid Rimu input.
	}
	return 0
}
