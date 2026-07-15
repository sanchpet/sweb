package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// partnerClientsCmd groups the referred-client roster and per-client views.
var partnerClientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "Referred-client roster, cards, comments and logs",
}

var partnerClientsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List referred clients",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		status, _ := cmd.Flags().GetInt("status")
		page, _ := cmd.Flags().GetInt("page")
		res, err := c.PartnerProgram.ClientsList(cmd.Context(), status, page)
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tLOGIN\tPLAN\tSTATUS\tREGISTERED\tPAYS(ALL)\tPAYS(MONTH)")
			for _, cl := range res.List {
				fmt.Fprintf(w, "%s\t%s\t%v\t%d\t%s\t%d\t%d\n",
					cl.ID, cl.CustLogin, cl.Plan, int64(cl.Status), cl.TS,
					int64(cl.PaysAll), int64(cl.PaysMonth))
			}
		})
	},
}

var partnerClientsCardCmd = &cobra.Command{
	Use:   "card <client-id>",
	Short: "Show a referred client's detailed card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		card, err := c.PartnerProgram.ClientCard(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return render(cmd, card, func(w io.Writer) {
			kv(w, "ID", card.ID)
			kv(w, "LOGIN", card.Login)
			fmt.Fprintf(w, "PLAN\t%v\n", card.PlanName)
			fmt.Fprintf(w, "STATUS\t%d\n", int64(card.Status))
			kv(w, "ATTRACTION", emptyDash(card.Attraction))
			kv(w, "COMMENT", emptyDash(card.Comment))
			kv(w, "CONTRACT", emptyDash(card.ContractNumber))
			kv(w, "REGISTERED", card.RegDate)
			fmt.Fprintf(w, "AMOUNT(PERIOD)\t%d\n", int64(card.AmountsPeriod))
			fmt.Fprintf(w, "AMOUNT(LAST MONTH)\t%d\n", int64(card.AmountsLastMonth))
		})
	},
}

var partnerClientsCommentCmd = &cobra.Command{
	Use:   "comment <client-id> <comment>",
	Short: "Set the partner's note on a referred client",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.PartnerProgram.SaveClientComment(cmd.Context(), args[0], args[1]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Comment saved on %s\n", args[0])
		return nil
	},
}

// partnerClientsLogCmd groups the paginated client event/finance logs.
var partnerClientsLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Referred-client event and finance logs",
}

var partnerClientsLogEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Show the client-event log",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		page, _ := cmd.Flags().GetInt("page")
		log, err := c.PartnerProgram.ClientLogEvents(cmd.Context(), page)
		if err != nil {
			return err
		}
		return render(cmd, log, func(w io.Writer) {
			fmt.Fprintln(w, "DATE\tEVENT")
			for _, e := range log.List {
				fmt.Fprintf(w, "%s\t%s\n", e.TS, e.EventName)
			}
		})
	},
}

var partnerClientsLogFinanceCmd = &cobra.Command{
	Use:   "finance",
	Short: "Show the client-finance log",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		page, _ := cmd.Flags().GetInt("page")
		log, err := c.PartnerProgram.ClientLogFinance(cmd.Context(), page)
		if err != nil {
			return err
		}
		return render(cmd, log, func(w io.Writer) {
			fmt.Fprintln(w, "DATE\tEVENT\tPAYMENT\tWITHDRAWAL\tLOCK")
			for _, e := range log.List {
				fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\n",
					e.TS, e.EventName, int64(e.Payment), int64(e.Withdrawal), int64(e.Lock))
			}
		})
	},
}

func init() {
	partnerClientsListCmd.Flags().Int("status", 0, "filter by client status (0 = all)")
	partnerClientsListCmd.Flags().Int("page", 0, "page of results (1-based; 0 lets the API default)")

	partnerClientsLogEventsCmd.Flags().Int("page", 0, "page of results (1-based; 0 lets the API default)")
	partnerClientsLogFinanceCmd.Flags().Int("page", 0, "page of results (1-based; 0 lets the API default)")
	partnerClientsLogCmd.AddCommand(partnerClientsLogEventsCmd, partnerClientsLogFinanceCmd)

	partnerClientsCmd.AddCommand(
		partnerClientsListCmd,
		partnerClientsCardCmd,
		partnerClientsCommentCmd,
		partnerClientsLogCmd,
	)
	partnerCmd.AddCommand(partnerClientsCmd)
}
