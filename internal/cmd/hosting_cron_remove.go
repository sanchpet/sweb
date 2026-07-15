package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cronRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a cron task — destructive",
	Long: `Remove a cron task via the "removeTask" method.

--task is the job id: the raw crontab line of the task as shown by
'sweb hosting cron list -o json' (the "task" field). This is DESTRUCTIVE; you
are asked to confirm unless --yes is given.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		task, _ := cmd.Flags().GetString("task")
		if task == "" {
			return errMissingTask
		}
		if !confirmed(cmd, fmt.Sprintf("Remove cron task %q? This cannot be undone.", task), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Cron.RemoveTask(cmd.Context(), task); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Removed cron task:", task)
		return nil
	},
}

func init() {
	cronRemoveCmd.Flags().String("task", "", "job id: the raw crontab line of the task to remove (from 'cron list -o json')")
	cronRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	cronCmd.AddCommand(cronRemoveCmd)
}
