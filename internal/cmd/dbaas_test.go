package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestDBaaSCommandTree(t *testing.T) {
	if !subNames(rootCmd)["dbaas"] {
		t.Fatal("root is missing subcommand \"dbaas\"")
	}

	var dbaas *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "dbaas" {
			dbaas = c
		}
	}
	if dbaas == nil {
		t.Fatal("dbaas command not registered")
	}
	sub := subNames(dbaas)
	for _, n := range []string{"list", "config", "create", "edit", "remove", "delete-database"} {
		if !sub[n] {
			t.Errorf("dbaas is missing subcommand %q", n)
		}
	}
}
