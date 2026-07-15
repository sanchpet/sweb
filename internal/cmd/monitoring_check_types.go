package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var checkTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List the available check types (Ping, HTTP, Port)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		types, err := c.MonitoringChecks.GetTypes(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, types, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tCODE\tNAME")
			for _, t := range types {
				fmt.Fprintf(w, "%s\t%s\t%s\n", t.ID, t.Code, t.Name)
			}
		})
	},
}

func init() {
	checkCmd.AddCommand(checkTypesCmd)
}
