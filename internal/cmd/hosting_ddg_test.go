package cmd

import "testing"

func TestHostingDDGCommandTree(t *testing.T) {
	hosting := findSub(rootCmd, "hosting")
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}
	ddg := findSub(hosting, "ddg")
	if ddg == nil {
		t.Fatal("hosting ddg command not registered")
	}

	ddgsub := subNames(ddg)
	for _, n := range []string{"list", "info", "count", "price", "enable", "disable"} {
		if !ddgsub[n] {
			t.Errorf("hosting ddg is missing subcommand %q", n)
		}
	}
}
