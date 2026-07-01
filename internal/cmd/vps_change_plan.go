package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var vpsChangePlanCmd = &cobra.Command{
	Use:   "change-plan <billing-id> <plan-id>",
	Short: "Change a VPS's tariff plan in place (resize)",
	Long: `Change a VPS's tariff plan via the "changePlan" method — a resize
without reprovisioning. The billing ID is the service identifier (login_vps_N),
shown in the BILLING_ID column of 'sweb vps list'. The plan ID comes from
'sweb vps config' (or a constructor plan id).

The resize is asynchronous — the command returns once the change is accepted;
the node may reboot while it applies.`,
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
		billingID := args[0]
		planID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("plan-id must be an integer: %w", err)
		}
		if err := c.VPS.ChangePlan(cmd.Context(), billingID, planID); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Changed plan of %s to %d (resize applying asynchronously)\n", billingID, planID)
		return nil
	},
}

func init() {
	vpsCmd.AddCommand(vpsChangePlanCmd)
}
