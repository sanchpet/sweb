package cmd

import "github.com/spf13/cobra"

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "Manage domains and subdomains",
}

// yesNo renders a boolean for table output.
func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// prolongModes are the auto-prolongation values changeProlong accepts.
var prolongModes = []string{"none", "manual", "bonus_money"}

func completeProlongModes(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 1 { // complete only the <mode> arg (after <domain>)
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return prolongModes, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(domainsCmd)
}
