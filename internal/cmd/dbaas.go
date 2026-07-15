package cmd

import "github.com/spf13/cobra"

// dbaasCmd groups the managed-database (DBaaS) service — clusters, the
// create-page catalog, and cluster/database lifecycle (endpoint /dbaas).
var dbaasCmd = &cobra.Command{
	Use:   "dbaas",
	Short: "Manage managed databases (DBaaS clusters)",
}

func init() {
	rootCmd.AddCommand(dbaasCmd)
}
