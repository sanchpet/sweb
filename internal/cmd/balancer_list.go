package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var balancerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's load balancers",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		bals, err := c.Balancer.List(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, bals, func(w io.Writer) {
			fmt.Fprintln(w, "BILLING_ID\tNAME\tTYPE\tPLAN\tIP\tPRICE/mo\tACTIVE\tACTION")
			for _, b := range bals {
				action := b.CurrentAction
				if action == "" {
					action = "-"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\t%t\t%s\n",
					b.BillingID, b.Name, b.Type, b.PlanName, b.IPBalancer,
					int64(b.Price), b.Active, action)
			}
		})
	},
}

func init() {
	balancerCmd.AddCommand(balancerListCmd)
}
