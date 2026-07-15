package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var domainsMoveCmd = &cobra.Command{
	Use:   "move <domain>",
	Short: "Add an existing domain to the account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		prolong, _ := cmd.Flags().GetString("prolong")
		dir, _ := cmd.Flags().GetString("dir")
		if err := c.Domains.Move(cmd.Context(), domain, prolong, dir); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Moved %s onto the account\n", domain)
		return nil
	},
}

func init() {
	// move/moveList use "no" (not "none") for the do-not-prolong token.
	domainsMoveCmd.Flags().String("prolong", "no", "auto-prolong mode: no|manual|bonus_money")
	domainsMoveCmd.Flags().String("dir", "", "home directory (optional)")
	domainsCmd.AddCommand(domainsMoveCmd)
}
