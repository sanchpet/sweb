package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var domainsAutoprolongCmd = &cobra.Command{
	Use:               "autoprolong <domain> <mode>",
	Short:             "Set a domain's auto-prolongation mode (none|manual|bonus_money)",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeProlongModes,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain, mode := args[0], args[1]
		if err := c.Domains.ChangeProlong(cmd.Context(), domain, mode); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Set auto-prolong for %s to %s\n", domain, mode)
		return nil
	},
}

func init() {
	domainsCmd.AddCommand(domainsAutoprolongCmd)
}
