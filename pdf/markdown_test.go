package pdf

import (
	"testing"
)

func TestParseMarkdownSegments_Headers(t *testing.T) {
	tests := []struct {
		input      string
		wantHeader bool
		wantLevel  int
		wantText   string
	}{
		{"# Heading 1", true, 1, "Heading 1"},
		{"## Heading 2", true, 2, "Heading 2"},
		{"### Heading 3", true, 3, "Heading 3"},
		{"#### Heading 4", true, 4, "Heading 4"},
		{"##### Heading 5", true, 5, "Heading 5"},
		{"###### Heading 6", true, 6, "Heading 6"},
		{"Regular text", false, 0, "Regular text"},
		{"#Not a header", false, 0, "#Not a header"}, // no space after #
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			segments := ParseMarkdownSegments(tt.input)
			if len(segments) != 1 {
				t.Fatalf("got %d segments, want 1", len(segments))
			}

			seg := segments[0]
			if seg.IsHeader != tt.wantHeader {
				t.Errorf("IsHeader = %v, want %v", seg.IsHeader, tt.wantHeader)
			}
			if seg.HeaderLevel != tt.wantLevel {
				t.Errorf("HeaderLevel = %d, want %d", seg.HeaderLevel, tt.wantLevel)
			}
			if seg.Text != tt.wantText {
				t.Errorf("Text = %q, want %q", seg.Text, tt.wantText)
			}
		})
	}
}

func TestParseMarkdownSegments_Bold(t *testing.T) {
	tests := []struct {
		input    string
		wantBold bool
		wantText string
	}{
		{"**bold text**", true, "bold text"},
		{"Some **bold** words", true, "Some bold words"},
		{"No bold here", false, "No bold here"},
		{"**multiple** bold **words**", true, "multiple bold words"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			segments := ParseMarkdownSegments(tt.input)
			if len(segments) != 1 {
				t.Fatalf("got %d segments, want 1", len(segments))
			}

			seg := segments[0]
			if seg.IsBold != tt.wantBold {
				t.Errorf("IsBold = %v, want %v", seg.IsBold, tt.wantBold)
			}
			if seg.Text != tt.wantText {
				t.Errorf("Text = %q, want %q", seg.Text, tt.wantText)
			}
		})
	}
}

func TestParseMarkdownSegments_Italic(t *testing.T) {
	tests := []struct {
		input      string
		wantItalic bool
		wantText   string
	}{
		{"*italic text*", true, "italic text"},
		{"Some *italic* words", true, "Some italic words"},
		{"No italic here", false, "No italic here"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			segments := ParseMarkdownSegments(tt.input)
			if len(segments) != 1 {
				t.Fatalf("got %d segments, want 1", len(segments))
			}

			seg := segments[0]
			if seg.IsItalic != tt.wantItalic {
				t.Errorf("IsItalic = %v, want %v", seg.IsItalic, tt.wantItalic)
			}
			if seg.Text != tt.wantText {
				t.Errorf("Text = %q, want %q", seg.Text, tt.wantText)
			}
		})
	}
}

func TestParseMarkdownSegments_MultiLine(t *testing.T) {
	input := `# Welcome

This is a paragraph.

## Section

More text here.`

	segments := ParseMarkdownSegments(input)

	expected := []struct {
		text     string
		isHeader bool
		level    int
	}{
		{"Welcome", true, 1},
		{"", false, 0},
		{"This is a paragraph.", false, 0},
		{"", false, 0},
		{"Section", true, 2},
		{"", false, 0},
		{"More text here.", false, 0},
	}

	if len(segments) != len(expected) {
		t.Fatalf("got %d segments, want %d", len(segments), len(expected))
	}

	for i, exp := range expected {
		seg := segments[i]
		if seg.Text != exp.text {
			t.Errorf("segment[%d].Text = %q, want %q", i, seg.Text, exp.text)
		}
		if seg.IsHeader != exp.isHeader {
			t.Errorf("segment[%d].IsHeader = %v, want %v", i, seg.IsHeader, exp.isHeader)
		}
		if seg.HeaderLevel != exp.level {
			t.Errorf("segment[%d].HeaderLevel = %d, want %d", i, seg.HeaderLevel, exp.level)
		}
	}
}

func TestStripMarkdown(t *testing.T) {
	input := `# Welcome

This is **bold** and *italic*.

## Section Two`

	result := StripMarkdown(input)
	expected := `Welcome
This is bold and italic.
Section Two`

	if result != expected {
		t.Errorf("StripMarkdown:\ngot:\n%s\nwant:\n%s", result, expected)
	}
}

func TestParseMarkdownSegments_EmptyInput(t *testing.T) {
	segments := ParseMarkdownSegments("")
	if len(segments) != 1 {
		t.Fatalf("got %d segments, want 1", len(segments))
	}
	if segments[0].Text != "" {
		t.Errorf("Text = %q, want empty string", segments[0].Text)
	}
}

func TestParseMarkdownSegments_WindowsLineEndings(t *testing.T) {
	input := "Line 1\r\nLine 2\r\n"
	segments := ParseMarkdownSegments(input)

	if len(segments) != 3 {
		t.Fatalf("got %d segments, want 3", len(segments))
	}

	if segments[0].Text != "Line 1" {
		t.Errorf("segment[0].Text = %q, want 'Line 1'", segments[0].Text)
	}
	if segments[1].Text != "Line 2" {
		t.Errorf("segment[1].Text = %q, want 'Line 2'", segments[1].Text)
	}
}
