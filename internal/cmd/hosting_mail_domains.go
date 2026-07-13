package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/vh/mail"
	"github.com/spf13/cobra"
)

var mailDomainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "List the account's mail domains and their settings",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		res, err := c.Mail.DomainsList(cmd.Context(), mail.ListOptions{})
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintln(w, "DOMAIN\tMAILBOXES\tQUOTA(MB)\tSPF\tDKIM\tSENDER-VERIFY\tAUTODISCOVER\tCOLLECTOR")
			for _, d := range res.List {
				fmt.Fprintf(w, "%s\t%d\t%d\t%s\t%s\t%s\t%s\t%s\n",
					d.FQDN, int64(d.MailboxesCnt), int64(d.Quota),
					onOff(d.SPF == 1), d.DKIM, onOff(d.SenderVerify == 1),
					onOff(d.AutoDiscover == 1), emptyDash(d.EmailCollector))
			}
		})
	},
}

var mailQuotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Show the total size (MB) of all mailboxes on the account",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		total, err := c.Mail.MailQuota(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, map[string]int64{"quotaMB": total}, func(w io.Writer) {
			kv(w, "QUOTA(MB)", fmt.Sprintf("%d", total))
		})
	},
}

// mailDomainCmd groups the domain-level mail toggles (SPF/sender-verify/
// autodiscover for every mailbox of a domain at once); DKIM lives under its own
// `mail dkim` group.
var mailDomainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Domain-level mail toggles (SPF, sender-verify, autodiscover)",
}

var mailDomainSpfCmd = &cobra.Command{
	Use:   "spf <domain> <on|off>",
	Short: "Toggle SPF filtering for every mailbox of a domain",
	Args:  cobra.ExactArgs(2),
	RunE:  domainToggle(func(c mailClient) toggleFn { return c.ChangeDomainSpf }, "SPF"),
}

var mailDomainSenderVerifyCmd = &cobra.Command{
	Use:   "sender-verify <domain> <on|off>",
	Short: "Toggle sender-address verification for a domain",
	Args:  cobra.ExactArgs(2),
	RunE:  domainToggle(func(c mailClient) toggleFn { return c.ChangeSenderVerify }, "sender-verify"),
}

var mailDomainAutoDiscoverCmd = &cobra.Command{
	Use:   "autodiscover <domain> <on|off>",
	Short: "Toggle mail-client auto-configuration for a domain",
	Args:  cobra.ExactArgs(2),
	RunE:  domainToggle(func(c mailClient) toggleFn { return c.ChangeAutoDiscover }, "autodiscover"),
}

func init() {
	mailCmd.AddCommand(mailDomainsCmd)
	mailCmd.AddCommand(mailQuotaCmd)
	mailDomainCmd.AddCommand(mailDomainSpfCmd)
	mailDomainCmd.AddCommand(mailDomainSenderVerifyCmd)
	mailDomainCmd.AddCommand(mailDomainAutoDiscoverCmd)
	mailCmd.AddCommand(mailDomainCmd)
}
