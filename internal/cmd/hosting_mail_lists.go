package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/vh/mail"
	"github.com/spf13/cobra"
)

// listPageOptions reads the shared --page/--limit flags into a ListOptions. The
// white/black-list reads require a non-zero page and limit, so both flags carry
// defaults (registered in each list command's init).
func listPageOptions(cmd *cobra.Command) mail.ListOptions {
	page, _ := cmd.Flags().GetInt("page")
	limit, _ := cmd.Flags().GetInt("limit")
	return mail.ListOptions{Page: page, Limit: limit}
}

// addPageFlags registers the --page/--limit flags a list-read command needs.
func addPageFlags(cmd *cobra.Command) {
	cmd.Flags().Int("page", 1, "page number (1-based)")
	cmd.Flags().Int("limit", 100, "rows per page")
}

// renderAddressList prints an {list, filterInfo} address envelope as a table.
func renderAddressList(cmd *cobra.Command, res *mail.AddressList, header string) error {
	return render(cmd, res, func(w io.Writer) {
		fmt.Fprintln(w, header)
		for _, a := range res.List {
			fmt.Fprintln(w, a)
		}
	})
}

// mailWhitelistCmd groups a mailbox's antispam whitelist (list/add/remove).
var mailWhitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Manage a mailbox's antispam whitelist",
}

var mailWhitelistListCmd = &cobra.Command{
	Use:   "list <domain> <mbox>",
	Short: "List a mailbox's antispam whitelist",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		res, err := c.Mail.Whitelist(cmd.Context(), args[0], args[1], listPageOptions(cmd))
		if err != nil {
			return err
		}
		return renderAddressList(cmd, res, "WHITELIST")
	},
}

var mailWhitelistAddCmd = &cobra.Command{
	Use:   "add <domain> <mbox> <address>",
	Short: "Add an address to a mailbox's whitelist",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		all, _ := cmd.Flags().GetBool("all")
		if err := c.Mail.AddToWhitelist(cmd.Context(), args[0], args[1], args[2], all); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s added to whitelist of %s@%s\n", args[2], args[1], args[0])
		return nil
	},
}

var mailWhitelistRemoveCmd = &cobra.Command{
	Use:   "remove <domain> <mbox> <address>",
	Short: "Remove an address from a mailbox's whitelist",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Mail.DropFromWhitelist(cmd.Context(), args[0], args[1], args[2]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s removed from whitelist of %s@%s\n", args[2], args[1], args[0])
		return nil
	},
}

// mailBlacklistCmd groups a mailbox's antispam blacklist (list/add/remove).
var mailBlacklistCmd = &cobra.Command{
	Use:   "blacklist",
	Short: "Manage a mailbox's antispam blacklist",
}

var mailBlacklistListCmd = &cobra.Command{
	Use:   "list <domain> <mbox>",
	Short: "List a mailbox's antispam blacklist",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		res, err := c.Mail.Blacklist(cmd.Context(), args[0], args[1], listPageOptions(cmd))
		if err != nil {
			return err
		}
		return renderAddressList(cmd, res, "BLACKLIST")
	},
}

var mailBlacklistAddCmd = &cobra.Command{
	Use:   "add <domain> <mbox> <address>",
	Short: "Add an address to a mailbox's blacklist",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		all, _ := cmd.Flags().GetBool("all")
		if err := c.Mail.AddToBlacklist(cmd.Context(), args[0], args[1], args[2], all); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s added to blacklist of %s@%s\n", args[2], args[1], args[0])
		return nil
	},
}

var mailBlacklistRemoveCmd = &cobra.Command{
	Use:   "remove <domain> <mbox> <address>",
	Short: "Remove an address from a mailbox's blacklist",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Mail.DropFromBlacklist(cmd.Context(), args[0], args[1], args[2]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s removed from blacklist of %s@%s\n", args[2], args[1], args[0])
		return nil
	},
}

func init() {
	addPageFlags(mailWhitelistListCmd)
	addPageFlags(mailBlacklistListCmd)
	mailWhitelistAddCmd.Flags().Bool("all", false, "apply the rule to every mailbox of the domain")
	mailBlacklistAddCmd.Flags().Bool("all", false, "apply the rule to every mailbox of the domain")

	mailWhitelistCmd.AddCommand(mailWhitelistListCmd, mailWhitelistAddCmd, mailWhitelistRemoveCmd)
	mailBlacklistCmd.AddCommand(mailBlacklistListCmd, mailBlacklistAddCmd, mailBlacklistRemoveCmd)
	mailCmd.AddCommand(mailWhitelistCmd)
	mailCmd.AddCommand(mailBlacklistCmd)
}
