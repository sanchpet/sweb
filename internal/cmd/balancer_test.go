package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestBalancerCommandTree(t *testing.T) {
	if !subNames(rootCmd)["balancer"] {
		t.Fatal("root is missing subcommand \"balancer\"")
	}

	var balancer *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "balancer" {
			balancer = c
		}
	}
	if balancer == nil {
		t.Fatal("balancer command not registered")
	}

	bsub := subNames(balancer)
	for _, n := range []string{"list", "config", "create", "edit", "remove"} {
		if !bsub[n] {
			t.Errorf("balancer is missing subcommand %q", n)
		}
	}
}
