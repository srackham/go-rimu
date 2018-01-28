# Rimu Markup for Go

_go-rimu_ is a port of the [Rimu Markup
language](http://rimumarkup.org) written in the Go language.


## Features

- Functionally identical to the [JavaScript
  implementation](https://github.com/srackham/rimu) version
  10.4.2 with the following exceptions:

  * Does not support _Expression macro values_.
  * Uses Go RE2 regular expressions which places some limitations
    on Replacements definitions and Inclusion/Exclusion macro invocations.

- Includes
  [rimucgo](http://rimumarkup.org/reference.html#rimuc-command)
  command-line compiler.
- Single build dependency: TODO.

Details:

- Github source repo: https://github.com/srackham/go-rimu


## Using the go-rimu library
Install with:

  go get

Example code: TODO


## rimuc compiler command
The executable is named `rimucgo`.


## Implementation
The largely one-to-one correspondence between the canonical
[TypeScript code](https://github.com/srackham/rimu) and the Go
code eased porting and debugging.  This will also make it easier to
cross-port new features and bug-fixes.

TypeScript-style namespaces are implemented as Go packages.

Both the Go and JavaScript implementations share the same JSON
driven test suites comprising over 250 compatibility checks.

