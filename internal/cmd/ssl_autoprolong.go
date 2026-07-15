package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var cloudSSLAutoprolongCmd = &cobra.Command{
	Use:   "autoprolong <id>",
	Short: "Enable or disable a certificate's auto-prolongation",
	Long: `Toggle a certificate's auto-prolongation via the "editAutoprolong" method.

<id> is a certificate id (see 'sweb ssl list'); pass exactly one of --enable or
--disable.`,
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
		enable, _ := cmd.Flags().GetBool("enable")
		disable, _ := cmd.Flags().GetBool("disable")
		if enable == disable {
			return fmt.Errorf("pass exactly one of --enable or --disable")
		}
		res, err := c.SSL.EditAutoprolong(cmd.Context(), id, enable)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Set auto-prolongation for certificate %d to %s (%s)\n", id, yesNo(enable), res)
		return nil
	},
}

func init() {
	cloudSSLAutoprolongCmd.Flags().Bool("enable", false, "enable auto-prolongation")
	cloudSSLAutoprolongCmd.Flags().Bool("disable", false, "disable auto-prolongation")
	cloudSSLAutoprolongCmd.MarkFlagsMutuallyExclusive("enable", "disable")

	cloudSSLCmd.AddCommand(cloudSSLAutoprolongCmd)
}
