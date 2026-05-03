package md2html

import (
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/grokify/mogo/os/osutil"
)

func NewMarkdownParserDefault() *parser.Parser {
	exts := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	return parser.NewWithExtensions(exts)
}

func NewHTMLRendererDefault() *html.Renderer {
	flags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: flags}
	return html.NewRenderer(opts)
}

func MarkdownToHTML(md []byte) []byte {
	p := NewMarkdownParserDefault()
	doc := p.Parse(md)
	r := NewHTMLRendererDefault()
	return markdown.Render(doc, r)
}

func MarkdownToHTMLFile(srcFilename, outFilename string, perm os.FileMode) error {
	bSrc, err := osutil.ReadFileSecure(srcFilename)
	if err != nil {
		return err
	}
	bOut := MarkdownToHTML(bSrc)
	return osutil.WriteFileSecure(outFilename, bOut, perm)
}
