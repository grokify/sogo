package md2html

import (
	"strings"
	"testing"
)

func TestRenderBasic(t *testing.T) {
	md := []byte(`# Hello World

This is a **test** document.

## Section One

Some content here.

` + "```go\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n```" + `

### Subsection

More content.
`)

	opts := DefaultOptions()
	result, err := Render(md, opts)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if result.HTML == "" {
		t.Error("Expected non-empty HTML output")
	}

	if !strings.Contains(result.HTML, "<h1") {
		t.Error("Expected h1 tag in output")
	}

	if !strings.Contains(result.HTML, "<strong>test</strong>") {
		t.Error("Expected strong tag in output")
	}
}

func TestRenderWithFrontmatter(t *testing.T) {
	md := []byte(`---
title: Test Document
subtitle: A subtitle here
author: Test Author
description: This is a test
version: v1.0.0
tags:
  - test
  - example
---

# Content

This is the content.
`)

	opts := DefaultOptions()
	result, err := Render(md, opts)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if result.Frontmatter == nil {
		t.Fatal("Expected frontmatter to be parsed")
	}

	if result.Frontmatter.Title != "Test Document" {
		t.Errorf("Expected title 'Test Document', got '%s'", result.Frontmatter.Title)
	}

	if result.Frontmatter.Subtitle != "A subtitle here" {
		t.Errorf("Expected subtitle 'A subtitle here', got '%s'", result.Frontmatter.Subtitle)
	}

	if result.Frontmatter.Author != "Test Author" {
		t.Errorf("Expected author 'Test Author', got '%s'", result.Frontmatter.Author)
	}

	if len(result.Frontmatter.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(result.Frontmatter.Tags))
	}
}

func TestRenderWithTOC(t *testing.T) {
	md := []byte(`# Title

## Section One

### Subsection A

## Section Two

### Subsection B

### Subsection C
`)

	opts := Options{
		GenerateTOC: true,
		TOCMaxDepth: 3,
	}
	result, err := Render(md, opts)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if result.TOC == nil {
		t.Fatal("Expected TOC to be generated")
	}

	// Should have 6 entries: 1 h1, 2 h2, 3 h3
	if len(result.TOC.Entries) != 6 {
		t.Errorf("Expected 6 TOC entries, got %d", len(result.TOC.Entries))
	}
}

func TestRenderWithWrapper(t *testing.T) {
	md := []byte(`# Hello

Some content.
`)

	opts := Options{
		Wrapper: true,
		Title:   "My Document",
		Lang:    "en",
		Theme:   "default",
	}
	result, err := Render(md, opts)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.HasPrefix(result.HTML, "<!DOCTYPE html>") {
		t.Error("Expected DOCTYPE declaration")
	}

	if !strings.Contains(result.HTML, `<html lang="en">`) {
		t.Error("Expected html tag with lang attribute")
	}

	if !strings.Contains(result.HTML, "<title>My Document</title>") {
		t.Error("Expected title tag")
	}

	if !strings.Contains(result.HTML, "<article>") {
		t.Error("Expected article tag")
	}

	if !strings.Contains(result.HTML, ":root {") {
		t.Error("Expected CSS custom properties")
	}
}

func TestThemes(t *testing.T) {
	themes := ListThemes()
	if len(themes) == 0 {
		t.Error("Expected at least one theme")
	}

	for _, name := range themes {
		theme, ok := GetTheme(name)
		if !ok {
			t.Errorf("Theme '%s' not found", name)
			continue
		}

		if theme.Name != name {
			t.Errorf("Theme name mismatch: expected '%s', got '%s'", name, theme.Name)
		}

		if len(theme.Tokens) == 0 {
			t.Errorf("Theme '%s' has no tokens", name)
		}

		if _, ok := theme.Tokens["color-bg"]; !ok {
			t.Errorf("Theme '%s' missing color-bg token", name)
		}
	}
}

func TestMergeTokens(t *testing.T) {
	theme := Theme{
		Tokens: map[string]string{
			"color-bg":   "#ffffff",
			"color-text": "#000000",
		},
	}

	custom := map[string]string{
		"color-bg":      "#f0f0f0",
		"color-primary": "#0000ff",
	}

	merged := MergeTokens(theme, custom)

	if merged["color-bg"] != "#f0f0f0" {
		t.Errorf("Custom token should override theme: expected '#f0f0f0', got '%s'", merged["color-bg"])
	}

	if merged["color-text"] != "#000000" {
		t.Errorf("Theme token should be preserved: expected '#000000', got '%s'", merged["color-text"])
	}

	if merged["color-primary"] != "#0000ff" {
		t.Errorf("Custom token should be added: expected '#0000ff', got '%s'", merged["color-primary"])
	}
}

func TestRenderCSSTokens(t *testing.T) {
	tokens := map[string]string{
		"color-bg":   "#ffffff",
		"color-text": "#000000",
	}

	css := RenderCSSTokens(tokens)

	if !strings.Contains(css, ":root {") {
		t.Error("Expected :root selector")
	}

	if !strings.Contains(css, "--color-bg: #ffffff;") {
		t.Error("Expected color-bg custom property")
	}

	if !strings.Contains(css, "--color-text: #000000;") {
		t.Error("Expected color-text custom property")
	}
}

func TestGenerateTOCIDsMatchRenderedHTML(t *testing.T) {
	md := []byte("# Hello, World!\n\n## Repeat\n\n## Repeat\n\n## `code` **bold**\n\n```\n# not a heading\n```\n")

	toc := GenerateTOC(md, 3)
	if len(toc.Entries) != 4 {
		t.Fatalf("expected 4 entries, got %d: %+v", len(toc.Entries), toc.Entries)
	}
	if toc.Entries[3].Text != "code bold" {
		t.Errorf("expected plain heading text, got %q", toc.Entries[3].Text)
	}

	result, err := Render(md, Options{})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	for _, e := range toc.Entries {
		if !strings.Contains(result.HTML, `id="`+e.ID+`"`) {
			t.Errorf("TOC id %q not found in rendered HTML:\n%s", e.ID, result.HTML)
		}
	}
}

func TestRenderTOCHTMLEscapes(t *testing.T) {
	toc := GenerateTOC([]byte("# <b>x</b> & y\n"), 3)
	out := toc.RenderTOCHTML()
	if strings.Contains(out, "<b>") {
		t.Errorf("TOC text not escaped: %s", out)
	}
}

func TestParseFrontmatterDelimiterInValue(t *testing.T) {
	fm, body, err := ParseFrontmatter([]byte("---\ntitle: a---b\n---\n\n# Body\n"))
	if err != nil {
		t.Fatalf("ParseFrontmatter: %v", err)
	}
	if fm == nil || fm.Title != "a---b" {
		t.Errorf("unexpected frontmatter: %+v", fm)
	}
	if string(body) != "# Body" {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestRenderDocumentEscapesLangAndVersion(t *testing.T) {
	out := RenderDocument(Document{
		Lang:        `en"><script>`,
		Frontmatter: &Frontmatter{Title: "T", Version: "<i>1</i>"},
	})
	if strings.Contains(out, "<script>") || strings.Contains(out, "<i>1") {
		t.Errorf("unescaped output:\n%s", out)
	}
}

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantFM      bool
		wantTitle   string
		wantContent string
	}{
		{
			name:        "no frontmatter",
			input:       "# Hello\n\nContent here.",
			wantFM:      false,
			wantContent: "# Hello\n\nContent here.",
		},
		{
			name:        "with frontmatter",
			input:       "---\ntitle: Test\n---\n\n# Hello",
			wantFM:      true,
			wantTitle:   "Test",
			wantContent: "# Hello",
		},
		{
			name:        "unclosed frontmatter",
			input:       "---\ntitle: Test\n\n# Hello",
			wantFM:      false,
			wantContent: "---\ntitle: Test\n\n# Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, content, err := ParseFrontmatter([]byte(tt.input))
			if err != nil {
				t.Fatalf("ParseFrontmatter failed: %v", err)
			}

			if tt.wantFM && fm == nil {
				t.Error("Expected frontmatter to be parsed")
			}

			if !tt.wantFM && fm != nil {
				t.Error("Expected no frontmatter")
			}

			if fm != nil && fm.Title != tt.wantTitle {
				t.Errorf("Expected title '%s', got '%s'", tt.wantTitle, fm.Title)
			}

			if strings.TrimSpace(string(content)) != strings.TrimSpace(tt.wantContent) {
				t.Errorf("Expected content '%s', got '%s'", tt.wantContent, string(content))
			}
		})
	}
}

func TestGetHighlightStyles(t *testing.T) {
	styles := GetHighlightStyles()
	if len(styles) == 0 {
		t.Error("Expected at least one highlight style")
	}

	// Check for some common styles
	found := make(map[string]bool)
	for _, s := range styles {
		found[s] = true
	}

	commonStyles := []string{"monokai", "github", "dracula"}
	for _, style := range commonStyles {
		if !found[style] {
			t.Errorf("Expected '%s' style to be available", style)
		}
	}
}

func TestGetHighlightCSS(t *testing.T) {
	css, err := GetHighlightCSS("github")
	if err != nil {
		t.Fatalf("GetHighlightCSS failed: %v", err)
	}

	if css == "" {
		t.Error("Expected non-empty CSS")
	}

	if !strings.Contains(css, ".chroma") {
		t.Error("Expected .chroma class in CSS")
	}
}
