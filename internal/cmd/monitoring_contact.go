package cmd

import "github.com/spf13/cobra"

// contactCmd groups the monitoring-contact operations (endpoint
// /monitoring/contacts): listing contacts, adding email/phone/Telegram
// contacts, editing and removing them, and the Telegram verification flow. It
// hangs off the monitoring group.
var contactCmd = &cobra.Command{
	Use:   "contact",
	Short: "Manage monitoring notification contacts",
}

func init() {
	monitoringCmd.AddCommand(contactCmd)
}
