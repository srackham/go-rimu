package main

import (
	"github.com/srackham/rimu-go/rimu"
)

func main() {
	println(rimu.Render("Hello *rimu-go!*", rimu.RenderOptions{}))
}
