package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var domainsRedirectCmd = &cobra.Command{
	Use:   "redirect <domain>",
	Short: "Show or set a domain's redirect URL",
	Long: `Show a domain's redirect URL, or set it with --url.

	sweb domains redirect example.com                       # show
	sweb domains redirect example.com --url https://x.org   # set`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client()
		if err != nil {
			return err
		}
		domain := args[0]
		if cmd.Flags().Changed("url") {
			url, _ := cmd.Flags().GetString("url")
			if err := c.Domains.SetRedirect(cmd.Context(), domain, url); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Set redirect for %s -> %s\n", domain, url)
			return nil
		}
		cur, err := c.Domains.Redirect(cmd.Context(), domain)
		if err != nil {
			return err
		}
		if cur == "" {
			fmt.Fprintf(cmd.OutOrStdout(), "%s has no redirect\n", domain)
			return nil
		}
		fmt.Fprintln(cmd.OutOrStdout(), cur)
		return nil
	},
}

func init() {
	domainsRedirectCmd.Flags().String("url", "", "set the redirect URL (empty to clear)")
	domainsCmd.AddCommand(domainsRedirectCmd)
}
