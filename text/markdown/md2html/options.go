package md2html

// Options configures the markdown to HTML conversion.
type Options struct {
	// Theme is a named theme preset (e.g., "default", "dark", "minimal").
	// If set, loads predefined CSS tokens for that theme.
	Theme string

	// CSSTokens are custom CSS properties to inject into the document.
	// Keys should not include the "--" prefix (e.g., "color-primary" not "--color-primary").
	CSSTokens map[string]string

	// CSSFile is a path to an external CSS file to link.
	CSSFile string

	// CSSInline is raw CSS to embed in a <style> tag.
	CSSInline string

	// GenerateTOC creates a table of contents from headings.
	GenerateTOC bool

	// TOCMaxDepth limits TOC to headings up to this level (default: 3 = h1-h3).
	TOCMaxDepth int

	// Highlight enables syntax highlighting for code blocks.
	Highlight bool

	// HighlightStyle is the chroma style name (e.g., "monokai", "github", "dracula").
	HighlightStyle string

	// Wrapper wraps output in a full HTML document with <html>, <head>, <body>.
	Wrapper bool

	// Title sets the document <title> (used with Wrapper).
	Title string

	// Lang sets the document language attribute (default: "en").
	Lang string

	// Meta contains additional <meta> tags as key-value pairs.
	Meta map[string]string
}

// DefaultOptions returns sensible defaults for article rendering.
func DefaultOptions() Options {
	return Options{
		Theme:          "default",
		CSSTokens:      make(map[string]string),
		GenerateTOC:    false,
		TOCMaxDepth:    3,
		Highlight:      true,
		HighlightStyle: "github",
		Wrapper:        false,
		Lang:           "en",
		Meta:           make(map[string]string),
	}
}

// ArticleOptions returns options optimized for long-form article rendering.
func ArticleOptions() Options {
	opts := DefaultOptions()
	opts.GenerateTOC = true
	opts.Wrapper = true
	opts.Highlight = true
	return opts
}

// FragmentOptions returns options for rendering HTML fragments (no wrapper).
func FragmentOptions() Options {
	opts := DefaultOptions()
	opts.Wrapper = false
	opts.GenerateTOC = false
	return opts
}
