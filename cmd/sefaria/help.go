package main

import (
	"github.com/ryanfaerman/go-sefaria/bidi"
	"github.com/spf13/cobra"
)

var (
	helpBidi = &cobra.Command{
		Use:   "bidirectional-text",
		Short: "information about bidirectional text handling",
		Long: `Bidirectional (bidi) text handling is essential when working with languages 
like Hebrew and Arabic that are read from right to left (RTL), especially when 
mixed with left-to-right (LTR) languages like English.

How It Works:
  This tool automatically detects RTL text and applies Unicode bidirectional 
  algorithm markers to ensure proper display. The process involves two steps:

  1. Detection and Marking:
   - Automatically detects Hebrew and Arabic characters
   - Wraps RTL text sequences with Unicode Right-to-Left Mark (RLM) and 
     Left-to-Right Mark (LRM) characters
   - Preserves LTR text unchanged

  2. Character Reordering:
   - Reorders characters within RLM/LRM boundaries for correct display
   - Ensures RTL text appears right-to-left in most terminals and viewers
   - Handles mixed LTR/RTL content appropriately

Supported Character Sets:
  Hebrew (עברית)
  Arabic (العربية)

  These character sets are generally used for Hebrew, Arabic, Farsi, Yiddish, and others.

Configuration:
  --no-bidi          Disable bidirectional text processing
                     Use when piping to programs that handle Unicode bidi correctly


Troubleshooting:
  If Hebrew or Arabic text appears incorrectly:

  - Ensure your terminal supports Unicode bidirectional text
  - Try using a modern terminal emulator (iTerm2, Terminal.app, etc.)
  - Check that your system has proper font support for Hebrew/Arabic
  - Use --no-bidi flag if output is being processed by another tool

Technical Details:

  The implementation uses Unicode bidirectional algorithm markers (RLM/LRM) 
  combined with character reordering to ensure proper display across different 
  environments. This approach provides maximum compatibility while maintaining 
  text integrity.
`,
	}

	helpOutput = &cobra.Command{
		Use:   "output",
		Short: "information about output formats",
		Long: `The following output formats are supported:

JSON Formats:
  json              JSON format (default)
  json-pretty       Pretty-printed JSON with indentation
  jsonl             JSON Lines format (one JSON object per line)
  json-lines        Alias for jsonl
  json-compact      Alias for jsonl

Structured Formats:
  yaml              YAML format
  yml               Alias for yaml
  xml               XML format

Tabular Formats:
  csv               CSV format (for flat data only)

Human-Readable Formats:
  text              Human-readable text format
  pretty            Alias for text
  human             Alias for text

Plain Text Formats:
  plain             Plain text (one item per line)
  shell             Alias for plain

Examples:
  sefaria text get "Genesis 1:1" --output-format=yaml
  sefaria terms completions "torah" --output-format=text
  sefaria index contents --output-format=csv
`,
	}

	helpLogging = &cobra.Command{
		Use:   "logging",
		Short: "information about logging options",
		Long: `The following logging options are supported:

Log Levels:
  debug             Most verbose logging
  info              Informational messages
  warn              Warning messages (default)
  error             Error messages only
  fatal             Fatal errors only
  panic             Panic-level errors only

Log Formats:
  json              JSON format for structured logging
  text              Plain text format
  console           Pretty-printed text with colors (default)
  human             Alias for console
  pretty            Alias for console

Examples:
  sefaria text get "Genesis 1:1" --log-level=debug --log-format=json
  sefaria terms completions "torah" --log-level=info
  sefaria index contents --log-level=error --log-format=text
`,
	}

	helpExamples = &cobra.Command{
		Use:   "examples",
		Short: "common usage examples",
		Long: `Common usage examples:

Basic Text Retrieval:
  sefaria text get "Genesis 1:1"
  sefaria text get "Berakhot 2a"
  sefaria text get "Mishnah Berakhot 1:1"

Get Text with Options:
  sefaria text get "Genesis 1:1" --output-format=yaml
  sefaria text get "Berakhot 2a" --lang=he

Explore Available Content:
  sefaria index contents
  sefaria text languages
  sefaria text versions "Genesis"

Find Related Content:
  sefaria related get "Genesis 1:1"
  sefaria topics recommended "Genesis 1:1" "Berakhot 2a"

Search and Discovery:
  sefaria terms completions "torah"
  sefaria topics all
  sefaria topics random
  sefaria text random

Calendar and Reading:
  sefaria calendar get
  sefaria calendar next-read "Bereshit"

Lexicon Lookup:
  sefaria lexicon get "תורה"
  sefaria lexicon completions "תור"

Output Formatting:
  sefaria text get "Genesis 1:1" --output-format=text
  sefaria index contents --output-format=csv
  sefaria topics all --output-format=yaml

Debugging and Logging:
  sefaria text get "Genesis 1:1" --log-level=debug
  sefaria index contents --log-format=json --log-level=info
`,
	}
)

func init() {
	helpBidi.SetOut(bidi.NewWriter(root.OutOrStdout(), true))
	root.AddCommand(
		helpOutput,
		helpLogging,
		helpExamples,
		helpBidi,
	)
}
