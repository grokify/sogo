# SoGo

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/sogo/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/sogo/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/sogo/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/sogo/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/sogo/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/sogo/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/sogo
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/sogo
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/sogo
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/sogo
 [viz-svg]: https://img.shields.io/badge/visualization-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fsogo
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/sogo/blob/main/LICENSE

## Overview

SoGo is a collection of Go utilities for common tasks including PDF generation, HTTP handling, database access, text processing, and more. The library wraps various Go packages to provide simpler, more consistent APIs.

Originally extracted from [`github.com/grokify/mogo`](https://github.com/grokify/mogo) to reduce dependency overhead, sogo focuses on providing useful wrappers with minimal transitive dependencies.

## Installation

```bash
go get github.com/grokify/sogo
```

## Packages

| Category | Package | Description |
|----------|---------|-------------|
| **Compression** | `compress/lzfutil` | LZF compression utilities |
| **Database** | `database/kvs` | Key-value store interface |
| | `database/kvs/files` | File-based key-value store |
| | `database/kvs/redis` | Redis client wrapper |
| | `database/kvs/ristretto` | Ristretto cache client |
| **Encoding** | `encoding/jsonutil` | JSON utilities (fastjson, jsonparser) |
| **Flags** | `flag/cobrautil` | Cobra command flag utilities |
| **Logging** | `log/logutil` | Logfmt utilities |
| | `log/slogutil` | slog with ANSI color support |
| **Network** | `net/http/anyhttp` | Unified net/http and fasthttp interface |
| | `net/http/fasthttputil` | Fasthttp utilities |
| | `net/http/httpsimple` | Simple HTTP server |
| | `net/imaputil` | IMAP client utilities |
| | `net/mailutil` | Email utilities |
| | `net/sftputil` | SFTP client utilities |
| | `net/urlutil` | URL parsing and manipulation |
| **PDF** | `pdf` | PDF generation with background images, text overlay, multi-page support |
| **Reflect** | `reflect/reflectutil` | Reflection utilities |
| **Text** | `text/currencyutil` | Currency formatting |
| | `text/markdown/md2html` | Markdown to HTML conversion |
| | `text/markdown/md2pdf` | Markdown to PDF conversion |
| | `text/markdown/remark` | Remark.js presentation utilities |
| | `text/mustacheutil` | Mustache template utilities |
| **Time** | `time/timezone` | Timezone abbreviation mapping |
| **Path** | `path/template` | Path template utilities |

## Usage Examples

### PDF with Background Image and Text Overlay

Create a PDF with a background image, logo, and styled text:

```go
import "github.com/grokify/sogo/pdf"

opts := pdf.BackgroundPDFOptions{
    PageSize: pdf.PageSizeLetter,
    TextStyle: pdf.TextStyle{
        FontName:   "Helvetica",
        FontSize:   14,
        Color:      "1 1 1", // white
        MarginTop:  200,
        MarginLeft: 72,
        LineHeight: 1.5,
    },
    Logo: &pdf.LogoOptions{
        Position: pdf.LogoPositionTopLeft,
        Scale:    0.3,
    },
    LogoPath: "logo.png",
}

err := pdf.CreateBackgroundPDFFile(
    "background.png",
    "# Title\n\nSubtitle text here",
    opts,
    "output.pdf",
)
```

### Multi-Page PDF with Per-Segment Styling

Create multi-page PDFs with different styles per text block:

```go
pages := []pdf.Page{
    {
        BackgroundPath: "cover-bg.png",
        LogoPath:       "logo.png",
        Logo:           &pdf.LogoOptions{Position: pdf.LogoPositionTopLeft, Scale: 0.3},
        TextBlocks: []pdf.TextBlock{
            {Text: "Document Title", FontSize: 28, Color: "1 1 1"},
            {Text: "Subtitle", FontSize: 20, Color: "1 1 1"},
        },
        TextStyle: pdf.TextStyle{MarginTop: 200, MarginLeft: 72},
    },
    {
        // White page (no background)
        TextBlocks: []pdf.TextBlock{
            {Text: "Chapter 1", FontSize: 18, Color: "0 0 0"},
            {Text: "Content goes here...", FontSize: 12, Color: "0 0 0"},
        },
        TextStyle: pdf.TextStyle{MarginTop: 72, MarginLeft: 72},
    },
}

err := pdf.CreateMultiPagePDFFile(pages, pdf.PageSizeLetter, "document.pdf")
```

### Markdown to PDF

Convert markdown files to PDF:

```go
import "github.com/grokify/sogo/text/markdown/md2pdf"

err := md2pdf.MarkdownToPDFFile("input.md", "output.pdf", 0600)
```

## Documentation

Full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/grokify/sogo).

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.

## License

MIT License. See [LICENSE](LICENSE) for details.
