package pdf

import (
	"regexp"
	"strings"
)

// TextSegment represents a parsed segment of markdown text with styling information.
type TextSegment struct {
	Text        string
	IsHeader    bool
	HeaderLevel int  // 1-6 for headers
	IsBold      bool // text was wrapped in **
	IsItalic    bool // text was wrapped in *
}

var (
	// headerPattern matches markdown headers (# to ######)
	headerPattern = regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
	// boldPattern matches **text**
	boldPattern = regexp.MustCompile(`\*\*([^*]+)\*\*`)
)

// ParseMarkdownSegments parses markdown text into a slice of TextSegments.
// This is a simple line-by-line parser that handles:
//   - Headers (# to ######)
//   - Bold (**text**)
//   - Italic (*text*)
//   - Plain paragraphs
//
// Note: This strips formatting markers for plain text rendering since pdfcpu's
// TextWatermark doesn't support rich text styling. The IsBold and IsItalic
// flags are preserved for potential future use with fonts that support it.
func ParseMarkdownSegments(md string) []TextSegment {
	lines := strings.Split(md, "\n")
	segments := make([]TextSegment, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		seg := parseLine(line)
		segments = append(segments, seg)
	}

	return segments
}

// parseLine parses a single line of markdown into a TextSegment.
func parseLine(line string) TextSegment {
	seg := TextSegment{}

	// Check for header
	if matches := headerPattern.FindStringSubmatch(line); len(matches) == 3 {
		seg.IsHeader = true
		seg.HeaderLevel = len(matches[1])
		line = matches[2]
	}

	// Strip bold markers and set flag
	if boldPattern.MatchString(line) {
		seg.IsBold = true
		line = boldPattern.ReplaceAllString(line, "$1")
	}

	// Strip italic markers and set flag
	// More careful handling to avoid false positives
	if strings.Contains(line, "*") && !strings.Contains(line, "**") {
		seg.IsItalic = true
		// Simple italic stripping: remove single * around words
		line = stripItalic(line)
	}

	seg.Text = strings.TrimSpace(line)
	return seg
}

// stripItalic removes single asterisks used for italic formatting.
func stripItalic(s string) string {
	result := strings.Builder{}
	inItalic := false
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '*' {
			// Check if this is a single asterisk (not double)
			isDouble := (i+1 < len(runes) && runes[i+1] == '*') ||
				(i > 0 && runes[i-1] == '*')
			if !isDouble {
				inItalic = !inItalic
				continue
			}
		}
		result.WriteRune(r)
	}

	return result.String()
}

// StripMarkdown removes all markdown formatting and returns plain text.
func StripMarkdown(md string) string {
	segments := ParseMarkdownSegments(md)
	lines := make([]string, 0, len(segments))
	for _, seg := range segments {
		if seg.Text != "" {
			lines = append(lines, seg.Text)
		}
	}
	return strings.Join(lines, "\n")
}
