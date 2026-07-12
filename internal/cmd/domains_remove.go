package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var domainsRemoveCmd = &cobra.Command{
	Use:   "remove <domain>",
	Short: "Remove a domain from the account — destructive",
	Long: `Remove a domain via the "remove" method.

This is DESTRUCTIVE. You are asked to confirm unless --yes is given.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		if !confirmed(cmd, fmt.Sprintf("Remove %s? This cannot be undone.", domain), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Domains.Remove(cmd.Context(), domain); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Removed", domain)
		return nil
	},
}

func init() {
	domainsRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	domainsCmd.AddCommand(domainsRemoveCmd)
}
