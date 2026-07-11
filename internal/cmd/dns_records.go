package cmd

import (
	"fmt"
	"io"
	"strings"

	sweb "github.com/sanchpet/sweb-go-sdk"
	"github.com/spf13/cobra"
)

var dnsRecordsCmd = &cobra.Command{
	Use:   "records <domain>",
	Short: "List a domain's DNS records",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		recs, err := c.DNS.Records(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return render(cmd, recs, func(w io.Writer) {
			fmt.Fprintln(w, "INDEX\tTYPE\tNAME\tVALUE\tDETAILS")
			for _, r := range recs {
				name := r.Name
				if name == "" {
					name = "@"
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", int64(r.Index), r.Type, name, truncateCell(r.Value, 60), recordDetails(r))
			}
		})
	},
}

// truncateCell shortens a value for table display so one long record (a DKIM
// TXT can run 250+ chars) does not blow the column width out for every row. The
// full value is always available via -o json or `dns export`.
func truncateCell(s string, limit int) string {
	r := []rune(s)
	if len(r) <= limit {
		return s
	}
	return string(r[:limit-1]) + "…"
}

// recordDetails renders the type-specific fields of a record into one column.
func recordDetails(r sweb.DNSRecord) string {
	var parts []string
	add := func(k string, v int64) {
		if v != 0 {
			parts = append(parts, fmt.Sprintf("%s=%d", k, v))
		}
	}
	switch r.Type {
	case "MX":
		add("priority", int64(r.Priority))
	case "SRV":
		add("priority", int64(r.Priority))
		add("weight", int64(r.Weight))
		add("port", int64(r.Port))
		add("ttl", int64(r.TTL))
		if r.Target != "" {
			parts = append(parts, "target="+r.Target)
		}
	case "TXT":
		add("ttl", int64(r.TTL))
	}
	return strings.Join(parts, " ")
}

func init() {
	dnsCmd.AddCommand(dnsRecordsCmd)
}
