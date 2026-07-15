package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkActivateCmd = &cobra.Command{
	Use:   "activate <id>",
	Short: "Enable a monitoring check",
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
		if err := c.MonitoringChecks.Activate(cmd.Context(), id); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Activated check", id)
		return nil
	},
}

var checkDeactivateCmd = &cobra.Command{
	Use:   "deactivate <id>",
	Short: "Disable a monitoring check",
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
		if err := c.MonitoringChecks.Deactivate(cmd.Context(), id); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Deactivated check", id)
		return nil
	},
}

var checkRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Remove a monitoring check — destructive",
	Long: `Delete a monitoring check (method "remove"). This is DESTRUCTIVE; you
are asked to confirm unless --yes.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := checkID(args[0])
		if err != nil {
			return err
		}
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Remove check %d? This cannot be undone.", id), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.MonitoringChecks.Remove(cmd.Context(), id); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Removed check", id)
		return nil
	},
}

func init() {
	checkRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	checkCmd.AddCommand(checkActivateCmd, checkDeactivateCmd, checkRemoveCmd)
}
