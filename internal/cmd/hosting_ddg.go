package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/vh/ddg"
	"github.com/spf13/cobra"
)

// ddgCmd groups the shared-hosting DDoS-Guard operations (SDK /vh/ddg): the read
// side (list protected domains, enable-page catalogue, domain count, price) plus
// the per-domain enable/disable mutations. It hangs off the hosting parent, so it
// inherits that group's profile binding.
var ddgCmd = &cobra.Command{
	Use:   "ddg",
	Short: "Manage DDoS-Guard protection",
}

var ddgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's domains and their DDoS-Guard state",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		list, err := c.DDoSGuard.List(cmd.Context(), &ddg.ListOptions{
			Page:        flagInt(cmd, "page"),
			PerPage:     flagInt(cmd, "per-page"),
			OrderField:  mustFlagString(cmd, "order-field"),
			OrderDirect: mustFlagString(cmd, "order-direct"),
		})
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "DOMAIN\tSTATUS\tIP\tEXPIRED\tBLOCKED")
			for _, d := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					d.FQDNReadable, d.Status, dash(d.IP), dash(d.Expired), dash(d.Blocked))
			}
		})
	},
}

var ddgInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show the enable-page catalogue: domains eligible to connect and the price",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		info, err := c.DDoSGuard.EnableInfo(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, info, func(w io.Writer) {
			fmt.Fprintf(w, "PRICE\t%g\n", float64(info.Price))
			fmt.Fprintln(w, "DOMAIN\tON OUR NS\tSSL")
			for _, d := range info.Domains {
				fmt.Fprintf(w, "%s\t%s\t%s\n", d.FQDNReadable, yesNo(d.IsOnOurNS), sslCell(d.SSL))
			}
		})
	},
}

var ddgCountCmd = &cobra.Command{
	Use:   "count",
	Short: "Count the account's domains (excludes technical ones)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		n, err := c.DDoSGuard.CountAllDomains(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, map[string]int64{"count": n}, func(w io.Writer) {
			fmt.Fprintf(w, "COUNT\t%d\n", n)
		})
	},
}

var ddgPriceCmd = &cobra.Command{
	Use:   "price",
	Short: "Show the DDoS-Guard service price",
	Long: `Show the DDoS-Guard service price.

By default this reports the price for your current tariff plan (method
"getPrice"). With --widget it reports the "Change tariff" widget pricing
(method "priceWidget"): the price on the current tariff and on a tariff from
another line.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if widget, _ := cmd.Flags().GetBool("widget"); widget {
			p, err := c.DDoSGuard.PriceWidget(cmd.Context())
			if err != nil {
				return err
			}
			return render(cmd, p, func(w io.Writer) {
				fmt.Fprintf(w, "CURRENT TARIFF\t%g\n", float64(p.Current))
				fmt.Fprintf(w, "OTHER TARIFF\t%g\n", float64(p.New))
			})
		}
		price, err := c.DDoSGuard.GetPrice(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, map[string]float64{"price": price}, func(w io.Writer) {
			fmt.Fprintf(w, "PRICE\t%g\n", price)
		})
	},
}

var ddgEnableCmd = &cobra.Command{
	Use:   "enable <domain>",
	Short: "Enable DDoS-Guard for a domain — bills",
	Long: `Enable (or unblock) DDoS-Guard for a domain via the "enable" method.

This CONNECTS a paid service and bills your account. You are asked to confirm
unless --yes is given. On success it reports the domain's new IP address.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		if !confirmed(cmd, fmt.Sprintf("Enable DDoS-Guard for %s? This connects a paid service and bills your account.", domain), "Enable") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.DDoSGuard.Enable(cmd.Context(), domain)
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintf(w, "DOMAIN\t%s\n", res.FQDNReadable)
			fmt.Fprintf(w, "ON OUR NS\t%s\n", yesNo(res.IsOnOurNS))
			fmt.Fprintf(w, "IP\t%s\n", dash(res.IP))
		})
	},
}

var ddgDisableCmd = &cobra.Command{
	Use:   "disable <domain>",
	Short: "Disable DDoS-Guard for a domain",
	Long: `Disable (or block) DDoS-Guard for a domain via the "disable" method.

You are asked to confirm unless --yes is given. On success it reports the
service's paid-through date.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		if !confirmed(cmd, fmt.Sprintf("Disable DDoS-Guard for %s?", domain), "Disable") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.DDoSGuard.Disable(cmd.Context(), domain)
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintf(w, "DOMAIN\t%s\n", res.FQDNReadable)
			fmt.Fprintf(w, "ON OUR NS\t%s\n", yesNo(res.IsOnOurNS))
			fmt.Fprintf(w, "PAID THROUGH\t%s\n", dash(res.Expire))
		})
	},
}

// dash renders an empty string as "—", used for the API's nullable date/IP fields.
func dash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

// sslCell renders an eligible domain's SSL summary for the info table.
func sslCell(s *ddg.SSLInfo) string {
	if s == nil {
		return "none"
	}
	if s.IsOur {
		return "ours"
	}
	if s.IsFilled {
		return "external"
	}
	return "none"
}

// mustFlagString reads a string flag, ignoring the lookup error (the flag is
// always registered on the command that reads it).
func mustFlagString(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetString(name)
	return v
}

func init() {
	ddgListCmd.Flags().Int("page", 0, "1-based page number")
	ddgListCmd.Flags().Int("per-page", 0, "domains per page")
	ddgListCmd.Flags().String("order-field", "", "sort field: status|fqdn")
	ddgListCmd.Flags().String("order-direct", "", "sort direction: ASC|DESC")

	ddgPriceCmd.Flags().Bool("widget", false, "show the 'Change tariff' widget pricing (priceWidget)")

	ddgEnableCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	ddgDisableCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	ddgCmd.AddCommand(
		ddgListCmd,
		ddgInfoCmd,
		ddgCountCmd,
		ddgPriceCmd,
		ddgEnableCmd,
		ddgDisableCmd,
	)
	hostingCmd.AddCommand(ddgCmd)
}
