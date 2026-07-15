package cmd

import (
	"fmt"
	"io"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var vpsIPCmd = &cobra.Command{
	Use:   "ip",
	Short: "Manage a VPS's public IPs and reverse-DNS (PTR)",
}

var vpsIPListCmd = &cobra.Command{
	Use:               "list <vps>",
	Short:             "List a VPS's public IPs",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeBillingIDs,
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
		return render(cmd, info.IPs, func(w io.Writer) {
			fmt.Fprintln(w, "IP\tGATEWAY\tNETMASK\tPTR\tPRICE/mo")
			for _, a := range info.IPs {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\n", a.IP, a.Gateway, a.Netmask, a.PTR, float64(a.Price))
			}
		})
	},
}

var vpsIPAddCmd = &cobra.Command{
	Use:               "add <vps> [count]",
	Short:             "Order additional public IPs — bills",
	Long:              "Order [count] (default 1) additional public IPs for a VPS via the \"add\" method. This bills; you are asked to confirm unless --yes.",
	Args:              cobra.RangeArgs(1, 2),
	ValidArgsFunction: completeBillingIDs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		count := 1
		if len(args) == 2 {
			if _, err := fmt.Sscanf(args[1], "%d", &count); err != nil || count < 1 {
				return fmt.Errorf("count must be a positive integer")
			}
		}
		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}
		if yes, _ := cmd.Flags().GetBool("yes"); !yes {
			confirm := false
			if err := huh.NewConfirm().
				Title(fmt.Sprintf("Order %d additional IP(s) for %q? This bills.", count, billingID)).
				Affirmative("Order").Negative("Cancel").Value(&confirm).Run(); err != nil {
				return err
			}
			if !confirm {
				fmt.Fprintln(cmd.OutOrStdout(), "aborted")
				return nil
			}
		}
		if _, err := c.IP.Add(cmd.Context(), billingID, count); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Ordered %d IP(s) for %s (assigned asynchronously — see 'sweb vps ip list %s')\n", count, billingID, args[0])
		return nil
	},
}

var vpsIPRemoveCmd = &cobra.Command{
	Use:               "remove <vps> <ip>",
	Short:             "Release a public IP from a VPS",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeBillingIDs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}
		ip := args[1]
		if yes, _ := cmd.Flags().GetBool("yes"); !yes {
			confirm := false
			if err := huh.NewConfirm().
				Title(fmt.Sprintf("Release IP %s from %q?", ip, billingID)).
				Affirmative("Release").Negative("Cancel").Value(&confirm).Run(); err != nil {
				return err
			}
			if !confirm {
				fmt.Fprintln(cmd.OutOrStdout(), "aborted")
				return nil
			}
		}
		if err := c.IP.Remove(cmd.Context(), billingID, ip); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Released %s from %s\n", ip, billingID)
		return nil
	},
}

var vpsIPMoveCmd = &cobra.Command{
	Use:   "move <ip>",
	Short: "Attach an IP to a VPS (--to) or detach it (--detach)",
	Long:  "Move a public IP via the \"move\" method: --to <vps> attaches it, --detach releases it from its current VPS.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		ip := args[0]
		to, _ := cmd.Flags().GetString("to")
		detach, _ := cmd.Flags().GetBool("detach")
		if (to == "") == !detach {
			return fmt.Errorf("give exactly one of --to <vps> or --detach")
		}
		var billingID string
		if !detach {
			if billingID, err = resolveVPS(cmd.Context(), c, to); err != nil {
				return err
			}
		}
		if err := c.IP.Move(cmd.Context(), ip, billingID); err != nil {
			return err
		}
		if detach {
			fmt.Fprintf(cmd.OutOrStdout(), "Detached %s\n", ip)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Attached %s to %s\n", ip, billingID)
		}
		return nil
	},
}

var vpsIPPtrCmd = &cobra.Command{
	Use:   "ptr",
	Short: "Get or set an IP's reverse-DNS (PTR) record",
}

var vpsIPPtrGetCmd = &cobra.Command{
	Use:   "get <ip>",
	Short: "Show an IP's PTR record",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		ptr, err := c.IP.GetPtr(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), ptr)
		return nil
	},
}

var vpsIPPtrSetCmd = &cobra.Command{
	Use:   "set <ip> <ptr>",
	Short: "Set an IP's PTR record (empty resets to default)",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		ptr := ""
		if len(args) == 2 {
			ptr = args[1]
		}
		if err := c.IP.EditPtr(cmd.Context(), args[0], ptr); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Set PTR of %s to %q\n", args[0], ptr)
		return nil
	},
}

func init() {
	vpsIPAddCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	vpsIPRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	vpsIPMoveCmd.Flags().String("to", "", "VPS (name or billing id) to attach the IP to")
	vpsIPMoveCmd.Flags().Bool("detach", false, "detach the IP from its current VPS")
	_ = vpsIPMoveCmd.RegisterFlagCompletionFunc("to", completeBillingIDs)

	vpsIPPtrCmd.AddCommand(vpsIPPtrGetCmd, vpsIPPtrSetCmd)
	vpsIPCmd.AddCommand(vpsIPListCmd, vpsIPAddCmd, vpsIPRemoveCmd, vpsIPMoveCmd, vpsIPPtrCmd)
	vpsCmd.AddCommand(vpsIPCmd)
}
