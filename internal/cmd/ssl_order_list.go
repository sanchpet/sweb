package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var cloudSSLOrderListCmd = &cobra.Command{
	Use:   "order-list",
	Short: "List the certificate products available for order",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		list, err := c.SSL.OrderList(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tNAME\tTYPE\tADVANTAGE")
			for _, o := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", o.ID, o.Name, o.Type, o.AdvantageText)
			}
		})
	},
}

func init() {
	cloudSSLCmd.AddCommand(cloudSSLOrderListCmd)
}
