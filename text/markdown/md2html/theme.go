package md2html

import (
	"bytes"
	"fmt"
	"sort"
)

// Theme represents a named set of CSS tokens.
type Theme struct {
	Name        string
	Description string
	Tokens      map[string]string
}

// Built-in themes
var builtinThemes = map[string]Theme{
	"default": {
		Name:        "default",
		Description: "Clean, readable light theme",
		Tokens: map[string]string{
			"color-bg":           "#ffffff",
			"color-bg-secondary": "#f8f9fa",
			"color-text":         "#212529",
			"color-text-muted":   "#6c757d",
			"color-primary":      "#0d6efd",
			"color-primary-dark": "#0a58ca",
			"color-accent":       "#ffc107",
			"color-border":       "#dee2e6",
			"color-code-bg":      "#f8f9fa",
			"font-body":          "system-ui, -apple-system, sans-serif",
			"font-heading":       "system-ui, -apple-system, sans-serif",
			"font-mono":          "ui-monospace, monospace",
			"line-height":        "1.6",
			"max-width":          "65ch",
		},
	},
	"dark": {
		Name:        "dark",
		Description: "Dark theme for reduced eye strain",
		Tokens: map[string]string{
			"color-bg":           "#1a1a2e",
			"color-bg-secondary": "#16213e",
			"color-text":         "#eaeaea",
			"color-text-muted":   "#a0a0a0",
			"color-primary":      "#4cc9f0",
			"color-primary-dark": "#3a0ca3",
			"color-accent":       "#f72585",
			"color-border":       "#2d2d44",
			"color-code-bg":      "#16213e",
			"font-body":          "system-ui, -apple-system, sans-serif",
			"font-heading":       "system-ui, -apple-system, sans-serif",
			"font-mono":          "ui-monospace, monospace",
			"line-height":        "1.6",
			"max-width":          "65ch",
		},
	},
	"minimal": {
		Name:        "minimal",
		Description: "Minimal, typography-focused theme",
		Tokens: map[string]string{
			"color-bg":           "#ffffff",
			"color-bg-secondary": "#fafafa",
			"color-text":         "#333333",
			"color-text-muted":   "#666666",
			"color-primary":      "#333333",
			"color-primary-dark": "#000000",
			"color-accent":       "#333333",
			"color-border":       "#eeeeee",
			"color-code-bg":      "#f5f5f5",
			"font-body":          "Georgia, serif",
			"font-heading":       "Georgia, serif",
			"font-mono":          "Menlo, monospace",
			"line-height":        "1.7",
			"max-width":          "60ch",
		},
	},
	"forest": {
		Name:        "forest",
		Description: "Professional green and gold theme",
		Tokens: map[string]string{
			"color-bg":           "#ffffff",
			"color-bg-secondary": "#f8fdf9",
			"color-text":         "#111827",
			"color-text-muted":   "#4b5563",
			"color-primary":      "#047857",
			"color-primary-dark": "#064e3b",
			"color-accent":       "#eab308",
			"color-border":       "#d1d5db",
			"color-code-bg":      "#f0fdf4",
			"font-body":          "system-ui, -apple-system, sans-serif",
			"font-heading":       "system-ui, -apple-system, sans-serif",
			"font-mono":          "ui-monospace, monospace",
			"line-height":        "1.6",
			"max-width":          "65ch",
		},
	},
}

// GetTheme returns a built-in theme by name.
func GetTheme(name string) (Theme, bool) {
	theme, ok := builtinThemes[name]
	return theme, ok
}

// ListThemes returns the names of all built-in themes.
func ListThemes() []string {
	names := make([]string, 0, len(builtinThemes))
	for name := range builtinThemes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// MergeTokens merges custom tokens with theme tokens.
// Custom tokens override theme tokens.
func MergeTokens(theme Theme, custom map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range theme.Tokens {
		merged[k] = v
	}
	for k, v := range custom {
		merged[k] = v
	}
	return merged
}

// RenderCSSTokens generates CSS custom properties from tokens.
func RenderCSSTokens(tokens map[string]string) string {
	if len(tokens) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString(":root {\n")

	// Sort keys for deterministic output
	keys := make([]string, 0, len(tokens))
	for k := range tokens {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		buf.WriteString(fmt.Sprintf("  --%s: %s;\n", k, tokens[k]))
	}

	buf.WriteString("}\n")
	return buf.String()
}
