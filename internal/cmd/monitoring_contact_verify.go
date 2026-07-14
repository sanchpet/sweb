package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var contactVerifyCmd = &cobra.Command{
	Use:   "verify <id> <code>",
	Short: "Confirm a contact with its verification code",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.MonitoringContacts.VerifyContact(cmd.Context(), args[0], args[1]); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Verified contact", args[0])
		return nil
	},
}

var contactVerifyStatusCmd = &cobra.Command{
	Use:   "verify-status <id>",
	Short: "Report whether a contact is verified",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		verified, err := c.MonitoringContacts.IsVerified(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return render(cmd, map[string]bool{"verified": verified}, func(w io.Writer) {
			fmt.Fprintf(w, "VERIFIED\t%t\n", verified)
		})
	},
}

func init() {
	contactCmd.AddCommand(contactVerifyCmd, contactVerifyStatusCmd)
}
