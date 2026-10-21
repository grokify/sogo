package md2html

import (
	"bytes"
	"fmt"
	"html"
	"sort"
)

// Document represents a full HTML document.
type Document struct {
	Title       string
	Lang        string
	Meta        map[string]string
	CSSTokens   string
	CSSInline   string
	CSSFile     string
	TOC         string
	Content     string
	Frontmatter *Frontmatter
}

// DefaultArticleCSS returns basic article styling that uses CSS tokens.
const DefaultArticleCSS = `
body {
  font-family: var(--font-body);
  line-height: var(--line-height);
  color: var(--color-text);
  background-color: var(--color-bg);
  margin: 0;
  padding: 2rem;
}

article {
  max-width: var(--max-width);
  margin: 0 auto;
}

h1, h2, h3, h4, h5, h6 {
  font-family: var(--font-heading);
  color: var(--color-text);
  margin-top: 2em;
  margin-bottom: 0.5em;
  line-height: 1.3;
}

h1 { font-size: 2.25rem; margin-top: 0; }
h2 { font-size: 1.75rem; }
h3 { font-size: 1.375rem; }
h4 { font-size: 1.125rem; }

a {
  color: var(--color-primary);
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

p {
  margin: 1em 0;
}

code {
  font-family: var(--font-mono);
  font-size: 0.9em;
  background-color: var(--color-code-bg);
  padding: 0.2em 0.4em;
  border-radius: 3px;
}

pre {
  background-color: var(--color-code-bg);
  padding: 1rem;
  border-radius: 6px;
  overflow-x: auto;
  border: 1px solid var(--color-border);
}

pre code {
  background: none;
  padding: 0;
}

blockquote {
  border-left: 4px solid var(--color-primary);
  margin: 1.5em 0;
  padding: 0.5em 1em;
  background-color: var(--color-bg-secondary);
  color: var(--color-text-muted);
}

table {
  width: 100%;
  border-collapse: collapse;
  margin: 1.5em 0;
}

th, td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

th {
  font-weight: 600;
  background-color: var(--color-bg-secondary);
}

hr {
  border: none;
  border-top: 1px solid var(--color-border);
  margin: 2em 0;
}

img {
  max-width: 100%;
  height: auto;
}

.toc {
  background-color: var(--color-bg-secondary);
  padding: 1rem 1.5rem;
  border-radius: 6px;
  margin-bottom: 2rem;
  border: 1px solid var(--color-border);
}

.toc ul {
  list-style: none;
  padding-left: 1rem;
  margin: 0.5rem 0;
}

.toc > ul {
  padding-left: 0;
}

.toc a {
  color: var(--color-text-muted);
}

.toc a:hover {
  color: var(--color-primary);
}

.article-meta {
  color: var(--color-text-muted);
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--color-border);
}

.article-meta h1 {
  color: var(--color-text);
  margin-bottom: 0.25rem;
}

.article-meta .subtitle {
  font-size: 1.25rem;
  color: var(--color-primary);
  margin-bottom: 1rem;
}

.article-meta .meta-info {
  font-size: 0.875rem;
}

@media print {
  body {
    padding: 0;
  }

  .toc {
    page-break-after: always;
  }
}
`

// RenderDocument generates a full HTML document.
func RenderDocument(doc Document) string {
	var buf bytes.Buffer

	// Doctype and html
	buf.WriteString("<!DOCTYPE html>\n")
	buf.WriteString(fmt.Sprintf(`<html lang="%s">`, html.EscapeString(doc.Lang)))
	buf.WriteString("\n<head>\n")

	// Meta charset
	buf.WriteString(`  <meta charset="UTF-8">`)
	buf.WriteString("\n")
	buf.WriteString(`  <meta name="viewport" content="width=device-width, initial-scale=1.0">`)
	buf.WriteString("\n")

	// Title
	title := doc.Title
	if title == "" && doc.Frontmatter != nil && doc.Frontmatter.Title != "" {
		title = doc.Frontmatter.Title
	}
	if title != "" {
		buf.WriteString(fmt.Sprintf("  <title>%s</title>\n", html.EscapeString(title)))
	}

	// Meta tags from frontmatter
	if doc.Frontmatter != nil {
		if doc.Frontmatter.Description != "" {
			buf.WriteString(fmt.Sprintf(`  <meta name="description" content="%s">`, html.EscapeString(doc.Frontmatter.Description)))
			buf.WriteString("\n")
		}
		if doc.Frontmatter.Author != "" {
			buf.WriteString(fmt.Sprintf(`  <meta name="author" content="%s">`, html.EscapeString(doc.Frontmatter.Author)))
			buf.WriteString("\n")
		}
	}

	// Additional meta tags
	if len(doc.Meta) > 0 {
		keys := make([]string, 0, len(doc.Meta))
		for k := range doc.Meta {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			buf.WriteString(fmt.Sprintf(`  <meta name="%s" content="%s">`, html.EscapeString(k), html.EscapeString(doc.Meta[k])))
			buf.WriteString("\n")
		}
	}

	// External CSS file
	if doc.CSSFile != "" {
		buf.WriteString(fmt.Sprintf(`  <link rel="stylesheet" href="%s">`, html.EscapeString(doc.CSSFile)))
		buf.WriteString("\n")
	}

	// CSS tokens and inline styles
	buf.WriteString("  <style>\n")
	if doc.CSSTokens != "" {
		buf.WriteString(doc.CSSTokens)
	}
	buf.WriteString(DefaultArticleCSS)
	if doc.CSSInline != "" {
		buf.WriteString(doc.CSSInline)
	}
	buf.WriteString("  </style>\n")

	buf.WriteString("</head>\n<body>\n")
	buf.WriteString("<article>\n")

	// Article header from frontmatter
	if doc.Frontmatter != nil && doc.Frontmatter.Title != "" {
		buf.WriteString(`<header class="article-meta">`)
		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf("<h1>%s</h1>\n", html.EscapeString(doc.Frontmatter.Title)))
		if doc.Frontmatter.Subtitle != "" {
			buf.WriteString(fmt.Sprintf(`<p class="subtitle">%s</p>`, html.EscapeString(doc.Frontmatter.Subtitle)))
			buf.WriteString("\n")
		}
		buf.WriteString(`<p class="meta-info">`)
		if doc.Frontmatter.Version != "" {
			buf.WriteString(html.EscapeString(doc.Frontmatter.Version))
		}
		if !doc.Frontmatter.Date.IsZero() {
			if doc.Frontmatter.Version != "" {
				buf.WriteString(" · ")
			}
			buf.WriteString(doc.Frontmatter.Date.Format("January 2006"))
		}
		buf.WriteString("</p>\n")
		buf.WriteString("</header>\n")
	}

	// Table of contents
	if doc.TOC != "" {
		buf.WriteString(doc.TOC)
		buf.WriteString("\n")
	}

	// Main content
	buf.WriteString(doc.Content)

	buf.WriteString("\n</article>\n")
	buf.WriteString("</body>\n</html>\n")

	return buf.String()
}
