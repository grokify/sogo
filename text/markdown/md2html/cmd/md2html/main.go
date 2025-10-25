package main

import (
	"fmt"
	"os"

	"github.com/grokify/sogo/text/markdown/md2html"
)

func main() {
	var srcfile string
	var outfile string

	if len(os.Args) < 2 {
		fmt.Println("Please provide a name of an input file.")
		os.Exit(1)
	} else {
		srcfile = os.Args[1]
	}
	if len(os.Args) >= 3 {
		outfile = os.Args[2]
	} else {
		outfile = srcfile + ".html"
	}

	err := md2html.MarkdownToHTMLFile(srcfile, outfile, 0600)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	fmt.Println("DONE")
	os.Exit(0)
}
