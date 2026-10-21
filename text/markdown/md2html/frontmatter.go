package md2html

import (
	"bytes"
	"time"

	"gopkg.in/yaml.v3"
)

// Frontmatter represents YAML frontmatter extracted from markdown.
type Frontmatter struct {
	Title       string            `yaml:"title" json:"title"`
	Subtitle    string            `yaml:"subtitle" json:"subtitle"`
	Description string            `yaml:"description" json:"description"`
	Author      string            `yaml:"author" json:"author"`
	Date        time.Time         `yaml:"date" json:"date"`
	Version     string            `yaml:"version" json:"version"`
	Tags        []string          `yaml:"tags" json:"tags"`
	Draft       bool              `yaml:"draft" json:"draft"`
	Extra       map[string]string `yaml:"extra" json:"extra"`
}

var frontmatterDelim = []byte("---")

// splitFrontmatter splits content into the YAML block and the remaining markdown.
// The opening and closing delimiters must each be a line consisting only of "---".
// ok is false if content has no complete frontmatter block.
func splitFrontmatter(content []byte) (yamlContent, markdown []byte, ok bool) {
	content = bytes.TrimSpace(content)
	first, rest, found := bytes.Cut(content, []byte("\n"))
	if !found || !bytes.Equal(bytes.TrimRight(first, " \t\r"), frontmatterDelim) {
		return nil, content, false
	}

	offset := 0
	for len(rest[offset:]) > 0 {
		line, _, _ := bytes.Cut(rest[offset:], []byte("\n"))
		if bytes.Equal(bytes.TrimRight(line, " \t\r"), frontmatterDelim) {
			return rest[:offset], bytes.TrimSpace(rest[offset+len(line):]), true
		}
		offset += len(line) + 1
		if offset > len(rest) {
			break
		}
	}
	return nil, content, false
}

// ParseFrontmatter extracts YAML frontmatter from markdown content.
// Returns the frontmatter, the remaining markdown content, and any error.
func ParseFrontmatter(content []byte) (*Frontmatter, []byte, error) {
	yamlContent, markdown, ok := splitFrontmatter(content)
	if !ok {
		return nil, bytes.TrimSpace(content), nil
	}

	var fm Frontmatter
	if err := yaml.Unmarshal(yamlContent, &fm); err != nil {
		return nil, bytes.TrimSpace(content), err
	}

	return &fm, markdown, nil
}

// HasFrontmatter checks if content starts with a complete YAML frontmatter block.
func HasFrontmatter(content []byte) bool {
	_, _, ok := splitFrontmatter(content)
	return ok
}
