package pdf

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"strings"

	"github.com/grokify/mogo/image/imageutil"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	pdfcolor "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
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

// TextAlign represents horizontal text alignment.
type TextAlign string

const (
	TextAlignLeft   TextAlign = "left"
	TextAlignCenter TextAlign = "center"
	TextAlignRight  TextAlign = "right"
)

// TextStyle defines styling for text overlay.
type TextStyle struct {
	FontName    string    // e.g., "Helvetica", "Courier"
	FontSize    float64   // points
	Color       string    // e.g., "0 0 0" (RGB) or "#000000"
	MarginTop   float64   // points from top
	MarginLeft  float64   // points from left
	MarginRight float64   // points from right (used to calculate line width)
	LineHeight  float64   // multiplier (e.g., 1.5)
	Align       TextAlign // horizontal alignment (default: left)
}

// TextBlock represents a single block of text with optional style overrides.
// Empty values inherit from the parent TextStyle.
type TextBlock struct {
	Text      string    // The text to render
	FontName  string    // Override font (empty = inherit from TextStyle)
	FontSize  float64   // Override font size in points (0 = inherit from TextStyle)
	Color     string    // Override color (empty = inherit from TextStyle)
	MarginTop float64   // Additional top margin before this block in points
	Align     TextAlign // Override alignment (empty = inherit from TextStyle)
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
	TextBlocks []TextBlock  // If set, used instead of markdown text for precise control
	Logo       *LogoOptions // Optional logo configuration (nil = no logo)
	LogoReader io.Reader    // Logo image reader (used with CreateBackgroundPDF)
	LogoPath   string       // Logo image path (used with CreateBackgroundPDFFile)
}

// Page represents a single page in a multi-page PDF.
type Page struct {
	BackgroundReader io.Reader    // Background image reader (nil = white page)
	BackgroundPath   string       // Background image path (file-based alternative)
	Logo             *LogoOptions // Optional logo configuration
	LogoReader       io.Reader    // Logo image reader
	LogoPath         string       // Logo image path
	TextStyle        TextStyle    // Base text style
	TextBlocks       []TextBlock  // Styled text blocks (takes precedence over MarkdownText)
	MarkdownText     string       // Markdown text (used if TextBlocks is empty)
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
// Text is rendered from either TextBlocks (if provided) or Markdown with basic formatting support.
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

	// Step 4: Add text overlay
	// TextBlocks take precedence over markdown text
	if len(opts.TextBlocks) > 0 {
		pdfBytes, err = addTextBlocksOverlay(pdfBytes, opts.TextBlocks, opts.TextStyle)
		if err != nil {
			return nil, fmt.Errorf("adding text blocks overlay: %w", err)
		}
	} else if strings.TrimSpace(markdownText) != "" {
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

	// Calculate scaled dimensions (ceil to ensure full coverage, crop handles overshoot)
	scaledWidth := int(math.Ceil(float64(srcWidth) * scale))
	scaledHeight := int(math.Ceil(float64(srcHeight) * scale))

	// Resize the image
	scaled := imageutil.Resize(scaledWidth, scaledHeight, img, imageutil.ScalerBest())

	// Center crop to exact target dimensions
	if scaled.Bounds().Dx() != targetWidth {
		scaled = imageutil.CropX(scaled, targetWidth, imageutil.AlignCenter)
	}
	if scaled.Bounds().Dy() != targetHeight {
		scaled = imageutil.CropY(scaled, targetHeight, imageutil.AlignCenter)
	}

	return scaled, nil
}

// averageEdgeColor samples pixels along the edges of an image and returns
// the average color. This is used to set a background color that matches
// the image edges, eliminating any visible gaps from sub-pixel rendering.
func averageEdgeColor(img image.Image) pdfcolor.SimpleColor {
	bounds := img.Bounds()
	var r, g, b uint64
	var count uint64

	// Sample all four edges
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		// Top edge
		c := img.At(x, bounds.Min.Y)
		cr, cg, cb, _ := c.RGBA()
		r += uint64(cr)
		g += uint64(cg)
		b += uint64(cb)
		count++

		// Bottom edge
		c = img.At(x, bounds.Max.Y-1)
		cr, cg, cb, _ = c.RGBA()
		r += uint64(cr)
		g += uint64(cg)
		b += uint64(cb)
		count++
	}
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		// Left edge
		c := img.At(bounds.Min.X, y)
		cr, cg, cb, _ := c.RGBA()
		r += uint64(cr)
		g += uint64(cg)
		b += uint64(cb)
		count++

		// Right edge
		c = img.At(bounds.Max.X-1, y)
		cr, cg, cb, _ = c.RGBA()
		r += uint64(cr)
		g += uint64(cg)
		b += uint64(cb)
		count++
	}

	if count == 0 {
		return pdfcolor.SimpleColor{R: 0, G: 0, B: 0}
	}

	// RGBA values are 16-bit (0-65535), convert to 0-1 range
	return pdfcolor.SimpleColor{
		R: float32(r/count) / 65535.0,
		G: float32(g/count) / 65535.0,
		B: float32(b/count) / 65535.0,
	}
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

	// Sample edge color to fill any sub-pixel gaps between image and page
	bgColor := averageEdgeColor(img)

	// Configure import to place image at full page size
	imp := &pdfcpu.Import{
		PageDim:  &types.Dim{Width: pageSize.Width, Height: pageSize.Height},
		PageSize: "",
		UserDim:  true,
		Pos:      types.Full,
		Scale:    1.0,
		InpUnit:  types.POINTS,
		BgColor:  &bgColor,
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
		yFromTop := currentY

		// Build watermark description string
		desc := buildTextWatermarkDesc(style.FontName, int(fontSize), style.Color, style.MarginLeft, yFromTop, style.Align, style.MarginRight)

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

// addTextBlocksOverlay adds styled text blocks onto the PDF.
func addTextBlocksOverlay(pdfBytes []byte, blocks []TextBlock, baseStyle TextStyle) ([]byte, error) {
	if len(blocks) == 0 {
		return pdfBytes, nil
	}

	currentY := baseStyle.MarginTop
	watermarks := make([]*model.Watermark, 0, len(blocks))

	for _, block := range blocks {
		// Apply any additional margin before this block
		if block.MarginTop > 0 {
			currentY += block.MarginTop
		}

		// Resolve effective values (block overrides or base style)
		fontName := baseStyle.FontName
		if block.FontName != "" {
			fontName = block.FontName
		}

		fontSize := baseStyle.FontSize
		if block.FontSize > 0 {
			fontSize = block.FontSize
		}

		color := baseStyle.Color
		if block.Color != "" {
			color = block.Color
		}

		align := baseStyle.Align
		if block.Align != "" {
			align = block.Align
		}

		// Skip empty text blocks (they're just spacers)
		if strings.TrimSpace(block.Text) == "" {
			// Still advance Y position for empty blocks with default line height
			currentY += fontSize * baseStyle.LineHeight
			continue
		}

		// Build watermark description
		desc := buildTextWatermarkDesc(fontName, int(fontSize), color, baseStyle.MarginLeft, currentY, align, baseStyle.MarginRight)

		wm, err := api.TextWatermark(block.Text, desc, true, false, types.POINTS)
		if err != nil {
			return nil, fmt.Errorf("creating text watermark: %w", err)
		}
		watermarks = append(watermarks, wm)

		// Advance Y position for next line
		currentY += fontSize * baseStyle.LineHeight
	}

	if len(watermarks) == 0 {
		return pdfBytes, nil
	}

	// Apply all watermarks to the PDF
	rs := bytes.NewReader(pdfBytes)
	var outBuf bytes.Buffer

	wmMap := make(map[int][]*model.Watermark)
	wmMap[1] = watermarks

	if err := api.AddWatermarksSliceMap(rs, &outBuf, wmMap, nil); err != nil {
		return nil, fmt.Errorf("adding watermarks: %w", err)
	}

	return outBuf.Bytes(), nil
}

// buildTextWatermarkDesc builds the pdfcpu watermark description string.
func buildTextWatermarkDesc(fontName string, fontSize int, color string, marginLeft, yFromTop float64, align TextAlign, marginRight float64) string {
	// Determine position and offset based on alignment
	var pos string
	var offsetX float64

	switch align {
	case TextAlignCenter:
		pos = "tc"
		offsetX = 0 // Centered, no horizontal offset
	case TextAlignRight:
		pos = "tr"
		offsetX = -marginRight // Offset from right edge
	default: // TextAlignLeft or empty
		pos = "tl"
		offsetX = marginLeft
	}

	return fmt.Sprintf(
		"font:%s, points:%d, fillcolor:%s, pos:%s, offset:%.1f -%.1f, scale:1 abs, rotation:0, opacity:1",
		fontName,
		fontSize,
		color,
		pos,
		offsetX,
		yFromTop,
	)
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
		"pos:%s, offset:%.1f %.1f, scale:%.3f %s, rotation:0, opacity:%.2f",
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

// CreateMultiPagePDF creates a PDF with multiple pages, each with its own background, logo, and text.
func CreateMultiPagePDF(pages []Page, pageSize PageSize) ([]byte, error) {
	if len(pages) == 0 {
		return nil, fmt.Errorf("no pages provided")
	}

	// Generate each page as a separate PDF
	var pagePDFs [][]byte
	for i, page := range pages {
		pageBytes, err := createSinglePage(page, pageSize)
		if err != nil {
			return nil, fmt.Errorf("creating page %d: %w", i+1, err)
		}
		pagePDFs = append(pagePDFs, pageBytes)
	}

	// If only one page, return it directly
	if len(pagePDFs) == 1 {
		return pagePDFs[0], nil
	}

	// Merge all pages into a single PDF
	return mergePDFs(pagePDFs)
}

// CreateMultiPagePDFFile creates a multi-page PDF and writes it to outPath.
func CreateMultiPagePDFFile(pages []Page, pageSize PageSize, outPath string) error {
	// Track opened files for cleanup
	var closers []io.Closer
	defer func() {
		for _, c := range closers {
			c.Close()
		}
	}()

	// Open any file-based resources
	for i := range pages {
		if pages[i].BackgroundPath != "" && pages[i].BackgroundReader == nil {
			f, err := os.Open(pages[i].BackgroundPath)
			if err != nil {
				return fmt.Errorf("opening background file for page %d: %w", i+1, err)
			}
			closers = append(closers, f)
			pages[i].BackgroundReader = f
		}
		if pages[i].LogoPath != "" && pages[i].Logo != nil && pages[i].LogoReader == nil {
			f, err := os.Open(pages[i].LogoPath)
			if err != nil {
				return fmt.Errorf("opening logo file for page %d: %w", i+1, err)
			}
			closers = append(closers, f)
			pages[i].LogoReader = f
		}
	}

	pdfBytes, err := CreateMultiPagePDF(pages, pageSize)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outPath, pdfBytes, 0600); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	return nil
}

// createSinglePage creates a single page PDF.
func createSinglePage(page Page, pageSize PageSize) ([]byte, error) {
	var pdfBytes []byte
	var err error

	// Create base page (with or without background)
	if page.BackgroundReader != nil {
		img, err := prepareBackgroundImage(page.BackgroundReader, pageSize)
		if err != nil {
			return nil, fmt.Errorf("preparing background image: %w", err)
		}
		pdfBytes, err = createImagePDF(img, pageSize)
		if err != nil {
			return nil, fmt.Errorf("creating image PDF: %w", err)
		}
	} else {
		// Create blank white page
		pdfBytes, err = createBlankPDF(pageSize)
		if err != nil {
			return nil, fmt.Errorf("creating blank PDF: %w", err)
		}
	}

	// Add logo if provided
	if page.Logo != nil && page.LogoReader != nil {
		pdfBytes, err = addLogoOverlay(pdfBytes, page.LogoReader, *page.Logo)
		if err != nil {
			return nil, fmt.Errorf("adding logo overlay: %w", err)
		}
	}

	// Add text (TextBlocks take precedence over markdown)
	if len(page.TextBlocks) > 0 {
		pdfBytes, err = addTextBlocksOverlay(pdfBytes, page.TextBlocks, page.TextStyle)
		if err != nil {
			return nil, fmt.Errorf("adding text blocks overlay: %w", err)
		}
	} else if strings.TrimSpace(page.MarkdownText) != "" {
		pdfBytes, err = addTextOverlay(pdfBytes, page.MarkdownText, page.TextStyle)
		if err != nil {
			return nil, fmt.Errorf("adding text overlay: %w", err)
		}
	}

	return pdfBytes, nil
}

// createBlankPDF creates a single blank white page PDF.
func createBlankPDF(pageSize PageSize) ([]byte, error) {
	// Create a white image the size of the page
	width := int(pageSize.Width)
	height := int(pageSize.Height)
	whiteImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with white
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			whiteImg.Set(x, y, white)
		}
	}

	return createImagePDF(whiteImg, pageSize)
}

// mergePDFs merges multiple PDF byte slices into a single PDF.
func mergePDFs(pdfs [][]byte) ([]byte, error) {
	if len(pdfs) == 0 {
		return nil, fmt.Errorf("no PDFs to merge")
	}
	if len(pdfs) == 1 {
		return pdfs[0], nil
	}

	// Create readers for all PDFs
	readers := make([]io.ReadSeeker, len(pdfs))
	for i, pdf := range pdfs {
		readers[i] = bytes.NewReader(pdf)
	}

	var outBuf bytes.Buffer
	if err := api.MergeRaw(readers, &outBuf, false, nil); err != nil {
		return nil, fmt.Errorf("merging PDFs: %w", err)
	}

	return outBuf.Bytes(), nil
}
