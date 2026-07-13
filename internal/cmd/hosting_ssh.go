package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sshCmd groups the shared-hosting SSH-access toggle (SDK /vh/utils): turning
// SSH access on for a fixed lease period and turning it off. It hangs off the
// hosting parent, so it inherits that group's profile binding.
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Toggle shared-hosting SSH access",
}

// sshPeriods are the lease durations the API documents for `ssh on` (hours).
var sshPeriods = []int{3, 8, 24}

var sshOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Enable SSH access for a fixed period",
	Long: `Enable SSH access to the shared-hosting account via the "sshOn" method.

--period is the lease duration in hours; it must be one of 3, 8, or 24.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		period := flagInt(cmd, "period")
		if !sshValidPeriod(period) {
			return fmt.Errorf("--period must be one of 3, 8, 24 (hours)")
		}
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.SSH.On(cmd.Context(), period); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "SSH access enabled for %d hours\n", period)
		return nil
	},
}

var sshOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Disable SSH access",
	Long:  `Disable SSH access to the shared-hosting account via the "sshOff" method.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.SSH.Off(cmd.Context()); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "SSH access disabled")
		return nil
	},
}

// sshValidPeriod reports whether h is one of the documented lease durations.
func sshValidPeriod(h int) bool {
	for _, p := range sshPeriods {
		if h == p {
			return true
		}
	}
	return false
}

func init() {
	sshOnCmd.Flags().Int("period", 24, "lease duration in hours: 3, 8, or 24")
	_ = sshOnCmd.RegisterFlagCompletionFunc("period", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"3", "8", "24"}, cobra.ShellCompDirectiveNoFileComp
	})

	sshCmd.AddCommand(sshOnCmd, sshOffCmd)
	hostingCmd.AddCommand(sshCmd)
}
