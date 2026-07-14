package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/monitoring/checks"
	"github.com/spf13/cobra"
)

var checkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's monitoring checks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		res, err := c.MonitoringChecks.Index(cmd.Context(), &checks.ListOptions{
			Page:    flagInt(cmd, "page"),
			PerPage: flagInt(cmd, "per-page"),
		})
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tNAME\tTYPE\tSTATUS\tDISABLED")
			for _, ch := range res.List {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%t\n",
					ch.ID, ch.Name, ch.Type, onOffStatus(ch.Status), ch.Disabled)
			}
		})
	},
}

// onOffStatus renders a check's active flag as active/disabled.
func onOffStatus(active bool) string {
	if active {
		return "active"
	}
	return "disabled"
}

func init() {
	checkListCmd.Flags().Int("page", 0, "page of results (1-based; 0 lets the API default)")
	checkListCmd.Flags().Int("per-page", 0, "rows per page (0 lets the API default)")
	checkCmd.AddCommand(checkListCmd)
}
