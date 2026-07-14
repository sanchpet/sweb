package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dbaasRemoveCmd = &cobra.Command{
	Use:   "remove <billing-id>",
	Short: "Remove a managed-database cluster — destructive",
	Long: `Remove a cluster via the "removeInstance" method. <billing-id> is from
'sweb dbaas list'. This removes the cluster and all of its databases.

This is DESTRUCTIVE. You are asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Remove cluster %q? This cannot be undone.", args[0]), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if _, err := c.DBaaS.RemoveInstance(cmd.Context(), args[0]); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Removed", args[0])
		return nil
	},
}

func init() {
	dbaasRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	dbaasCmd.AddCommand(dbaasRemoveCmd)
}
