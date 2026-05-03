package md2pdf

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/mogo/os/osutil"
	"github.com/phpdave11/gofpdf"

	"github.com/grokify/sogo/text/markdown/md2html"
)

const (
	FontArial = "Arial"
)

// convertHTMLForGofpdf converts HTML to a format compatible with gofpdf's HTMLBasic.
// HTMLBasic only supports <b>, <i>, <u>, and <a> tags. This function converts
// other common HTML elements to compatible equivalents.
func convertHTMLForGofpdf(htmlStr string) string {
	s := htmlStr

	// Convert content inside <code> tags
	// Note: gofpdf's HTMLBasic strips anything that looks like <tagname>, so users
	// should use {braces} instead of <angles> for placeholder text in code blocks
	codePattern := regexp.MustCompile(`<code>(.*?)</code>`)
	s = codePattern.ReplaceAllStringFunc(s, func(match string) string {
		content := codePattern.FindStringSubmatch(match)[1]
		// Unescape common HTML entities
		content = strings.ReplaceAll(content, "&amp;", "&")
		content = strings.ReplaceAll(content, "&quot;", "\"")
		// Note: &lt; and &gt; are intentionally NOT unescaped because gofpdf
		// would interpret <text> as an HTML tag and strip it
		// Wrap in bold as a visual indicator (since gofpdf doesn't support monospace)
		return "<b>[" + content + "]</b>"
	})

	// Convert headers to bold with line breaks
	for i := 6; i >= 1; i-- {
		openTag := fmt.Sprintf("<h%d[^>]*>", i)
		closeTag := fmt.Sprintf("</h%d>", i)
		s = regexp.MustCompile(openTag).ReplaceAllString(s, "<br><br><b>")
		s = strings.ReplaceAll(s, closeTag, "</b><br><br>")
	}

	// Convert paragraphs to line breaks
	s = regexp.MustCompile(`<p[^>]*>`).ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "</p>", "<br><br>")

	// Convert ordered lists with numbers
	olPattern := regexp.MustCompile(`(?s)<ol[^>]*>(.*?)</ol>`)
	s = olPattern.ReplaceAllStringFunc(s, func(match string) string {
		inner := olPattern.FindStringSubmatch(match)[1]
		counter := 1
		liPattern := regexp.MustCompile(`<li>(.*?)</li>`)
		inner = liPattern.ReplaceAllStringFunc(inner, func(li string) string {
			content := liPattern.FindStringSubmatch(li)[1]
			result := fmt.Sprintf("%d. %s<br>", counter, content)
			counter++
			return result
		})
		return "<br>" + inner + "<br>"
	})

	// Convert unordered lists with dashes
	ulPattern := regexp.MustCompile(`(?s)<ul[^>]*>(.*?)</ul>`)
	s = ulPattern.ReplaceAllStringFunc(s, func(match string) string {
		inner := ulPattern.FindStringSubmatch(match)[1]
		inner = strings.ReplaceAll(inner, "<li>", "- ")
		inner = strings.ReplaceAll(inner, "</li>", "<br>")
		return "<br>" + inner + "<br>"
	})

	// Remove other unsupported tags (html, body, pre, etc.)
	s = regexp.MustCompile(`<pre[^>]*>`).ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "</pre>", "")
	s = regexp.MustCompile(`<html[^>]*>`).ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "</html>", "")
	s = regexp.MustCompile(`<body[^>]*>`).ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "</body>", "")

	// Clean up excessive line breaks
	s = regexp.MustCompile(`(<br>\s*){3,}`).ReplaceAllString(s, "<br><br>")
	s = strings.TrimPrefix(s, "<br>")
	s = strings.TrimPrefix(s, "<br>")

	return s
}

func HTMLToPDFBytes(b []byte) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont(FontArial, "", 12)
	pdf.AddPage()

	// Convert HTML to gofpdf-compatible format
	content := convertHTMLForGofpdf(string(b))

	htmlWriter := pdf.HTMLBasicNew()
	htmlWriter.Write(5, content)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}

func HTMLToPDFFile(srcFile, outFile string, perm os.FileMode) error {
	if b, err := osutil.ReadFileSecure(srcFile); err != nil {
		return err
	} else if bPDF, err := HTMLToPDFBytes(b); err != nil {
		return err
	} else {
		return osutil.WriteFileSecure(outFile, bPDF, perm)
	}
}

func MarkdownToPDFBytes(b []byte, wrapPage bool) ([]byte, error) {
	htBytes := md2html.MarkdownToHTML(b)
	if wrapPage {
		htBytes = append([]byte("<html><body>"), htBytes...)
		htBytes = append(htBytes, []byte("</body></html>")...)
	}
	return HTMLToPDFBytes(htBytes)
}

func MarkdownToPDFFile(srcFile, outFile string, perm os.FileMode) error {
	if b, err := osutil.ReadFileSecure(srcFile); err != nil {
		return err
	} else if bPDF, err := MarkdownToPDFBytes(b, true); err != nil {
		return err
	} else {
		return osutil.WriteFileSecure(outFile, bPDF, perm)
	}
}
