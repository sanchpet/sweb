package cmd

import (
	"fmt"
	"io"

	domainapi "github.com/sanchpet/sweb-go-sdk/domains"
	"github.com/spf13/cobra"
)

var domainsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's domains",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		typ, _ := cmd.Flags().GetString("type")
		filter, _ := cmd.Flags().GetString("filter")
		packages, _ := cmd.Flags().GetBool("packages")
		domains, err := c.Domains.List(cmd.Context(), &domainapi.ListOptions{
			Type:         typ,
			Filter:       filter,
			ShowPackages: packages,
		})
		if err != nil {
			return err
		}
		return render(cmd, domains, func(w io.Writer) {
			fmt.Fprintln(w, "DOMAIN\tTECH\tDOCROOT\tSUBS")
			for _, d := range domains {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", d.FQDNReadable, d.FQDNTech, d.Docroot, len(d.Subdomains))
			}
		})
	},
}

func init() {
	domainsListCmd.Flags().String("type", "all", "domain type: all|sweb|free|other")
	domainsListCmd.Flags().String("filter", "", "substring filter on the domain name")
	domainsListCmd.Flags().Bool("packages", false, "include domain packages")
	domainsCmd.AddCommand(domainsListCmd)
}
