package main

import (
	"fmt"

	"github.com/srackham/go-rimu/v11/rimu"
)

func main() {
	// Prints "<p><em>Hello Rimu</em>!</p>"
	fmt.Println(rimu.Render("*Hello Rimu*!", rimu.RenderOptions{}))
}
