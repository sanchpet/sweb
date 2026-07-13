package cmd

import "github.com/spf13/cobra"

// hostingCmd groups the shared-hosting panel services — mail, databases, sites,
// SSL, backups, cron and the rest. SpaceWeb serves the hosting panel from a
// separate account than the cloud/VPS one, so bind this group to that profile
// once and every subcommand inherits it: `sweb profile bind hosting <profile>`.
var hostingCmd = &cobra.Command{
	Use:   "hosting",
	Short: "Manage shared-hosting services (mail, databases, sites, SSL, …)",
}

func init() {
	rootCmd.AddCommand(hostingCmd)
}
