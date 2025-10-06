package main

import (
	"fmt"

	"github.com/ryanfaerman/go-sefaria"
	"github.com/spf13/cobra"
	"github.com/urfave/sflags/gen/gpflag"
)

var (
	cmdTerms = &cobra.Command{
		Use:   "terms",
		Short: "Search and get completions for terms in Sefaria's database",
		Long: `Terms commands allow you to search and explore Sefaria's term database.

The terms system provides autocomplete functionality for finding texts, topics,
and other entities in Sefaria. Use these commands to discover available terms
and get suggestions for partial searches.

Examples:
  sefaria terms completions "torah"
  sefaria terms completions "berakhot"
  sefaria terms completions "תורה"
`,
	}

	cmdTermsGet = &cobra.Command{
		Use:   "get [term]",
		Short: "Get a term by its exact name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			term, err := client.Terms.Get(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("cannot get term: %w", err)
			}
			renderer.Render(term)
			return nil
		},
	}

	cmdTermsNameOptions = &struct {
		*sefaria.TermNameOptions
		Full bool `flag:"full" desc:"display full output"`
	}{}
	cmdTermsName = &cobra.Command{
		Use:   "completions [term]",
		Short: "Get term completions for partial search terms",
		Long: `Get autocomplete suggestions for partial term searches.

This command searches Sefaria's term database and returns suggestions that
match the provided partial term. It's useful for discovering available texts,
topics, and other entities when you're not sure of the exact spelling or
when you want to explore what's available.

The search works with both English and Hebrew terms, and supports partial
matches. By default, only completion titles are returned, but you can use
the --full flag to get complete term information.

Arguments:
  term    The partial term to search for (e.g., "torah", "berakhot", "תורה")

Options:
  --full  Display complete term information instead of just titles

Examples:
  sefaria terms completions "torah"
  sefaria terms completions "berakhot" --full
  sefaria terms completions "תורה"
  sefaria terms completions "genesis"
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			completions, err := client.Terms.Name(cmd.Context(), args[0], cmdTermsNameOptions.TermNameOptions)
			if err != nil {
				return fmt.Errorf("cannot get term completions: %w", err)
			}
			if !cmdTermsNameOptions.Full {
				renderer.Render(completions.CompletionTitles)
			} else {
				renderer.Render(completions)
			}
			return nil
		},
	}
)

func init() {
	if err := gpflag.ParseTo(cmdTermsNameOptions, cmdTermsName.Flags()); err != nil {
		panic("cannot activate command flags")
	}
	cmdTerms.AddCommand(cmdTermsGet, cmdTermsName)

	root.AddCommand(cmdTerms)
}
