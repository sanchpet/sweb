package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/vh/partner"
	"github.com/spf13/cobra"
)

// partnerCmd groups the partner-program service (endpoint /vh/partnerProgram):
// the referral hosting/VPS catalog, becoming a partner and filling requisites,
// the referred-client roster with per-client card / event / finance logs,
// advertising materials, referral-site and link statistics, and reward
// withdrawal. It hangs off `hosting`, so it inherits that group's profile
// binding.
var partnerCmd = &cobra.Command{
	Use:   "partner",
	Short: "Partner-program: clients, stats, payouts",
}

// partnerStatusCmd checks whether a candidate referral login is free to use
// (method "checkLogin") — the pre-flight before placing a referral order.
var partnerStatusCmd = &cobra.Command{
	Use:   "status <login>",
	Short: "Check whether a referral login is available",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		free, err := c.PartnerProgram.CheckLogin(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return render(cmd, map[string]bool{"available": free}, func(w io.Writer) {
			kv(w, "LOGIN", args[0])
			kv(w, "AVAILABLE", onOff(free))
		})
	},
}

// partnerJoinCmd enrolls the account in the partner program (method
// "startPartnership"). It is a one-time account change, so it confirms.
var partnerJoinCmd = &cobra.Command{
	Use:   "join",
	Short: "Enroll the account in the partner program",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, "Enroll this account in the partner program?", "Join") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.PartnerProgram.StartPartnership(cmd.Context()); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Joined the partner program")
		return nil
	},
}

// partnerRequisitesCmd groups the legal-requisites subcommands.
var partnerRequisitesCmd = &cobra.Command{
	Use:   "requisites",
	Short: "Partner legal requisites (INN/SNILS/address)",
}

// partnerRequisitesSetCmd saves the partner's legal identifiers (method
// "fillPartnerRequisites").
var partnerRequisitesSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the partner's legal requisites",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		inn, _ := cmd.Flags().GetString("inn")
		snils, _ := cmd.Flags().GetString("snils")
		addr, _ := cmd.Flags().GetString("address")
		if inn == "" || snils == "" || addr == "" {
			return fmt.Errorf("--inn, --snils and --address are all required")
		}
		if err := c.PartnerProgram.FillRequisites(cmd.Context(), partner.Requisites{
			INN:        inn,
			SNILS:      snils,
			RegAddress: addr,
		}); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Requisites saved")
		return nil
	},
}

// partnerOSConfigCmd prints the VPS ordering catalog (method "vpsOsConfig"):
// datacenters and the selectable OS distributions, used to fill an order.
var partnerOSConfigCmd = &cobra.Command{
	Use:   "os-config",
	Short: "Show the VPS ordering catalog (datacenters, OS list)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		cfg, err := c.PartnerProgram.OSConfig(cmd.Context())
		if err != nil {
			return err
		}
		return render(cmd, cfg, func(w io.Writer) {
			fmt.Fprintln(w, "DATACENTERS")
			fmt.Fprintln(w, "ID\tNAME\tLOCATION")
			for _, dc := range cfg.Datacenters {
				fmt.Fprintf(w, "%s\t%s\t%s\n", dc.ID, dc.Name, dc.Location)
			}
			fmt.Fprintln(w)
			fmt.Fprintln(w, "OPERATING SYSTEMS")
			fmt.Fprintln(w, "ID\tNAME\tPLAN\tPANELS")
			for _, os := range cfg.SelectOS {
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", os.ID, os.Name, os.PlanID, os.PanelType)
			}
		})
	},
}

func init() {
	partnerRequisitesSetCmd.Flags().String("inn", "", "taxpayer identification number (INN)")
	partnerRequisitesSetCmd.Flags().String("snils", "", "individual insurance account number (SNILS)")
	partnerRequisitesSetCmd.Flags().String("address", "", "registration address")
	partnerRequisitesCmd.AddCommand(partnerRequisitesSetCmd)

	partnerJoinCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	partnerCmd.AddCommand(
		partnerStatusCmd,
		partnerJoinCmd,
		partnerRequisitesCmd,
		partnerOSConfigCmd,
	)
	hostingCmd.AddCommand(partnerCmd)
}
