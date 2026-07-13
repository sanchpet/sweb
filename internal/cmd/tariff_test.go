package cmd

import "testing"

// TestTariffCommandTree guards the top-level tariff group and its read-only
// subcommands. It reuses findSub (hosting_db_test.go).
func TestTariffCommandTree(t *testing.T) {
	tariff := findSub(rootCmd, "tariff")
	if tariff == nil {
		t.Fatal("tariff command not registered")
	}
	for _, n := range []string{"show", "server"} {
		if findSub(tariff, n) == nil {
			t.Errorf("tariff is missing subcommand %q", n)
		}
	}
}
