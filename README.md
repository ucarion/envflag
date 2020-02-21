# envflag [![godoc][godoc-badge]][godoc-url] [![ci][ci-badge]][ci-url]

[godoc-badge]: https://godoc.org/github.com/ucarion/envflag?status.svg
[godoc-url]: https://godoc.org/github.com/ucarion/envflag
[ci-badge]: https://github.com/ucarion/envflag/workflows/CI/badge.svg?branch=master
[ci-url]: https://github.com/ucarion/envflag/actions

`envflag` is a Golang package that enhances the standard library's `flag`
package with the ability to read from environment variables. Just change:

```golang
flag.Parse()
```

To:

```golang
envflag.Parse()
flag.Parse()
```

And you'll get env variable goodness. No further changes required.

## Installation

You can add this package using `go get` as follows:

```bash
go get github.com/ucarion/envflag
```

## Example

Here is your typical example usage of `flag`:

```golang
package main

import (
  "flag"
  "fmt"
)

func main() {
  foo := flag.String("foo", "asdf", "some string param")
  bar := flag.Int("bar", 123, "some int param")

  flag.Parse()

  fmt.Println("foo", *foo)
  fmt.Println("bar", *bar)
}
```

Here is how you convert that into also using `envflag`:

```golang
package main

import (
  "flag"
  "fmt"

  "github.com/ucarion/envflag"
)

func main() {
  foo := flag.String("foo", "asdf", "some string param")
  bar := flag.Int("bar", 123, "some int param")

  envflag.Parse()
  flag.Parse()

  fmt.Println("foo", *foo)
  fmt.Println("bar", *bar)
}
```

Assuming you put this in a file called `./examples/simple/main.go` ([like the
one you can find in this repo](./examples/simple/main.go)), you can invoke it as
so:

```text
$ go run ./examples/simple/...
foo asdf
bar 123

$ SIMPLE_FOO=from-env go run ./examples/simple/...
foo from-env
bar 123

$ SIMPLE_FOO=from-env go run ./examples/simple/... --foo=from-args
foo from-args
bar 123

$ SIMPLE_FOO=from-env SIMPLE_BAR=456 go run ./examples/simple/...
foo from-env
bar 456
```

The env variables are prefixed with `SIMPLE_`, because that's the basename of
`os.Args[0]`. Prefixing env variables like this helps you keep your config
separate from others.

If you would prefer to disable this prefixing, instead of doing:

```golang
envflag.Parse()
```

Do:

```golang
// The first parameter to ParseFlagSet is a prefix for all env variables.
//
// The empty string disables prefixing env variables.
envflag.ParseFlagSet("", flag.CommandLine)
```
