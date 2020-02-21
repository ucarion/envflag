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
