package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var dbaasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's managed-database clusters",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		idx, err := c.DBaaS.List(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, idx, func(w io.Writer) {
			fmt.Fprintln(w, "BILLING_ID\tNAME\tENGINE\tPLAN\tIP\tPRICE/mo\tSTATUS\tACTION")
			for _, in := range idx.Instances {
				name := in.DisplayName
				if name == "" {
					name = in.Name
				}
				action := in.CurrentAction
				if action == "" {
					action = "-"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%.2f\t%s\t%s\n",
					in.BillingID, name, in.Engine, in.Plan.Name, in.IP,
					float64(in.Price), in.Status, action)
			}
		})
	},
}

func init() {
	dbaasCmd.AddCommand(dbaasListCmd)
}
