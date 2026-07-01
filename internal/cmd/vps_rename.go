package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var vpsRenameCmd = &cobra.Command{
	Use:   "rename <billing-id> <new-name>",
	Short: "Rename a VPS (change its alias) in place",
	Long: `Rename a VPS via the "rename" method. The billing ID is the service
identifier (login_vps_N), shown in the BILLING_ID column of 'sweb vps list'.

This is an in-place label change — it does not reprovision or bill.`,
	Args: cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 { // complete billing IDs for the first arg only
			return completeBillingIDs(cmd, args, toComplete)
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		billingID, newName := args[0], args[1]
		if err := c.VPS.Rename(cmd.Context(), billingID, newName); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Renamed %s to %q\n", billingID, newName)
		return nil
	},
}

func init() {
	vpsCmd.AddCommand(vpsRenameCmd)
}
