# Rimu Markup for Go

**IMPORTANT**: `go-rimu` has not yet reached 1.0 and there are no
stability guarantees.

_go-rimu_ is a port of the [Rimu Markup
language](http://rimumarkup.org) written in the Go language.


## Features
Functionally identical to the [JavaScript
implementation](https://github.com/srackham/rimu) version 10.4.2 with
the following exceptions:

  * Does not support _Expression macro values_.
  * Because the Go `regexp` package uses RE2 regular expressions there
    are some limitations on Replacements definitions and
    Inclusion/Exclusion macro invocations.


## Using the go-rimu library
Install with:

    go get -u github.com/srackham/go-rimu/...

Example usage:

``` go
package main

import (
    "fmt"

    "github.com/srackham/go-rimu/rimu"
)

func main() {
    // Prints "<p><em>Hello Rimu!</em></p>""
    fmt.Println(rimu.Render("*Hello Rimu!*", rimu.RenderOptions{}))
}
```


## rimuc compiler command
The executable is named `rimucgo` and is functionally identical to the
[JavaScript rimuc](http://rimumarkup.org/reference.html#rimuc-command)
command-line compiler.


## Implementation
- The largely one-to-one correspondence between the canonical
  [TypeScript code](https://github.com/srackham/rimu) and the Go code
  eased porting and debugging.  This will also make it easier to
  cross-port new features and bug-fixes.

- TypeScript-style namespaces are implemented as Go packages.

- Both the Go and JavaScript implementations share the same JSON
  driven test suites comprising over 250 compatibility checks.

