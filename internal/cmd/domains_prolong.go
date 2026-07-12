package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var domainsProlongCmd = &cobra.Command{
	Use:   "prolong <domain>",
	Short: "Prolong (renew) a domain's registration — bills",
	Long: `Prolong a domain via the "prolong" method. This CHARGES the account (money or
bonus points per --pay). You are asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		if !confirmed(cmd, fmt.Sprintf("Prolong %s? This charges the account.", domain), "Prolong") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		pay, _ := cmd.Flags().GetString("pay")
		if err := c.Domains.Prolong(cmd.Context(), domain, pay); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Prolonged %s\n", domain)
		return nil
	},
}

func init() {
	domainsProlongCmd.Flags().String("pay", "balance", "payment source: balance|bonus")
	domainsProlongCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	domainsCmd.AddCommand(domainsProlongCmd)
}
