package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// mailDeliveryCmd groups a mailbox's mailing (delivery) list: the quota summary,
// the address list, and add/remove.
var mailDeliveryCmd = &cobra.Command{
	Use:   "delivery",
	Short: "Manage a mailbox's mailing (delivery) list",
}

var mailDeliveryInfoCmd = &cobra.Command{
	Use:   "info <domain> <mbox>",
	Short: "Show mailing-list and delivery-address quota usage for a mailbox",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		info, err := c.Mail.DeliveryInfo(cmd.Context(), args[0], args[1])
		if err != nil {
			return err
		}
		return render(cmd, info, func(w io.Writer) {
			fmt.Fprintln(w, "RESOURCE\tCURRENT\tMAX")
			fmt.Fprintf(w, "groups\t%d\t%d\n", int64(info.Groups.Current), int64(info.Groups.Max))
			fmt.Fprintf(w, "addresses\t%d\t%d\n", int64(info.Addresses.Current), int64(info.Addresses.Max))
		})
	},
}

var mailDeliveryListCmd = &cobra.Command{
	Use:   "list <domain> <mbox>",
	Short: "List a mailbox's mailing (delivery) addresses",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		res, err := c.Mail.DeliveryAddressesList(cmd.Context(), args[0], args[1], listPageOptions(cmd))
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintln(w, "ADDRESS")
			for _, a := range res.List {
				fmt.Fprintln(w, a)
			}
		})
	},
}

var mailDeliveryAddCmd = &cobra.Command{
	Use:   "add <domain> <mbox> <email>",
	Short: "Add an address to a mailbox's mailing list",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Mail.AddDeliveryAddress(cmd.Context(), args[0], args[1], args[2]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s added to mailing list of %s@%s\n", args[2], args[1], args[0])
		return nil
	},
}

var mailDeliveryRemoveCmd = &cobra.Command{
	Use:   "remove <domain> <mbox> <email>",
	Short: "Remove an address from a mailbox's mailing list",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Mail.DropDeliveryAddress(cmd.Context(), args[0], args[1], args[2]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s removed from mailing list of %s@%s\n", args[2], args[1], args[0])
		return nil
	},
}

// mailCollectorCmd groups the domain-level mail collector (get/set/remove and
// the cross-domain confirmation step).
var mailCollectorCmd = &cobra.Command{
	Use:   "collector",
	Short: "Manage a domain's mail collector",
}

var mailCollectorGetCmd = &cobra.Command{
	Use:   "get <domain>",
	Short: "Show a domain's mail-collector address",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		addr, err := c.Mail.MailsCollector(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return render(cmd, map[string]string{"collector": addr}, func(w io.Writer) {
			kv(w, "COLLECTOR", emptyDash(addr))
		})
	},
}

var mailCollectorSetCmd = &cobra.Command{
	Use:   "set <domain> <email>",
	Short: "Set a domain's mail-collector address",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		res, err := c.Mail.ChangeMailsCollector(cmd.Context(), args[0], args[1])
		if err != nil {
			return err
		}
		if res == 2 {
			fmt.Fprintf(cmd.OutOrStdout(),
				"collector for %s set to %s — the target is on a domain you do not own; confirm with the emailed token via `sweb hosting mail collector confirm %s <token>`\n",
				args[0], args[1], args[0])
			return nil
		}
		fmt.Fprintf(cmd.OutOrStdout(), "collector for %s set to %s\n", args[0], args[1])
		return nil
	},
}

var mailCollectorRemoveCmd = &cobra.Command{
	Use:   "remove <domain>",
	Short: "Remove a domain's mail-collector address",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Remove the mail collector for %s?", args[0]), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Mail.RemoveMailsCollector(cmd.Context(), args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "collector removed for %s\n", args[0])
		return nil
	},
}

// mailCollectorConfirmCmd completes a cross-domain collector setup: when
// `collector set` targets a domain not on the account, the API emails a token
// that this command submits. Kept as a leaf of the collector group rather than a
// top-level command since it is only reachable from that flow.
var mailCollectorConfirmCmd = &cobra.Command{
	Use:   "confirm <domain> <token>",
	Short: "Confirm a cross-domain mail collector with the emailed token",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Mail.ConfirmMailsCollectorEmail(cmd.Context(), args[0], args[1]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "collector confirmed for %s\n", args[0])
		return nil
	},
}

func init() {
	addPageFlags(mailDeliveryListCmd)
	mailDeliveryCmd.AddCommand(mailDeliveryInfoCmd, mailDeliveryListCmd, mailDeliveryAddCmd, mailDeliveryRemoveCmd)
	mailCmd.AddCommand(mailDeliveryCmd)

	mailCollectorRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	mailCollectorCmd.AddCommand(mailCollectorGetCmd, mailCollectorSetCmd, mailCollectorRemoveCmd, mailCollectorConfirmCmd)
	mailCmd.AddCommand(mailCollectorCmd)
}
