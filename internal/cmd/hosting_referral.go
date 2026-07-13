package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/sanchpet/sweb-go-sdk/vh/referral"
	"github.com/spf13/cobra"
)

// referralCmd groups the referral-program operations (SDK /vh/referralProgram):
// the read side (list) plus the add/confirm/remove referral-site lifecycle. It
// hangs off the hosting parent, so it inherits that group's profile binding.
var referralCmd = &cobra.Command{
	Use:   "referral",
	Short: "Manage referral sites",
}

var referralListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's referral sites",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		list, err := c.ReferralProgram.List(cmd.Context(), &referral.ListOptions{
			Page:  flagInt(cmd, "page"),
			Limit: flagInt(cmd, "limit"),
		})
		if err != nil {
			return err
		}
		return render(cmd, list, func(w io.Writer) {
			fmt.Fprintln(w, "ID\tDOMAIN\tCONFIRMED\tCLIENTS\tCREATED")
			for _, s := range list.List {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
					s.ID, s.Domain, yesNo(s.Confirmed), int64(s.ClientsCount), s.Created)
			}
		})
	},
}

var referralAddCmd = &cobra.Command{
	Use:   "add <site>",
	Short: "Register a referral site",
	Long: `Register a new referral site via the "addReferralSite" method.

<site> is the domain to register (e.g. example.com).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		if err := c.ReferralProgram.Add(cmd.Context(), domain); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Added referral site:", domain)
		return nil
	},
}

var referralConfirmCmd = &cobra.Command{
	Use:   "confirm <site>",
	Short: "Confirm ownership of a referral site",
	Long: `Confirm ownership of a referral site via the "confirmReferralSite" method.

<site> is the numeric site id from 'sweb hosting referral list' (the ID column).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := referralSiteID(args[0])
		if err != nil {
			return err
		}
		if err := c.ReferralProgram.Confirm(cmd.Context(), id); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Confirmed referral site:", args[0])
		return nil
	},
}

var referralRemoveCmd = &cobra.Command{
	Use:   "remove <site>",
	Short: "Remove a referral site — destructive",
	Long: `Remove a referral site via the "removeReferralSite" method.

<site> is the numeric site id from 'sweb hosting referral list' (the ID column).
This is DESTRUCTIVE; you are asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := referralSiteID(args[0])
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Remove referral site %q? This cannot be undone.", args[0]), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.ReferralProgram.Remove(cmd.Context(), id); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Removed referral site:", args[0])
		return nil
	},
}

// referralSiteID parses the numeric referral-site id that confirm/remove take
// (a Site.ID from 'referral list'), surfacing a friendly error on non-numerics.
func referralSiteID(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return 0, fmt.Errorf("site id must be an integer (from 'sweb hosting referral list'): %w", err)
	}
	return id, nil
}

func init() {
	referralListCmd.Flags().Int("page", 0, "1-based page number")
	referralListCmd.Flags().Int("limit", 0, "records per page")

	referralRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	referralCmd.AddCommand(
		referralListCmd,
		referralAddCmd,
		referralConfirmCmd,
		referralRemoveCmd,
	)
	hostingCmd.AddCommand(referralCmd)
}
