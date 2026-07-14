package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var contactEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Edit a contact's value and name",
	Long: `Update a contact's value and name (method "editContact"). For a
Telegram contact only --name applies (its value is set by verification).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		value, _ := cmd.Flags().GetString("value")
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name is required")
		}
		c, err := client()
		if err != nil {
			return err
		}
		if err := c.MonitoringContacts.EditContact(cmd.Context(), args[0], value, name); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Edited contact", args[0])
		return nil
	},
}

var contactRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Remove a notification contact — destructive",
	Long: `Delete a monitoring contact (method "deleteContact"). This is
DESTRUCTIVE; you are asked to confirm unless --yes.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		if !confirmed(cmd, fmt.Sprintf("Remove contact %s? This cannot be undone.", args[0]), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.MonitoringContacts.DeleteContact(cmd.Context(), args[0]); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Removed contact", args[0])
		return nil
	},
}

func init() {
	contactEditCmd.Flags().String("value", "", "new contact value (email address or phone)")
	contactEditCmd.Flags().String("name", "", "new display name")
	contactRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	contactCmd.AddCommand(contactEditCmd, contactRemoveCmd)
}
