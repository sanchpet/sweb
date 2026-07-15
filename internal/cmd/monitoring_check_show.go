package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var checkShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a check's full configuration and attached contacts",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := checkID(args[0])
		if err != nil {
			return err
		}
		c, err := client()
		if err != nil {
			return err
		}
		info, err := c.MonitoringChecks.GetFullCheckInfo(cmd.Context(), id)
		if err != nil {
			return err
		}
		return render(cmd, info, func(w io.Writer) {
			fmt.Fprintf(w, "ID\t%d\n", int64(info.ID))
			fmt.Fprintf(w, "NAME\t%s\n", info.Name)
			fmt.Fprintf(w, "TYPE\t%d\n", int64(info.Type))
			fmt.Fprintf(w, "STATUS\t%s\n", onOffStatus(info.Status))
			for _, s := range info.Settings {
				fmt.Fprintf(w, "SETTING/%s\t%s\n", s.Type, s.Value)
			}
			for _, ct := range info.Contacts {
				fmt.Fprintf(w, "CONTACT\t%d\t%s\t%s\t%s\tverified=%t\n",
					int64(ct.ID), ct.Type, ct.Name, ct.Value, ct.Verified)
			}
		})
	},
}

func init() {
	checkCmd.AddCommand(checkShowCmd)
}
