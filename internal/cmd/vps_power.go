package cmd

import (
	"context"
	"fmt"
	"time"

	sweb "github.com/sanchpet/sweb-go-sdk"
	"github.com/spf13/cobra"
)

// powerCmd builds a start/stop/reboot command. The three share the same shape —
// resolve <vps>, issue the SDK power call, then either return (the change applies
// asynchronously) or, with --wait, poll until the node settles. action is the SDK
// call; settled is the past-tense verb for the success line.
func powerCmd(use, short string, action func(context.Context, *sweb.Client, string) error, settled string) *cobra.Command {
	c := &cobra.Command{
		Use:               use + " <vps>",
		Short:             short,
		Long:              short + `. <vps> is the VPS name (alias) or its billing ID (login_vps_N), from 'sweb vps list'.` + "\n\nThe change is asynchronous; pass --wait to block until the node settles.",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeBillingIDs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cl, err := client()
			if err != nil {
				return err
			}
			billingID, err := resolveVPS(cmd.Context(), cl, args[0])
			if err != nil {
				return err
			}
			if err := action(cmd.Context(), cl, billingID); err != nil {
				return err
			}

			if wait, _ := cmd.Flags().GetBool("wait"); !wait {
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s (applying asynchronously; --wait to block)\n", settled, billingID)
				return nil
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Minute)
			defer cancel()
			fmt.Fprintf(cmd.ErrOrStderr(), "accepted; waiting for %s to settle…\n", billingID)
			last := ""
			node, err := cl.VPS.WaitForIdle(ctx, billingID, 5*time.Second, func(a string) {
				if a != "" && a != last {
					fmt.Fprintf(cmd.ErrOrStderr(), "  → %s\n", a)
					last = a
				}
			})
			if err != nil {
				return fmt.Errorf("waiting for %s to settle: %w", billingID, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s (now %s)\n", settled, billingID, runState(node))
			return nil
		},
	}
	c.Flags().Bool("wait", false, "block until the operation settles instead of returning immediately")
	return c
}

var (
	vpsStartCmd = powerCmd("start", "Power on a VPS",
		func(ctx context.Context, c *sweb.Client, id string) error { return c.VPS.PowerOn(ctx, id) }, "Powered on")
	vpsStopCmd = powerCmd("stop", "Power off a VPS",
		func(ctx context.Context, c *sweb.Client, id string) error { return c.VPS.PowerOff(ctx, id) }, "Powered off")
	vpsRebootCmd = powerCmd("reboot", "Reboot a VPS",
		func(ctx context.Context, c *sweb.Client, id string) error { return c.VPS.Reboot(ctx, id) }, "Rebooted")
)

// runState renders a VPS's power state as a word.
func runState(v *sweb.VPS) string {
	if v.IsRunning == 1 {
		return "running"
	}
	return "stopped"
}

func init() {
	vpsCmd.AddCommand(vpsStartCmd, vpsStopCmd, vpsRebootCmd)
}
