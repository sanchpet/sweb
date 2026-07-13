package cmd

import (
	"fmt"

	"github.com/sanchpet/sweb-go-sdk/dns"
	"github.com/spf13/cobra"
)

var dnsSrvCmd = &cobra.Command{
	Use:   "srv <domain>",
	Short: "Add, edit, or remove an SRV record",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		action, ok, err := resolveDNSAction(cmd, fmt.Sprintf("SRV record index %d on %s", flagInt(cmd, "index"), domain))
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		f := cmd.Flags()
		priority, _ := f.GetInt("priority")
		ttl, _ := f.GetInt("ttl")
		weight, _ := f.GetInt("weight")
		port, _ := f.GetInt("port")
		target, _ := f.GetString("target")
		service, _ := f.GetString("service")
		protocol, _ := f.GetString("protocol")
		sub, _ := f.GetString("subdomain")
		if err := c.DNS.EditSRV(cmd.Context(), domain, action, dns.SRVRecord{
			Index:     flagInt(cmd, "index"),
			Priority:  priority,
			TTL:       ttl,
			Weight:    weight,
			Target:    target,
			Service:   service,
			Protocol:  protocol,
			Port:      port,
			SubDomain: sub,
		}); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s SRV record on %s\n", pastTense(action), domain)
		return nil
	},
}

func init() {
	addDNSEditFlags(dnsSrvCmd)
	f := dnsSrvCmd.Flags()
	f.Int("priority", 0, "SRV priority")
	f.Int("ttl", 0, "record TTL (seconds)")
	f.Int("weight", 0, "SRV weight")
	f.Int("port", 0, "SRV target port")
	f.String("target", "", "SRV target host")
	f.String("service", "", "service name (e.g. sip)")
	f.String("protocol", "tcp", "protocol (tcp|udp)")
	f.String("subdomain", "", "subdomain (default: the main domain)")
	dnsCmd.AddCommand(dnsSrvCmd)
}
