package cmd

import "testing"

func TestHostingDiskUsageTree(t *testing.T) {
	hosting := findSub(rootCmd, "hosting")
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}
	du := findSub(hosting, "disk-usage")
	if du == nil {
		t.Fatal("hosting disk-usage command not registered")
	}

	dusub := subNames(du)
	for _, n := range []string{"report", "tasks", "scan", "email", "set-email"} {
		if !dusub[n] {
			t.Errorf("hosting disk-usage is missing subcommand %q", n)
		}
	}
}
