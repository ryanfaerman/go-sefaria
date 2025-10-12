# go-sefaria

[![Go Report Card](https://goreportcard.com/badge/github.com/ryanfaerman/go-sefaria)](https://goreportcard.com/report/github.com/ryanfaerman/go-sefaria)
[![GoDoc](https://godoc.org/github.com/ryanfaerman/go-sefaria?status.svg)](https://godoc.org/github.com/ryanfaerman/go-sefaria)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org/)

A Go client library for the [Sefaria API](https://www.sefaria.org/api), providing programmatic access to Sefaria's vast library of Jewish texts and resources.

## About Sefaria

Sefaria is a non-profit organization dedicated to building the future of Jewish learning in an open and participatory way. Their platform provides access to thousands of Jewish texts including Tanakh, Talmud, Mishnah, and many other sources in multiple languages.

This client library (and the CLI tool) is not affiliated with Sefaria but is built to facilitate easy access to their API for Go developers.

## Features

- **Comprehensive API Coverage**: Access to all major Sefaria API endpoints
- **Multiple Services**: Text retrieval, index exploration, calendar information, lexicon lookups, topic discovery, and term completions
- **Bidirectional Text Support**: Built-in handling for Hebrew and Arabic text with proper RTL/LTR rendering
- **Robust HTTP Client**: Retry logic, validation, and comprehensive error handling
- **Flexible Configuration**: Customizable endpoints, logging, and HTTP clients
- **CLI Tool**: Command-line interface for interactive use and scripting

## Installation

```bash
go get github.com/ryanfaerman/go-sefaria
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ryanfaerman/go-sefaria"
)

func main() {
    // Create a new client
    client := sefaria.NewClient()

    // Get text from Genesis 1:1
    text, err := client.Text.Get(context.Background(), "Genesis 1:1", nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Reference: %s\n", text.Ref)
    fmt.Printf("Text: %s\n", text.Text[0])
}
```

## Data normalization and alterations

The data from sefaria comes from many sources and has some inconsistencies. This library attempts to normalize the data as much 
as possible, but some alterations are made to the data to improve usability:

- Replace typographic punctuation with ASCII equivalents (e.g., curly quotes to straight quotes)
- Add bidirectional unicode markers to Hebrew and Arabic text for proper rendering (see `bidi` notes below)


### Bidirectional Text Support

The library includes comprehensive support for Hebrew and Arabic text through the `bidi` package:

```go
import "github.com/ryanfaerman/go-sefaria/bidi"

// Create a bidirectional-aware writer
writer := bidi.NewWriter(os.Stdout, true)

// Use with Hebrew text
hebrewText := bidi.String("בראשית ברא אלהים")
fmt.Fprintf(writer, "Hebrew: %s\n", hebrewText)
```

## Configuration

Customize the client with various options:

```go
client := sefaria.NewClient(
    sefaria.WithAPIEndpoint("https://custom-sefaria.org/api"),
    sefaria.WithLogger(logger),
    sefaria.WithHTTPClient(customHTTPClient),
)
```

## CLI Tool

The package includes a command-line tool for interactive use:

```bash
# Install the CLI
go install github.com/ryanfaerman/go-sefaria/cmd/sefaria@latest

# Get text
sefaria text get "Genesis 1:1"

# Search terms
sefaria terms completions "torah"

# Get calendar info
sefaria calendar get

# Multiple output formats
sefaria text get "Genesis 1:1" --output-format=yaml
```

## Requirements

- Go 1.25.0 or later
- Internet connection for API access

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Links

- [Sefaria Website](https://www.sefaria.org)
- [Sefaria API Documentation](https://www.sefaria.org/api)
- [Go Documentation](https://pkg.go.dev/github.com/ryanfaerman/go-sefaria)
