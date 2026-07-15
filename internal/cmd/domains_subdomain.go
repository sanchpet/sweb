package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var domainsSubdomainCmd = &cobra.Command{
	Use:   "subdomain",
	Short: "Create or remove subdomains",
}

var domainsSubdomainCreateCmd = &cobra.Command{
	Use:   "create <domain> <machine>",
	Short: "Create a subdomain",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain, machine := args[0], args[1]
		dir, _ := cmd.Flags().GetString("dir")
		if err := c.Domains.CreateSubdomain(cmd.Context(), domain, machine, dir); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Created subdomain %s.%s\n", machine, domain)
		return nil
	},
}

var domainsSubdomainRemoveCmd = &cobra.Command{
	Use:   "remove <domain> <machine>",
	Short: "Remove a subdomain — destructive",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain, machine := args[0], args[1]
		if !confirmed(cmd, fmt.Sprintf("Remove subdomain %s.%s?", machine, domain), "Remove") {
			fmt.Fprintln(cmd.OutOrStdout(), "aborted")
			return nil
		}
		if err := c.Domains.RemoveSubdomain(cmd.Context(), domain, machine); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Removed subdomain %s.%s\n", machine, domain)
		return nil
	},
}

func init() {
	domainsSubdomainCreateCmd.Flags().String("dir", "", "site directory")
	domainsSubdomainRemoveCmd.Flags().Bool("yes", false, "skip the confirmation prompt")
	domainsSubdomainCmd.AddCommand(domainsSubdomainCreateCmd)
	domainsSubdomainCmd.AddCommand(domainsSubdomainRemoveCmd)
	domainsCmd.AddCommand(domainsSubdomainCmd)
}
