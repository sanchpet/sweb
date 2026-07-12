package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// domainPricing is the combined read behind `domains price`.
type domainPricing struct {
	Domain            string  `json:"domain"`
	RegistrationPrice float64 `json:"registration_price"`
	RegisterAvailable bool    `json:"register_available"`
	TransferAvailable bool    `json:"transfer_available"`
}

var domainsPriceCmd = &cobra.Command{
	Use:   "price <domain>",
	Short: "Show registration price and registration/transfer availability",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		pay, _ := cmd.Flags().GetString("pay")

		price, err := c.Domains.RegistrationPrice(cmd.Context(), domain)
		if err != nil {
			return err
		}
		regOK, err := c.Domains.RegAvailable(cmd.Context(), domain, pay)
		if err != nil {
			return err
		}
		xferOK, err := c.Domains.TransferAvailable(cmd.Context(), domain)
		if err != nil {
			return err
		}

		p := domainPricing{
			Domain:            domain,
			RegistrationPrice: price,
			RegisterAvailable: regOK,
			TransferAvailable: xferOK,
		}
		return render(cmd, p, func(w io.Writer) {
			fmt.Fprintf(w, "REGISTRATION PRICE\t%g\n", p.RegistrationPrice)
			fmt.Fprintf(w, "REGISTER AVAILABLE\t%s\n", yesNo(p.RegisterAvailable))
			fmt.Fprintf(w, "TRANSFER AVAILABLE\t%s\n", yesNo(p.TransferAvailable))
		})
	},
}

func init() {
	domainsPriceCmd.Flags().String("pay", "balance", "payment source for the availability check: balance|bonus")
	domainsCmd.AddCommand(domainsPriceCmd)
}
