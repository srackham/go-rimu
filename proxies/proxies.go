// Package proxies exists solely to eliminate import cycles stemming from the importation of the api package.
package proxies

// Function pointers initialised by api.init().
var ApiInit func()
var ApiRender func(source string) string

func init() {
	// So we can use these functions in imported packages without incuring import cycle errors.
	// ApiInit = api.Init
	// ApiRender = api.Render
}
