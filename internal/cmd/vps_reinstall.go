package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var vpsReinstallCmd = &cobra.Command{
	Use:   "reinstall <vps>",
	Short: "Reinstall a VPS's OS — destructive",
	Long: `Reinstall a VPS's operating system via the "reinstallOs" method. <vps> is the
VPS name (alias) or its billing ID (login_vps_N), from 'sweb vps list'.

--os is a distributive id (see 'sweb vps config'). This is DESTRUCTIVE: the
system disk is wiped unless --keep-disk is given. You are asked to confirm unless
--yes. The rebuild is asynchronous; pass --wait to block until it settles.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeBillingIDs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		f := cmd.Flags()
		distr, _ := f.GetInt("os")
		if distr == 0 {
			return fmt.Errorf("--os is required (a distributive id — see 'sweb vps config')")
		}
		keepDisk, _ := f.GetBool("keep-disk")

		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}

		if yes, _ := f.GetBool("yes"); !yes {
			title := fmt.Sprintf("Reinstall OS on %q to distributive %d?", billingID, distr)
			if !keepDisk {
				title += " This WIPES the system disk."
			}
			confirm := false
			if err := huh.NewConfirm().Title(title).
				Affirmative("Reinstall").Negative("Cancel").Value(&confirm).Run(); err != nil {
				return err
			}
			if !confirm {
				fmt.Fprintln(cmd.OutOrStdout(), "aborted")
				return nil
			}
		}

		if err := c.VPS.ReinstallOS(cmd.Context(), billingID, distr, keepDisk); err != nil {
			return err
		}

		if wait, _ := f.GetBool("wait"); !wait {
			fmt.Fprintf(cmd.OutOrStdout(), "Reinstalling %s (applying asynchronously; --wait to block)\n", billingID)
			return nil
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Minute)
		defer cancel()
		fmt.Fprintf(cmd.ErrOrStderr(), "accepted; waiting for %s to rebuild…\n", billingID)
		last := ""
		node, err := c.VPS.WaitForIdle(ctx, billingID, 10*time.Second, func(a string) {
			if a != "" && a != last {
				fmt.Fprintf(cmd.ErrOrStderr(), "  → %s\n", a)
				last = a
			}
		})
		if err != nil {
			return fmt.Errorf("waiting for reinstall to settle: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Reinstalled %s (OS %s, now %s)\n", billingID, node.OSDistribution, runState(node))
		return nil
	},
}

func init() {
	f := vpsReinstallCmd.Flags()
	f.Int("os", 0, "distributive id to install (see `sweb vps config`)")
	f.Bool("keep-disk", false, "keep the data disk instead of wiping it (save_disk)")
	f.Bool("wait", false, "block until the rebuild settles instead of returning immediately")
	f.Bool("yes", false, "skip the confirmation prompt")
	vpsCmd.AddCommand(vpsReinstallCmd)
}
