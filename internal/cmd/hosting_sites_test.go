package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestHostingSitesTree(t *testing.T) {
	var hosting *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "hosting" {
			hosting = c
		}
	}
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}

	var sites *cobra.Command
	for _, c := range hosting.Commands() {
		if c.Name() == "sites" {
			sites = c
		}
	}
	if sites == nil {
		t.Fatal("hosting sites command not registered")
	}

	ssub := subNames(sites)
	for _, n := range []string{"list", "info", "backends", "add", "edit", "remove", "set-domain", "set-backend"} {
		if !ssub[n] {
			t.Errorf("hosting sites is missing subcommand %q", n)
		}
	}
}

// TestHostingSitesRequiredFlags guards that the mutations declaring required
// inputs actually mark them required, so a bare invocation errors rather than
// sending an empty field to the API.
func TestHostingSitesRequiredFlags(t *testing.T) {
	for _, tc := range []struct {
		cmd   *cobra.Command
		flags []string
	}{
		{sitesAddCmd, []string{"alias", "docroot", "domain"}},
		{sitesEditCmd, []string{"alias"}},
		{sitesSetDomainCmd, []string{"docroot"}},
		{sitesSetBackendCmd, []string{"backend"}},
	} {
		for _, name := range tc.flags {
			f := tc.cmd.Flags().Lookup(name)
			if f == nil {
				t.Errorf("%s: flag %q not registered", tc.cmd.Name(), name)
				continue
			}
			if f.Annotations[cobra.BashCompOneRequiredFlag] == nil {
				t.Errorf("%s: flag %q is not marked required", tc.cmd.Name(), name)
			}
		}
	}
}
