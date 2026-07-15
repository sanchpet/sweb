package cmd

import "testing"

// Asserts the `domains persons` group and its subcommands are wired.
// It reuses findSub (hosting_db_test.go).
func TestDomainsPersonsCommandTree(t *testing.T) {
	domains := findSub(rootCmd, "domains")
	if domains == nil {
		t.Fatal("domains command not registered")
	}
	persons := findSub(domains, "persons")
	if persons == nil {
		t.Fatal("domains persons command not registered")
	}
	for _, n := range []string{
		"list", "info", "create-individual", "create-company",
	} {
		if findSub(persons, n) == nil {
			t.Errorf("domains persons is missing subcommand %q", n)
		}
	}
}

func TestPersonTypeLabel(t *testing.T) {
	cases := map[string]string{
		"f":     "individual",
		"ip":    "entrepreneur",
		"u":     "legal",
		"other": "other",
	}
	for in, want := range cases {
		if got := personTypeLabel(in); got != want {
			t.Errorf("personTypeLabel(%q) = %q, want %q", in, got, want)
		}
	}
}
