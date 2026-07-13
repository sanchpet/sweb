package cmd

import (
	"fmt"

	"github.com/sanchpet/sweb-go-sdk/dns"
	"github.com/spf13/cobra"
)

var dnsMxCmd = &cobra.Command{
	Use:   "mx <domain>",
	Short: "Add, edit, or remove an MX record",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		action, ok, err := resolveDNSAction(cmd, fmt.Sprintf("MX record index %d on %s", flagInt(cmd, "index"), domain))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		priority, _ := cmd.Flags().GetInt("priority")
		value, _ := cmd.Flags().GetString("value")
		sub, _ := cmd.Flags().GetString("subdomain")
		if err := c.DNS.EditMX(cmd.Context(), domain, action, dns.MXRecord{
			Index:     flagInt(cmd, "index"),
			Priority:  priority,
			Value:     value,
			SubDomain: sub,
		}); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s MX record on %s\n", pastTense(action), domain)
		return nil
	},
}

func init() {
	addDNSEditFlags(dnsMxCmd)
	dnsMxCmd.Flags().Int("priority", 10, "MX priority")
	dnsMxCmd.Flags().String("value", "", "mail server host")
	dnsMxCmd.Flags().String("subdomain", "", "subdomain (default: the main domain)")
	dnsCmd.AddCommand(dnsMxCmd)
}
