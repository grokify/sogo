package md2html

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func MarkdownToHTMLGoldmark(md []byte) ([]byte, error) {
	var buf bytes.Buffer
	mdParser := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	if err := mdParser.Convert(md, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
