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
	for _, n := range []string{"configure", "vps", "token"} {
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
	for _, n := range []string{"list", "create", "config", "delete", "rename", "change-plan", "local-ip", "start", "stop", "reboot", "status", "reinstall", "clone", "logs", "ip"} {
		if !vsub[n] {
			t.Errorf("vps is missing subcommand %q", n)
		}
	}

	// ip carries list/add/remove/move + a ptr subgroup.
	var ipCmd *cobra.Command
	for _, c := range vps.Commands() {
		if c.Name() == "ip" {
			ipCmd = c
		}
	}
	if ipCmd == nil {
		t.Fatal("vps ip command not registered")
	}
	isub := subNames(ipCmd)
	for _, n := range []string{"list", "add", "remove", "move", "ptr"} {
		if !isub[n] {
			t.Errorf("vps ip is missing subcommand %q", n)
		}
	}

	// local-ip carries its own show/add/remove subcommands.
	var localIP *cobra.Command
	for _, c := range vps.Commands() {
		if c.Name() == "local-ip" {
			localIP = c
		}
	}
	if localIP == nil {
		t.Fatal("vps local-ip command not registered")
	}
	lsub := subNames(localIP)
	for _, n := range []string{"show", "add", "remove"} {
		if !lsub[n] {
			t.Errorf("vps local-ip is missing subcommand %q", n)
		}
	}
}
