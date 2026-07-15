package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/vh/partner"
	"github.com/spf13/cobra"
)

// partnerWithdrawCmd groups reward-withdrawal: the saved requisites/balance and
// placing a payout order.
var partnerWithdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Reward withdrawal (requisites, balance, payout order)",
}

var partnerWithdrawRequisitesCmd = &cobra.Command{
	Use:   "requisites",
	Short: "Show the balance, payout methods and saved bank requisites",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		req, err := c.PartnerProgram.WithdrawalRequisites(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, req, func(w io.Writer) {
			fmt.Fprintf(w, "BALANCE\t%g\n", float64(req.Balance))
			fmt.Fprintln(w)
			fmt.Fprintln(w, "PAYOUT METHODS")
			fmt.Fprintln(w, "TYPE\tNAME\tENABLED\tMIN\tMONTH-MAX")
			for _, way := range req.OrderTypes {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\n",
					way.Type, way.Name, onOff(way.Enable),
					int64(way.MinimumAmount), int64(way.MaximumMonthAmount))
			}
		})
	},
}

// partnerWithdrawSendCmd places a payout order (method "sendWithdrawalOrder").
// It moves money, so it confirms unless --yes.
var partnerWithdrawSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Place a reward-payout order — moves money",
	Long: `Place a reward-payout order via "sendWithdrawalOrder". --type is a payout
order type from 'partner withdraw requisites'; the bank --req-* fields are
required only for a bank-account payout (type 1).

This MOVES MONEY. You are asked to confirm unless --yes is given.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		orderType, _ := cmd.Flags().GetInt("type")
		amount, _ := cmd.Flags().GetFloat64("amount")
		if amount <= 0 {
			return fmt.Errorf("--amount must be greater than 0")
		}
		if !confirmed(cmd, fmt.Sprintf("Withdraw %g via order type %d?", amount, orderType), "Withdraw") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.PartnerProgram.SendWithdrawalOrder(cmd.Context(), partner.WithdrawalOrder{
			OrderType:       orderType,
			CountMoney:      amount,
			ReqUserName:     flagStr(cmd, "req-name"),
			ReqPayPurpose:   flagStr(cmd, "req-purpose"),
			ReqCheckAccount: flagStr(cmd, "req-account"),
			ReqBankName:     flagStr(cmd, "req-bank"),
			ReqBIC:          flagStr(cmd, "req-bic"),
			ReqCorrAccount:  flagStr(cmd, "req-corr"),
		}); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Withdrawal order placed")
		return nil
	},
}

// flagStr reads a string flag, ignoring the (never-set) lookup error.
func flagStr(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetString(name)
	return v
}

func init() {
	partnerWithdrawSendCmd.Flags().Int("type", 0, "payout order type (from `withdraw requisites`)")
	partnerWithdrawSendCmd.Flags().Float64("amount", 0, "amount to withdraw")
	partnerWithdrawSendCmd.Flags().String("req-name", "", "account holder name (bank payout)")
	partnerWithdrawSendCmd.Flags().String("req-purpose", "", "payment purpose (bank payout)")
	partnerWithdrawSendCmd.Flags().String("req-account", "", "settlement account (bank payout)")
	partnerWithdrawSendCmd.Flags().String("req-bank", "", "bank name (bank payout)")
	partnerWithdrawSendCmd.Flags().String("req-bic", "", "bank BIC (bank payout)")
	partnerWithdrawSendCmd.Flags().String("req-corr", "", "correspondent account (bank payout)")
	partnerWithdrawSendCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	partnerWithdrawCmd.AddCommand(partnerWithdrawRequisitesCmd, partnerWithdrawSendCmd)
	partnerCmd.AddCommand(partnerWithdrawCmd)
}
