package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var vpsLogsCmd = &cobra.Command{
	Use:   "logs <vps>",
	Short: "Show a VPS's operation log",
	Long: `Show the record of lifecycle actions (create, reinstall, resize, …) run
against a VPS, via the "logs" method. <vps> is the VPS name (alias) or its
billing ID (login_vps_N). Read-only; use -o json for the raw entries.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeBillingIDs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}
		logs, err := c.VPS.Logs(cmd.Context(), billingID)
		if err != nil {
			return err
		}
		return render(cmd, logs, func(w io.Writer) {
			fmt.Fprintln(w, "TYPE\tSTATUS\tSTARTED\tENDED")
			for _, e := range logs {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.Type, e.Status, e.StartedAt, e.EndedAt)
			}
		})
	},
}

func init() {
	vpsCmd.AddCommand(vpsLogsCmd)
}
