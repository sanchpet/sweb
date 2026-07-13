package cmd

import (
	"errors"

	"github.com/sanchpet/sweb-go-sdk/vh/cron"
	"github.com/spf13/cobra"
)

// errMissingCommand is returned when a cron add/edit omits the required
// --command.
var errMissingCommand = errors.New("--command is required")

// errMissingTask is returned when a cron edit/remove omits the required --task
// job id (the raw crontab line from `cron list`).
var errMissingTask = errors.New("--task is required (the raw crontab line from 'cron list -o json')")

// cronCmd groups the shared-hosting crontab service (endpoint /vh/cron): listing
// the account's cron tasks plus the add/edit/remove lifecycle. It hangs off the
// hosting group, so it inherits that group's profile binding
// (`sweb profile bind hosting <profile>`).
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Manage shared-hosting cron tasks",
}

// addScheduleFlags registers the five cron-position flags plus --command shared
// by `cron add` and `cron edit`. The wildcard convention (`*`) is folded into
// numeric ranges: the API takes concrete positions, so we default to the first
// valid value of each field (minute/hour/day 0/0/1, month/weekday 0 = every).
func addScheduleFlags(cmd *cobra.Command) {
	cmd.Flags().Int("minute", 0, "minute (0..59)")
	cmd.Flags().Int("hour", 0, "hour (0..23)")
	cmd.Flags().Int("day", 1, "day of month (1..31)")
	cmd.Flags().Int("month", 0, "month (0..12; 0 = every month)")
	cmd.Flags().Int("weekday", 0, "day of week (0..7; 0 and 7 = Sunday)")
	cmd.Flags().String("command", "", "command line to run")
}

// scheduleFromFlags builds a cron.Schedule from the flags registered by
// addScheduleFlags. --command is required and its absence is reported.
func scheduleFromFlags(cmd *cobra.Command) (cron.Schedule, error) {
	command, _ := cmd.Flags().GetString("command")
	if command == "" {
		return cron.Schedule{}, errMissingCommand
	}
	return cron.Schedule{
		Minute:  flagInt(cmd, "minute"),
		Hour:    flagInt(cmd, "hour"),
		Day:     flagInt(cmd, "day"),
		Month:   flagInt(cmd, "month"),
		Weekday: flagInt(cmd, "weekday"),
		Command: command,
	}, nil
}

func init() {
	hostingCmd.AddCommand(cronCmd)
}
