package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var domainsSubdomainsCmd = &cobra.Command{
	Use:   "subdomains <domain>",
	Short: "List a domain's subdomains",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		subs, err := c.Domains.Subdomains(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return render(cmd, subs, func(w io.Writer) {
			fmt.Fprintln(w, "NAME\tVALUE")
			for _, s := range subs {
				fmt.Fprintf(w, "%s\t%s\n", s.Name, s.Value)
			}
		})
	},
}

func init() {
	domainsCmd.AddCommand(domainsSubdomainsCmd)
}
