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
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fsogo
 [loc-svg]: https://tokei.rs/b1/github/grokify/sogo
 [repo-url]: https://github.com/grokify/sogo
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/sogo/blob/main/LICENSE

## Description

SoGo is a collection of Go wrappers to make using various Go libaries easier. Much of the original code for this library was originally in the [`github.com/grokify/mogo`](https://github.com/grokify/mogo) package, however, the code was extracted in the interest of reducing the dependencies associated with that module. This module uses additional dependencies, with the goal of limiting dependencies to specific packages.
