package main

import (
	"fmt"

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
		Use:   "next-read [parsha]",
		Short: "Get the next scheduled reading for a given parsha",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reading, err := client.Calendar.NextRead(cmd.Context(), args[0])
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
