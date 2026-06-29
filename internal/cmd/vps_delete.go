package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var vpsDeleteCmd = &cobra.Command{
	Use:   "delete <billing-id>",
	Short: "Delete (cancel) a VPS by its billing ID — destructive",
	Long: `Delete a VPS via the "remove" method. The billing ID is the service
identifier (login_vps_N), shown in the BILLING_ID column of 'sweb vps list'.

This is DESTRUCTIVE. You are asked to confirm unless --yes is given.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeBillingIDs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		billingID := args[0]

		if yes, _ := cmd.Flags().GetBool("yes"); !yes {
			confirm := false
			if err := huh.NewConfirm().
				Title(fmt.Sprintf("Delete VPS %q? This cannot be undone.", billingID)).
				Affirmative("Delete").
				Negative("Cancel").
				Value(&confirm).
				Run(); err != nil {
				return err
			}
			if !confirm {
				fmt.Fprintln(cmd.OutOrStdout(), "aborted")
				return nil
			}
		}

		if _, err := c.VPS.Remove(cmd.Context(), billingID); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Deleted", billingID)
		return nil
	},
}

func init() {
	vpsDeleteCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	vpsCmd.AddCommand(vpsDeleteCmd)
}
