package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/flex"
	"github.com/spf13/cobra"
)

// payCmd groups the account-level billing operations (endpoint /pay): balance,
// autopayment/deferment state, payment recommendations and upcoming payments,
// the runway to blocking, and the reserves holding funds. Billing is
// account-level (not per-hosting-service), so this is a top-level group; bind it
// to a profile once with `sweb profile bind pay <profile>`.
var payCmd = &cobra.Command{
	Use:   "pay",
	Short: "Account billing: balance, payments, reserves",
}

// payMoney formats a flex.Float money amount for table output.
func payMoney(f flex.Float) string { return fmt.Sprintf("%.2f", float64(f)) }

var payBalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Show the account balance broken down by pocket",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		b, err := c.Pay.GetBalance(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, b, func(w io.Writer) {
			row := func(k, v string) { fmt.Fprintf(w, "%s\t%s\n", k, v) }
			row("REAL", payMoney(b.RealBalance))
			row("BONUS", payMoney(b.BonusBalance))
			row("CLOUD", payMoney(b.CloudBalance))
			row("OTHER", payMoney(b.OtherBalance))
			row("CREDIT", payMoney(b.CreditBalance))
			row("CREDIT CLOUD", payMoney(b.CreditCloudBalance))
			row("CREDIT OTHER", payMoney(b.CreditOtherBalance))
			row("MULTIPLE", yesNo(b.MultipleBalanceEnabled))
			for id, amt := range b.VATBalance {
				row("VAT["+id+"]", amt)
			}
		})
	},
}

var payIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Billing overview: balance, autopayment, deferment, block countdown",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		a, err := c.Pay.Index(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, a, func(w io.Writer) {
			row := func(k, v string) { fmt.Fprintf(w, "%s\t%s\n", k, v) }
			row("STATUS", a.Status)
			row("REAL BALANCE", payMoney(a.Balance.RealBalance))
			row("BONUS BALANCE", payMoney(a.Balance.BonusBalance))
			row("CREDIT BALANCE", payMoney(a.Balance.CreditBalance))
			row("AUTOPAYMENT", yesNo(int64(a.AutoPaymentEnable) == 1))
			row("BLOCKED MONEY", payMoney(a.BlockedMoney))
			row("DOMAIN BONUSES", fmt.Sprintf("%d", int64(a.DomainBonuses)))
			if a.BlockInfo.DaysWord != "" {
				row("BLOCK IN", fmt.Sprintf("%d days (%s)", int64(a.BlockInfo.Days), a.BlockInfo.DaysWord))
			} else if int64(a.BlockInfo.Days) != 0 {
				row("BLOCK IN", fmt.Sprintf("%d days", int64(a.BlockInfo.Days)))
			}
			if a.BlockInfo.DaysDate != "" {
				row("BLOCK DATE", a.BlockInfo.DaysDate)
			}
			row("DEFERMENT", fmt.Sprintf("%d days", int64(a.Deferment.Value)))
			if a.EdgeDate != "" {
				row("DOCS FROM", a.EdgeDate)
			}
		})
	},
}

var payAutopayCmd = &cobra.Command{
	Use:   "autopay",
	Short: "Report whether autopayment is available on the account",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		on, err := c.Pay.IsAutopaymentEnable(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, map[string]bool{"autopayment_enable": on}, func(w io.Writer) {
			fmt.Fprintf(w, "AUTOPAYMENT\t%s\n", yesNo(on))
		})
	},
}

var payRecommendationsCmd = &cobra.Command{
	Use:   "recommendations",
	Short: "List services recommended for payment",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		balance, _ := cmd.Flags().GetBool("balance")
		r, err := c.Pay.GetPayRecommendations(cmd.Context(), balance)
		if err != nil {
			return err
		}
		return render(cmd, r, func(w io.Writer) {
			fmt.Fprintln(w, "BUCKET\tID\tNAME\tDATE\tCOST")
			for _, rec := range r.RecommendedForPay {
				fmt.Fprintf(w, "pay\t%d\t%s\t%s\t%s\n", int64(rec.ID), rec.Name, rec.Date, payMoney(rec.Cost))
			}
			for _, rec := range r.RecommendedForPayBalance {
				fmt.Fprintf(w, "balance\t%d\t%s\t%s\t%s\n", int64(rec.ID), rec.Name, rec.Date, payMoney(rec.Cost))
			}
		})
	},
}

var payRecommendationCostCmd = &cobra.Command{
	Use:   "recommendation-cost",
	Short: "Show the total cost of the recommended payments",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		cost, err := c.Pay.GetRecommendationTotalCost(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, map[string]float64{"total_cost": float64(cost)}, func(w io.Writer) {
			fmt.Fprintf(w, "TOTAL COST\t%s\n", payMoney(cost))
		})
	},
}

var payUpcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "List upcoming payments for the hosting account",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		ps, err := c.Pay.GetUpcomingPaymentsVh(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, ps, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tNAME\tDATE\tCOST\tTYPE")
			for _, p := range ps {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", int64(p.ID), p.Name, p.Date, payMoney(p.Cost), p.Type)
			}
		})
	},
}

var payRemainsCmd = &cobra.Command{
	Use:   "remains",
	Short: "Show the runway to account blocking (date and days remaining)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		date, err := c.Pay.GetRemainsDate(cmd.Context())
		if err != nil {
			return err
		}
		days, err := c.Pay.GetRemainsDays(cmd.Context())
		if err != nil {
			return err
		}
		out := struct {
			Date string `json:"date"`
			Days int64  `json:"days"`
		}{Date: date, Days: int64(days)}
		return render(cmd, out, func(w io.Writer) {
			fmt.Fprintf(w, "DATE\t%s\n", date)
			fmt.Fprintf(w, "DAYS\t%d\n", int64(days))
		})
	},
}

var payReservesCmd = &cobra.Command{
	Use:   "reserves",
	Short: "List active reserves currently holding funds",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		rs, err := c.Pay.GetActiveReserves(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, rs, func(w io.Writer) {
			fmt.Fprintln(w, "CHARGE\tTYPE\tBALANCE\tTITLE\tEND DATE")
			for _, r := range rs {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", payMoney(r.Charge), r.Type, r.BalanceType, r.Info.Title, r.Info.EndDate)
			}
		})
	},
}

var payDefermentCmd = &cobra.Command{
	Use:   "deferment",
	Short: "Toggle the payment deferment — mutating",
	Long: `Turn the account's payment deferment on or off via the "changeDeferment"
method. This changes billing state on the account; you are asked to confirm
unless --yes is given.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		on, _ := cmd.Flags().GetBool("on")
		verb := "Disable"
		if on {
			verb = "Enable"
		}
		if !confirmed(cmd, fmt.Sprintf("%s payment deferment on this account?", verb), verb) {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Pay.ChangeDeferment(cmd.Context(), on); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Payment deferment %sd\n", verb)
		return nil
	},
}

func init() {
	payRecommendationsCmd.Flags().Bool("balance", false, "also include top-up (balance) recommendations")
	payDefermentCmd.Flags().Bool("on", false, "turn the deferment on (default off)")
	payDefermentCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	payCmd.AddCommand(
		payBalanceCmd,
		payIndexCmd,
		payAutopayCmd,
		payRecommendationsCmd,
		payRecommendationCostCmd,
		payUpcomingCmd,
		payRemainsCmd,
		payReservesCmd,
		payDefermentCmd,
	)
	rootCmd.AddCommand(payCmd)
}
