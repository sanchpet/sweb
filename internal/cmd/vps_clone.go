package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var vpsCloneCmd = &cobra.Command{
	Use:   "clone <vps> [plan-id]",
	Short: "Clone a VPS into a new one — bills",
	Long: `Clone a VPS via the "copy" method — provisions a NEW, billed VPS from the
source. <vps> is the source VPS name (alias) or its billing ID (login_vps_N).

[plan-id] is the new VPS's plan (see 'sweb vps config'); if omitted, the source's
current plan is used. This creates a billed resource — you are asked to confirm
unless --yes. Provisioning is asynchronous; find the new node with 'sweb vps list'.`,
	Args:              cobra.RangeArgs(1, 2),
	ValidArgsFunction: completeBillingIDs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		billingID, err := resolveVPS(cmd.Context(), c, args[0])
		if err != nil {
			return err
		}

		var planID int
		if len(args) == 2 {
			if planID, err = strconv.Atoi(args[1]); err != nil {
				return fmt.Errorf("plan-id must be an integer: %w", err)
			}
		} else {
			// Default to the source's current plan — look it up from the listing.
			list, lerr := c.VPS.List(cmd.Context())
			if lerr != nil {
				return lerr
			}
			for i := range list {
				if list[i].BillingID == billingID {
					planID = int(list[i].PlanID)
					break
				}
			}
			if planID == 0 {
				return fmt.Errorf("could not read the source plan of %s — pass an explicit [plan-id]", billingID)
			}
		}

		if yes, _ := cmd.Flags().GetBool("yes"); !yes {
			confirm := false
			if err := huh.NewConfirm().
				Title(fmt.Sprintf("Clone %q onto plan %d? This provisions a new, billed VPS.", billingID, planID)).
				Affirmative("Clone").Negative("Cancel").Value(&confirm).Run(); err != nil {
				return err
			}
			if !confirm {
				fmt.Fprintln(cmd.OutOrStdout(), "aborted")
				return nil
			}
		}

		if _, err := c.VPS.Copy(cmd.Context(), billingID, planID); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Cloning %s onto plan %d (provisioning asynchronously — see 'sweb vps list')\n", billingID, planID)
		return nil
	},
}

func init() {
	vpsCloneCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	vpsCmd.AddCommand(vpsCloneCmd)
}
