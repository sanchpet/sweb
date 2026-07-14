package cmd

import "github.com/spf13/cobra"

// balancerCmd groups the cloud load-balancer service (endpoint /balancer):
// list the account's balancers, show the order catalog, and the
// create/edit/remove lifecycle. Like the VPS group, it lives on the cloud
// account, so bind it to that profile once and every subcommand inherits it.
var balancerCmd = &cobra.Command{
	Use:   "balancer",
	Short: "Manage cloud load balancers (list, config, create, edit, remove)",
}

func init() {
	rootCmd.AddCommand(balancerCmd)
}
