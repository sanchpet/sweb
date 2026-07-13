package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cronAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a cron task",
	Long: `Add a cron task via the "addTask" method.

The schedule is the five cron positions (minute, hour, day, month, weekday) plus
the command to run. --command is required; the schedule flags default to a
run at 00:00 on day 1 of every month.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		sc, err := scheduleFromFlags(cmd)
		if err != nil {
			return err
		}
		if err := c.Cron.AddTask(cmd.Context(), sc); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Added cron task: %s\n", sc.Command)
		return nil
	},
}

func init() {
	addScheduleFlags(cronAddCmd)
	cronCmd.AddCommand(cronAddCmd)
}
