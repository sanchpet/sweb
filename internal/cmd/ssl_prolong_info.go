package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/spf13/cobra"
)

var cloudSSLProlongInfoCmd = &cobra.Command{
	Use:   "prolong-info <id>",
	Short: "Show the prolongation options for a certificate",
	Long: `Show a certificate's prolongation options via the "getProlongInfo" method.

<id> is a certificate id (see 'sweb ssl list'); the result lists the per-period
prices and the product ids.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("id must be an integer: %w", err)
		}
		info, err := c.SSL.ProlongInfo(cmd.Context(), id)
		if err != nil {
			return err
		}
		return render(cmd, info, func(w io.Writer) {
			if info == nil {
				fmt.Fprintln(w, "no prolongation offered")
				return
			}
			fmt.Fprintf(w, "TITLE\t%s\n", info.Title)
			fmt.Fprintf(w, "CURRENT ID\t%d\n", int64(info.CurrentCertificateID))
			fmt.Fprintf(w, "FREE\t%s\n", yesNo(info.IsFreeCertificate))
			fmt.Fprintln(w, "PERIOD(MONTHS)\tPRICE\tPRODUCT ID")
			for period, price := range info.Prices {
				fmt.Fprintf(w, "%s\t%.2f\t%s\n", period, float64(price), info.IDs[period])
			}
		})
	},
}

func init() {
	cloudSSLCmd.AddCommand(cloudSSLProlongInfoCmd)
}
