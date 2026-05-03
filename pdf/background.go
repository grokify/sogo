package pdf

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/grokify/mogo/image/imageutil"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// PageSize represents standard page dimensions in points (72 dpi).
type PageSize struct {
	Width  float64
	Height float64
}

var (
	// PageSizeLetter is 8.5x11 inches at 72 dpi.
	PageSizeLetter = PageSize{Width: 612, Height: 792}
	// PageSizeA4 is the standard A4 paper size at 72 dpi.
	PageSizeA4 = PageSize{Width: 595, Height: 842}
)

// TextStyle defines styling for text overlay.
type TextStyle struct {
	FontName    string  // e.g., "Helvetica", "Courier"
	FontSize    float64 // points
	Color       string  // e.g., "0 0 0" (RGB) or "#000000"
	MarginTop   float64 // points from top
	MarginLeft  float64 // points from left
	MarginRight float64 // points from right (used to calculate line width)
	LineHeight  float64 // multiplier (e.g., 1.5)
}

// LogoPosition represents anchor positions for logo placement.
type LogoPosition string

const (
	LogoPositionTopLeft      LogoPosition = "tl"
	LogoPositionTopCenter    LogoPosition = "tc"
	LogoPositionTopRight     LogoPosition = "tr"
	LogoPositionLeft         LogoPosition = "l"
	LogoPositionCenter       LogoPosition = "c"
	LogoPositionRight        LogoPosition = "r"
	LogoPositionBottomLeft   LogoPosition = "bl"
	LogoPositionBottomCenter LogoPosition = "bc"
	LogoPositionBottomRight  LogoPosition = "br"
)

// LogoOptions configures logo placement on the PDF.
type LogoOptions struct {
	Position LogoPosition // Anchor position (e.g., "tl", "br", "c")
	OffsetX  float64      // Horizontal offset from anchor in points
	OffsetY  float64      // Vertical offset from anchor in points
	Scale    float64      // Scale factor (0.0-1.0 relative, or absolute if ScaleAbs is true)
	ScaleAbs bool         // If true, Scale represents absolute width in points
	Opacity  float64      // Opacity (0.0-1.0), defaults to 1.0 if zero
}

// DefaultLogoOptions returns sensible defaults for logo placement.
func DefaultLogoOptions() LogoOptions {
	return LogoOptions{
		Position: LogoPositionBottomRight,
		OffsetX:  -20,
		OffsetY:  20,
		Scale:    0.15,
		ScaleAbs: false,
		Opacity:  1.0,
	}
}

// DefaultTextStyle returns a reasonable default text style.
func DefaultTextStyle() TextStyle {
	return TextStyle{
		FontName:    "Helvetica",
		FontSize:    12,
		Color:       "0 0 0",
		MarginTop:   72,
		MarginLeft:  72,
		MarginRight: 72,
		LineHeight:  1.5,
	}
}

// BackgroundPDFOptions configures PDF generation.
type BackgroundPDFOptions struct {
	PageSize   PageSize
	TextStyle  TextStyle
	Logo       *LogoOptions // Optional logo configuration (nil = no logo)
	LogoReader io.Reader    // Logo image reader (used with CreateBackgroundPDF)
	LogoPath   string       // Logo image path (used with CreateBackgroundPDFFile)
}

// DefaultBackgroundPDFOptions returns default options for letter size with standard text styling.
func DefaultBackgroundPDFOptions() BackgroundPDFOptions {
	return BackgroundPDFOptions{
		PageSize:  PageSizeLetter,
		TextStyle: DefaultTextStyle(),
	}
}

// CreateBackgroundPDF creates a PDF with a background image, optional logo, and text overlay.
// The background image is scaled to cover the page (minimum scaling to fill, then center-cropped).
// Text is rendered from Markdown with basic formatting support (headers, paragraphs).
// Logo is placed according to LogoOptions if provided.
func CreateBackgroundPDF(
	imageReader io.Reader,
	markdownText string,
	opts BackgroundPDFOptions,
) ([]byte, error) {
	// Step 1: Prepare the background image (scale to cover, center crop)
	img, err := prepareBackgroundImage(imageReader, opts.PageSize)
	if err != nil {
		return nil, fmt.Errorf("preparing background image: %w", err)
	}

	// Step 2: Create PDF with the background image
	pdfBytes, err := createImagePDF(img, opts.PageSize)
	if err != nil {
		return nil, fmt.Errorf("creating image PDF: %w", err)
	}

	// Step 3: Add logo overlay if provided
	if opts.Logo != nil && opts.LogoReader != nil {
		pdfBytes, err = addLogoOverlay(pdfBytes, opts.LogoReader, *opts.Logo)
		if err != nil {
			return nil, fmt.Errorf("adding logo overlay: %w", err)
		}
	}

	// Step 4: Add text overlay if markdown text is provided
	if strings.TrimSpace(markdownText) != "" {
		pdfBytes, err = addTextOverlay(pdfBytes, markdownText, opts.TextStyle)
		if err != nil {
			return nil, fmt.Errorf("adding text overlay: %w", err)
		}
	}

	return pdfBytes, nil
}

// CreateBackgroundPDFFile is a convenience wrapper that reads an image file,
// creates the PDF, and writes it to an output file.
// If opts.LogoPath is set and opts.Logo is configured, the logo will be added.
func CreateBackgroundPDFFile(
	imagePath string,
	markdownText string,
	opts BackgroundPDFOptions,
	outPath string,
) error {
	imgFile, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("opening image file: %w", err)
	}
	defer imgFile.Close()

	// If logo path is provided and logo options are set, open the logo file
	if opts.LogoPath != "" && opts.Logo != nil {
		logoFile, err := os.Open(opts.LogoPath)
		if err != nil {
			return fmt.Errorf("opening logo file: %w", err)
		}
		defer logoFile.Close()
		opts.LogoReader = logoFile
	}

	pdfBytes, err := CreateBackgroundPDF(imgFile, markdownText, opts)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outPath, pdfBytes, 0600); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	return nil
}

// prepareBackgroundImage scales an image to cover the page dimensions and center-crops
// to the exact page size. This implements "cover" scaling behavior.
func prepareBackgroundImage(r io.Reader, pageSize PageSize) (image.Image, error) {
	// Decode the image
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	// Target dimensions in pixels (at 72 dpi, points == pixels)
	targetWidth := int(pageSize.Width)
	targetHeight := int(pageSize.Height)

	srcWidth := img.Bounds().Dx()
	srcHeight := img.Bounds().Dy()

	// Calculate scale factors for cover behavior
	// Cover means we scale to fill the entire target area (some parts may be cropped)
	scaleX := float64(targetWidth) / float64(srcWidth)
	scaleY := float64(targetHeight) / float64(srcHeight)

	// Use the larger scale factor to ensure the image covers the target area
	scale := scaleX
	if scaleY > scaleX {
		scale = scaleY
	}

	// Calculate scaled dimensions
	scaledWidth := int(float64(srcWidth) * scale)
	scaledHeight := int(float64(srcHeight) * scale)

	// Resize the image
	scaled := imageutil.Resize(scaledWidth, scaledHeight, img, imageutil.ScalerBest())

	// Center crop to exact target dimensions
	if scaled.Bounds().Dx() > targetWidth {
		scaled = imageutil.CropX(scaled, targetWidth, imageutil.AlignCenter)
	}
	if scaled.Bounds().Dy() > targetHeight {
		scaled = imageutil.CropY(scaled, targetHeight, imageutil.AlignCenter)
	}

	return scaled, nil
}

// createImagePDF creates a single-page PDF with the image filling the page.
func createImagePDF(img image.Image, pageSize PageSize) ([]byte, error) {
	// Write image to a temporary file (pdfcpu needs a reader for ImportImages)
	tmpFile, err := os.CreateTemp("", "bg-*.png")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if err := png.Encode(tmpFile, img); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("encoding PNG: %w", err)
	}
	tmpFile.Close()

	// Configure import to place image at full page size
	imp := &pdfcpu.Import{
		PageDim:  &types.Dim{Width: pageSize.Width, Height: pageSize.Height},
		PageSize: "",
		UserDim:  true,
		Pos:      types.Full,
		Scale:    1.0,
		InpUnit:  types.POINTS,
	}

	// Create PDF with the image
	var pdfBuf bytes.Buffer
	imgFile, err := os.Open(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("reopening temp file: %w", err)
	}
	defer imgFile.Close()

	if err := api.ImportImages(nil, &pdfBuf, []io.Reader{imgFile}, imp, nil); err != nil {
		return nil, fmt.Errorf("importing image to PDF: %w", err)
	}

	return pdfBuf.Bytes(), nil
}

// addTextOverlay adds text from markdown content onto the PDF.
func addTextOverlay(pdfBytes []byte, markdown string, style TextStyle) ([]byte, error) {
	// Parse markdown into text segments
	segments := ParseMarkdownSegments(markdown)
	if len(segments) == 0 {
		return pdfBytes, nil
	}

	// Build watermarks map for all text segments
	// We'll position text from top-left, working downward
	currentY := style.MarginTop
	watermarks := make([]*model.Watermark, 0, len(segments))

	for _, seg := range segments {
		fontSize := style.FontSize
		if seg.IsHeader {
			// Scale font size based on header level (h1 = largest)
			// h1: 2.0x, h2: 1.8x, h3: 1.6x, h4: 1.4x, h5: 1.2x, h6: 1.0x
			headerScale := 2.0 - float64(seg.HeaderLevel-1)*0.2
			fontSize = style.FontSize * headerScale
		}

		// Text to render (strip markdown formatting for plain text output)
		text := seg.Text
		if text == "" {
			// Empty line - just advance Y position
			currentY += fontSize * style.LineHeight
			continue
		}

		// Calculate Y offset from top of page
		// pdfcpu uses bottom-left origin, so we need to convert
		// offset is relative to the anchor position
		yFromTop := currentY

		// Build watermark description string
		// Using position tl (top-left) with offset
		// Note: pdfcpu expects integer points
		desc := fmt.Sprintf(
			"font:%s, points:%d, fillcolor:%s, pos:tl, offset:%.1f -%.1f, scale:1 abs, opacity:1",
			style.FontName,
			int(fontSize),
			style.Color,
			style.MarginLeft,
			yFromTop,
		)

		wm, err := api.TextWatermark(text, desc, true, false, types.POINTS)
		if err != nil {
			return nil, fmt.Errorf("creating text watermark: %w", err)
		}
		watermarks = append(watermarks, wm)

		// Advance Y position for next line
		currentY += fontSize * style.LineHeight
	}

	if len(watermarks) == 0 {
		return pdfBytes, nil
	}

	// Apply all watermarks to the PDF
	rs := bytes.NewReader(pdfBytes)
	var outBuf bytes.Buffer

	// Create a watermarks map for page 1
	wmMap := make(map[int][]*model.Watermark)
	wmMap[1] = watermarks

	if err := api.AddWatermarksSliceMap(rs, &outBuf, wmMap, nil); err != nil {
		return nil, fmt.Errorf("adding watermarks: %w", err)
	}

	return outBuf.Bytes(), nil
}

// addLogoOverlay adds a logo image onto the PDF at the specified position.
func addLogoOverlay(pdfBytes []byte, logoReader io.Reader, opts LogoOptions) ([]byte, error) {
	// Read logo data into buffer (we may need to use it with the file-based API)
	logoData, err := io.ReadAll(logoReader)
	if err != nil {
		return nil, fmt.Errorf("reading logo data: %w", err)
	}

	// Write logo to temp file (pdfcpu ImageWatermark needs a file path with proper extension)
	tmpFile, err := os.CreateTemp("", "logo-*.png")
	if err != nil {
		return nil, fmt.Errorf("creating temp logo file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(logoData); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("writing temp logo file: %w", err)
	}
	tmpFile.Close()

	// Build watermark description
	opacity := opts.Opacity
	if opacity <= 0 {
		opacity = 1.0
	}

	scaleType := "rel"
	if opts.ScaleAbs {
		scaleType = "abs"
	}

	desc := fmt.Sprintf(
		"pos:%s, offset:%.1f %.1f, scale:%.3f %s, opacity:%.2f",
		opts.Position,
		opts.OffsetX,
		opts.OffsetY,
		opts.Scale,
		scaleType,
		opacity,
	)

	wm, err := api.ImageWatermark(tmpPath, desc, true, false, types.POINTS)
	if err != nil {
		return nil, fmt.Errorf("creating image watermark: %w", err)
	}

	// Apply watermark to PDF
	rs := bytes.NewReader(pdfBytes)
	var outBuf bytes.Buffer

	if err := api.AddWatermarks(rs, &outBuf, nil, wm, nil); err != nil {
		return nil, fmt.Errorf("adding logo watermark: %w", err)
	}

	return outBuf.Bytes(), nil
}
