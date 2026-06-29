package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var vpsConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the catalog of selectable VPS options (plans, OS, datacenters)",
	Long:  "Lists plans, OS images, and datacenters with the numeric IDs used by `sweb vps create`.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		cfg, err := c.VPS.AvailableConfig(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, cfg, func(w io.Writer) {
			fmt.Fprintln(w, "PLANS")
			fmt.Fprintln(w, "ID\tNAME\tCPU\tRAM\tDISK\tPRICE/mo\tSOLD_OUT")
			for _, p := range cfg.VPSPlans {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s %s\t%d\t%t\n",
					p.ID, p.Name, p.CPUCores, p.RAM, p.VolumeDisk, p.DiskType, p.PricePerMonth, p.SoldOut)
			}
			fmt.Fprintln(w, "\nDATACENTERS")
			fmt.Fprintln(w, "ID\tNAME\tLOCATION")
			for _, d := range cfg.Datacenters {
				fmt.Fprintf(w, "%s\t%s\t%s\n", d.ID, d.Name, d.Location)
			}
			fmt.Fprintln(w, "\nOS IMAGES")
			fmt.Fprintln(w, "DISTR_ID\tNAME\tVERSION")
			for _, o := range cfg.SelectOS {
				fmt.Fprintf(w, "%s\t%s\t%s\n", o.OSDistributionID, o.Name, o.Version)
			}
		})
	},
}

func init() {
	vpsCmd.AddCommand(vpsConfigCmd)
}
