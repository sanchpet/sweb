package cmd

import (
	"fmt"
	"io"
	"sort"

	"github.com/spf13/cobra"
)

var dbaasConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the cluster-creation catalog (engines and plans)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		cfg, err := c.DBaaS.AvailableConfig(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, cfg, func(w io.Writer) {
			fmt.Fprintln(w, "ENGINE\tVERSION")
			engines := make([]string, 0, len(cfg.Engines))
			for name := range cfg.Engines {
				engines = append(engines, name)
			}
			sort.Strings(engines)
			for _, name := range engines {
				for _, e := range cfg.Engines[name] {
					fmt.Fprintf(w, "%s\t%s\n", name, e.Version)
				}
			}
			fmt.Fprintln(w, "\nPLAN_ID\tNAME\tCPU\tMEM_GB\tSTORAGE_GB")
			for _, p := range cfg.Plans {
				fmt.Fprintf(w, "%d\t%s\t%d\t%d\t%d\n",
					int64(p.ID), p.Name, int64(p.CPU), int64(p.Memory), int64(p.Storage))
			}
		})
	},
}

func init() {
	dbaasCmd.AddCommand(dbaasConfigCmd)
}
