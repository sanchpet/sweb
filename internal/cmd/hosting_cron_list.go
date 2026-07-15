package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var cronListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's cron tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		tasks, err := c.Cron.GetTasks(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, tasks, func(w io.Writer) {
			fmt.Fprintln(w, "MIN\tHOUR\tDAY\tMON\tWDAY\tCOMMAND")
			for _, t := range tasks {
				fmt.Fprintf(w, "%d\t%d\t%d\t%d\t%d\t%s\n",
					int64(t.Minute), int64(t.Hour), int64(t.Day),
					int64(t.Month), int64(t.Weekday), t.Command)
			}
		})
	},
}

func init() {
	cronCmd.AddCommand(cronListCmd)
}
