package md2html

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// TOCEntry represents a single entry in the table of contents.
type TOCEntry struct {
	Level    int
	Text     string
	ID       string
	Children []*TOCEntry
}

// TOC represents a table of contents.
type TOC struct {
	Entries []*TOCEntry
}

// GenerateTOC extracts headings from markdown and builds a TOC.
// Heading IDs come from goldmark's auto-ID generation, so TOC anchors match
// the IDs in the rendered HTML (including non-ASCII and duplicate headings).
func GenerateTOC(md []byte, maxDepth int) *TOC {
	if maxDepth <= 0 {
		maxDepth = 6
	}

	toc := &TOC{Entries: make([]*TOCEntry, 0)}

	p := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID())).Parser()
	doc := p.Parse(text.NewReader(md))

	// Walk only returns an error if the walker does; ours never does.
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		h, ok := n.(*ast.Heading)
		if !ok || !entering {
			return ast.WalkContinue, nil
		}
		if h.Level <= maxDepth {
			id, _ := h.AttributeString("id")
			idBytes, _ := id.([]byte)
			toc.Entries = append(toc.Entries, &TOCEntry{
				Level:    h.Level,
				Text:     headingText(h, md),
				ID:       string(idBytes),
				Children: make([]*TOCEntry, 0),
			})
		}
		return ast.WalkSkipChildren, nil
	})

	return toc
}

// headingText returns the plain text of a heading, without inline markup.
func headingText(h *ast.Heading, source []byte) string {
	var buf bytes.Buffer
	_ = ast.Walk(h, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch t := n.(type) {
		case *ast.Text:
			buf.Write(t.Segment.Value(source))
		case *ast.String:
			buf.Write(t.Value)
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(buf.String())
}

// RenderTOCHTML renders the TOC as an HTML unordered list.
func (toc *TOC) RenderTOCHTML() string {
	if len(toc.Entries) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString(`<nav class="toc" aria-label="Table of Contents">`)
	buf.WriteString("\n<ul>\n")

	prevLevel := 0
	for _, entry := range toc.Entries {
		// Handle nesting
		if entry.Level > prevLevel {
			for i := prevLevel; i < entry.Level; i++ {
				if i > 0 {
					buf.WriteString("<ul>\n")
				}
			}
		} else if entry.Level < prevLevel {
			for i := entry.Level; i < prevLevel; i++ {
				buf.WriteString("</li>\n</ul>\n")
			}
			buf.WriteString("</li>\n")
		} else if prevLevel > 0 {
			buf.WriteString("</li>\n")
		}

		buf.WriteString(fmt.Sprintf(`<li><a href="#%s">%s</a>`, html.EscapeString(entry.ID), html.EscapeString(entry.Text)))
		prevLevel = entry.Level
	}

	// Close remaining tags
	for i := 0; i < prevLevel; i++ {
		buf.WriteString("</li>\n")
		if i < prevLevel-1 {
			buf.WriteString("</ul>\n")
		}
	}

	buf.WriteString("</ul>\n</nav>")
	return buf.String()
}

// RenderTOCMarkdown renders the TOC as markdown.
func (toc *TOC) RenderTOCMarkdown() string {
	var buf bytes.Buffer

	for _, entry := range toc.Entries {
		indent := strings.Repeat("  ", entry.Level-1)
		buf.WriteString(fmt.Sprintf("%s- [%s](#%s)\n", indent, entry.Text, entry.ID))
	}

	return buf.String()
}
