package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func subNames(c *cobra.Command) map[string]bool {
	m := map[string]bool{}
	for _, sub := range c.Commands() {
		m[sub.Name()] = true
	}
	return m
}

func TestCommandTree(t *testing.T) {
	root := subNames(rootCmd)
	for _, n := range []string{"configure", "vps"} {
		if !root[n] {
			t.Errorf("root is missing subcommand %q", n)
		}
	}

	var vps *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "vps" {
			vps = c
		}
	}
	if vps == nil {
		t.Fatal("vps command not registered")
	}
	vsub := subNames(vps)
	for _, n := range []string{"list", "create", "config"} {
		if !vsub[n] {
			t.Errorf("vps is missing subcommand %q", n)
		}
	}
}
