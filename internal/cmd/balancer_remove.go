package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var balancerRemoveCmd = &cobra.Command{
	Use:   "remove <billing-id>",
	Short: "Remove (cancel) a load balancer — destructive",
	Long: `Remove a load balancer via the "remove" method. <billing-id> is a
Balancer.BillingID from 'sweb balancer list'.

This is DESTRUCTIVE. You are asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		billingID := args[0]
		if !confirmed(cmd, fmt.Sprintf("Remove balancer %q? This cannot be undone.", billingID), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}

		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Balancer.Remove(cmd.Context(), billingID); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Removed", billingID)
		return nil
	},
}

func init() {
	balancerRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	balancerCmd.AddCommand(balancerRemoveCmd)
}
