// Package proxies exists solely to eliminate import cycles stemming from the importation of the api package.
package proxies

// Function pointers initialised by api.init().
var ApiInit func()
var ApiRender func(source string) string
