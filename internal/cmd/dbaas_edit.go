package cmd

import (
	"fmt"

	"github.com/sanchpet/sweb-go-sdk/dbaas"
	"github.com/spf13/cobra"
)

var dbaasEditCmd = &cobra.Command{
	Use:   "edit <billing-id>",
	Short: "Edit a managed-database cluster (users, plan, name)",
	Long: `Edit a cluster's users, plan or display name via the "editInstance" method.
<billing-id> is from 'sweb dbaas list'. An omitted flag leaves that facet
unchanged.

--user follows the edit semantics: a user with a password is created, and any
existing user absent from the given list is removed; pass no --user to leave the
user set untouched. Each --user is "name:password" and may be repeated.

This MUTATES the cluster. You are asked to confirm unless --yes.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		plan, _ := f.GetInt("plan")
		displayName, _ := f.GetString("name")
		userSpecs, _ := f.GetStringArray("user")

		req := dbaas.EditInstanceRequest{
			BillingID:   args[0],
			PlanID:      plan,
			DisplayName: displayName,
		}
		if len(userSpecs) > 0 {
			users, err := parseDBaaSUsers(userSpecs)
			if err != nil {
				return err
			}
			req.Users = users
		}

		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Edit cluster %q?", args[0]), "Edit") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.DBaaS.EditInstance(cmd.Context(), req); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Edited", args[0])
		return nil
	},
}

func init() {
	f := dbaasEditCmd.Flags()
	f.Int("plan", 0, "new plan ID (planId) — see `sweb dbaas config`")
	f.String("name", "", "new display name for the cluster")
	f.StringArray("user", nil, "cluster user as name:password (repeatable; edit semantics)")
	f.Bool("yes", false, "skip the confirmation prompt")

	dbaasCmd.AddCommand(dbaasEditCmd)
}
