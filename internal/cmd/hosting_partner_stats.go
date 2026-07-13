package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

// partnerMaterialsCmd groups the advertising-material catalog.
var partnerMaterialsCmd = &cobra.Command{
	Use:   "materials",
	Short: "Advertising banners (types and per-type listing)",
}

var partnerMaterialsTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List the selectable advertising-material types",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		types, err := c.PartnerProgram.AdvertMaterialTypes(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, types, func(w io.Writer) {
			fmt.Fprintln(w, "VALUE\tNAME")
			for _, t := range types {
				fmt.Fprintf(w, "%s\t%s\n", t.Value, t.Name)
			}
		})
	},
}

var partnerMaterialsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List advertising banners (--type filters; all = every type)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		typ, _ := cmd.Flags().GetString("type")
		mats, err := c.PartnerProgram.AdvertMaterials(cmd.Context(), typ)
		if err != nil {
			return err
		}
		return render(cmd, mats, func(w io.Writer) {
			fmt.Fprintln(w, "SIZES\tFILESIZE\tCODE")
			for _, m := range mats {
				fmt.Fprintf(w, "%s\t%s\t%s\n", m.Sizes, m.FileSize, m.Code)
			}
		})
	},
}

// yearMonth resolves the --year/--month flags, defaulting to the current month.
func yearMonth(cmd *cobra.Command) (year, month int) {
	year, _ = cmd.Flags().GetInt("year")
	month, _ = cmd.Flags().GetInt("month")
	now := time.Now()
	if year == 0 {
		year = now.Year()
	}
	if month == 0 {
		month = int(now.Month())
	}
	return year, month
}

var partnerStatsCmd = &cobra.Command{
	Use:   "stats <site>",
	Short: "Per-referral-site daily statistics",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		year, month := yearMonth(cmd)
		stat, err := c.PartnerProgram.GetStatistic(cmd.Context(), args[0], year, month)
		if err != nil {
			return err
		}
		return render(cmd, stat, func(w io.Writer) { statDataTable(w, stat.Data) })
	},
}

var partnerLinksCmd = &cobra.Command{
	Use:   "links",
	Short: "Referral-link daily statistics",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		year, month := yearMonth(cmd)
		stat, err := c.PartnerProgram.GetLinkStatistics(cmd.Context(), year, month)
		if err != nil {
			return err
		}
		return render(cmd, stat, func(w io.Writer) { statDataTable(w, stat.Data) })
	},
}

// statDataTable prints the positional daily-series rows verbatim, tab-joining
// each row's raw JSON cells so both the site and link series render without the
// table having to know their (differing) column layouts.
func statDataTable(w io.Writer, data [][]json.RawMessage) {
	for _, row := range data {
		for i, cell := range row {
			if i > 0 {
				fmt.Fprint(w, "\t")
			}
			fmt.Fprint(w, string(cell))
		}
		fmt.Fprintln(w)
	}
}

func init() {
	partnerMaterialsListCmd.Flags().String("type", "all", "material type value (from `materials types`; all = every type)")
	partnerMaterialsCmd.AddCommand(partnerMaterialsTypesCmd, partnerMaterialsListCmd)

	for _, c := range []*cobra.Command{partnerStatsCmd, partnerLinksCmd} {
		c.Flags().Int("year", 0, "year (defaults to the current year)")
		c.Flags().Int("month", 0, "month 1-12 (defaults to the current month)")
	}

	partnerCmd.AddCommand(partnerMaterialsCmd, partnerStatsCmd, partnerLinksCmd)
}
