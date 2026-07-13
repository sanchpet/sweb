package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cronEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing cron task",
	Long: `Replace an existing cron task via the "editTask" method.

--task identifies the entry to change: it is the raw crontab line (the job id)
of the task as shown by 'sweb hosting cron list -o json' (the "task" field). The
schedule flags plus --command describe the replacement.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		oldTask, _ := cmd.Flags().GetString("task")
		if oldTask == "" {
			return errMissingTask
		}
		sc, err := scheduleFromFlags(cmd)
		if err != nil {
			return err
		}
		if err := c.Cron.EditTask(cmd.Context(), oldTask, sc); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Edited cron task: %s\n", sc.Command)
		return nil
	},
}

func init() {
	addScheduleFlags(cronEditCmd)
	cronEditCmd.Flags().String("task", "", "job id: the raw crontab line of the task to replace (from 'cron list -o json')")
	cronCmd.AddCommand(cronEditCmd)
}
