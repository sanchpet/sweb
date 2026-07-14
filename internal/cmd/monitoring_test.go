package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestMonitoringCommandTree(t *testing.T) {
	var monitoring *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "monitoring" {
			monitoring = c
		}
	}
	if monitoring == nil {
		t.Fatal("monitoring command not registered")
	}

	// monitoring carries the plan/subscription commands plus the check and
	// contact subgroups.
	msub := subNames(monitoring)
	for _, n := range []string{"plans", "enable", "disable", "change", "check", "contact"} {
		if !msub[n] {
			t.Errorf("monitoring is missing subcommand %q", n)
		}
	}

	// check carries the CRUD + toggle + read commands.
	var check *cobra.Command
	for _, c := range monitoring.Commands() {
		if c.Name() == "check" {
			check = c
		}
	}
	if check == nil {
		t.Fatal("monitoring check command not registered")
	}
	csub := subNames(check)
	for _, n := range []string{"list", "types", "show", "create", "edit", "activate", "deactivate", "remove", "history"} {
		if !csub[n] {
			t.Errorf("monitoring check is missing subcommand %q", n)
		}
	}

	// contact carries the add/edit/remove lifecycle plus verification.
	var contact *cobra.Command
	for _, c := range monitoring.Commands() {
		if c.Name() == "contact" {
			contact = c
		}
	}
	if contact == nil {
		t.Fatal("monitoring contact command not registered")
	}
	ctsub := subNames(contact)
	for _, n := range []string{"list", "add-email", "add-phone", "add-telegram", "edit", "remove", "verify", "verify-status"} {
		if !ctsub[n] {
			t.Errorf("monitoring contact is missing subcommand %q", n)
		}
	}
}
