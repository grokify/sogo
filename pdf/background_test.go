package pdf

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// createTestImage creates a simple colored test image.
func createTestImage(width, height int, clr color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, clr)
		}
	}
	return img
}

func TestPrepareBackgroundImage_Portrait(t *testing.T) {
	// Create a portrait image (taller than wide)
	img := createTestImage(400, 600, color.RGBA{R: 128, G: 128, B: 255, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encoding test image: %v", err)
	}

	result, err := prepareBackgroundImage(&buf, PageSizeLetter)
	if err != nil {
		t.Fatalf("prepareBackgroundImage: %v", err)
	}

	// Result should match letter page dimensions
	if result.Bounds().Dx() != int(PageSizeLetter.Width) {
		t.Errorf("width = %d, want %d", result.Bounds().Dx(), int(PageSizeLetter.Width))
	}
	if result.Bounds().Dy() != int(PageSizeLetter.Height) {
		t.Errorf("height = %d, want %d", result.Bounds().Dy(), int(PageSizeLetter.Height))
	}
}

func TestPrepareBackgroundImage_Landscape(t *testing.T) {
	// Create a landscape image (wider than tall)
	img := createTestImage(800, 400, color.RGBA{R: 0, G: 128, B: 255, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encoding test image: %v", err)
	}

	result, err := prepareBackgroundImage(&buf, PageSizeLetter)
	if err != nil {
		t.Fatalf("prepareBackgroundImage: %v", err)
	}

	// Result should match letter page dimensions
	if result.Bounds().Dx() != int(PageSizeLetter.Width) {
		t.Errorf("width = %d, want %d", result.Bounds().Dx(), int(PageSizeLetter.Width))
	}
	if result.Bounds().Dy() != int(PageSizeLetter.Height) {
		t.Errorf("height = %d, want %d", result.Bounds().Dy(), int(PageSizeLetter.Height))
	}
}

func TestPrepareBackgroundImage_Square(t *testing.T) {
	// Create a square image
	img := createTestImage(500, 500, color.RGBA{R: 255, G: 128, B: 0, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encoding test image: %v", err)
	}

	result, err := prepareBackgroundImage(&buf, PageSizeLetter)
	if err != nil {
		t.Fatalf("prepareBackgroundImage: %v", err)
	}

	// Result should match letter page dimensions
	if result.Bounds().Dx() != int(PageSizeLetter.Width) {
		t.Errorf("width = %d, want %d", result.Bounds().Dx(), int(PageSizeLetter.Width))
	}
	if result.Bounds().Dy() != int(PageSizeLetter.Height) {
		t.Errorf("height = %d, want %d", result.Bounds().Dy(), int(PageSizeLetter.Height))
	}
}

func TestCreateBackgroundPDF_Basic(t *testing.T) {
	// Create a simple test image
	img := createTestImage(800, 1000, color.RGBA{R: 200, G: 220, B: 240, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encoding test image: %v", err)
	}

	markdown := `# Welcome

This is a test document with **important** text.

## Section One

Some paragraph content here.
`

	opts := DefaultBackgroundPDFOptions()
	pdfBytes, err := CreateBackgroundPDF(&buf, markdown, opts)
	if err != nil {
		t.Fatalf("CreateBackgroundPDF: %v", err)
	}

	// Basic sanity check: PDF should start with %PDF-
	if len(pdfBytes) < 5 || string(pdfBytes[:5]) != "%PDF-" {
		t.Error("output doesn't appear to be a valid PDF")
	}
}

func TestCreateBackgroundPDF_EmptyText(t *testing.T) {
	// Create a simple test image
	img := createTestImage(600, 800, color.RGBA{R: 100, G: 150, B: 200, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encoding test image: %v", err)
	}

	opts := DefaultBackgroundPDFOptions()
	pdfBytes, err := CreateBackgroundPDF(&buf, "", opts)
	if err != nil {
		t.Fatalf("CreateBackgroundPDF with empty text: %v", err)
	}

	// Should still produce valid PDF
	if len(pdfBytes) < 5 || string(pdfBytes[:5]) != "%PDF-" {
		t.Error("output doesn't appear to be a valid PDF")
	}
}

func TestCreateBackgroundPDFFile(t *testing.T) {
	// Create temp directory for test files
	tmpDir := t.TempDir()

	// Create a test image file
	img := createTestImage(800, 1000, color.RGBA{R: 180, G: 200, B: 220, A: 255})
	imgPath := filepath.Join(tmpDir, "test_bg.png")

	imgFile, err := os.Create(imgPath)
	if err != nil {
		t.Fatalf("creating image file: %v", err)
	}
	if err := png.Encode(imgFile, img); err != nil {
		imgFile.Close()
		t.Fatalf("encoding image: %v", err)
	}
	imgFile.Close()

	// Create PDF
	outPath := filepath.Join(tmpDir, "output.pdf")
	markdown := `# Test Document

This is **bold** and *italic* text.
`

	opts := BackgroundPDFOptions{
		PageSize: PageSizeLetter,
		TextStyle: TextStyle{
			FontName:    "Helvetica",
			FontSize:    14,
			Color:       "0 0 0",
			MarginTop:   100,
			MarginLeft:  72,
			MarginRight: 72,
			LineHeight:  1.5,
		},
	}

	if err := CreateBackgroundPDFFile(imgPath, markdown, opts, outPath); err != nil {
		t.Fatalf("CreateBackgroundPDFFile: %v", err)
	}

	// Verify output file exists and is valid PDF
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Error("output doesn't appear to be a valid PDF")
	}
}

func TestPageSizeConstants(t *testing.T) {
	// Letter: 8.5 x 11 inches at 72 dpi
	if PageSizeLetter.Width != 612 {
		t.Errorf("PageSizeLetter.Width = %f, want 612", PageSizeLetter.Width)
	}
	if PageSizeLetter.Height != 792 {
		t.Errorf("PageSizeLetter.Height = %f, want 792", PageSizeLetter.Height)
	}

	// A4: 210 x 297 mm at 72 dpi
	// 210mm = 8.27in = 595.44pt (rounded to 595)
	// 297mm = 11.69in = 841.68pt (rounded to 842)
	if PageSizeA4.Width != 595 {
		t.Errorf("PageSizeA4.Width = %f, want 595", PageSizeA4.Width)
	}
	if PageSizeA4.Height != 842 {
		t.Errorf("PageSizeA4.Height = %f, want 842", PageSizeA4.Height)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultBackgroundPDFOptions()

	if opts.PageSize != PageSizeLetter {
		t.Errorf("default page size is not Letter")
	}

	style := opts.TextStyle
	if style.FontName != "Helvetica" {
		t.Errorf("default font = %q, want Helvetica", style.FontName)
	}
	if style.FontSize != 12 {
		t.Errorf("default font size = %f, want 12", style.FontSize)
	}
	if style.MarginTop != 72 {
		t.Errorf("default margin top = %f, want 72", style.MarginTop)
	}
}

func TestCreateBackgroundPDF_WithLogo(t *testing.T) {
	// Create background image
	bgImg := createTestImage(800, 1000, color.RGBA{R: 200, G: 220, B: 240, A: 255})
	var bgBuf bytes.Buffer
	if err := png.Encode(&bgBuf, bgImg); err != nil {
		t.Fatalf("encoding background image: %v", err)
	}

	// Create logo image (small, different color)
	logoImg := createTestImage(100, 50, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	var logoBuf bytes.Buffer
	if err := png.Encode(&logoBuf, logoImg); err != nil {
		t.Fatalf("encoding logo image: %v", err)
	}

	logoOpts := DefaultLogoOptions()
	opts := BackgroundPDFOptions{
		PageSize:   PageSizeLetter,
		TextStyle:  DefaultTextStyle(),
		Logo:       &logoOpts,
		LogoReader: &logoBuf,
	}

	pdfBytes, err := CreateBackgroundPDF(&bgBuf, "# Document with Logo", opts)
	if err != nil {
		t.Fatalf("CreateBackgroundPDF with logo: %v", err)
	}

	if len(pdfBytes) < 5 || string(pdfBytes[:5]) != "%PDF-" {
		t.Error("output doesn't appear to be a valid PDF")
	}
}

func TestCreateBackgroundPDFFile_WithLogo(t *testing.T) {
	tmpDir := t.TempDir()

	// Create background image file
	bgImg := createTestImage(800, 1000, color.RGBA{R: 180, G: 200, B: 220, A: 255})
	bgPath := filepath.Join(tmpDir, "background.png")
	bgFile, err := os.Create(bgPath)
	if err != nil {
		t.Fatalf("creating background file: %v", err)
	}
	if err := png.Encode(bgFile, bgImg); err != nil {
		bgFile.Close()
		t.Fatalf("encoding background: %v", err)
	}
	bgFile.Close()

	// Create logo image file
	logoImg := createTestImage(80, 80, color.RGBA{R: 0, G: 100, B: 200, A: 255})
	logoPath := filepath.Join(tmpDir, "logo.png")
	logoFile, err := os.Create(logoPath)
	if err != nil {
		t.Fatalf("creating logo file: %v", err)
	}
	if err := png.Encode(logoFile, logoImg); err != nil {
		logoFile.Close()
		t.Fatalf("encoding logo: %v", err)
	}
	logoFile.Close()

	outPath := filepath.Join(tmpDir, "output_with_logo.pdf")

	logoOpts := LogoOptions{
		Position: LogoPositionBottomRight,
		OffsetX:  -30,
		OffsetY:  30,
		Scale:    0.1,
		ScaleAbs: false,
		Opacity:  0.9,
	}

	opts := BackgroundPDFOptions{
		PageSize:  PageSizeLetter,
		TextStyle: DefaultTextStyle(),
		Logo:      &logoOpts,
		LogoPath:  logoPath,
	}

	if err := CreateBackgroundPDFFile(bgPath, "# Test with Logo\n\nSome content here.", opts, outPath); err != nil {
		t.Fatalf("CreateBackgroundPDFFile with logo: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Error("output doesn't appear to be a valid PDF")
	}
}

func TestLogoPositions(t *testing.T) {
	positions := []LogoPosition{
		LogoPositionTopLeft,
		LogoPositionTopCenter,
		LogoPositionTopRight,
		LogoPositionLeft,
		LogoPositionCenter,
		LogoPositionRight,
		LogoPositionBottomLeft,
		LogoPositionBottomCenter,
		LogoPositionBottomRight,
	}

	// Create background and logo images once
	bgImg := createTestImage(612, 792, color.RGBA{R: 240, G: 240, B: 240, A: 255})
	logoImg := createTestImage(50, 50, color.RGBA{R: 0, G: 128, B: 0, A: 255})

	for _, pos := range positions {
		t.Run(string(pos), func(t *testing.T) {
			var bgBuf, logoBuf bytes.Buffer
			if err := png.Encode(&bgBuf, bgImg); err != nil {
				t.Fatalf("encoding background: %v", err)
			}
			if err := png.Encode(&logoBuf, logoImg); err != nil {
				t.Fatalf("encoding logo: %v", err)
			}

			logoOpts := LogoOptions{
				Position: pos,
				Scale:    0.1,
				Opacity:  1.0,
			}

			opts := BackgroundPDFOptions{
				PageSize:   PageSizeLetter,
				TextStyle:  DefaultTextStyle(),
				Logo:       &logoOpts,
				LogoReader: &logoBuf,
			}

			pdfBytes, err := CreateBackgroundPDF(&bgBuf, "", opts)
			if err != nil {
				t.Fatalf("CreateBackgroundPDF with position %s: %v", pos, err)
			}

			if len(pdfBytes) < 5 || string(pdfBytes[:5]) != "%PDF-" {
				t.Error("output doesn't appear to be a valid PDF")
			}
		})
	}
}

func TestDefaultLogoOptions(t *testing.T) {
	opts := DefaultLogoOptions()

	if opts.Position != LogoPositionBottomRight {
		t.Errorf("Position = %q, want %q", opts.Position, LogoPositionBottomRight)
	}
	if opts.Scale != 0.15 {
		t.Errorf("Scale = %f, want 0.15", opts.Scale)
	}
	if opts.Opacity != 1.0 {
		t.Errorf("Opacity = %f, want 1.0", opts.Opacity)
	}
	if opts.ScaleAbs {
		t.Error("ScaleAbs should be false by default")
	}
}
