package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/ssl"
	"github.com/spf13/cobra"
)

var cloudSSLListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's SSL certificates",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		res, err := c.SSL.List(cmd.Context(), &ssl.ListOptions{
			Page:        flagInt(cmd, "page"),
			PerPage:     flagInt(cmd, "per-page"),
			OrderField:  cloudSSLFlagString(cmd, "order-field"),
			OrderDirect: cloudSSLFlagString(cmd, "order-direct"),
		})
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tDOMAIN\tNAME\tSTATUS\tIP\tVALID TO\tAUTOPROLONG")
			if res == nil {
				return
			}
			for _, cert := range res.List {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n",
					int64(cert.ID), cert.Domain, cert.Name, cert.Status, cert.IP,
					cert.ValidTo, yesNo(cert.Autoprolong))
			}
		})
	},
}

func init() {
	cloudSSLListCmd.Flags().Int("page", 0, "1-based page number")
	cloudSSLListCmd.Flags().Int("per-page", 0, "records per page")
	cloudSSLListCmd.Flags().String("order-field", "", "sort field: id|valid_to|fqdn|status|ip")
	cloudSSLListCmd.Flags().String("order-direct", "", "sort direction: asc|desc")

	cloudSSLCmd.AddCommand(cloudSSLListCmd)
}
