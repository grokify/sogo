// Package mkdn2pdf provides Markdown to PDF conversion using the maroto library.
// It supports tables, headers, lists, and styled text blocks.
package mkdn2pdf

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// Options configures PDF generation from Markdown.
type Options struct {
	// PageSize is the page size (e.g., "A4", "Letter"). Default: "A4".
	PageSize string
	// FontFamily is the font family to use. Default: "Arial".
	FontFamily string
	// FontSize is the base font size. Default: 12.
	FontSize float64
	// Margins for the page (top, left, right, bottom).
	MarginTop    float64
	MarginLeft   float64
	MarginRight  float64
	MarginBottom float64
	// Title is an optional document title.
	Title string
	// Author is an optional document author.
	Author string
}

// DefaultOptions returns sensible default options.
func DefaultOptions() Options {
	return Options{
		PageSize:     "A4",
		FontFamily:   "Arial",
		FontSize:     12,
		MarginTop:    15,
		MarginLeft:   15,
		MarginRight:  15,
		MarginBottom: 15,
	}
}

// withDefaults returns a copy of o with zero-valued fields set to defaults.
// Default margins apply only when all four margins are zero, so a caller can
// still ask for a single zero margin by setting another one.
func (o Options) withDefaults() Options {
	d := DefaultOptions()
	if o.PageSize == "" {
		o.PageSize = d.PageSize
	}
	if o.FontFamily == "" {
		o.FontFamily = d.FontFamily
	}
	if o.FontSize == 0 {
		o.FontSize = d.FontSize
	}
	if o.MarginTop == 0 && o.MarginLeft == 0 && o.MarginRight == 0 && o.MarginBottom == 0 {
		o.MarginTop, o.MarginLeft = d.MarginTop, d.MarginLeft
		o.MarginRight, o.MarginBottom = d.MarginRight, d.MarginBottom
	}
	return o
}

// pageSizes maps Options.PageSize values (case-insensitive) to maroto page sizes.
var pageSizes = map[string]pagesize.Type{
	"a3":     pagesize.A3,
	"a4":     pagesize.A4,
	"a5":     pagesize.A5,
	"letter": pagesize.Letter,
	"legal":  pagesize.Legal,
}

// newMaroto builds a maroto document from options that have had defaults applied.
func newMaroto(opts Options) core.Maroto {
	b := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(opts.MarginLeft).
		WithTopMargin(opts.MarginTop).
		WithRightMargin(opts.MarginRight).
		WithBottomMargin(opts.MarginBottom).
		WithDefaultFont(&props.Font{Family: opts.FontFamily, Size: opts.FontSize})
	if size, ok := pageSizes[strings.ToLower(opts.PageSize)]; ok {
		b = b.WithPageSize(size)
	}
	if opts.Title != "" {
		b = b.WithTitle(opts.Title, true)
	}
	if opts.Author != "" {
		b = b.WithAuthor(opts.Author, true)
	}
	return maroto.New(b.Build())
}

// MarkdownToPDF converts Markdown content to PDF bytes.
func MarkdownToPDF(markdown string, opts Options) ([]byte, error) {
	opts = opts.withDefaults()
	m := newMaroto(opts)

	if opts.Title != "" {
		m.AddRow(10, text.NewCol(12, opts.Title, props.Text{
			Size:  16,
			Style: fontstyle.Bold,
			Align: align.Center,
		}))
		m.AddRow(5) // spacing
	}

	// Parse and render markdown content
	if err := renderMarkdown(m, markdown, opts); err != nil {
		return nil, fmt.Errorf("rendering markdown: %w", err)
	}

	doc, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("generating PDF: %w", err)
	}

	return doc.GetBytes(), nil
}

// MarkdownToPDFFile converts a Markdown file to a PDF file.
func MarkdownToPDFFile(srcPath, outPath string, opts Options) error {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("reading source file: %w", err)
	}

	pdfBytes, err := MarkdownToPDF(string(content), opts)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outPath, pdfBytes, 0600); err != nil {
		return fmt.Errorf("writing PDF file: %w", err)
	}

	return nil
}

// MarkdownBytesToPDF converts Markdown bytes to PDF bytes.
func MarkdownBytesToPDF(markdown []byte, opts Options) ([]byte, error) {
	return MarkdownToPDF(string(markdown), opts)
}

// renderMarkdown parses markdown and adds content to the maroto document.
func renderMarkdown(m core.Maroto, markdown string, opts Options) error {
	lines := strings.Split(markdown, "\n")
	var inCodeBlock bool
	var codeBlockContent strings.Builder
	var inTable bool
	var tableRows [][]string
	var tableHeader []string

	fontSize := opts.FontSize
	if fontSize == 0 {
		fontSize = 12
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Handle code blocks
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				// End code block
				addCodeBlock(m, codeBlockContent.String(), fontSize)
				codeBlockContent.Reset()
				inCodeBlock = false
			} else {
				// Start code block
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			codeBlockContent.WriteString(line + "\n")
			continue
		}

		// Handle tables
		if strings.HasPrefix(strings.TrimSpace(line), "|") {
			if !inTable {
				inTable = true
				tableHeader = parseTableRow(line)
			} else if isTableSeparator(line) {
				// Skip separator line
				continue
			} else {
				tableRows = append(tableRows, parseTableRow(line))
			}
			continue
		} else if inTable {
			// End of table
			addTable(m, tableHeader, tableRows, fontSize)
			tableHeader = nil
			tableRows = nil
			inTable = false
		}

		// Handle headers
		if strings.HasPrefix(line, "######") {
			addHeader(m, strings.TrimPrefix(line, "######"), 6, fontSize)
		} else if strings.HasPrefix(line, "#####") {
			addHeader(m, strings.TrimPrefix(line, "#####"), 5, fontSize)
		} else if strings.HasPrefix(line, "####") {
			addHeader(m, strings.TrimPrefix(line, "####"), 4, fontSize)
		} else if strings.HasPrefix(line, "###") {
			addHeader(m, strings.TrimPrefix(line, "###"), 3, fontSize)
		} else if strings.HasPrefix(line, "##") {
			addHeader(m, strings.TrimPrefix(line, "##"), 2, fontSize)
		} else if strings.HasPrefix(line, "#") {
			addHeader(m, strings.TrimPrefix(line, "#"), 1, fontSize)
		} else if strings.HasPrefix(strings.TrimSpace(line), "- ") || strings.HasPrefix(strings.TrimSpace(line), "* ") {
			// Unordered list item
			addListItem(m, strings.TrimSpace(line)[2:], false, 0, fontSize)
		} else if match := regexp.MustCompile(`^\s*(\d+)\.\s+(.*)$`).FindStringSubmatch(line); match != nil {
			// Ordered list item
			addListItem(m, match[2], true, 0, fontSize)
		} else if strings.TrimSpace(line) == "" {
			// Empty line - add spacing
			m.AddRow(3)
		} else {
			// Regular paragraph
			addParagraph(m, line, fontSize)
		}
	}

	// Handle any remaining table
	if inTable && len(tableHeader) > 0 {
		addTable(m, tableHeader, tableRows, fontSize)
	}

	return nil
}

func addHeader(m core.Maroto, content string, level int, baseSize float64) {
	content = strings.TrimSpace(content)
	sizes := map[int]float64{
		1: baseSize + 8,
		2: baseSize + 6,
		3: baseSize + 4,
		4: baseSize + 2,
		5: baseSize + 1,
		6: baseSize,
	}
	size := sizes[level]
	if size == 0 {
		size = baseSize
	}

	m.AddRow(size+2, text.NewCol(12, content, props.Text{
		Size:  size,
		Style: fontstyle.Bold,
		Top:   2,
	}))
	m.AddRow(2) // spacing after header
}

func addParagraph(m core.Maroto, content string, fontSize float64) {
	content = strings.TrimSpace(content)
	if content == "" {
		return
	}

	// Handle inline formatting (bold, italic)
	content = processInlineFormatting(content)

	m.AddRow(6, text.NewCol(12, content, props.Text{
		Size: fontSize,
	}))
}

func addListItem(m core.Maroto, content string, ordered bool, index int, fontSize float64) {
	prefix := "  \u2022 " // bullet
	if ordered {
		prefix = fmt.Sprintf("  %d. ", index+1)
	}

	m.AddRow(6, text.NewCol(12, prefix+content, props.Text{
		Size: fontSize,
	}))
}

func addCodeBlock(m core.Maroto, content string, fontSize float64) {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	for _, line := range lines {
		m.AddRow(5, text.NewCol(12, "  "+line, props.Text{
			Size:  fontSize - 1,
			Style: fontstyle.Italic,
		}))
	}
	m.AddRow(2) // spacing after code block
}

func addTable(m core.Maroto, header []string, rows [][]string, fontSize float64) {
	if len(header) == 0 {
		return
	}

	colWidth := 12 / len(header)
	if colWidth < 1 {
		colWidth = 1
	}

	// Add header row
	var headerCols []core.Col
	for _, h := range header {
		headerCols = append(headerCols, text.NewCol(colWidth, h, props.Text{
			Size:  fontSize,
			Style: fontstyle.Bold,
			Align: align.Center,
		}))
	}
	m.AddRow(8, headerCols...)

	// Add data rows
	for _, row := range rows {
		var cols []core.Col
		for i, cell := range row {
			if i >= len(header) {
				break
			}
			cols = append(cols, text.NewCol(colWidth, cell, props.Text{
				Size:  fontSize,
				Align: align.Center,
			}))
		}
		// Pad with empty cols if needed
		for len(cols) < len(header) {
			cols = append(cols, text.NewCol(colWidth, "", props.Text{Size: fontSize}))
		}
		m.AddRow(7, cols...)
	}

	m.AddRow(3) // spacing after table
}

func parseTableRow(line string) []string {
	line = strings.Trim(line, "|")
	parts := strings.Split(line, "|")
	var cells []string
	for _, p := range parts {
		cells = append(cells, strings.TrimSpace(p))
	}
	return cells
}

func isTableSeparator(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "|") {
		return false
	}
	// Check if it's a separator line like |---|---|
	return regexp.MustCompile(`^\|[\s\-:|]+\|$`).MatchString(line)
}

func processInlineFormatting(content string) string {
	// Remove markdown bold/italic markers for plain text output
	// In a more advanced implementation, we could track formatting
	content = regexp.MustCompile(`\*\*([^*]+)\*\*`).ReplaceAllString(content, "$1")
	content = regexp.MustCompile(`\*([^*]+)\*`).ReplaceAllString(content, "$1")
	content = regexp.MustCompile(`__([^_]+)__`).ReplaceAllString(content, "$1")
	content = regexp.MustCompile(`_([^_]+)_`).ReplaceAllString(content, "$1")
	content = regexp.MustCompile("`([^`]+)`").ReplaceAllString(content, "$1")
	return content
}

// TableToPDF renders a simple table to PDF.
func TableToPDF(header []string, rows [][]string, opts Options) ([]byte, error) {
	opts = opts.withDefaults()
	m := newMaroto(opts)

	if opts.Title != "" {
		m.AddRow(10, text.NewCol(12, opts.Title, props.Text{
			Size:  16,
			Style: fontstyle.Bold,
			Align: align.Center,
		}))
		m.AddRow(5)
	}

	fontSize := opts.FontSize
	if fontSize == 0 {
		fontSize = 12
	}

	addTable(m, header, rows, fontSize)

	doc, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("generating PDF: %w", err)
	}

	return doc.GetBytes(), nil
}

// Writer wraps a maroto instance for building PDFs programmatically.
type Writer struct {
	m        core.Maroto
	opts     Options
	fontSize float64
}

// NewWriter creates a new PDF writer with the given options.
func NewWriter(opts Options) *Writer {
	opts = opts.withDefaults()

	return &Writer{
		m:        newMaroto(opts),
		opts:     opts,
		fontSize: opts.FontSize,
	}
}

// AddTitle adds a centered title.
func (w *Writer) AddTitle(title string) {
	w.m.AddRow(10, text.NewCol(12, title, props.Text{
		Size:  w.fontSize + 6,
		Style: fontstyle.Bold,
		Align: align.Center,
	}))
	w.m.AddRow(5)
}

// AddHeader adds a section header.
func (w *Writer) AddHeader(content string, level int) {
	addHeader(w.m, content, level, w.fontSize)
}

// AddParagraph adds a paragraph of text.
func (w *Writer) AddParagraph(content string) {
	addParagraph(w.m, content, w.fontSize)
}

// AddTable adds a table with header and rows.
func (w *Writer) AddTable(header []string, rows [][]string) {
	addTable(w.m, header, rows, w.fontSize)
}

// AddSpace adds vertical spacing.
func (w *Writer) AddSpace(height float64) {
	w.m.AddRow(height)
}

// AddMarkdown parses and adds markdown content.
func (w *Writer) AddMarkdown(markdown string) error {
	return renderMarkdown(w.m, markdown, w.opts)
}

// Generate produces the final PDF bytes.
func (w *Writer) Generate() ([]byte, error) {
	doc, err := w.m.Generate()
	if err != nil {
		return nil, fmt.Errorf("generating PDF: %w", err)
	}
	return doc.GetBytes(), nil
}

// WriteFile generates and writes the PDF to a file.
func (w *Writer) WriteFile(path string) error {
	data, err := w.Generate()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
