package main

import (
	"flag"
	"os"
)

var fileFlag string
var extractFlag bool

func init() {
	flag.StringVar(&fileFlag, "f", "out.tar", "tar file name")
	flag.BoolVar(&extractFlag, "x", false, "extract tar")

	flag.Parse()
}

func main() {
	if len(flag.Args()) == 0 {
		println("files required")
		os.Exit(1)
	}

	if extractFlag {
		if err := toUntar(); err != nil {
			println("failed to extract tar file")
			os.Exit(1)
		}
	} else {
		if err := toTar(); err != nil {
			println("failed to create tar")
			os.Exit(1)
		}
	}
	println("done")
}

func toTar() error {
	//todo realize
	return nil
}

func toUntar() error {
	//todo realuze
	return nil
}
