package md2html

import (
	"bytes"
	"os"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

// Result contains the rendered HTML and extracted metadata.
type Result struct {
	HTML        string
	Frontmatter *Frontmatter
	TOC         *TOC
}

// Render converts markdown to HTML with the given options.
//
// Raw HTML in the markdown, Options.CSSInline and Options.CSSTokens are passed
// through unescaped, so only render trusted input.
func Render(md []byte, opts Options) (*Result, error) {
	result := &Result{}

	// Parse frontmatter
	fm, content, err := ParseFrontmatter(md)
	if err != nil {
		return nil, err
	}
	result.Frontmatter = fm

	// Generate TOC if requested
	if opts.GenerateTOC {
		maxDepth := opts.TOCMaxDepth
		if maxDepth <= 0 {
			maxDepth = 3
		}
		result.TOC = GenerateTOC(content, maxDepth)
	}

	// Build goldmark instance
	gmOpts := []goldmark.Option{
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithUnsafe(),
		),
	}

	// Add syntax highlighting
	if opts.Highlight {
		style := opts.HighlightStyle
		if style == "" {
			style = "github"
		}
		gmOpts = append(gmOpts, goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle(style),
				highlighting.WithFormatOptions(
					html.WithClasses(false),
					html.WithLineNumbers(false),
				),
			),
		))
	}

	mdParser := goldmark.New(gmOpts...)

	// Convert markdown to HTML
	var buf bytes.Buffer
	if err := mdParser.Convert(content, &buf); err != nil {
		return nil, err
	}

	htmlContent := buf.String()

	// If wrapper requested, build full document
	if opts.Wrapper {
		// Get theme tokens
		tokens := make(map[string]string)
		if opts.Theme != "" {
			if theme, ok := GetTheme(opts.Theme); ok {
				tokens = MergeTokens(theme, opts.CSSTokens)
			}
		}
		if len(opts.CSSTokens) > 0 && len(tokens) == 0 {
			tokens = opts.CSSTokens
		}

		doc := Document{
			Title:       opts.Title,
			Lang:        opts.Lang,
			Meta:        opts.Meta,
			CSSTokens:   RenderCSSTokens(tokens),
			CSSInline:   opts.CSSInline,
			CSSFile:     opts.CSSFile,
			Content:     htmlContent,
			Frontmatter: fm,
		}

		if result.TOC != nil {
			doc.TOC = result.TOC.RenderTOCHTML()
		}

		if doc.Lang == "" {
			doc.Lang = "en"
		}

		result.HTML = RenderDocument(doc)
	} else {
		result.HTML = htmlContent
	}

	return result, nil
}

// RenderFile reads a markdown file and converts it to HTML.
func RenderFile(srcPath string, opts Options) (*Result, error) {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}
	return Render(content, opts)
}

// RenderToFile converts markdown to HTML and writes to a file.
func RenderToFile(md []byte, outPath string, opts Options) error {
	result, err := Render(md, opts)
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(result.HTML), 0600)
}

// RenderFileToFile reads a markdown file and writes HTML to another file.
func RenderFileToFile(srcPath, outPath string, opts Options) error {
	result, err := RenderFile(srcPath, opts)
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(result.HTML), 0600)
}

// GetHighlightStyles returns available syntax highlighting style names.
func GetHighlightStyles() []string {
	return styles.Names()
}

// GetHighlightCSS returns the CSS for a syntax highlighting style.
func GetHighlightCSS(styleName string) (string, error) {
	style := styles.Get(styleName)
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.WithClasses(true))
	var buf bytes.Buffer
	if err := formatter.WriteCSS(&buf, style); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetLanguages returns available syntax highlighting languages.
func GetLanguages() []string {
	names := make([]string, 0)
	for _, l := range lexers.GlobalLexerRegistry.Lexers {
		config := l.Config()
		names = append(names, config.Name)
	}
	return names
}
