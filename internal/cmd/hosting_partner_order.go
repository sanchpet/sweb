package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/vh/partner"
	"github.com/spf13/cobra"
)

// partnerOrderCmd groups referral order placement. Every subcommand here places
// a BILLED order for a referred client, so each confirms unless --yes.
var partnerOrderCmd = &cobra.Command{
	Use:   "order",
	Short: "Place a referral order (hosting/VIP/VPS) — bills",
}

// orderResultTable prints the credentials the API returns for a placed order.
func orderResultTable(w io.Writer, r *partner.OrderResult) {
	kv(w, "LOGIN", r.Login)
	kv(w, "PASSWORD", r.Password)
}

// standardOrderFromFlags builds a shared-hosting order from the common flags.
func standardOrderFromFlags(cmd *cobra.Command) partner.StandardOrder {
	tariff, _ := cmd.Flags().GetInt("tariff")
	period, _ := cmd.Flags().GetInt("period")
	return partner.StandardOrder{
		Email:    flagStr(cmd, "email"),
		TariffID: tariff,
		Period:   period,
		Login:    flagStr(cmd, "login"),
		Password: flagStr(cmd, "password"),
	}
}

var partnerOrderVhCmd = &cobra.Command{
	Use:   "vh",
	Short: "Place a standard hosting order for a referred client — bills",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, "Place a standard hosting order? This bills the referred client.", "Order") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.PartnerProgram.CreateStandardOrder(cmd.Context(), standardOrderFromFlags(cmd))
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) { orderResultTable(w, res) })
	},
}

var partnerOrderVipCmd = &cobra.Command{
	Use:   "vip",
	Short: "Place a VIP hosting order for a referred client — bills",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, "Place a VIP hosting order? This bills the referred client.", "Order") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.PartnerProgram.CreateVIPOrder(cmd.Context(), standardOrderFromFlags(cmd))
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) { orderResultTable(w, res) })
	},
}

var partnerOrderVpsCmd = &cobra.Command{
	Use:   "vps",
	Short: "Place a VPS order for a referred client — bills",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, "Place a VPS order? This bills the referred client.", "Order") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		tariff, _ := cmd.Flags().GetInt("tariff")
		distro, _ := cmd.Flags().GetInt("distro")
		period, _ := cmd.Flags().GetInt("period")
		dc, _ := cmd.Flags().GetInt("datacenter")
		res, err := c.PartnerProgram.CreateVPSOrder(cmd.Context(), partner.VPSOrder{
			Email:          flagStr(cmd, "email"),
			TariffID:       tariff,
			DistributiveID: distro,
			Period:         period,
			Login:          flagStr(cmd, "login"),
			Password:       flagStr(cmd, "password"),
			Datacenter:     dc,
		})
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) { orderResultTable(w, res) })
	},
}

// addStandardOrderFlags registers the flags shared by the vh/vip/vps orders.
func addStandardOrderFlags(cmd *cobra.Command) {
	cmd.Flags().String("email", "", "referred client's email")
	cmd.Flags().Int("tariff", 0, "tariff (plan) id, from `partner plans`")
	cmd.Flags().Int("period", 0, "billing period length in months")
	cmd.Flags().String("login", "", "new account login")
	cmd.Flags().String("password", "", "new account password")
	cmd.Flags().Bool("yes", false, "skip the confirmation prompt")
}

func init() {
	addStandardOrderFlags(partnerOrderVhCmd)
	addStandardOrderFlags(partnerOrderVipCmd)
	addStandardOrderFlags(partnerOrderVpsCmd)
	partnerOrderVpsCmd.Flags().Int("distro", 0, "OS distribution id, from `partner os-config`")
	partnerOrderVpsCmd.Flags().Int("datacenter", 0, "datacenter id, from `partner os-config`")

	partnerOrderCmd.AddCommand(partnerOrderVhCmd, partnerOrderVipCmd, partnerOrderVpsCmd)
	partnerCmd.AddCommand(partnerOrderCmd)
}
