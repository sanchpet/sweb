package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// TestCloudSSLCommandTree asserts the cloud/account SSL group (`sweb ssl`,
// SDK /vps/ssl) is registered on the root with its full leaf set — distinct
// from the shared-hosting `hosting ssl` group (SDK /vh/ssl).
func TestCloudSSLCommandTree(t *testing.T) {
	var sslGroup *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "ssl" {
			sslGroup = c
		}
	}
	if sslGroup == nil {
		t.Fatal("ssl command not registered on root")
	}
	sub := subNames(sslGroup)
	for _, n := range []string{"list", "order-list", "download", "prolong-info", "autoprolong", "order", "remove"} {
		if !sub[n] {
			t.Errorf("ssl is missing subcommand %q", n)
		}
	}
}
