package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var vpsChangePlanCmd = &cobra.Command{
	Use:   "change-plan <billing-id> [plan-id]",
	Short: "Change a VPS's tariff plan in place (resize)",
	Long: `Change a VPS's tariff plan via the "changePlan" method — a resize
without reprovisioning. The billing ID is the service identifier (login_vps_N),
shown in the BILLING_ID column of 'sweb vps list'.

Target the new plan one of two ways (like 'sweb vps create'):

  • a stock plan:      change-plan <billing-id> <plan-id>   (see 'sweb vps config')
  • the configurator:  change-plan <billing-id> --cpu N --ram N --disk N [--category id]
                       (resolves a custom plan; ram and disk are in GB; default
                        category is NVMe)

The resize runs as a sequence of async actions (Modify → ExtIpAdd → …). By
default the command waits for it to settle, printing each phase; pass --async to
return as soon as the change is accepted. Note: shrinking the disk is refused by
the API, so the target's disk must be >= the current one.`,
	Args: cobra.RangeArgs(1, 2),
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

		f := cmd.Flags()
		cpu, _ := f.GetInt("cpu")
		ram, _ := f.GetInt("ram")
		disk, _ := f.GetInt("disk")
		category, _ := f.GetInt("category")
		configurator := cpu != 0 || ram != 0 || disk != 0

		var planID int
		switch {
		case len(args) == 2:
			if configurator {
				return fmt.Errorf("specify either a <plan-id> or --cpu/--ram/--disk, not both")
			}
			if planID, err = strconv.Atoi(args[1]); err != nil {
				return fmt.Errorf("plan-id must be an integer: %w", err)
			}
		case configurator:
			if cpu == 0 || ram == 0 || disk == 0 {
				return fmt.Errorf("the configurator needs all of --cpu, --ram and --disk")
			}
			if category == 0 {
				category = 1 // NVMe ("Быстрые") — see `sweb vps config`
			}
			if planID, err = c.VPS.GetConstructorPlanID(cmd.Context(), cpu, ram, disk, category); err != nil {
				return fmt.Errorf("resolve configurator plan: %w", err)
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "configurator %dcpu/%dGB/%dGB (category %d) → plan %d\n",
				cpu, ram, disk, category, planID)
		default:
			return fmt.Errorf("provide a <plan-id>, or --cpu/--ram/--disk to resolve one")
		}

		if err := c.VPS.ChangePlan(cmd.Context(), billingID, planID); err != nil {
			return err
		}

		if async, _ := f.GetBool("async"); async {
			fmt.Fprintf(cmd.OutOrStdout(), "Changed plan of %s to %d (resize applying asynchronously)\n", billingID, planID)
			return nil
		}

		// Default: wait until the node settles. A resize is a sequence of async
		// actions (Modify → ExtIpAdd → …) with is_running staying 1, so we poll
		// current_action via the SDK and print each phase.
		ctx, cancel := context.WithTimeout(cmd.Context(), 15*time.Minute)
		defer cancel()
		fmt.Fprintf(cmd.ErrOrStderr(), "resize accepted; waiting for %s to settle…\n", billingID)
		last := ""
		node, err := c.VPS.WaitForIdle(ctx, billingID, 10*time.Second, func(action string) {
			if action != "" && action != last {
				fmt.Fprintf(cmd.ErrOrStderr(), "  → %s\n", action)
				last = action
			}
		})
		if err != nil {
			return fmt.Errorf("waiting for resize to settle: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Resized %s → %s (%dcpu / %dMB / %s)\n",
			billingID, node.PlanName, int64(node.CPU), int64(node.RAM), node.Disk)
		return nil
	},
}

func init() {
	f := vpsChangePlanCmd.Flags()
	f.Int("cpu", 0, "configurator: CPU cores (with --ram/--disk, instead of a plan-id)")
	f.Int("ram", 0, "configurator: RAM in GB")
	f.Int("disk", 0, "configurator: disk in GB")
	f.Int("category", 0, "configurator: category id (default 1 = NVMe) — see `sweb vps config`")
	f.Bool("async", false, "return immediately instead of waiting for the resize to settle")
	_ = vpsChangePlanCmd.RegisterFlagCompletionFunc("category", completeCategories)

	vpsCmd.AddCommand(vpsChangePlanCmd)
}
