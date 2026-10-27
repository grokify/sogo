module github.com/grokify/sogo

go 1.26.1

require (
	github.com/alecthomas/chroma/v2 v2.27.0
	github.com/apex/gateway v1.1.2
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/buaazp/fasthttprouter v0.1.1
	github.com/buger/jsonparser v1.6.1
	github.com/cbroglie/mustache v1.4.2
	github.com/derekstavis/go-qs v0.0.0-20250518184349-717ef4cb7534
	github.com/dgraph-io/ristretto v0.2.0
	github.com/emersion/go-imap v1.2.1
	github.com/fatih/color v1.19.0
	github.com/go-logfmt/logfmt v0.6.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/gomarkdown/markdown v0.0.0-20261011113952-f83953680332
	github.com/gorilla/mux v1.8.1
	github.com/grokify/mogo v0.75.0
	github.com/jessevdk/go-flags v1.6.1
	github.com/johnfercher/maroto/v2 v2.4.3
	github.com/leekchan/accounting v1.0.0
	github.com/oleiade/reflections v1.1.0
	github.com/pdfcpu/pdfcpu v0.15.0
	github.com/phpdave11/gofpdf v1.4.3
	github.com/pkg/sftp v1.13.12
	github.com/shopspring/decimal v1.6.0
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10
	github.com/tkuchiki/go-timezone v0.2.3
	github.com/valyala/fasthttp v1.76.0
	github.com/valyala/fastjson v1.6.10
	github.com/valyala/quicktemplate v1.8.0
	github.com/yuin/goldmark v1.8.6
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20230729083705-37449abec8cc
	github.com/zhuyie/golzf v0.0.0-20161112031142-8387b0307ade
	golang.org/x/crypto v0.58.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/aws/aws-lambda-go v1.55.1 // indirect
	github.com/boombuler/barcode v1.1.0 // indirect
	github.com/caarlos0/env/v11 v11.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2/v2 v2.9.1 // indirect
	github.com/dustin/go-humanize v1.1.0 // indirect
	github.com/emersion/go-sasl v0.0.0-20241020182733-b788ff22d5a6 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hhrutter/tiff v1.0.7 // indirect
	github.com/huandu/xstrings v1.6.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/johnfercher/go-tree v1.1.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.20.3 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.16 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mattn/go-runewidth v0.0.31 // indirect
	github.com/molecule-man/go-brrr v1.3.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/exp v0.0.0-20261014143641-f0251eb4aa3e // indirect
	golang.org/x/image v0.47.0 // indirect
	golang.org/x/net v0.61.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
)

// pdfcpu v0.16+ added context.Context to its api package, which breaks
// maroto/v2 (built against pdfcpu v0.11). Remove once maroto supports it.
exclude (
	github.com/pdfcpu/pdfcpu v0.16.0-rc.1
	github.com/pdfcpu/pdfcpu v0.16.0
	github.com/pdfcpu/pdfcpu v0.16.1
	github.com/pdfcpu/pdfcpu v0.16.2
	github.com/pdfcpu/pdfcpu v0.16.3
	github.com/pdfcpu/pdfcpu v1.0.0-test
)
