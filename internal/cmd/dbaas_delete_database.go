package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dbaasDeleteDatabaseCmd = &cobra.Command{
	Use:   "delete-database <billing-id> <db-name>",
	Short: "Delete a single database from a cluster — destructive",
	Long: `Delete one database from a cluster via the "deleteDatabase" method.
<billing-id> and <db-name> are from 'sweb dbaas list' (db-name is the technical
Name, not the display name).

This is DESTRUCTIVE. You are asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		billingID, dbName := args[0], args[1]
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Delete database %q from cluster %q? This cannot be undone.", dbName, billingID), "Delete") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if _, err := c.DBaaS.DeleteDatabase(cmd.Context(), billingID, dbName); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Deleted database", dbName, "from", billingID)
		return nil
	},
}

func init() {
	dbaasDeleteDatabaseCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	dbaasCmd.AddCommand(dbaasDeleteDatabaseCmd)
}
