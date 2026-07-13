package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// diskUsageCmd groups the shared-hosting disk-usage services under
// `sweb hosting disk-usage` (SDK /vh/utils/diskUsage): the per-backend quota
// breakdown and scan-task state, triggering a new scan, and the over-quota
// notification email. It hangs off the hosting parent, so it inherits that
// group's profile binding.
var diskUsageCmd = &cobra.Command{
	Use:   "disk-usage",
	Short: "Shared-hosting disk-usage reports",
}

var diskUsageReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Show the per-backend disk-usage (quota) breakdown",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		list, err := c.DiskUsage.List(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "TARIFF(MB)\tREAL(MB)\tDB(MB)\tMAIL(MB)\tFILES(MB)\tFILES")
			for _, u := range list {
				fmt.Fprintf(w, "%.1f\t%.1f\t%.1f\t%.1f\t%.1f\t%d\n",
					float64(u.TariffQuota), float64(u.RealQuota), float64(u.DBQuota),
					float64(u.MailQuota), float64(u.FilesQuota), int64(u.FilesNum))
			}
		})
	},
}

var diskUsageTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Show the disk-usage scan-task state",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		info, err := c.DiskUsage.TasksInfo(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, info, func(w io.Writer) {
			fmt.Fprintf(w, "ACTIVE TASKS\t%d\n", int64(info.ActiveTasksCount))
			fmt.Fprintf(w, "LAST DONE\t%s\n", info.LastDoneTaskDate)
		})
	},
}

var diskUsageScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Queue a new disk-usage scan",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.DiskUsage.StartTask(cmd.Context()); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Scan queued")
		return nil
	},
}

var diskUsageEmailCmd = &cobra.Command{
	Use:   "email",
	Short: "Show the over-quota notification email",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		email, err := c.DiskUsage.Email(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, map[string]string{"email": email}, func(w io.Writer) {
			fmt.Fprintf(w, "EMAIL\t%s\n", email)
		})
	},
}

var diskUsageSetEmailCmd = &cobra.Command{
	Use:   "set-email <email>",
	Short: "Set the over-quota notification email",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.DiskUsage.ChangeEmail(cmd.Context(), args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Notification email set to %s\n", args[0])
		return nil
	},
}

func init() {
	diskUsageCmd.AddCommand(
		diskUsageReportCmd,
		diskUsageTasksCmd,
		diskUsageScanCmd,
		diskUsageEmailCmd,
		diskUsageSetEmailCmd,
	)
	hostingCmd.AddCommand(diskUsageCmd)
}
