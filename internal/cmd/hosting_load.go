package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// loadCmd groups the shared-hosting server-load statistics (SDK /vh/load): the
// available reporting periods and a period's per-day load table. Both sides are
// read-only. It hangs off the hosting parent, so it inherits that group's
// profile binding.
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Shared-hosting server-load stats",
}

var loadPeriodsCmd = &cobra.Command{
	Use:   "periods",
	Short: "List the months for which load statistics exist",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		periods, err := c.HostingLoad.Periods(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, periods, func(w io.Writer) {
			fmt.Fprintln(w, "YEAR\tMONTH")
			for _, p := range periods {
				fmt.Fprintf(w, "%s\t%s\n", p.Year, p.Month)
			}
		})
	},
}

var loadTableCmd = &cobra.Command{
	Use:   "table",
	Short: "Show a period's per-day server-load table",
	Long: `Show a period's load table via the "getLoadTable" method.

--year, --month and --type are all optional: an omitted --year/--month lets the
API pick its default period, and an omitted --type returns every kind. --type
filters by load kind (cpu or mysql).`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		loadType, _ := cmd.Flags().GetString("type")
		tbl, err := c.HostingLoad.LoadTable(cmd.Context(), flagInt(cmd, "year"), flagInt(cmd, "month"), loadType)
		if err != nil {
			return err
		}
		return render(cmd, tbl, func(w io.Writer) {
			fmt.Fprintln(w, "DATE\tCPU\tMYSQL")
			for _, d := range tbl.List {
				fmt.Fprintf(w, "%s\t%g\t%d\n", d.Date, float64(d.CPU), int64(d.Mysql))
			}
		})
	},
}

func init() {
	loadTableCmd.Flags().Int("year", 0, "reporting year (0 lets the API pick its default period)")
	loadTableCmd.Flags().Int("month", 0, "reporting month 1-12 (0 lets the API pick its default period)")
	loadTableCmd.Flags().String("type", "", "load kind to filter by: cpu|mysql (empty returns every kind)")
	_ = loadTableCmd.RegisterFlagCompletionFunc("type", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"cpu", "mysql"}, cobra.ShellCompDirectiveNoFileComp
	})

	loadCmd.AddCommand(
		loadPeriodsCmd,
		loadTableCmd,
	)
	hostingCmd.AddCommand(loadCmd)
}
