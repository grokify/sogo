package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/grokify/sogo/text/markdown/md2html"
	"github.com/spf13/cobra"
)

var (
	theme          string
	cssFile        string
	cssInline      string
	generateTOC    bool
	tocMaxDepth    int
	highlight      bool
	highlightStyle string
	wrapper        bool
	title          string
	lang           string
	listThemes     bool
	listStyles     bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "md2html [input.md] [output.html]",
		Short: "Convert markdown to HTML with theming and syntax highlighting",
		Long: `A markdown to HTML converter with support for:
  - CSS token-based theming
  - YAML frontmatter parsing
  - Table of contents generation
  - Syntax highlighting with chroma
  - Full HTML document wrapper`,
		Args: cobra.RangeArgs(0, 2),
		Run:  run,
	}

	rootCmd.Flags().StringVar(&theme, "theme", "", "Theme preset (default, dark, minimal, forest)")
	rootCmd.Flags().StringVar(&cssFile, "css-file", "", "External CSS file to link")
	rootCmd.Flags().StringVar(&cssInline, "css-inline", "", "Inline CSS to embed")
	rootCmd.Flags().BoolVar(&generateTOC, "toc", false, "Generate table of contents")
	rootCmd.Flags().IntVar(&tocMaxDepth, "toc-depth", 3, "Maximum heading depth for TOC")
	rootCmd.Flags().BoolVar(&highlight, "highlight", true, "Enable syntax highlighting")
	rootCmd.Flags().StringVar(&highlightStyle, "highlight-style", "github", "Syntax highlighting style")
	rootCmd.Flags().BoolVar(&wrapper, "wrapper", false, "Wrap output in full HTML document")
	rootCmd.Flags().StringVar(&title, "title", "", "Document title (used with --wrapper)")
	rootCmd.Flags().StringVar(&lang, "lang", "en", "Document language attribute")
	rootCmd.Flags().BoolVar(&listThemes, "list-themes", false, "List available themes and exit")
	rootCmd.Flags().BoolVar(&listStyles, "list-styles", false, "List available syntax highlighting styles and exit")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Handle list commands
	if listThemes {
		fmt.Println("Available themes:")
		for _, name := range md2html.ListThemes() {
			theme, _ := md2html.GetTheme(name)
			fmt.Printf("  %s - %s\n", name, theme.Description)
		}
		return
	}

	if listStyles {
		fmt.Println("Available syntax highlighting styles:")
		styles := md2html.GetHighlightStyles()
		for i, name := range styles {
			if i > 0 && i%8 == 0 {
				fmt.Println()
			}
			fmt.Printf("  %s", name)
		}
		fmt.Println()
		return
	}

	// Require input file
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: input file required")
		fmt.Fprintln(os.Stderr, "Usage: md2html [input.md] [output.html]")
		os.Exit(1)
	}

	srcFile := args[0]
	outFile := srcFile + ".html"
	if len(args) >= 2 {
		outFile = args[1]
	} else if strings.HasSuffix(srcFile, ".md") {
		outFile = strings.TrimSuffix(srcFile, ".md") + ".html"
	}

	// Build options
	opts := md2html.Options{
		Theme:          theme,
		CSSFile:        cssFile,
		CSSInline:      cssInline,
		GenerateTOC:    generateTOC,
		TOCMaxDepth:    tocMaxDepth,
		Highlight:      highlight,
		HighlightStyle: highlightStyle,
		Wrapper:        wrapper,
		Title:          title,
		Lang:           lang,
		CSSTokens:      make(map[string]string),
		Meta:           make(map[string]string),
	}

	// Convert file
	result, err := md2html.RenderFile(srcFile, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Write output
	if err := os.WriteFile(outFile, []byte(result.HTML), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	// Print summary
	fmt.Printf("Converted: %s -> %s\n", srcFile, outFile)
	if result.Frontmatter != nil && result.Frontmatter.Title != "" {
		fmt.Printf("Title: %s\n", result.Frontmatter.Title)
	}
	if result.TOC != nil && len(result.TOC.Entries) > 0 {
		fmt.Printf("TOC: %d entries\n", len(result.TOC.Entries))
	}
}
