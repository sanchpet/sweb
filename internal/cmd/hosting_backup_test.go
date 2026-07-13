package cmd

import "testing"

func TestHostingBackupCommandTree(t *testing.T) {
	hosting := findSub(rootCmd, "hosting")
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}
	backup := findSub(hosting, "backup")
	if backup == nil {
		t.Fatal("hosting backup command not registered")
	}

	bsub := subNames(backup)
	for _, n := range []string{
		"dates", "files", "mysql", "snapshot",
		"restore-files", "restore-mysql", "receive-files", "receive-mysql", "download",
	} {
		if !bsub[n] {
			t.Errorf("hosting backup is missing subcommand %q", n)
		}
	}
}

// TestHostingBackupConfirmFlags checks that every destructive/mutating backup
// command registers the --yes escape hatch that confirmed() reads.
func TestHostingBackupConfirmFlags(t *testing.T) {
	for _, path := range [][]string{
		{"hosting", "backup", "snapshot"},
		{"hosting", "backup", "restore-files"},
		{"hosting", "backup", "restore-mysql"},
	} {
		cmd, _, err := rootCmd.Find(path)
		if err != nil {
			t.Errorf("%v: not found: %v", path, err)
			continue
		}
		if cmd.Flags().Lookup("yes") == nil {
			t.Errorf("%v: missing --yes flag", path)
		}
	}
}
