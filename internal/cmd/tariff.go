package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// tariffCmd groups the account tariff and server-info reads (SDK /tariff): the
// current plan with real usage (show) and the node the account is hosted on
// (server). Both are read-only and hang off root, so tariff is a bindable
// top-level profile group.
var tariffCmd = &cobra.Command{
	Use:   "tariff",
	Short: "Show the account tariff and server info",
}

var tariffShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current tariff plan and real resource usage",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		t, err := c.Tariff.Index(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, t, func(w io.Writer) {
			row := func(k, v string) { fmt.Fprintf(w, "%s\t%s\n", k, v) }
			info := t.Info
			row("PLAN", info.Name)
			row("PLAN ID", fmt.Sprintf("%d", int64(info.PlanID)))
			row("CATEGORY", info.Category)
			row("PRICE", fmt.Sprintf("%d", int64(info.Price)))
			row("PRICE 6MO", fmt.Sprintf("%d", int64(info.Price6)))
			row("PRICE 12MO", fmt.Sprintf("%d", int64(info.Price12)))
			row("RENEWAL (MONTHS)", fmt.Sprintf("%d", int64(info.Duration)))

			usage := t.Usage
			row("DISK QUOTA", fmt.Sprintf("%s / %d", usage.Quota, int64(info.Quota)))
			row("MAIL QUOTA", fmt.Sprintf("%s / %d", usage.MailQuota, int64(info.MailQuota)))
			row("MYSQL", fmt.Sprintf("%d / %d", int64(usage.MySQL), int64(info.MySQL)))
			row("POSTGRESQL", fmt.Sprintf("%d / %d", int64(usage.PostgreSQL), int64(info.PostgreSQL)))
			row("SITES", fmt.Sprintf("%d / %d", int64(usage.Site), int64(info.Site)))
			row("MAILBOXES", fmt.Sprintf("%d", int64(usage.Mailbox)))
		})
	},
}

var tariffServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Show the server the account is hosted on and its software stack",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		s, err := c.Tariff.ServerInfo(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, s, func(w io.Writer) {
			row := func(k, v string) {
				if v != "" {
					fmt.Fprintf(w, "%s\t%s\n", k, v)
				}
			}
			row("SERVER", s.Name)
			row("IP", s.IP)
			row("OS", s.OS)
			row("APACHE", s.Apache)
			row("MYSQL", s.MySQL)
			row("PERL", s.Perl)
			row("PYTHON", s.Python)
			row("RUBY", s.Ruby)
			if len(s.Backend) > 0 {
				fmt.Fprintln(w, "")
				fmt.Fprintln(w, "BACKEND ID\tTYPE\tPORT\tDESCRIPTION")
				for _, b := range s.Backend {
					fmt.Fprintf(w, "%d\t%s\t%d\t%s\n", int64(b.ID), b.Type, int64(b.Port), b.Descr)
				}
			}
		})
	},
}

func init() {
	tariffCmd.AddCommand(tariffShowCmd, tariffServerCmd)
	rootCmd.AddCommand(tariffCmd)
}
