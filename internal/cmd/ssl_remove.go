package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var cloudSSLRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Delete a certificate — destructive",
	Long: `Delete a certificate via the "removeCertificate" method.

<id> is a certificate id (see 'sweb ssl list'). This is DESTRUCTIVE. You are
asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("id must be an integer: %w", err)
		}
		if !confirmed(cmd, fmt.Sprintf("Remove certificate %d? This cannot be undone.", id), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		res, err := c.SSL.RemoveCertificate(cmd.Context(), id)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Removed certificate %d (%s)\n", id, res)
		return nil
	},
}

func init() {
	cloudSSLRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")

	cloudSSLCmd.AddCommand(cloudSSLRemoveCmd)
}
