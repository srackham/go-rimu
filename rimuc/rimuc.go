package main

import (
	"fmt"
	"os"

	"github.com/srackham/rimu-go/rimu"
)

var osExit = os.Exit // Mocked by tests.

func die(message string) {
	if message != "" {
		fmt.Fprint(os.Stderr, message)
	}
	osExit(1)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--illegal" {
		die("illegal option: --illegal")
	} else {
		output := rimu.Render("*Hello World!*", rimu.RenderOptions{})
		fmt.Print(output)
		osExit(0)
	}
}
