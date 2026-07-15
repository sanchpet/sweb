package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/apierr"
	"github.com/spf13/cobra"
)

// domainPricing is the combined read behind `domains price`. Each check is
// independent: the API answers only some of them depending on the domain's
// state (see the command's Long help), so a field is nil when its call was
// refused, with the reason collected in Notes.
type domainPricing struct {
	Domain            string   `json:"domain"`
	RegistrationPrice *float64 `json:"registration_price,omitempty"`
	RegisterAvailable *bool    `json:"register_available,omitempty"`
	TransferAvailable *bool    `json:"transfer_available,omitempty"`
	Notes             []string `json:"notes,omitempty"`
}

var domainsPriceCmd = &cobra.Command{
	Use:   "price <domain>",
	Short: "Show registration price and registration/transfer availability",
	Long: `Show a domain's registration price and registration/transfer availability.

The three underlying API methods answer for different domain states:
regAvailable works for any domain (reporting whether it can be registered),
while priceForRegistration and priceForTrasfer only answer for a domain already
on your account. Each check is reported independently — one the API refuses is
shown as "—", with the reason listed under the table.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		pay, _ := cmd.Flags().GetString("pay")
		ctx := cmd.Context()

		p := domainPricing{Domain: domain}
		note := func(what string, err error) {
			p.Notes = append(p.Notes, fmt.Sprintf("%s — %s", what, apiReason(err)))
		}
		if v, err := c.Domains.RegistrationPrice(ctx, domain); err == nil {
			p.RegistrationPrice = &v
		} else {
			note("registration price", err)
		}
		if v, err := c.Domains.RegAvailable(ctx, domain, pay); err == nil {
			p.RegisterAvailable = &v
		} else {
			note("register available", err)
		}
		if v, err := c.Domains.TransferAvailable(ctx, domain); err == nil {
			p.TransferAvailable = &v
		} else {
			note("transfer available", err)
		}

		return render(cmd, p, func(w io.Writer) {
			fmt.Fprintf(w, "REGISTRATION PRICE\t%s\n", floatCell(p.RegistrationPrice))
			fmt.Fprintf(w, "REGISTER AVAILABLE\t%s\n", boolCell(p.RegisterAvailable))
			fmt.Fprintf(w, "TRANSFER AVAILABLE\t%s\n", boolCell(p.TransferAvailable))
			for _, n := range p.Notes {
				fmt.Fprintf(w, "\t%s\n", n)
			}
		})
	},
}

// floatCell renders an optional price, "—" when the API did not answer.
func floatCell(f *float64) string {
	if f == nil {
		return "—"
	}
	return fmt.Sprintf("%g", *f)
}

// boolCell renders an optional yes/no, "—" when the API did not answer.
func boolCell(b *bool) string {
	if b == nil {
		return "—"
	}
	return yesNo(*b)
}

// apiReason extracts the human-readable message from a SpaceWeb API error.
func apiReason(err error) string {
	var e *apierr.Error
	if errors.As(err, &e) {
		return e.Message
	}
	return err.Error()
}

func init() {
	domainsPriceCmd.Flags().String("pay", "balance", "payment source for the availability check: balance|bonus")
	domainsCmd.AddCommand(domainsPriceCmd)
}
