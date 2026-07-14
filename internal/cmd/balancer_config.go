package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

var balancerConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the load-balancer order catalog (plans, protocols)",
	Long: `Lists the plans and front-end protocols used by 'sweb balancer create', and
reports whether ordering a new balancer is currently available (isCreateEnable).`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		cfg, err := c.Balancer.AvailableConfig(cmd.Context())
		if err != nil {
			return err
		}
		enabled, err := c.Balancer.IsCreateEnable(cmd.Context())
		if err != nil {
			return err
		}
		// Carry both the catalog and the availability flag on the json path.
		data := struct {
			CreateEnabled bool `json:"createEnabled"`
			*balancerConfig
		}{enabled, (*balancerConfig)(cfg)}
		return render(cmd, data, func(w io.Writer) {
			fmt.Fprintf(w, "CREATE AVAILABLE\t%t\n\n", enabled)
			fmt.Fprintln(w, "PLANS")
			fmt.Fprintln(w, "ID\tTAG\tTITLE\tPRICE/mo")
			for _, p := range cfg.Plans {
				fmt.Fprintf(w, "%s\t%s\t%s\t%g\n", p.ID, p.Tag, p.Title, float64(p.Price))
			}
			fmt.Fprintln(w, "\nPROTOCOLS (front-end → allowed back-ends)")
			fmt.Fprintln(w, "ID\tNAME\tRESTRICTIONS")
			for _, pr := range cfg.Protocols {
				r := strings.Join(pr.Restrictions, ",")
				if r == "" {
					r = "-"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", pr.ID, pr.Name, r)
			}
		})
	},
}

func init() {
	balancerCmd.AddCommand(balancerConfigCmd)
}
