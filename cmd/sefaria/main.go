package main

import (
	"fmt"
	"log/slog"

	"github.com/ryanfaerman/go-sefaria"
	"github.com/ryanfaerman/go-sefaria/bidi"
	"github.com/ryanfaerman/go-sefaria/cmd/sefaria/internal/render"
	"github.com/ryanfaerman/go-sefaria/cmd/sefaria/internal/version"
	"github.com/spf13/cobra"
	"github.com/urfave/sflags/gen/gpflag"
)

type Config struct {
	LogLevel  string `flag:"log-level" desc:"(debug, info, warn, error, fatal, panic)"`
	LogFormat string `flag:"log-format" desc:"log format (text, json, console)"`

	OutputFormat string `flag:"output-format f" desc:"output format (text, json, yaml, xml, csv)"`

	DisableBidi bool `flag:"no-bidi" desc:"disable bidi text handling"`
}

var (
	client *sefaria.Client
	config = &Config{
		LogLevel:     "warn",
		LogFormat:    "console",
		OutputFormat: "json",
	}
	renderer render.Renderer
	logger   *slog.Logger

	root = &cobra.Command{
		Use:     "sefaria",
		Version: version.String(),
		Short:   "A command line tool for interacting with the Sefaria API",
		Long: `Sefaria CLI is a command-line interface for accessing Sefaria's vast library
of Jewish texts and resources.

Sefaria is a non-profit organization dedicated to building the future of Jewish
learning in an open and participatory way. This CLI tool provides programmatic
access to Sefaria's API, allowing you to:

• Retrieve texts from Tanakh, Talmud, Mishnah, and thousands of other sources
• Search and explore Sefaria's comprehensive index of Jewish texts
• Access translations in multiple languages
• Find related content and topics
• Get calendar and reading schedule information
• Look up terms and get autocomplete suggestions

The tool supports multiple output formats (JSON, YAML, XML, CSV, text) and
includes comprehensive logging options for debugging and monitoring.

Examples:
  sefaria terms completions "torah"
  sefaria --output-format=yaml terms completions "berakhot" --full
  sefaria --log-level=debug terms completions "תורה"

For more information about specific commands, use:
  sefaria help <command>

For help topics, use:
  sefaria help <topic>
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			logger, err = NewLogger(config.LogFormat, config.LogLevel, "")
			if err != nil {
				return fmt.Errorf("cannot create logger: %w", err)
			}

			w := cmd.OutOrStdout()
			if !config.DisableBidi {
				w = bidi.NewWriter(w, true)
			}
			renderer = render.NewRenderer(config.OutputFormat, w)

			client = sefaria.NewClient(sefaria.WithLogger(logger))

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return renderer.Flush()
		},
	}
)

func init() {
	if err := gpflag.ParseTo(config, root.PersistentFlags()); err != nil {
		panic("cannot activate command flags")
	}
}

func main() {
	root.Execute()
}
