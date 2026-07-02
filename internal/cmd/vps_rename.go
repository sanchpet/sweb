package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var vpsRenameCmd = &cobra.Command{
	Use:   "rename <vps> <new-name>",
	Short: "Rename a VPS (change its alias) in place",
	Long: `Rename a VPS via the "rename" method. <vps> is the VPS name (alias) or its
billing ID (login_vps_N) — both from 'sweb vps list'.

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
		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}
		newName := args[1]
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
