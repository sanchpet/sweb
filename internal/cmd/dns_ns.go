package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dnsNsCmd = &cobra.Command{
	Use:   "ns <domain>",
	Short: "Add, edit, or remove an NS record",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		action, ok, err := resolveDNSAction(cmd, fmt.Sprintf("NS record index %d on %s", flagInt(cmd, "index"), domain))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		sub, _ := cmd.Flags().GetString("subdomain")
		value, _ := cmd.Flags().GetString("value")
		if err := c.DNS.EditNS(cmd.Context(), domain, action, flagInt(cmd, "index"), sub, value); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s NS record on %s\n", pastTense(action), domain)
		return nil
	},
}

func init() {
	addDNSEditFlags(dnsNsCmd)
	dnsNsCmd.Flags().String("subdomain", "", "subdomain")
	dnsNsCmd.Flags().String("value", "", "nameserver host")
	dnsCmd.AddCommand(dnsNsCmd)
}
