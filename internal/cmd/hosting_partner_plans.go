package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/vh/partner"
	"github.com/spf13/cobra"
)

// partnerPlansCmd groups the referral hosting-plan catalog (standard/VIP).
var partnerPlansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Referral hosting-plan catalog (standard/VIP)",
}

// plansTable renders a plan list: one row per plan showing its resource limits,
// then a nested price line per billing period.
func plansTable(w io.Writer, plans []partner.Plan) {
	fmt.Fprintln(w, "ID\tNAME\tDISK(GB)\tSITES\tDB\tFTP\tMAIL")
	for _, p := range plans {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%d\t%d\n",
			p.ID, p.Name, int64(p.Disk), int64(p.Sites),
			int64(p.DBCount), int64(p.FTPCount), int64(p.MailCount))
		for _, per := range p.Period {
			fmt.Fprintf(w, "  %dmo\t%d руб\tSSL:%s\tdomains:%d\t%s\t\n",
				int64(per.Length), int64(per.Price), onOffInt(per.SSL),
				int64(per.Domain), per.DomainZone)
		}
	}
}

// onOffInt renders a 0/1 flag int as off/on.
func onOffInt[T ~int64](v T) string {
	if v != 0 {
		return "on"
	}
	return "off"
}

// partnerPlansFor builds the RunE for a plans subcommand over the selected SDK
// call (StandardPlans/VipPlans).
func partnerPlansFor(sel func(*partner.Service) func(cmd *cobra.Command) ([]partner.Plan, error)) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		plans, err := sel(c.PartnerProgram)(cmd)
		if err != nil {
			return err
		}
		return render(cmd, plans, func(w io.Writer) { plansTable(w, plans) })
	}
}

var partnerPlansStandardCmd = &cobra.Command{
	Use:   "standard",
	Short: "List standard referral hosting plans",
	Args:  cobra.NoArgs,
	RunE: partnerPlansFor(func(s *partner.Service) func(*cobra.Command) ([]partner.Plan, error) {
		return func(cmd *cobra.Command) ([]partner.Plan, error) { return s.StandardPlans(cmd.Context()) }
	}),
}

var partnerPlansVipCmd = &cobra.Command{
	Use:   "vip",
	Short: "List VIP referral hosting plans",
	Args:  cobra.NoArgs,
	RunE: partnerPlansFor(func(s *partner.Service) func(*cobra.Command) ([]partner.Plan, error) {
		return func(cmd *cobra.Command) ([]partner.Plan, error) { return s.VipPlans(cmd.Context()) }
	}),
}

func init() {
	partnerPlansCmd.AddCommand(partnerPlansStandardCmd, partnerPlansVipCmd)
	partnerCmd.AddCommand(partnerPlansCmd)
}
