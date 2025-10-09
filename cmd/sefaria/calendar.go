package main

import (
	"fmt"

	"github.com/ryanfaerman/go-sefaria"
	"github.com/ryanfaerman/go-sefaria/cmd/sefaria/internal/phonetic"
	"github.com/spf13/cobra"
)

var (
	cmdCalendar = &cobra.Command{
		Use:   "calendar",
		Short: "Get calendar information from Sefaria",
	}

	cmdCalendarGet = &cobra.Command{
		Use:   "get",
		Short: "Get the daily or weekly learning schedule for a given date",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmdCalendarNextRead = &cobra.Command{
		Use:        "next-read [parsha]",
		Short:      "Get the next scheduled reading for a given parsha",
		Args:       cobra.ExactArgs(1),
		SuggestFor: []string{"next"},
		RunE: func(cmd *cobra.Command, args []string) error {
			matches := phonetic.Matches(sefaria.Parshiot, args[0])
			if len(matches) == 0 {
				return fmt.Errorf("no close matches found for parsha: %s", args[0])
			}

			if len(matches) > 1 {
				fmt.Printf("Error: Unknown parsha '%s'\n\n", args[0])
				fmt.Println("Which parsha did you mean?")
				for _, m := range matches {
					fmt.Printf(" - %s\n", m.Candidate)
				}
				return nil
			}

			reading, err := client.Calendar.NextRead(cmd.Context(), matches[0].Candidate)
			if err != nil {
				return fmt.Errorf("cannot get next reading: %w", err)
			}

			renderer.Render(reading)
			return nil
		},
	}
)

func init() {
	cmdCalendar.AddCommand(cmdCalendarGet, cmdCalendarNextRead)
	root.AddCommand(cmdCalendar)
}
