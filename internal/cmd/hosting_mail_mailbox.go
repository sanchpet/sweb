package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/vh/mail"
	"github.com/spf13/cobra"
)

// mailMailboxCmd groups per-mailbox operations: list, lifecycle (create/delete),
// password/comment, the antispam filter level, per-mailbox SPF, and purging old
// mail.
var mailMailboxCmd = &cobra.Command{
	Use:     "mailbox",
	Aliases: []string{"mbox"},
	Short:   "Manage a domain's mailboxes",
}

var mailMailboxListCmd = &cobra.Command{
	Use:   "list <domain>",
	Short: "List the mailboxes of a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		search, _ := cmd.Flags().GetString("search")
		res, err := c.Mail.MailboxesList(cmd.Context(), args[0], search, mail.ListOptions{})
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintln(w, "MAILBOX\tQUOTA(MB)\tSPF\tANTISPAM\tPURPOSE\tCOMMENT")
			for _, m := range res.List {
				fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\n",
					m.Mbox, int64(m.Quota), onOff(m.SPF == 1),
					antispamLabel(int(m.Antispam)), emptyDash(m.Purpose), emptyDash(m.Comment))
			}
		})
	},
}

var mailMailboxCreateCmd = &cobra.Command{
	Use:   "create <domain> <mbox>",
	Short: "Create a mailbox (billable) and print its credentials",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		password, _ := cmd.Flags().GetString("password")
		comment, _ := cmd.Flags().GetString("comment")
		if password == "" {
			return fmt.Errorf("--password is required")
		}
		if !confirmed(cmd, fmt.Sprintf("Create mailbox %s@%s? This is billable.", args[1], args[0]), "Create") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.Mail.CreateMbox(cmd.Context(), args[0], args[1], password, comment)
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			kv(w, "LOGIN", res.Login)
			kv(w, "PASSWORD", res.Password)
			kv(w, "WEBMAIL", res.WebMail)
			for _, s := range res.MailProgramSettings {
				kv(w, s.Name, fmt.Sprintf("%s:%s", s.Server, s.Port))
			}
		})
	},
}

var mailMailboxDeleteCmd = &cobra.Command{
	Use:   "delete <domain> <mbox>",
	Short: "Delete a mailbox",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Delete mailbox %s@%s? This is irreversible.", args[1], args[0]), "Delete") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Mail.DropMbox(cmd.Context(), args[0], args[1]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "deleted %s@%s\n", args[1], args[0])
		return nil
	},
}

var mailMailboxPasswordCmd = &cobra.Command{
	Use:   "password <domain> <mbox>",
	Short: "Change a mailbox's password",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		password, _ := cmd.Flags().GetString("password")
		if password == "" {
			return fmt.Errorf("--password is required")
		}
		if err := c.Mail.ChangeMailboxPassword(cmd.Context(), args[0], args[1], password); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "password changed for %s@%s\n", args[1], args[0])
		return nil
	},
}

var mailMailboxCommentCmd = &cobra.Command{
	Use:   "comment <domain> <mbox> <comment>",
	Short: "Set a mailbox's comment",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.Mail.UpdateComment(cmd.Context(), args[0], args[1], args[2]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "comment set for %s@%s\n", args[1], args[0])
		return nil
	},
}

var mailMailboxAntispamCmd = &cobra.Command{
	Use:   "antispam <domain> <mbox> <hard|medium|soft|off>",
	Short: "Set a mailbox's antispam filter level",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		level, err := antispamValue(args[2])
		if err != nil {
			return err
		}
		if err := c.Mail.UpdateAntispamState(cmd.Context(), args[0], args[1], level); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "antispam set to %s for %s@%s\n", args[2], args[1], args[0])
		return nil
	},
}

var mailMailboxSpfCmd = &cobra.Command{
	Use:   "spf <domain> <mbox> <on|off>",
	Short: "Toggle SPF filtering for a single mailbox",
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
		if err := c.Mail.ChangeMailboxSpf(cmd.Context(), args[0], args[1], on); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "SPF %s for %s@%s\n", onOff(on), args[1], args[0])
		return nil
	},
}

var mailMailboxPurgeCmd = &cobra.Command{
	Use:   "purge <domain> <mbox>",
	Short: "Delete a mailbox's messages older than N days",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		days, _ := cmd.Flags().GetInt("days")
		if !confirmed(cmd, fmt.Sprintf("Delete messages older than %d days in %s@%s? This is irreversible.", days, args[1], args[0]), "Delete") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Mail.DeleteMails(cmd.Context(), args[0], args[1], days); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "purged messages older than %d days in %s@%s\n", days, args[1], args[0])
		return nil
	},
}

var mailMailboxRequisitesCmd = &cobra.Command{
	Use:   "requisites <mbox-login> <email>",
	Short: "Email a mailbox's connection requisites to an address",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		password, _ := cmd.Flags().GetString("password")
		if password == "" {
			return fmt.Errorf("--password is required")
		}
		if err := c.Mail.SendRequisitesToEmail(cmd.Context(), args[1], args[0], password); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "requisites for %s sent to %s\n", args[0], args[1])
		return nil
	},
}

func init() {
	mailMailboxListCmd.Flags().String("search", "", "substring filter on the mailbox name")
	mailMailboxCreateCmd.Flags().String("password", "", "mailbox password (required)")
	mailMailboxCreateCmd.Flags().String("comment", "", "mailbox comment")
	mailMailboxCreateCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	mailMailboxDeleteCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	mailMailboxPasswordCmd.Flags().String("password", "", "new mailbox password (required)")
	mailMailboxPurgeCmd.Flags().Int("days", 30, "delete messages older than this many days")
	mailMailboxPurgeCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	mailMailboxRequisitesCmd.Flags().String("password", "", "the mailbox's password (required)")

	for _, sub := range []*cobra.Command{
		mailMailboxListCmd, mailMailboxCreateCmd, mailMailboxDeleteCmd,
		mailMailboxPasswordCmd, mailMailboxCommentCmd, mailMailboxAntispamCmd,
		mailMailboxSpfCmd, mailMailboxPurgeCmd, mailMailboxRequisitesCmd,
	} {
		mailMailboxCmd.AddCommand(sub)
	}
	mailCmd.AddCommand(mailMailboxCmd)
}
