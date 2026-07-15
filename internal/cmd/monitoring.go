package cmd

import "github.com/spf13/cobra"

// monitoringCmd groups the uptime-monitoring services. Three SDK services back
// it: the monitoring subscription/tariff (endpoint /monitoring), the checks
// (`monitoring check`, endpoint /monitoring/checks), and the notification
// contacts (`monitoring contact`, endpoint /monitoring/contacts).
var monitoringCmd = &cobra.Command{
	Use:   "monitoring",
	Short: "Manage uptime monitoring (plans, checks, contacts)",
}

func init() {
	rootCmd.AddCommand(monitoringCmd)
}
