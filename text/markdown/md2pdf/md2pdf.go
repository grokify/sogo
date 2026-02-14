package md2pdf

import (
	"bytes"
	"fmt"
	"os"

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

	html := pdf.HTMLBasicNew()
	html.Write(5, string(b))

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}

func HTMLToPDFFile(srcFile, outFile string, perm os.FileMode) error {
	if b, err := os.ReadFile(srcFile); err != nil {
		return err
	} else if bPDF, err := HTMLToPDFBytes(b); err != nil {
		return err
	} else {
		return os.WriteFile(outFile, bPDF, perm)
	}
}

func MarkdownToPDFBytes(b []byte, wrapPage bool) ([]byte, error) {
	htBytes := md2html.MarkdownToHTML(b)
	if wrapPage {
		htBytes = append([]byte("<html><body>"), htBytes...)
		htBytes = append(htBytes, []byte("</body></html>")...)
	}
	fmt.Println(string(htBytes))
	return HTMLToPDFBytes(htBytes)
}

func MarkdownToPDFFile(srcFile, outFile string, perm os.FileMode) error {
	if b, err := os.ReadFile(srcFile); err != nil {
		return err
	} else if bPDF, err := MarkdownToPDFBytes(b, true); err != nil {
		return err
	} else {
		return os.WriteFile(outFile, bPDF, perm)
	}
}
