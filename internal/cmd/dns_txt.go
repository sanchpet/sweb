package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dnsTxtCmd = &cobra.Command{
	Use:   "txt <domain>",
	Short: "Add, edit, or remove a TXT record",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		action, ok, err := resolveDNSAction(cmd, fmt.Sprintf("TXT record index %d on %s", flagInt(cmd, "index"), domain))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		sub, _ := cmd.Flags().GetString("subdomain")
		value, _ := cmd.Flags().GetString("value")
		if err := c.DNS.EditTXT(cmd.Context(), domain, action, flagInt(cmd, "index"), sub, value); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s TXT record on %s\n", pastTense(action), domain)
		return nil
	},
}

func init() {
	addDNSEditFlags(dnsTxtCmd)
	dnsTxtCmd.Flags().String("subdomain", "", "subdomain (default: the main domain)")
	dnsTxtCmd.Flags().String("value", "", "TXT value (e.g. an SPF string)")
	dnsCmd.AddCommand(dnsTxtCmd)
}
