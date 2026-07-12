package cmd

import (
	"fmt"

	sweb "github.com/sanchpet/sweb-go-sdk"
	"github.com/spf13/cobra"
)

var domainsRegisterCmd = &cobra.Command{
	Use:   "register <domain>",
	Short: "Register a domain on the account — bills",
	Long: `Register a domain via the "reg" method. This CHARGES the account (money or
bonus points per --pay). You are asked to confirm unless --yes is given.

--person is the domain-person (registrant contact) id; most registrations
require it.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		if !confirmed(cmd, fmt.Sprintf("Register %s? This charges the account.", domain), "Register") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		pay, _ := cmd.Flags().GetString("pay")
		prolong, _ := cmd.Flags().GetString("prolong")
		dir, _ := cmd.Flags().GetString("dir")
		shield, _ := cmd.Flags().GetBool("shield")
		if err := c.Domains.Register(cmd.Context(), sweb.RegisterOptions{
			Domain:      domain,
			PayType:     pay,
			DomPerson:   flagInt(cmd, "person"),
			ProlongType: prolong,
			AutoReg:     flagInt(cmd, "auto"),
			Dir:         dir,
			IDShield:    shield,
		}); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Registered %s\n", domain)
		return nil
	},
}

func init() {
	domainsRegisterCmd.Flags().String("pay", "balance", "payment source: balance|bonus")
	domainsRegisterCmd.Flags().Int("person", 0, "domain-person (registrant contact) id")
	domainsRegisterCmd.Flags().String("prolong", "none", "auto-prolong mode: none|manual|bonus_money")
	domainsRegisterCmd.Flags().String("dir", "", "relative site directory")
	domainsRegisterCmd.Flags().Int("auto", 0, "auto-registration flag")
	domainsRegisterCmd.Flags().Bool("shield", false, "hide WHOIS")
	domainsRegisterCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	domainsCmd.AddCommand(domainsRegisterCmd)
}
