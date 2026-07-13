package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// mailAutoreplyCmd groups a mailbox's autoresponder (get/set).
var mailAutoreplyCmd = &cobra.Command{
	Use:   "autoreply",
	Short: "Manage a mailbox's autoresponder",
}

var mailAutoreplyGetCmd = &cobra.Command{
	Use:   "get <domain> <mbox>",
	Short: "Show a mailbox's autoresponder text",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		text, err := c.Mail.Autoreply(cmd.Context(), args[0], args[1])
		if err != nil {
			return err
		}
		return render(cmd, map[string]string{"autoreply": text}, func(w io.Writer) {
			fmt.Fprintln(w, text)
		})
	},
}

var mailAutoreplySetCmd = &cobra.Command{
	Use:   "set <domain> <mbox> <text>",
	Short: "Set a mailbox's autoresponder text (empty text disables it)",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Mail.ChangeAutoreply(cmd.Context(), args[0], args[1], args[2]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "autoreply set for %s@%s\n", args[1], args[0])
		return nil
	},
}

// mailDkimCmd groups the domain-level DKIM signing toggle.
var mailDkimCmd = &cobra.Command{
	Use:   "dkim",
	Short: "Manage a domain's DKIM signing",
}

var mailDkimEnableCmd = &cobra.Command{
	Use:   "enable <domain>",
	Short: "Enable DKIM signing for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Mail.EnableDkim(cmd.Context(), args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "DKIM enabled for %s\n", args[0])
		return nil
	},
}

var mailDkimDisableCmd = &cobra.Command{
	Use:   "disable <domain>",
	Short: "Disable DKIM signing for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Mail.DisableDkim(cmd.Context(), args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "DKIM disabled for %s\n", args[0])
		return nil
	},
}

func init() {
	mailAutoreplyCmd.AddCommand(mailAutoreplyGetCmd)
	mailAutoreplyCmd.AddCommand(mailAutoreplySetCmd)
	mailCmd.AddCommand(mailAutoreplyCmd)

	mailDkimCmd.AddCommand(mailDkimEnableCmd)
	mailDkimCmd.AddCommand(mailDkimDisableCmd)
	mailCmd.AddCommand(mailDkimCmd)
}
