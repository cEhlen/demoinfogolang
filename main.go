package main

import (
	"flag"
)

func main() {
	filename := flag.String("f", "", "File Path")
	flag.Parse()

	var demo Dumper
	demo.Open(*filename)
	demo.PrintHeader()
	demo.Dump()
}
