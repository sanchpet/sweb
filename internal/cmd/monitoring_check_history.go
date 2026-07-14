package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/monitoring/checks"
	"github.com/spf13/cobra"
)

var checkHistoryCmd = &cobra.Command{
	Use:   "history <id>",
	Short: "Show a check's event history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := checkID(args[0])
		if err != nil {
			return err
		}
		c, err := client()
		if err != nil {
			return err
		}
		start, _ := cmd.Flags().GetString("start")
		finish, _ := cmd.Flags().GetString("finish")
		res, err := c.MonitoringChecks.History(cmd.Context(), id, &checks.HistoryOptions{
			StartDate:  start,
			FinishDate: finish,
			Page:       flagInt(cmd, "page"),
			PerPage:    flagInt(cmd, "per-page"),
		})
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tCHECK\tTIMESTAMP\tSUCCESS")
			for _, e := range res.List {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.ID, e.CheckID, e.TS, e.Success)
			}
		})
	},
}

func init() {
	checkHistoryCmd.Flags().String("start", "", "inclusive lower bound on the event date")
	checkHistoryCmd.Flags().String("finish", "", "inclusive upper bound on the event date")
	checkHistoryCmd.Flags().Int("page", 0, "page of results (1-based; 0 lets the API default)")
	checkHistoryCmd.Flags().Int("per-page", 0, "rows per page (0 lets the API default)")
	checkCmd.AddCommand(checkHistoryCmd)
}
