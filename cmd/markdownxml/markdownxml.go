package main

import (
	"io/ioutil"
	"os"

	"github.com/jsternberg/markdownxml"
	flag "github.com/spf13/pflag"
	"gopkg.in/russross/blackfriday.v2"
)

func main() {
	flag.Parse()
	args := flag.Args()
	in, err := ioutil.ReadFile(args[0])
	if err != nil {
		panic(err)
	}

	out := blackfriday.Run(in,
		blackfriday.WithRenderer(markdownxml.NewRenderer()),
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
	)
	os.Stdout.Write(out)
}
