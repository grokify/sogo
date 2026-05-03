package md2pdf

import (
	"bytes"
	"fmt"
	"html"
	"os"

	"github.com/grokify/mogo/os/osutil"
	"github.com/phpdave11/gofpdf"

	"github.com/grokify/sogo/text/markdown/md2html"
)

const (
	FontArial = "Arial"
)

func HTMLToPDFBytes(b []byte) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont(FontArial, "", 12)
	pdf.AddPage()

	// Unescape HTML entities (e.g., &lt; → <) since gofpdf's HTMLBasic
	// renders them as literal text instead of interpreting them
	content := html.UnescapeString(string(b))

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
