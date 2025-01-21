# Rimu Markup for Go

_go-rimu_ is a port of the [Rimu Markup
language](https://github.com/srackham/rimu) written in the Go language.

## Features

Functionally identical to the [JavaScript
implementation](https://github.com/srackham/rimu) version 11.4 with
the following exceptions:

-   Because the Go `regexp` package uses RE2 regular expressions there are
    [some limitations](https://srackham.github.io/rimu/reference.html#regular-expressions)
    on the regular expressions used in Replacements definitions and
    Inclusion/Exclusion macro invocations.

## Installation

Download, build, test and install (requires Go 1.17 or better):

    git clone https://github.com/srackham/go-rimu.git
    cd go-rimu
    make

## Using the go-rimu library

Example usage:

```go
package main

import (
    "fmt"

    "github.com/srackham/go-rimu/v11/rimu"
)

func main() {
    // Prints "<p><em>Hello Rimu</em>!</p>"
    fmt.Println(rimu.Render("*Hello Rimu*!", rimu.RenderOptions{}))
}
```

To compile and run this simple application:

1.  Copy the code above to a file named `hello-rimu.go` and put it in an empty
    directory.
2.  Change to the directory and run the following Go commands:

        go mod init example.com/hello-rimu
        go mod tidy
        go run hello-rimu.go

See also Rimu
[API documentation](https://srackham.github.io/rimu/reference.html#api).

## Rimu CLI command

The [Rimu CLI command](https://srackham.github.io/rimu/reference.html#rimuc-command) is named
`rimugo`.

## Learn more

Read the [documentation](https://srackham.github.io/rimu/reference.html) and
experiment with Rimu in the [Rimu
Playground](http://srackham.github.io/rimu/rimuplayground.html).

See the Rimu [Change
Log](http://srackham.github.io/rimu/changelog.html) for the latest
changes.

## Implementation

-   The largely one-to-one correspondence between the canonical
    [TypeScript code](https://github.com/srackham/rimu) and the Go code
    eased porting and debugging. This will also make it easier to
    cross-port new features and bug-fixes.

-   All Rimu implementations share the same JSON driven test suites
    comprising over 300 compatibility checks.
