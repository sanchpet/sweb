package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/spf13/cobra"
)

var monitoringPlansCmd = &cobra.Command{
	Use:   "plans",
	Short: "List the available monitoring tariff plans",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		plans, err := c.Monitoring.Plans(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, plans, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tNAME\tCHECKS\tSMS\tPRICE")
			for _, p := range plans {
				fmt.Fprintf(w, "%d\t%s\t%d\t%d\t%.2f\n",
					int64(p.ID), p.Name, int64(p.Checks), int64(p.SMS), float64(p.Price))
			}
		})
	},
}

// planID parses the shared <plan-id> positional argument.
func planID(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return 0, fmt.Errorf("plan id must be an integer: %q", arg)
	}
	return id, nil
}

var monitoringEnableCmd = &cobra.Command{
	Use:   "enable <plan-id>",
	Short: "Subscribe to a monitoring plan — this BILLS",
	Long: `Subscribe the account to the monitoring tariff with the given plan id
(method "enable"). This BILLS; you are asked to confirm unless --yes.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := planID(args[0])
		if err != nil {
			return err
		}
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Enable monitoring plan %d? This will bill your account.", id), "Enable") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Monitoring.Enable(cmd.Context(), id); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Enabled monitoring plan", id)
		return nil
	},
}

var monitoringDisableCmd = &cobra.Command{
	Use:   "disable <plan-id>",
	Short: "Cancel the monitoring subscription",
	Long: `Cancel the monitoring tariff subscription (method "disable"). You are
asked to confirm unless --yes.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := planID(args[0])
		if err != nil {
			return err
		}
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Disable monitoring plan %d?", id), "Disable") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Monitoring.Disable(cmd.Context(), id); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Disabled monitoring plan", id)
		return nil
	},
}

var monitoringChangeCmd = &cobra.Command{
	Use:   "change <plan-id>",
	Short: "Switch the monitoring subscription to another plan — may BILL",
	Long: `Switch the monitoring subscription to a different tariff plan (method
"change"). This may BILL; you are asked to confirm unless --yes.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := planID(args[0])
		if err != nil {
			return err
		}
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Change monitoring subscription to plan %d? This may bill your account.", id), "Change") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Monitoring.Change(cmd.Context(), id); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Changed monitoring subscription to plan", id)
		return nil
	},
}

func init() {
	for _, c := range []*cobra.Command{monitoringEnableCmd, monitoringDisableCmd, monitoringChangeCmd} {
		c.Flags().Bool("yes", false, "skip the confirmation prompt")
	}
	monitoringCmd.AddCommand(monitoringPlansCmd, monitoringEnableCmd, monitoringDisableCmd, monitoringChangeCmd)
}
