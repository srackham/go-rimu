package main

import (
	"fmt"
	"os"

	"github.com/srackham/rimu-go/rimu"
)

// Mocked by tests.
var osExit = os.Exit

func main() {
	fmt.Print(rimu.Render("*Hello World!*", rimu.RenderOptions{}))
	osExit(0)
}
