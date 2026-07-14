package cmd

import (
	"fmt"
	"strconv"

	"github.com/sanchpet/sweb-go-sdk/ssl"
	"github.com/spf13/cobra"
)

var cloudSSLOrderCmd = &cobra.Command{
	Use:   "order <domain> <certificate-id> <confirm-mail>",
	Short: "Order a certificate — bills the account",
	Long: `Order a certificate via the "orderSubmit" method.

<domain> is the fully-qualified domain to cover, <certificate-id> is a product
id from 'sweb ssl order-list', and <confirm-mail> is the confirmation mailbox.
Flags carry the optional order fields (--person-id, --company-link, --subdomain,
--autoprolong, --old-certificate-id, --from-prolongation).

This orders a PAID certificate and BILLS the account — you are asked to confirm
unless --yes is given.`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		certificateID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("certificate-id must be an integer: %w", err)
		}
		confirmMail := args[2]
		opts := &ssl.OrderSubmitOptions{
			PersonID:         flagInt(cmd, "person-id"),
			CompanyPageLink:  cloudSSLFlagString(cmd, "company-link"),
			Subdomain:        cloudSSLFlagString(cmd, "subdomain"),
			OldCertificateID: flagInt(cmd, "old-certificate-id"),
		}
		opts.Autoprolong, _ = cmd.Flags().GetBool("autoprolong")
		opts.FromProlongation, _ = cmd.Flags().GetBool("from-prolongation")
		if !confirmed(cmd, fmt.Sprintf("Order a paid certificate for %s? This bills the account.", domain), "Order") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.SSL.OrderSubmit(cmd.Context(), domain, certificateID, confirmMail, opts)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Ordered a certificate for %s (%s)\n", domain, res)
		return nil
	},
}

func init() {
	cloudSSLOrderCmd.Flags().Int("person-id", 0, "domain-person id")
	cloudSSLOrderCmd.Flags().String("company-link", "", `"about the company" URL (EV/OV)`)
	cloudSSLOrderCmd.Flags().String("subdomain", "", "subdomain the certificate covers")
	cloudSSLOrderCmd.Flags().Bool("autoprolong", false, "enable auto-prolongation on the new certificate")
	cloudSSLOrderCmd.Flags().Int("old-certificate-id", 0, "prior domain-person id (prolongation)")
	cloudSSLOrderCmd.Flags().Bool("from-prolongation", false, "order originates from a prolongation confirmation")
	cloudSSLOrderCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	cloudSSLCmd.AddCommand(cloudSSLOrderCmd)
}
