package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var checkEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit an existing monitoring check",
	Long: `Update an existing check (method "edit"). The flags mirror
'monitoring check create' minus --type (edit is keyed by id).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := checkID(args[0])
		if err != nil {
			return err
		}
		spec, err := checkSpecFromFlags(cmd)
		if err != nil {
			return err
		}
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.MonitoringChecks.Edit(cmd.Context(), id, spec); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Edited check", id)
		return nil
	},
}

func init() {
	addCheckSpecFlags(checkEditCmd)
	checkCmd.AddCommand(checkEditCmd)
}
