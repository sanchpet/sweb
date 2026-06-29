package cmd

import "github.com/spf13/cobra"

var vpsCmd = &cobra.Command{
	Use:   "vps",
	Short: "Manage VPS instances",
}

func init() {
	rootCmd.AddCommand(vpsCmd)
}
