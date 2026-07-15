package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// mailForwardingCmd groups a mailbox's forwarding addresses and the
// delete-after-forwarding toggle.
var mailForwardingCmd = &cobra.Command{
	Use:   "forwarding",
	Short: "Manage a mailbox's forwarding addresses",
}

var mailForwardingListCmd = &cobra.Command{
	Use:   "list <domain> <mbox>",
	Short: "List a mailbox's forwarding addresses",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		addrs, err := c.Mail.ForwardingEmailsList(cmd.Context(), args[0], args[1])
		if err != nil {
			return err
		}
		deleteAfter, err := c.Mail.IsDeletingAfterForwarding(cmd.Context(), args[0], args[1])
		if err != nil {
			return err
		}
		out := map[string]any{"forwarding": addrs, "deleteAfterForwarding": deleteAfter}
		return render(cmd, out, func(w io.Writer) {
			kv(w, "DELETE-AFTER-FORWARDING", onOff(deleteAfter))
			fmt.Fprintln(w, "FORWARDING")
			for _, a := range addrs {
				fmt.Fprintln(w, a)
			}
		})
	},
}

var mailForwardingAddCmd = &cobra.Command{
	Use:   "add <domain> <mbox> <email>",
	Short: "Add a forwarding address to a mailbox",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Mail.AddForwardingEmail(cmd.Context(), args[0], args[1], args[2]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "forwarding %s@%s -> %s added\n", args[1], args[0], args[2])
		return nil
	},
}

var mailForwardingRemoveCmd = &cobra.Command{
	Use:   "remove <domain> <mbox> <email>",
	Short: "Remove a forwarding address from a mailbox",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Remove forwarding %s@%s -> %s?", args[1], args[0], args[2]), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Mail.RemoveForwardingEmail(cmd.Context(), args[0], args[1], args[2]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "forwarding %s@%s -> %s removed\n", args[1], args[0], args[2])
		return nil
	},
}

var mailForwardingDeleteAfterCmd = &cobra.Command{
	Use:   "delete-after <domain> <mbox> <on|off>",
	Short: "Toggle deleting messages from the source mailbox after forwarding",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		on, err := parseOnOff(args[2])
		if err != nil {
			return err
		}
		if err := c.Mail.ChangeDeletingAfterForwarding(cmd.Context(), args[0], args[1], on); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "delete-after-forwarding %s for %s@%s\n", onOff(on), args[1], args[0])
		return nil
	},
}

func init() {
	mailForwardingRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	for _, sub := range []*cobra.Command{
		mailForwardingListCmd, mailForwardingAddCmd,
		mailForwardingRemoveCmd, mailForwardingDeleteAfterCmd,
	} {
		mailForwardingCmd.AddCommand(sub)
	}
	mailCmd.AddCommand(mailForwardingCmd)
}
