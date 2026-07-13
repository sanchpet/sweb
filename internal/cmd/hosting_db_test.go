package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// findSub returns a named direct child of c, or nil.
func findSub(c *cobra.Command, name string) *cobra.Command {
	for _, sub := range c.Commands() {
		if sub.Name() == name {
			return sub
		}
	}
	return nil
}

func TestHostingDBCommandTree(t *testing.T) {
	hosting := findSub(rootCmd, "hosting")
	if hosting == nil {
		t.Fatal("hosting command not registered")
	}
	db := findSub(hosting, "db")
	if db == nil {
		t.Fatal("hosting db command not registered")
	}

	dbsub := subNames(db)
	for _, n := range []string{"list", "pma-user", "mysql", "pgsql"} {
		if !dbsub[n] {
			t.Errorf("hosting db is missing subcommand %q", n)
		}
	}

	mysql := findSub(db, "mysql")
	if mysql == nil {
		t.Fatal("hosting db mysql command not registered")
	}
	msub := subNames(mysql)
	for _, n := range []string{"create", "delete", "password", "import", "copy", "comment", "access"} {
		if !msub[n] {
			t.Errorf("hosting db mysql is missing subcommand %q", n)
		}
	}

	access := findSub(mysql, "access")
	if access == nil {
		t.Fatal("hosting db mysql access command not registered")
	}
	asub := subNames(access)
	for _, n := range []string{"list", "grant", "revoke"} {
		if !asub[n] {
			t.Errorf("hosting db mysql access is missing subcommand %q", n)
		}
	}

	pgsql := findSub(db, "pgsql")
	if pgsql == nil {
		t.Fatal("hosting db pgsql command not registered")
	}
	psub := subNames(pgsql)
	for _, n := range []string{"create", "delete", "password"} {
		if !psub[n] {
			t.Errorf("hosting db pgsql is missing subcommand %q", n)
		}
	}
}

// TestHostingDBConfirmFlags checks that every destructive/mutating db command
// registers the --yes escape hatch that confirmed() reads.
func TestHostingDBConfirmFlags(t *testing.T) {
	for _, path := range [][]string{
		{"hosting", "db", "mysql", "delete"},
		{"hosting", "db", "mysql", "password"},
		{"hosting", "db", "mysql", "import"},
		{"hosting", "db", "mysql", "access", "revoke"},
		{"hosting", "db", "pgsql", "delete"},
		{"hosting", "db", "pgsql", "password"},
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
