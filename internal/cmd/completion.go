package cmd

import (
	"fmt"

	sweb "github.com/sanchpet/sweb-go-sdk"
	"github.com/spf13/cobra"
)

// completeOutput offers the static -o values.
func completeOutput(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"table", "json"}, cobra.ShellCompDirectiveNoFileComp
}

// completePlans / completeDatacenters / completeDistributives offer dynamic
// values for `vps create` flags, fetched from the live catalog
// (getAvailableConfig). They degrade to no suggestions when there is no token
// or the API call fails — completion must never error out the shell.
//
// Each suggestion is "value\tdescription"; the shell shows the description.

func completePlans(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	cfg, ok := completionConfig(cmd)
	if !ok {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var out []string
	for _, p := range cfg.VPSPlans {
		if p.SoldOut {
			continue
		}
		out = append(out, fmt.Sprintf("%d\t%s", p.ID, p.Name))
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

func completeDatacenters(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	cfg, ok := completionConfig(cmd)
	if !ok {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var out []string
	for _, d := range cfg.Datacenters {
		out = append(out, fmt.Sprintf("%s\t%s (%s)", d.ID, d.Name, d.Location))
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

func completeDistributives(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	cfg, ok := completionConfig(cmd)
	if !ok {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var out []string
	for _, o := range cfg.SelectOS {
		out = append(out, fmt.Sprintf("%s\t%s %s", o.OSDistributionID, o.Name, o.Version))
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

func completeCategories(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	cfg, ok := completionConfig(cmd)
	if !ok {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var out []string
	for _, ct := range cfg.Categories {
		out = append(out, fmt.Sprintf("%s\t%s", ct.ID, ct.Name))
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

// completeBillingIDs offers existing VPS billing IDs (for `vps delete`), each
// labelled with the VPS name. Degrades to no suggestions on any failure.
func completeBillingIDs(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	c, err := client()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	list, err := c.VPS.List(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var out []string
	for _, v := range list {
		if v.BillingID == "" {
			continue
		}
		out = append(out, fmt.Sprintf("%s\t%s", v.BillingID, v.Name))
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

// completionConfig builds a client and fetches the catalog for a completion
// callback, returning ok=false on any failure (no token, network, etc.) so the
// shell falls back to no suggestions rather than an error.
func completionConfig(cmd *cobra.Command) (*sweb.AvailableConfig, bool) {
	c, err := client()
	if err != nil {
		return nil, false
	}
	ac, err := c.VPS.AvailableConfig(cmd.Context())
	if err != nil {
		return nil, false
	}
	return ac, true
}
