package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var domainsInfoCmd = &cobra.Command{
	Use:   "info <domain>",
	Short: "Show full information for a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		info, err := c.Domains.Info(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return render(cmd, info, func(w io.Writer) {
			row := func(k, v string) { fmt.Fprintf(w, "%s\t%s\n", k, v) }
			row("DOMAIN", args[0])
			row("OURS", yesNo(info.IsOur == 1))
			row("EXPIRED", info.Expired)
			row("AUTOPROLONG", info.Autoprolong)
			row("CAN PROLONG", yesNo(info.CanProlong == 1))
			row("PROLONG PRICE", fmt.Sprintf("%d", int64(info.ProlongPrice)))
			row("REG PRICE", fmt.Sprintf("%d", int64(info.RegPrice)))
			if info.TransferPrice >= 0 { // -1 means transfer not offered
				row("TRANSFER PRICE", fmt.Sprintf("%d", int64(info.TransferPrice)))
			}
			if info.Registrar != "" {
				row("REGISTRAR", info.Registrar)
			}
			row("DOCROOT", info.DocRoot)
			if info.RedirectURL != "" {
				row("REDIRECT", info.RedirectURL)
			}
		})
	},
}

func init() {
	domainsCmd.AddCommand(domainsInfoCmd)
}
