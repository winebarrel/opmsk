package main

import (
	"io"
	"log"
	"os"

	"github.com/winebarrel/opmsk"
)

func init() {
	log.SetFlags(0)
}

func main() {
	opts := parseOptions()

	var in io.Reader

	if opts.file == "" || opts.file == "-" {
		in = os.Stdin
	} else {
		file, err := os.OpenFile(opts.file, os.O_RDONLY, 0666)

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
		in = file
	}

	err := opmsk.Mask(in, os.Stdout)

	if err != nil {
		log.Fatal(err)
	}
}
