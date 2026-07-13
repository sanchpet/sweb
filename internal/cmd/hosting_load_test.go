package cmd

import "testing"

func TestHostingLoadTree(t *testing.T) {
	hosting := findSub(rootCmd, "hosting")
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}

	loadGrp := findSub(hosting, "load")
	if loadGrp == nil {
		t.Fatal("hosting load command not registered")
	}

	lsub := subNames(loadGrp)
	for _, n := range []string{"periods", "table"} {
		if !lsub[n] {
			t.Errorf("hosting load is missing subcommand %q", n)
		}
	}
}
