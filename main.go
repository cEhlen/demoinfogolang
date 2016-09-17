package main

import (
	"flag"

	"github.com/cehlen/demoinfogolang/demo"
)

func main() {
	filename := flag.String("f", "", "File Path")
	flag.Parse()

	var d demo.Dumper
	d.Open(*filename)
	d.PrintHeader()
	d.Dump()
}
