package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grokify/sogo/text/markdown/remark"
)

func main() {
	slides := remark.PresentationData{
		Slides: []remark.RemarkSlideData{
			{
				Layout:   "middle, center, inverse",
				Class:    "false",
				Markdown: "# Test Slide\n\nTest Remark Slide",
			},
			{
				Markdown: "# Test Slide\n\nTest Remark Slide",
			},
		},
	}
	html := remark.RemarkHTML(slides)
	fmt.Println(html)

	err := os.WriteFile("test_slides.html", []byte(html), 0600)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DONE")
}
