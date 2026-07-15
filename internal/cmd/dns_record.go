package cmd

import (
	"fmt"

	"github.com/sanchpet/sweb-go-sdk/dns"
	"github.com/spf13/cobra"
)

var dnsRecordCmd = &cobra.Command{
	Use:   "record <domain>",
	Short: "Add, edit, or remove a general record (A, AAAA, CNAME, …)",
	Long: `Add, edit, or remove a general DNS record via the "editMain" method —
the record types without a dedicated command (A, AAAA, CNAME, …). Use --type to
pick the record type and --index (from 'sweb dns records') to target an existing
record for edit/remove.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		rtype, _ := cmd.Flags().GetString("type")
		action, ok, err := resolveDNSAction(cmd, fmt.Sprintf("%s record index %d on %s", rtype, flagInt(cmd, "index"), domain))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		f := cmd.Flags()
		name, _ := f.GetString("name")
		value, _ := f.GetString("value")
		if err := c.DNS.EditMain(cmd.Context(), domain, action, dns.MainRecord{
			Index: flagInt(cmd, "index"),
			Name:  name,
			Type:  rtype,
			Value: value,
		}); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s %s record on %s\n", pastTense(action), rtype, domain)
		return nil
	},
}

func init() {
	addDNSEditFlags(dnsRecordCmd)
	f := dnsRecordCmd.Flags()
	f.String("name", "", "subdomain name, or empty for the apex")
	f.String("type", "A", "record type: A, AAAA, CNAME, …")
	f.String("value", "", "record value")
	dnsCmd.AddCommand(dnsRecordCmd)
}
