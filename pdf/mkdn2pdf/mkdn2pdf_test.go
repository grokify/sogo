package mkdn2pdf

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

const sampleMarkdown = `# Heading

Some **bold** text and a paragraph.

- item one
- item two

| A | B |
|---|---|
| 1 | 2 |

` + "```" + `
code block
` + "```" + `
`

func TestMarkdownToPDF(t *testing.T) {
	out, err := MarkdownToPDF(sampleMarkdown, DefaultOptions())
	if err != nil {
		t.Fatalf("MarkdownToPDF: %v", err)
	}
	if !bytes.HasPrefix(out, []byte("%PDF-")) {
		t.Errorf("output is not a PDF: %q", out[:min(len(out), 8)])
	}
}

func TestMarkdownToPDFKeepsPartialOptions(t *testing.T) {
	opts := Options{Title: "Report", Author: "Someone", PageSize: "Letter", MarginLeft: 5}
	out, err := MarkdownToPDF(sampleMarkdown, opts)
	if err != nil {
		t.Fatalf("MarkdownToPDF: %v", err)
	}
	if !bytes.HasPrefix(out, []byte("%PDF-")) {
		t.Error("output is not a PDF")
	}
}

func TestWithDefaults(t *testing.T) {
	got := Options{Title: "T"}.withDefaults()
	if got.Title != "T" {
		t.Errorf("Title lost: %+v", got)
	}
	if got.PageSize != "A4" || got.FontSize != 12 || got.MarginTop != 15 {
		t.Errorf("defaults not applied: %+v", got)
	}

	got = Options{MarginLeft: 5}.withDefaults()
	if got.MarginLeft != 5 || got.MarginTop != 0 {
		t.Errorf("explicit margins not preserved: %+v", got)
	}
}

func TestMarkdownToPDFFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "in.md")
	dst := filepath.Join(dir, "out.pdf")
	if err := os.WriteFile(src, []byte(sampleMarkdown), 0600); err != nil {
		t.Fatal(err)
	}
	if err := MarkdownToPDFFile(src, dst, DefaultOptions()); err != nil {
		t.Fatalf("MarkdownToPDFFile: %v", err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		t.Error("output file is not a PDF")
	}
}

func TestTableToPDFAndWriter(t *testing.T) {
	out, err := TableToPDF([]string{"A", "B"}, [][]string{{"1", "2"}}, Options{})
	if err != nil || !bytes.HasPrefix(out, []byte("%PDF-")) {
		t.Fatalf("TableToPDF: %v", err)
	}

	w := NewWriter(Options{})
	w.AddTitle("Title")
	w.AddHeader("Section", 2)
	w.AddParagraph("text")
	if err := w.AddMarkdown(sampleMarkdown); err != nil {
		t.Fatalf("AddMarkdown: %v", err)
	}
	if out, err = w.Generate(); err != nil || !bytes.HasPrefix(out, []byte("%PDF-")) {
		t.Fatalf("Generate: %v", err)
	}
}
