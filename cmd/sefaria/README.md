# Sefaria CLI

A command-line interface for accessing Sefaria's vast library of Jewish texts and resources.

## About

Sefaria is a non-profit organization dedicated to building the future of Jewish learning in an open and participatory way. This CLI tool provides programmatic access to Sefaria's API, allowing you to:

- Retrieve texts from Tanakh, Talmud, Mishnah, and thousands of other sources
- Search and explore Sefaria's comprehensive index of Jewish texts
- Access translations in multiple languages
- Find related content and topics
- Get calendar and reading schedule information
- Look up terms and get autocomplete suggestions

This CLI tool is not affiliated with Sefaria but is built to facilitate easy access to their API for developers and researchers.

## Installation

### From Source

```bash
go install github.com/ryanfaerman/go-sefaria/cmd/sefaria@latest
```

### Building from Source

```bash
git clone https://github.com/ryanfaerman/go-sefaria.git
cd go-sefaria/cmd/sefaria
go build -o sefaria
```

## Quick Start

```bash
# Search for terms
sefaria terms completions "torah"

# Get detailed term information
sefaria terms completions "berakhot" --full

# Use different output formats
sefaria terms completions "תורה" --output-format=yaml

# Enable debug logging
sefaria terms completions "genesis" --log-level=debug
```

## Global Options

The CLI supports several global options that apply to all commands:

### Output Formats

- `--output-format` or `-f`: Choose output format (default: `json`)
  - `json`: JSON format (default)
  - `json-pretty`: Pretty-printed JSON with indentation
  - `jsonl`: JSON Lines format (one JSON object per line)
  - `yaml`/`yml`: YAML format
  - `xml`: XML format
  - `csv`: CSV format (for flat data only)
  - `text`/`pretty`/`human`: Human-readable text format
  - `plain`/`shell`: Plain text (one item per line)

### Logging Options

- `--log-level`: Set logging verbosity (default: `warn`)
  - `debug`: Most verbose logging
  - `info`: Informational messages
  - `warn`: Warning messages (default)
  - `error`: Error messages only
  - `fatal`: Fatal errors only
  - `panic`: Panic-level errors only

- `--log-format`: Set log format (default: `console`)
  - `json`: JSON format for structured logging
  - `text`: Plain text format
  - `console`/`human`/`pretty`: Pretty-printed text with colors

### Bidirectional Text Support

- `--no-bidi`: Disable bidirectional text processing

The CLI automatically handles Hebrew and Arabic text with proper RTL/LTR rendering. Use `--no-bidi` when piping to programs that handle Unicode bidi correctly.

## Commands

### Terms

Search and explore Sefaria's term database for autocomplete functionality.

#### `sefaria terms completions [term]`

Get autocomplete suggestions for partial term searches.

**Arguments:**
- `term`: The partial term to search for (e.g., "torah", "berakhot", "תורה")

**Options:**
- `--full`: Display complete term information instead of just titles

**Examples:**
```bash
# Basic search
sefaria terms completions "torah"

# Get full information
sefaria terms completions "berakhot" --full

# Search in Hebrew
sefaria terms completions "תורה"

# Different output formats
sefaria terms completions "genesis" --output-format=yaml
sefaria terms completions "mishnah" --output-format=text
```

#### `sefaria terms get [term]`

Get a term by its exact name.

**Arguments:**
- `term`: The exact term name to retrieve

**Examples:**
```bash
sefaria terms get "Torah"
sefaria terms get "Berakhot"
```

## Help Topics

The CLI includes several help topics for detailed information:

### `sefaria help bidirectional-text`

Information about bidirectional text handling for Hebrew and Arabic text.

### `sefaria help output`

Detailed information about all supported output formats.

### `sefaria help logging`

Information about logging options and formats.

### `sefaria help examples`

Common usage examples and patterns.

## Examples

### Basic Usage

```bash
# Search for Torah-related terms
sefaria terms completions "torah"

# Get detailed information about Berakhot
sefaria terms completions "berakhot" --full

# Search in Hebrew
sefaria terms completions "תורה"
```

### Output Formatting

```bash
# YAML output
sefaria terms completions "genesis" --output-format=yaml

# Human-readable text
sefaria terms completions "mishnah" --output-format=text

# CSV for data processing
sefaria terms completions "talmud" --output-format=csv
```

### Debugging and Logging

```bash
# Enable debug logging
sefaria terms completions "torah" --log-level=debug

# JSON logging for structured output
sefaria terms completions "berakhot" --log-format=json --log-level=info

# Disable bidirectional text processing
sefaria terms completions "תורה" --no-bidi
```

### Scripting and Automation

```bash
# Pipe to other tools
sefaria terms completions "torah" --output-format=plain | grep -i "genesis"

# Save to file
sefaria terms completions "berakhot" --full --output-format=yaml > berakhot.yaml

# Process with jq
sefaria terms completions "mishnah" | jq '.[] | select(.title | contains("Berakhot"))'
```

## Planned Features

The CLI is actively developed and the following commands are planned:

- **Text Commands**: Retrieve specific texts and translations
- **Index Commands**: Explore Sefaria's text index
- **Calendar Commands**: Access Jewish calendar and reading schedules
- **Lexicon Commands**: Look up Hebrew/Aramaic terms
- **Topics Commands**: Discover and explore topics
- **Related Commands**: Find related content

## Troubleshooting

### Hebrew/Arabic Text Display Issues

If Hebrew or Arabic text appears incorrectly:

1. Ensure your terminal supports Unicode bidirectional text
2. Try using a modern terminal emulator (iTerm2, Terminal.app, etc.)
3. Check that your system has proper font support for Hebrew/Arabic
4. Use `--no-bidi` flag if output is being processed by another tool

### Network Issues

The CLI requires internet access to connect to Sefaria's API. If you encounter network issues:

1. Check your internet connection
2. Verify that `https://www.sefaria.org/api` is accessible
3. Use `--log-level=debug` to see detailed request information

### Performance

For large result sets:

1. Use `--output-format=plain` for faster processing
2. Pipe output to `head` or `tail` to limit results
3. Use `--log-level=error` to reduce logging overhead

## Contributing

Contributions are welcome! The CLI is part of the larger go-sefaria project. Please see the main project README for contribution guidelines.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Links

- [Sefaria Website](https://www.sefaria.org)
- [Sefaria API Documentation](https://www.sefaria.org/api)
- [Main go-sefaria Library](https://github.com/ryanfaerman/go-sefaria)
- [Go Documentation](https://pkg.go.dev/github.com/ryanfaerman/go-sefaria)
