package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var vpsLocalIPCmd = &cobra.Command{
	Use:   "local-ip",
	Short: "Manage a VPS's private (local) network attachment",
	Long: `Attach or detach a VPS to the account private (local) network, or show its
local IP (methods addLocal / removeLocal / index on /vps/ip).

The private network lets VPS talk over a private L2 (10.0.0.0/x) without going
over the public internet — e.g. cluster/etcd traffic. Attaching an existing VPS
is declarative (no re-create). The guest OS still needs the private interface
configured (netplan/ifcfg) with the assigned local IP — see 'show'.`,
}

var vpsLocalIPShowCmd = &cobra.Command{
	Use:               "show <vps>",
	Short:             "Show a VPS's local (private) IP",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeBillingIDArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}
		info, err := c.IP.Info(cmd.Context(), billingID)
		if err != nil {
			return err
		}
		if len(info.LocalIP) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "%s is not attached to the local network\n", args[0])
			return nil
		}
		for _, l := range info.LocalIP {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\tmask %s\tmac %s\n", l.IP, l.Mask, l.MAC)
		}
		return nil
	},
}

var vpsLocalIPAddCmd = &cobra.Command{
	Use:   "add <vps>",
	Short: "Attach a VPS to the private (local) network",
	Long: `Attach an existing VPS to the account private network via "addLocal".
<vps> is the VPS name (alias) or its billing ID. SpaceWeb assigns the local IP;
by default the command waits for it and prints it (pass --async to return as soon
as the attach is accepted). Configure the private NIC in the guest OS afterwards.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeBillingIDArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}
		if err := c.IP.AddLocal(cmd.Context(), billingID); err != nil {
			return err
		}
		if async, _ := cmd.Flags().GetBool("async"); async {
			fmt.Fprintf(cmd.OutOrStdout(), "Attach of %s to the local network started\n", billingID)
			return nil
		}
		ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Minute)
		defer cancel()
		fmt.Fprintf(cmd.ErrOrStderr(), "attach accepted; waiting for the local IP on %s…\n", billingID)
		lip, err := c.IP.WaitForLocalIP(ctx, billingID, 5*time.Second)
		if err != nil {
			return fmt.Errorf("waiting for local IP: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Attached %s to the local network: %s (mask %s)\n", billingID, lip.IP, lip.Mask)
		return nil
	},
}

var vpsLocalIPRemoveCmd = &cobra.Command{
	Use:               "remove <vps>",
	Short:             "Detach a VPS from the private (local) network",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeBillingIDArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}
		if err := c.IP.RemoveLocal(cmd.Context(), billingID); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Detached %s from the local network\n", billingID)
		return nil
	},
}

// completeBillingIDArg completes billing IDs for the first positional arg only.
func completeBillingIDArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completeBillingIDs(cmd, args, toComplete)
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	vpsLocalIPAddCmd.Flags().Bool("async", false, "return immediately instead of waiting for the local IP")
	vpsLocalIPCmd.AddCommand(vpsLocalIPShowCmd, vpsLocalIPAddCmd, vpsLocalIPRemoveCmd)
	vpsCmd.AddCommand(vpsLocalIPCmd)
}
