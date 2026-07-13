package cmd

import (
	"fmt"
	"io"

	"github.com/charmbracelet/huh"
	"github.com/sanchpet/sweb-go-sdk/vh/hosting"
	"github.com/spf13/cobra"
)

// dbCmd groups the shared-hosting database services under `sweb hosting db`:
// the account's MySQL/PgSQL databases, their create/delete/password lifecycle,
// MySQL import/copy/comment, remote-access rules, and the PhpMyAdmin user.
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage shared-hosting databases (MySQL/PgSQL)",
}

// resolveDBPassword returns the password for a create/change-password command:
// the --password flag when set, otherwise a masked interactive prompt. An empty
// result (non-interactive with no flag) is an error, since the API requires one.
func resolveDBPassword(cmd *cobra.Command) (string, error) {
	if pw, _ := cmd.Flags().GetString("password"); pw != "" {
		return pw, nil
	}
	var pw string
	if err := huh.NewInput().
		Title("Database password").
		EchoMode(huh.EchoModePassword).
		Value(&pw).
		Run(); err != nil {
		return "", err
	}
	if pw == "" {
		return "", fmt.Errorf("a password is required: pass --password or enter one at the prompt")
	}
	return pw, nil
}

var dbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the account's MySQL and PgSQL databases",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		filter, _ := cmd.Flags().GetString("filter")
		page, _ := cmd.Flags().GetInt("page")
		res, err := c.HostingDB.DatabaseList(cmd.Context(), hosting.ListOptions{
			Page:   page,
			Filter: filter,
		})
		if err != nil {
			return err
		}
		return render(cmd, res, func(w io.Writer) {
			fmt.Fprintln(w, "NAME\tTYPE\tVERSION\tLOGIN\tTABLES\tSIZE(MB)\tCHARSET\tCOMMENT")
			for _, d := range res.List {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%.1f\t%s\t%s\n",
					d.Name, d.Type, d.Version, d.Login,
					int64(d.CountTables), float64(d.SizeTables), d.Charset, d.Comment)
			}
		})
	},
}

var dbPmaUserCmd = &cobra.Command{
	Use:   "pma-user <db>",
	Short: "Provision a temporary PhpMyAdmin login for a database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		u, err := c.HostingDB.GetPmaUser(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return render(cmd, u, func(w io.Writer) {
			fmt.Fprintf(w, "URL\t%s\n", u.URL)
			fmt.Fprintf(w, "DB\t%s\n", u.DB)
			fmt.Fprintf(w, "USER\t%s\n", u.User)
			fmt.Fprintf(w, "PASS\t%s\n", u.Pass)
		})
	},
}

func init() {
	dbListCmd.Flags().String("filter", "", "substring filter on the database name")
	dbListCmd.Flags().Int("page", 0, "page of results (1-based; 0 lets the API default)")

	dbCmd.AddCommand(dbListCmd, dbPmaUserCmd)
	hostingCmd.AddCommand(dbCmd)
}
