package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var vpsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List VPS instances",
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		list, err := c.VPS.List(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "NAME\tUID\tPLAN\tCPU\tRAM(MB)\tIP\tRUNNING")
			for _, v := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%s\t%t\n",
					v.Name, v.UID, v.PlanName, v.CPU, v.RAM, v.IP, v.IsRunning == 1)
			}
		})
	},
}

func init() {
	vpsCmd.AddCommand(vpsListCmd)
}
