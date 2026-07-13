package cmd

import (
	"fmt"
	"io"

	"github.com/sanchpet/sweb-go-sdk/vh/hosting"
	"github.com/spf13/cobra"
)

// dbMysqlCmd groups the MySQL database lifecycle under `sweb hosting db mysql`.
var dbMysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Manage MySQL databases",
}

var dbMysqlCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a MySQL database",
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
		comment, _ := cmd.Flags().GetString("comment")
		version, _ := cmd.Flags().GetString("version")
		if err := c.HostingDB.MysqlCreate(cmd.Context(), hosting.MysqlCreateOptions{
			Name:     args[0],
			Password: password,
			Comment:  comment,
			Version:  version,
		}); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Created MySQL database %s\n", args[0])
		return nil
	},
}

var dbMysqlDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a MySQL database — destructive",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Delete MySQL database %q? This is irreversible.", args[0]), "Delete") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.HostingDB.MysqlDelete(cmd.Context(), args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Deleted MySQL database %s\n", args[0])
		return nil
	},
}

var dbMysqlPasswordCmd = &cobra.Command{
	Use:   "password <name>",
	Short: "Change a MySQL database password",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Change the password of MySQL database %q?", args[0]), "Change") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		password, err := resolveDBPassword(cmd)
		if err != nil {
			return err
		}
		if err := c.HostingDB.MysqlChangePass(cmd.Context(), args[0], password); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Changed password of MySQL database %s\n", args[0])
		return nil
	},
}

var dbMysqlImportCmd = &cobra.Command{
	Use:   "import <name> <file>",
	Short: "Import a MySQL database from a file in your home directory — destructive",
	Long:  "Import a MySQL database from a dump at <file> (a path in your hosting home directory) via the \"databaseMysqlImport\" method. DESTRUCTIVE: overwrites the target database. Confirms unless --yes.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Import into MySQL database %q from %q? This overwrites it.", args[0], args[1]), "Import") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.HostingDB.MysqlImport(cmd.Context(), args[0], args[1]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Imported %s into MySQL database %s\n", args[1], args[0])
		return nil
	},
}

var dbMysqlCopyCmd = &cobra.Command{
	Use:   "copy <name>",
	Short: "Queue an archive (backup copy) of a MySQL database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.HostingDB.MysqlMakeCopy(cmd.Context(), args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Queued a copy of MySQL database %s\n", args[0])
		return nil
	},
}

var dbMysqlCommentCmd = &cobra.Command{
	Use:   "comment <name> <comment>",
	Short: "Set the comment on a MySQL database",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.HostingDB.EditComment(cmd.Context(), "mysql", args[0], args[1]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Set comment on MySQL database %s\n", args[0])
		return nil
	},
}

// dbMysqlAccessCmd groups the MySQL remote-access rules under
// `sweb hosting db mysql access`.
var dbMysqlAccessCmd = &cobra.Command{
	Use:   "access",
	Short: "Manage a MySQL database's remote-access rules",
}

var dbMysqlAccessListCmd = &cobra.Command{
	Use:   "list [db]",
	Short: "List the remote-access rules for a MySQL database",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		var dbName string
		if len(args) == 1 {
			dbName = args[0]
		}
		res, err := c.HostingDB.MysqlAccessList(cmd.Context(), dbName)
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintln(w, "RULE")
			for _, r := range res.List {
				fmt.Fprintf(w, "%s\n", r)
			}
		})
	},
}

var dbMysqlAccessGrantCmd = &cobra.Command{
	Use:   "grant <db> <rule>",
	Short: "Add a remote-access rule to a MySQL database",
	Long:  "Add a remote-access rule (a \"localhost\", an IP, or a subnet) to a MySQL database.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.HostingDB.MysqlAccessCreate(cmd.Context(), args[0], args[1]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Granted %s access to MySQL database %s\n", args[1], args[0])
		return nil
	},
}

var dbMysqlAccessRevokeCmd = &cobra.Command{
	Use:   "revoke <db> <rule>",
	Short: "Remove a remote-access rule from a MySQL database",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Revoke %q access from MySQL database %q?", args[1], args[0]), "Revoke") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.HostingDB.MysqlAccessDelete(cmd.Context(), args[0], args[1]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Revoked %s access from MySQL database %s\n", args[1], args[0])
		return nil
	},
}

func init() {
	dbMysqlCreateCmd.Flags().String("password", "", "database password (prompted securely when omitted)")
	dbMysqlCreateCmd.Flags().String("comment", "", "optional comment")
	dbMysqlCreateCmd.Flags().String("version", "", "MySQL version (API default when omitted)")

	dbMysqlDeleteCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	dbMysqlPasswordCmd.Flags().String("password", "", "new database password (prompted securely when omitted)")
	dbMysqlPasswordCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	dbMysqlImportCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	dbMysqlAccessRevokeCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	dbMysqlAccessCmd.AddCommand(dbMysqlAccessListCmd, dbMysqlAccessGrantCmd, dbMysqlAccessRevokeCmd)
	dbMysqlCmd.AddCommand(dbMysqlCreateCmd, dbMysqlDeleteCmd, dbMysqlPasswordCmd,
		dbMysqlImportCmd, dbMysqlCopyCmd, dbMysqlCommentCmd, dbMysqlAccessCmd)
	dbCmd.AddCommand(dbMysqlCmd)
}
