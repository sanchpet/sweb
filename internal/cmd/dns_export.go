package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var dnsExportCmd = &cobra.Command{
	Use:   "export <domain>",
	Short: "Print the raw BIND-style zone file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		zf, err := c.DNS.GetFile(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		// -o json emits the full ZoneFile object; the table path prints the raw
		// zone content (what you'd redirect to a .zone file).
		return render(cmd, zf, func(w io.Writer) {
			fmt.Fprintln(w, zf.Content)
		})
	},
}

func init() {
	dnsCmd.AddCommand(dnsExportCmd)
}
