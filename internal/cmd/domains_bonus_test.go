package cmd

import "testing"

// TestBonusCommandTree asserts the domains bonus subgroup hangs off domains
// with its list/status/buy subcommands. It reuses findSub (hosting_db_test.go).
func TestBonusCommandTree(t *testing.T) {
	domains := findSub(rootCmd, "domains")
	if domains == nil {
		t.Fatal("domains command not registered")
	}
	bonusGrp := findSub(domains, "bonus")
	if bonusGrp == nil {
		t.Fatal("domains bonus command not registered")
	}
	for _, n := range []string{"list", "status", "buy"} {
		if findSub(bonusGrp, n) == nil {
			t.Errorf("domains bonus is missing subcommand %q", n)
		}
	}
}
