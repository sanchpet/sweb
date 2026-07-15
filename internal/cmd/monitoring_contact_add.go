package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var contactAddEmailCmd = &cobra.Command{
	Use:   "add-email <email> <name>",
	Short: "Add an email notification contact",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := c.MonitoringContacts.AddEmail(cmd.Context(), args[0], args[1])
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Added email contact", id)
		return nil
	},
}

var contactAddPhoneCmd = &cobra.Command{
	Use:   "add-phone <phone> <name>",
	Short: "Add a phone notification contact",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := c.MonitoringContacts.AddPhone(cmd.Context(), args[0], args[1])
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Added phone contact", id)
		return nil
	},
}

var contactAddTelegramCmd = &cobra.Command{
	Use:   "add-telegram <name>",
	Short: "Add a Telegram notification contact and request its verification code",
	Long: `Add a Telegram contact (method "addTelegram"), then request its
verification code (method "requestTelegramVerifyCode"). Send the printed code to
the SpaceWeb bot, then confirm with 'monitoring contact verify <id> <code>'.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		id, err := c.MonitoringContacts.AddTelegram(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Added Telegram contact", id)
		code, err := c.MonitoringContacts.RequestTelegramVerifyCode(cmd.Context(), strconv.FormatInt(id, 10))
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(),
			"Send this code to the SpaceWeb bot, then run 'monitoring contact verify %d <code>':\n%s\n", id, code)
		return nil
	},
}

func init() {
	contactCmd.AddCommand(contactAddEmailCmd, contactAddPhoneCmd, contactAddTelegramCmd)
}
