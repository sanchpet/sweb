package cmd

import (
	"fmt"

	"github.com/sanchpet/sweb-go-sdk/vh/hosting"
	"github.com/spf13/cobra"
)

// dbPgsqlCmd groups the PostgreSQL database lifecycle under
// `sweb hosting db pgsql`.
var dbPgsqlCmd = &cobra.Command{
	Use:   "pgsql",
	Short: "Manage PostgreSQL databases",
}

var dbPgsqlCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a PostgreSQL database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		password, err := resolveDBPassword(cmd)
		if err != nil {
			return err
		}
		charset, _ := cmd.Flags().GetString("charset")
		comment, _ := cmd.Flags().GetString("comment")
		if err := c.HostingDB.PgsqlCreate(cmd.Context(), hosting.PgsqlCreateOptions{
			Name:     args[0],
			Password: password,
			Charset:  charset,
			Comment:  comment,
		}); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Created PostgreSQL database %s\n", args[0])
		return nil
	},
}

var dbPgsqlDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a PostgreSQL database — destructive",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Delete PostgreSQL database %q? This is irreversible.", args[0]), "Delete") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.HostingDB.PgsqlDelete(cmd.Context(), args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Deleted PostgreSQL database %s\n", args[0])
		return nil
	},
}

var dbPgsqlPasswordCmd = &cobra.Command{
	Use:   "password <name>",
	Short: "Change a PostgreSQL database password",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Change the password of PostgreSQL database %q?", args[0]), "Change") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		password, err := resolveDBPassword(cmd)
		if err != nil {
			return err
		}
		if err := c.HostingDB.PgsqlChangePass(cmd.Context(), args[0], password); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Changed password of PostgreSQL database %s\n", args[0])
		return nil
	},
}

func init() {
	dbPgsqlCreateCmd.Flags().String("password", "", "database password (prompted securely when omitted)")
	dbPgsqlCreateCmd.Flags().String("charset", "", "database charset")
	dbPgsqlCreateCmd.Flags().String("comment", "", "optional comment")

	dbPgsqlDeleteCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	dbPgsqlPasswordCmd.Flags().String("password", "", "new database password (prompted securely when omitted)")
	dbPgsqlPasswordCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	dbPgsqlCmd.AddCommand(dbPgsqlCreateCmd, dbPgsqlDeleteCmd, dbPgsqlPasswordCmd)
	dbCmd.AddCommand(dbPgsqlCmd)
}
