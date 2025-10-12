package main

import (
	"fmt"

	"github.com/araddon/dateparse"
	"github.com/ryanfaerman/go-sefaria"
	"github.com/ryanfaerman/go-sefaria/cmd/sefaria/internal/phonetic"
	"github.com/ryanfaerman/go-sefaria/param"
	"github.com/ryanfaerman/go-sefaria/tz"
	"github.com/spf13/cobra"
	"github.com/urfave/sflags/gen/gpflag"
)

var (
	cmdCalendar = &cobra.Command{
		Use:   "calendar",
		Short: "Get calendar information from Sefaria",
		Long: `Calendar commands provide access to Jewish calendar information and reading schedules.

The calendar system in Sefaria provides information about:
• Daily and weekly Torah reading schedules
• Holiday and festival dates
• Custom reading traditions (Ashkenaz, Sefard, Mizrahi)
• Diaspora vs. Israel calendar differences
• Time zone support for accurate local times

These commands are essential for:
• Finding what Torah portion is read on a specific date
• Discovering upcoming reading schedules
• Planning study sessions around the Jewish calendar
• Getting accurate calendar information for different locations

Examples:
  sefaria calendar get
  sefaria calendar get --date="2024-01-15"
  sefaria calendar get --diaspora --tradition=ashkenaz
  sefaria calendar next-read "bereshit"
`,
	}

	optsCalendarGet = &struct {
		Diaspora  bool   `flag:"diaspora" desc:"enable or disable diaspora mode"`
		Date      string `flag:"date" desc:"the date to get the calendar for (default: today)"`
		TimeZone  string `flag:"timezone" desc:"the timezone to use (default: auto-detected)"`
		Tradition string `flag:"tradition" desc:"the Jewish tradition to use (ashkenaz, sefard, or mizrahi)"`
	}{
		TimeZone: tz.Detect(),
		Diaspora: tz.Detect() != "" && tz.Detect() != "Asia/Jerusalem",
	}

	cmdCalendarGet = &cobra.Command{
		Use:   "get",
		Short: "Get the daily or weekly learning schedule for a given date",
		Long: `Get the Torah reading schedule and calendar information for a specific date.

This command retrieves the complete learning schedule for a given date, including:
• Torah portion (parsha) readings
• Haftarah readings
• Holiday information
• Special calendar events
• Reading traditions and customs

The command automatically detects your timezone and location to provide accurate
calendar information. You can override these settings with command-line flags.

Arguments:
  None (uses today's date by default)

Options:
  --date        The date to get calendar information for (default: today)
                Supports various date formats: "2024-01-15", "Jan 15, 2024", "15/01/2024"
  --diaspora    Enable diaspora mode (default: auto-detected based on timezone)
                Diaspora mode affects holiday observance and calendar calculations
  --timezone    Specify timezone (default: auto-detected)
                Examples: "America/New_York", "Europe/London", "Asia/Jerusalem"
  --tradition   Specify reading tradition (ashkenaz, sefard, mizrahi)
                Different traditions may have different reading schedules

Examples:
  # Get today's reading schedule
  sefaria calendar get

  # Get schedule for a specific date
  sefaria calendar get --date="2024-01-15"

  # Get schedule with specific tradition
  sefaria calendar get --tradition=ashkenaz

  # Get schedule for diaspora with custom timezone
  sefaria calendar get --diaspora --timezone="America/New_York"

  # Get schedule for a specific date with tradition
  sefaria calendar get --date="2024-03-15" --tradition=sefard
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := &sefaria.CalendarGetOptions{
				Diaspora: param.BoolInt(optsCalendarGet.Diaspora),
				TimeZone: optsCalendarGet.TimeZone,
				Custom:   optsCalendarGet.Tradition,
			}

			if optsCalendarGet.Date != "" {
				t, err := dateparse.ParseAny(optsCalendarGet.Date)
				if err != nil {
					return fmt.Errorf("cannot parse date: %w", err)
				}

				opts.Year, opts.Month, opts.Day = t.Year(), int(t.Month()), t.Day()
			}

			schedule, err := client.Calendar.Get(cmd.Context(), opts)
			if err != nil {
				return fmt.Errorf("cannot get calendar: %w", err)
			}
			renderer.Render(schedule.Learnings)
			return nil
		},
	}

	cmdCalendarNextRead = &cobra.Command{
		Use:   "next-read [parsha]",
		Short: "Get the next scheduled reading for a given parsha",
		Long: `Find the next scheduled reading date for a specific Torah portion (parsha).

This command helps you discover when a particular Torah portion will next be read
according to the Jewish calendar. It uses phonetic matching to find the correct
parsha even if you don't spell it exactly right.

The command searches through all available Torah portions and finds the closest
match to your input. If multiple matches are found, it will list them for you
to choose from.

Arguments:
  parsha    The name of the Torah portion to find the next reading for
            Examples: "bereshit", "noach", "lech-lecha", "vayera"
            Supports partial matches and phonetic matching

Examples:
  # Find next reading for Bereshit
  sefaria calendar next-read "bereshit"

  # Find next reading for Noah (phonetic matching)
  sefaria calendar next-read "noach"

  # Find next reading for Lech Lecha
  sefaria calendar next-read "lech-lecha"

  # Find next reading with partial match
  sefaria calendar next-read "vay"

Notes:
  • The command uses phonetic matching, so approximate spellings work
  • If multiple matches are found, you'll be shown options to choose from
  • The search is case-insensitive
  • Hebrew names are also supported
`,
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
	if err := gpflag.ParseTo(optsCalendarGet, cmdCalendarGet.Flags()); err != nil {
		panic("cannot activate command flags")
	}
	cmdCalendar.AddCommand(cmdCalendarGet, cmdCalendarNextRead)
	root.AddCommand(cmdCalendar)
}
